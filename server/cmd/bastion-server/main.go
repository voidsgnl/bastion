package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/voidsgnl/bastion/server/internal/api"
	"github.com/voidsgnl/bastion/server/internal/store"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	uiDir := flag.String("ui-dir", "../ui", "directory containing the static UI")
	flag.Parse()

	memory := store.NewMemory()
	handler := api.NewHandler(memory)

	uiPath, err := filepath.Abs(*uiDir)
	if err != nil {
		log.Fatalf("resolve UI directory: %v", err)
	}
	handler.SetStatic(http.FileServer(http.Dir(uiPath)))

	log.Printf("Bastion server listening on %s", *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal(err)
	}
}
