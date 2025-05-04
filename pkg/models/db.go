package models

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// DBConfig holds the database configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DBManager manages database connections
type DBManager struct {
	db     *sql.DB
	config DBConfig
	mu     sync.RWMutex
}

var (
	instance *DBManager
	once     sync.Once
)

// NewDBManager creates a new DBManager instance
func NewDBManager(config DBConfig) (*DBManager, error) {
	var err error
	once.Do(func() {
		instance = &DBManager{
			config: config,
		}
		err = instance.connect()
	})
	return instance, err
}

// GetDBManager returns the singleton DBManager instance
func GetDBManager() *DBManager {
	if instance == nil {
		panic("DBManager not initialized. Call NewDBManager first.")
	}
	return instance
}

// connect establishes a connection to the database
func (m *DBManager) connect() error {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		m.config.Host,
		m.config.Port,
		m.config.User,
		m.config.Password,
		m.config.DBName,
		m.config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	m.mu.Lock()
	m.db = db
	m.mu.Unlock()

	return nil
}

// Close closes the database connection
func (m *DBManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// Begin starts a new transaction
func (m *DBManager) Begin(ctx context.Context) (*sql.Tx, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.db.BeginTx(ctx, nil)
}

// Exec executes a query without returning any rows
func (m *DBManager) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.db.ExecContext(ctx, query, args...)
}

// Query executes a query that returns rows
func (m *DBManager) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that returns a single row
func (m *DBManager) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.db.QueryRowContext(ctx, query, args...)
}

// Ping verifies the connection to the database is still alive
func (m *DBManager) Ping(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.db.PingContext(ctx)
}

// Stats returns database statistics
func (m *DBManager) Stats() sql.DBStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.db.Stats()
}
