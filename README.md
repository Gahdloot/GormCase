# GORM-ORM

A Django-like ORM for Golang that is compatible with PostgreSQL.

## Features

- Model definitions with fields and relationships
- Query building (filter, exclude, annotate)
- Migrations
- Model validation
- Model methods and properties
- Transaction management
- Raw SQL support
- And much more...

## Installation

```bash
go get github.com/Gahdloot/GormCase
```

## Quick Start

1. First, define your model:

```go
package models

import (
    "github.com/Gahdloot/GormCase/pkg/fields"
    "github.com/Gahdloot/GormCase/pkg/models"
)

type User struct {
    *models.BaseModel
    Username    *fields.CharField
    Email       *fields.CharField
    Password    *fields.CharField
    FirstName   *fields.CharField
    LastName    *fields.CharField
    IsActive    bool
    LastLogin   time.Time
    DateJoined  time.Time
}

func NewUser() *User {
    return &User{
        BaseModel:  &models.BaseModel{},
        Username:   fields.NewCharField("username", 150),
        Email:      fields.NewCharField("email", 254),
        Password:   fields.NewCharField("password", 128),
        FirstName:  fields.NewCharField("first_name", 150),
        LastName:   fields.NewCharField("last_name", 150),
        IsActive:   true,
        DateJoined: time.Now(),
    }
}
```

2. Create a migration:

```go
package migrations

import (
    "context"
    "database/sql"
)

type Migration001CreateUsersTable struct{}

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
    `
    _, err := tx.ExecContext(ctx, query)
    return err
}

func (m *Migration001CreateUsersTable) Down(ctx context.Context, tx *sql.Tx) error {
    query := `DROP TABLE IF EXISTS users CASCADE;`
    _, err := tx.ExecContext(ctx, query)
    return err
}

func (m *Migration001CreateUsersTable) Name() string {
    return "001_create_users_table"
}
```

3. Use the ORM:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Gahdloot/GormCase/pkg/migration"
    "github.com/Gahdloot/GormCase/pkg/migration/migrations"
    "github.com/Gahdloot/GormCase/pkg/models"
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
        log.Fatal(err)
    }
    defer db.Close()

    // Run migrations
    migrationManager := migration.NewMigrationManager(db)
    migrationManager.AddMigration(&migrations.Migration001CreateUsersTable{})

    ctx := context.Background()
    if err := migrationManager.Migrate(ctx); err != nil {
        log.Fatal(err)
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
        log.Fatal(err)
    }

    fmt.Printf("Created user with ID: %v\n", user.GetID())
}
```

## Features in Detail

### Model Fields

The ORM supports various field types:

- CharField
- IntegerField (coming soon)
- BooleanField (coming soon)
- DateTimeField (coming soon)
- ForeignKey (coming soon)
- ManyToMany (coming soon)
- And more...

### Querying

```go
// Find by ID
user := models.NewUser()
err := user.FindByID(ctx, 1)

// Find by username
user := models.NewUser()
err := user.FindByUsername(ctx, "john_doe")

// Custom queries coming soon:
// users := models.User{}.Objects().Filter("age__gt", 18).OrderBy("-created_at").All()
```

### Migrations

```go
// Run migrations
migrationManager := migration.NewMigrationManager(db)
migrationManager.AddMigration(&migrations.Migration001CreateUsersTable{})
err := migrationManager.Migrate(ctx)

// Rollback last migration
err := migrationManager.Rollback(ctx)

// Reset all migrations
err := migrationManager.Reset(ctx)
```

### Validation

```go
// Built-in validation
user := models.NewUser()
err := user.Validate()

// Custom validation coming soon:
// func (u *User) Validate() error {
//     if err := u.BaseModel.Validate(); err != nil {
//         return err
//     }
//     // Custom validation logic
//     return nil
// }
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT
