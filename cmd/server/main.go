package main

import (
	"log"
	"the-unified-document-viewer/internal/api"
	"the-unified-document-viewer/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	server := api.NewServer(cfg)

	log.Printf("Starting server on %s:%s", cfg.Host, cfg.Port)
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
