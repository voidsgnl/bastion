package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"path/filepath"

	"github.com/voidsgnl/bastion/server/internal/api"
	"github.com/voidsgnl/bastion/server/internal/engine"
	"github.com/voidsgnl/bastion/server/internal/store"
)

func main() {
	addr := flag.String("addr", "", "HTTP listen address")
	uiDir := flag.String("ui-dir", "../ui", "directory containing the static UI")
	engineMode := flag.String("engine", "fixture", "scan runner: fixture or nmap-lab")
	flag.Parse()
	listenAddress, err := resolveListenAddress(*addr, *engineMode)
	if err != nil {
		log.Fatal(err)
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
	handler := api.NewHandlerWithRunner(memory, runner)

	uiPath, err := filepath.Abs(*uiDir)
	if err != nil {
		log.Fatalf("resolve UI directory: %v", err)
	}
	handler.SetStatic(http.FileServer(http.Dir(uiPath)))

	log.Printf("Bastion server listening on %s", listenAddress)
	if err := http.ListenAndServe(listenAddress, handler); err != nil {
		log.Fatal(err)
	}
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
