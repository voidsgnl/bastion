package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/voidsgnl/bastion/server/internal/api"
	"github.com/voidsgnl/bastion/server/internal/auth"
	"github.com/voidsgnl/bastion/server/internal/engine"
	"github.com/voidsgnl/bastion/server/internal/store"
)

func main() {
	addr := flag.String("addr", "", "HTTP listen address")
	uiDir := flag.String("ui-dir", "../ui", "directory containing the static UI")
	engineMode := flag.String("engine", "fixture", "scan runner: fixture or nmap-lab")
	storeMode := flag.String("store", "memory", "data store: memory or postgres")
	databaseURL := flag.String("database-url", os.Getenv("BASTION_DATABASE_URL"), "PostgreSQL connection URL")
	adminUser := flag.String("auth-user", envOr("BASTION_ADMIN_USER", "admin"), "bootstrap operator username")
	tlsCert := flag.String("tls-cert", "", "TLS certificate for non-loopback listening")
	tlsKey := flag.String("tls-key", "", "TLS private key for non-loopback listening")
	flag.Parse()
	configuredAddr := *addr
	if configuredAddr == "" && *engineMode != "nmap-lab" {
		configuredAddr = "127.0.0.1:8080"
	}
	listenAddress, err := resolveListenAddress(configuredAddr, *engineMode)
	if err != nil {
		log.Fatal(err)
	}
	if (*tlsCert == "") != (*tlsKey == "") {
		log.Fatal("-tls-cert and -tls-key must be provided together")
	}
	if !loopbackListenAddress(listenAddress) && *tlsCert == "" {
		log.Fatal("non-loopback listening requires -tls-cert and -tls-key")
	}
	adminPassword := os.Getenv("BASTION_ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("set BASTION_ADMIN_PASSWORD (minimum 12 characters)")
	}
	authManager, err := auth.NewManager(*adminUser, adminPassword, 12*time.Hour)
	if err != nil {
		log.Fatalf("configure authentication: %v", err)
	}

	memory := store.NewMemory()
	var runner engine.Runner
	switch *engineMode {
	case "fixture":
		runner = engine.FixtureRunner{}
	case "nmap-lab":
		runner = engine.NewNmapRunner()
	default:
		log.Fatalf("unsupported engine %q (choose fixture or nmap-lab)", *engineMode)
	}

	var repository store.Repository = memory
	var postgres *store.Postgres
	if *storeMode == "postgres" {
		if *databaseURL == "" {
			log.Fatal("PostgreSQL store requires -database-url or BASTION_DATABASE_URL")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		postgres, err = store.NewPostgres(ctx, *databaseURL)
		cancel()
		if err != nil {
			log.Fatalf("initialize PostgreSQL store: %v", err)
		}
		defer func() {
			if err := postgres.Close(); err != nil {
				log.Printf("close PostgreSQL store: %v", err)
			}
		}()
		repository = postgres
	} else if *storeMode != "memory" {
		log.Fatalf("unsupported store %q (choose memory or postgres)", *storeMode)
	}
	apiHandler := api.NewHandlerWithRepository(repository, runner)

	uiPath, err := filepath.Abs(*uiDir)
	if err != nil {
		log.Fatalf("resolve UI directory: %v", err)
	}
	apiHandler.SetStatic(http.FileServer(http.Dir(uiPath)))
	handler := authManager.Middleware(apiHandler)

	log.Printf("Bastion server listening on %s", listenAddress)
	server := &http.Server{
		Addr: listenAddress, Handler: handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
		TLSConfig:         &tls.Config{MinVersion: tls.VersionTLS12},
	}
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	serveErrors := make(chan error, 1)
	go func() {
		if *tlsCert != "" {
			serveErrors <- server.ServeTLS(listener, *tlsCert, *tlsKey)
			return
		}
		serveErrors <- server.Serve(listener)
	}()

	stop, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()
	select {
	case err := <-serveErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve HTTP: %v", err)
		}
	case <-stop.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown: %v", err)
			_ = server.Close()
		}
	}
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func loopbackListenAddress(address string) bool {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return false
	}
	if strings.EqualFold(host, "localhost") {
		return true
	}
	parsed, err := netip.ParseAddr(host)
	return err == nil && parsed.IsLoopback()
}

func resolveListenAddress(configured, engineMode string) (string, error) {
	if configured == "" {
		if engineMode == "nmap-lab" {
			return "127.0.0.1:8080", nil
		}
		return ":8080", nil
	}
	if engineMode != "nmap-lab" {
		return configured, nil
	}
	host, _, err := net.SplitHostPort(configured)
	if err != nil {
		return "", err
	}
	address, err := netip.ParseAddr(host)
	if err != nil || !address.IsLoopback() {
		return "", fmt.Errorf("nmap-lab mode requires a loopback HTTP listen address")
	}
	return configured, nil
}
