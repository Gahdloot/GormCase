package migration

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/Gahdloot/GormCase/pkg/models"
)

// Migration represents a database migration
type Migration interface {
	Up(ctx context.Context, tx *sql.Tx) error
	Down(ctx context.Context, tx *sql.Tx) error
	Name() string
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db         *models.DBManager
	migrations []Migration
}

// NewMigrationManager creates a new MigrationManager
func NewMigrationManager(db *models.DBManager) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// AddMigration adds a migration to the manager
func (m *MigrationManager) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// createMigrationsTable creates the migrations table if it doesn't exist
func (m *MigrationManager) createMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := m.db.Exec(ctx, query)
	return err
}

// getMigratedNames returns a list of applied migration names
func (m *MigrationManager) getMigratedNames(ctx context.Context) (map[string]bool, error) {
	query := "SELECT name FROM migrations"
	rows, err := m.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	migrated := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrated[name] = true
	}

	return migrated, rows.Err()
}

// Migrate runs all pending migrations
func (m *MigrationManager) Migrate(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get list of applied migrations
	migrated, err := m.getMigratedNames(ctx)
	if err != nil {
		return fmt.Errorf("failed to get migrated names: %v", err)
	}

	// Sort migrations by name
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Name() < m.migrations[j].Name()
	})

	// Run pending migrations
	for _, migration := range m.migrations {
		if !migrated[migration.Name()] {
			// Start transaction
			tx, err := m.db.Begin(ctx)
			if err != nil {
				return fmt.Errorf("failed to start transaction: %v", err)
			}

			// Run migration
			if err := migration.Up(ctx, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to run migration %s: %v", migration.Name(), err)
			}

			// Record migration
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO migrations (name, applied_at) VALUES ($1, $2)",
				migration.Name(),
				time.Now(),
			); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %s: %v", migration.Name(), err)
			}

			// Commit transaction
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit migration %s: %v", migration.Name(), err)
			}

			fmt.Printf("Applied migration: %s\n", migration.Name())
		}
	}

	return nil
}

// Rollback rolls back the last migration
func (m *MigrationManager) Rollback(ctx context.Context) error {
	// Get last applied migration
	var name string
	var appliedAt time.Time
	err := m.db.QueryRow(ctx,
		"SELECT name, applied_at FROM migrations ORDER BY applied_at DESC LIMIT 1",
	).Scan(&name, &appliedAt)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no migrations to roll back")
	}
	if err != nil {
		return fmt.Errorf("failed to get last migration: %v", err)
	}

	// Find migration
	var migration Migration
	for _, m := range m.migrations {
		if m.Name() == name {
			migration = m
			break
		}
	}
	if migration == nil {
		return fmt.Errorf("migration %s not found", name)
	}

	// Start transaction
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Run rollback
	if err := migration.Down(ctx, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to roll back migration %s: %v", name, err)
	}

	// Remove migration record
	if _, err := tx.ExecContext(ctx,
		"DELETE FROM migrations WHERE name = $1",
		name,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %s: %v", name, err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback %s: %v", name, err)
	}

	fmt.Printf("Rolled back migration: %s\n", name)
	return nil
}

// Reset rolls back all migrations
func (m *MigrationManager) Reset(ctx context.Context) error {
	for {
		err := m.Rollback(ctx)
		if err != nil {
			if err.Error() == "no migrations to roll back" {
				break
			}
			return err
		}
	}
	return nil
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Name   string
	Status string
}

// Status returns the status of all migrations
func (m *MigrationManager) Status(ctx context.Context) ([]MigrationStatus, error) {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(ctx); err != nil {
		return nil, fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get list of applied migrations
	migrated, err := m.getMigratedNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get migrated names: %v", err)
	}

	// Sort migrations by name
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Name() < m.migrations[j].Name()
	})

	// Create status list
	status := make([]MigrationStatus, len(m.migrations))
	for i, migration := range m.migrations {
		status[i] = MigrationStatus{
			Name:   migration.Name(),
			Status: "pending",
		}
		if migrated[migration.Name()] {
			status[i].Status = "applied"
		}
	}

	return status, nil
}
