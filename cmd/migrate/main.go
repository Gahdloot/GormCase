package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Gahdloot/GormCase/pkg/migration"
	"github.com/Gahdloot/GormCase/pkg/models"
)

func main() {
	// Parse command line arguments
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: migrate <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  up     - Run all pending migrations")
		fmt.Println("  down   - Rollback the last migration")
		fmt.Println("  reset  - Rollback all migrations")
		fmt.Println("  status - Show migration status")
		os.Exit(1)
	}

	// Initialize database connection
	dbConfig := models.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "gorm_orm",
		SSLMode:  "disable",
	}

	db, err := models.NewDBManager(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize migration manager
	manager := migration.NewMigrationManager(db)

	// Create context
	ctx := context.Background()

	// Execute command
	command := args[0]
	switch command {
	case "up":
		if err := manager.Migrate(ctx); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		if err := manager.Rollback(ctx); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		fmt.Println("Rolled back last migration")

	case "reset":
		if err := manager.Reset(ctx); err != nil {
			log.Fatalf("Failed to reset migrations: %v", err)
		}
		fmt.Println("Reset all migrations")

	case "status":
		status, err := manager.Status(ctx)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
		fmt.Println("Migration Status:")
		for _, s := range status {
			fmt.Printf("- %s: %s\n", s.Name, s.Status)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
