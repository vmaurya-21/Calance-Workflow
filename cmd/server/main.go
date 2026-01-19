package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/vmaurya-21/Calance-Workflow/internal/config"
	"github.com/vmaurya-21/Calance-Workflow/internal/database"
	"github.com/vmaurya-21/Calance-Workflow/internal/logger"
	"github.com/vmaurya-21/Calance-Workflow/internal/router"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// Use fmt for pre-logger errors
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.InitLogger(cfg.Log.Level, cfg.Log.Format)
	logger.Info().Msg("Configuration loaded successfully")

	// Debug: Log GitHub config
	logger.Debug().
		Str("client_id", cfg.GitHub.ClientID).
		Str("redirect_url", cfg.GitHub.RedirectURL).
		Msg("GitHub configuration")

	// Debug: Log database connection details
	logger.Debug().
		Str("host", cfg.Database.Host).
		Str("port", cfg.Database.Port).
		Str("user", cfg.Database.User).
		Str("dbname", cfg.Database.DBName).
		Msg("Database connection details")

	// Initialize database
	if err := database.InitDatabase(cfg); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database")
	}
	logger.Info().Msg("Database initialized successfully")

	// Set up router
	r := router.SetupRouter(database.GetDB(), cfg)
	logger.Info().Msg("Router configured successfully")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		logger.Info().Str("address", addr).Msg("Starting server")
		if err := r.Run(addr); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	<-quit
	logger.Info().Msg("Shutting down server...")

	// Close database connection
	if err := database.CloseDatabase(); err != nil {
		logger.Error().Err(err).Msg("Error closing database")
	}

	logger.Info().Msg("Server stopped")
}
