package database

import (
	"fmt"
	"log"

	"github.com/vmaurya-21/Calance-Workflow/internal/config"
	"github.com/vmaurya-21/Calance-Workflow/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase(cfg *config.Config) error {
	var err error

	dsn := cfg.GetDatabaseDSN()

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully")

	// Run migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// runMigrations runs database migrations
func runMigrations() error {
	log.Println("Running database migrations...")

	// GORM AutoMigrate for development convenience
	// For production, use golang-migrate CLI tool with SQL migration files in db/migrations/
	// Example: migrate -path db/migrations -database "postgresql://user:pass@host:port/dbname?sslmode=disable" up
	err := DB.AutoMigrate(
		&models.User{},
		&models.Token{},
		// Add other models here as needed
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
