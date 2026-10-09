package main

import (
	"fmt"
	"log"

	"qr-file-hub/backend/internal/config"
	"qr-file-hub/backend/internal/repository"
	"qr-file-hub/backend/internal/router"
	"qr-file-hub/backend/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	storage, err := service.NewStorageService(cfg)
	if err != nil {
		log.Fatalf("Storage service initialization failed: %v", err)
	}

	bridge := service.NewBridgeService(db)

	r := router.SetupRouter(cfg, db, storage, bridge)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
