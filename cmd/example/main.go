package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mac/gorm-orm/pkg/migration"
	"github.com/mac/gorm-orm/pkg/migration/migrations"
	"github.com/mac/gorm-orm/pkg/models"
)

func main() {
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
	migrationManager := migration.NewMigrationManager(db)
	migrationManager.AddMigration(&migrations.Migration001CreateUsersTable{})

	// Run migrations
	ctx := context.Background()
	if err := migrationManager.Migrate(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create a new user
	user := models.NewUser()
	user.Username.SetValue("john_doe")
	user.Email.SetValue("john@example.com")
	user.Password.SetValue("hashed_password")
	user.FirstName.SetValue("John")
	user.LastName.SetValue("Doe")

	// Save the user
	if err := user.Save(ctx); err != nil {
		log.Fatalf("Failed to save user: %v", err)
	}

	fmt.Printf("Created user with ID: %v\n", user.GetID())

	// Find the user by ID
	foundUser := models.NewUser()
	if err := foundUser.FindByID(ctx, user.GetID().(int)); err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}

	fmt.Printf("Found user: %+v\n", foundUser)

	// Update the user
	foundUser.FirstName.SetValue("Johnny")
	if err := foundUser.Save(ctx); err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

	fmt.Println("Updated user successfully")

	// Delete the user
	if err := foundUser.Delete(ctx); err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}

	fmt.Println("Deleted user successfully")

	// Rollback migrations (for cleanup)
	if err := migrationManager.Reset(ctx); err != nil {
		log.Fatalf("Failed to rollback migrations: %v", err)
	}
}
