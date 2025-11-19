package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vmaurya-21/Calance-Workflow/internal/config"
	// "github.com/vmaurya-21/Calance-Workflow/internal/database"
	"github.com/vmaurya-21/Calance-Workflow/internal/router"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("Configuration loaded successfully")

	// Debug: Log GitHub config
	log.Printf("DEBUG - GitHub ClientID: %s", cfg.GitHub.ClientID)
	log.Printf("DEBUG - GitHub RedirectURL: %s", cfg.GitHub.RedirectURL)

	// TODO: Initialize database
	// if err := database.InitDatabase(cfg); err != nil {
	// 	log.Fatalf("Failed to initialize database: %v", err)
	// }
	// log.Println("Database initialized successfully")

	// Set up router
	r := router.SetupRouter(nil, cfg)
	log.Println("Router configured successfully")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		log.Printf("Starting server on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// TODO: Close database connection
	// if err := database.CloseDatabase(); err != nil {
	// 	log.Printf("Error closing database: %v", err)
	// }

	log.Println("Server stopped")
}
