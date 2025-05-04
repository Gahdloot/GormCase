package migrations

import (
	"context"
	"database/sql"
)

// Migration001CreateUsersTable represents the first migration
type Migration001CreateUsersTable struct{}

// Up creates the users table
func (m *Migration001CreateUsersTable) Up(ctx context.Context, tx *sql.Tx) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(150) NOT NULL UNIQUE,
			email VARCHAR(254) NOT NULL UNIQUE,
			password VARCHAR(128) NOT NULL,
			first_name VARCHAR(150),
			last_name VARCHAR(150),
			is_active BOOLEAN NOT NULL DEFAULT true,
			last_login TIMESTAMP WITH TIME ZONE,
			date_joined TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX idx_users_username ON users(username);
		CREATE INDEX idx_users_email ON users(email);
		CREATE INDEX idx_users_is_active ON users(is_active);
	`

	_, err := tx.ExecContext(ctx, query)
	return err
}

// Down drops the users table
func (m *Migration001CreateUsersTable) Down(ctx context.Context, tx *sql.Tx) error {
	query := `
		DROP TABLE IF EXISTS users CASCADE;
	`

	_, err := tx.ExecContext(ctx, query)
	return err
}

// Name returns the name of the migration
func (m *Migration001CreateUsersTable) Name() string {
	return "001_create_users_table"
}
