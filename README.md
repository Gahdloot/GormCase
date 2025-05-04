# GormCase

A Django-inspired ORM for Go, built on top of GORM.

## Installation

```bash
go get github.com/Gahdloot/GormCase
```

## Features

- Django-style model definitions
- Automatic migration generation and management
- Field types with validation and constraints
- Relationship handling (One-to-Many, Many-to-Many)
- Transaction support
- Query building

## Field Types

### Basic Fields

#### CharField

```go
name := fields.NewCharField("name").SetMaxLength(100)
```

#### TextField

```go
description := fields.NewTextField("description")
```

#### IntegerField

```go
age := fields.NewIntegerField("age").
    SetMinValue(0).
    SetMaxValue(150)
```

#### DecimalField

```go
price := fields.NewDecimalField("price", 10, 2).  // 10 total digits, 2 decimal places
    SetMinValue(0)
```

#### BooleanField

```go
isActive := fields.NewBooleanField("is_active").
    SetDefault(true)
```

#### DateTimeField

```go
createdAt := fields.NewDateTimeField("created_at").
    SetAutoNowAdd()  // Set to current time on creation
updatedAt := fields.NewDateTimeField("updated_at").
    SetAutoNow()     // Update to current time on save
```

#### UUIDField

```go
id := fields.NewUUIDField("id").
    SetAutoGenerate()  // Generate UUID automatically
```

#### JSONField

```go
metadata := fields.NewJSONField("metadata")
// Get as map
data, err := metadata.GetMap()
// Get as slice
items, err := metadata.GetSlice()
```

### Relationship Fields

#### ForeignKey (One-to-Many)

```go
profile := fields.NewForeignKey("profile_id", &Profile{}).
    SetOnDelete("CASCADE").
    SetOnUpdate("CASCADE").
    SetRelatedName("user")
```

#### ManyToMany

```go
groups := fields.NewManyToMany("groups", &Group{}).
    SetThrough(&UserGroup{}).  // Optional through model
    SetRelatedName("users")
```

## Model Definition

```go
type User struct {
    *models.BaseModel
    ID        *fields.UUIDField
    Name      *fields.CharField
    Age       *fields.IntegerField
    IsActive  *fields.BooleanField
    CreatedAt *fields.DateTimeField
    Profile   *fields.ForeignKey
    Groups    *fields.ManyToMany
}

func NewUser() *User {
    return &User{
        BaseModel: models.NewBaseModel(),
        ID:        fields.NewUUIDField("id").SetAutoGenerate(),
        Name:      fields.NewCharField("name").SetMaxLength(100),
        Age:       fields.NewIntegerField("age").SetMinValue(0).SetMaxValue(150),
        IsActive:  fields.NewBooleanField("is_active").SetDefault(true),
        CreatedAt: fields.NewDateTimeField("created_at").SetAutoNowAdd(),
        Profile:   fields.NewForeignKey("profile_id", &Profile{}).SetOnDelete("CASCADE"),
        Groups:    fields.NewManyToMany("groups", &Group{}),
    }
}
```

## Migrations

### Generate Migrations

```bash
go run cmd/makemigrations/main.go --model=User
```

### Run Migrations

```bash
# Run all pending migrations
go run cmd/migrate/main.go up

# Rollback last migration
go run cmd/migrate/main.go down

# Reset all migrations
go run cmd/migrate/main.go reset

# Show migration status
go run cmd/migrate/main.go status
```

## Usage Example

```go
package main

import (
    "context"
    "log"

    "github.com/Gahdloot/GormCase/pkg/models"
    "github.com/Gahdloot/GormCase/pkg/migration"
)

func main() {
    // Initialize database
    db, err := models.NewDBManager(models.DBConfig{
        Host:     "localhost",
        Port:     5432,
        User:     "postgres",
        Password: "postgres",
        DBName:   "testdb",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create migration manager
    mgr := migration.NewMigrationManager(db)

    // Add migrations
    mgr.AddMigration(&CreateUserTable{})
    mgr.AddMigration(&CreateProfileTable{})

    // Run migrations
    ctx := context.Background()
    if err := mgr.Up(ctx); err != nil {
        log.Fatal(err)
    }

    // Create a new user
    user := models.NewUser()
    user.Name.SetValue("John Doe")
    user.Age.SetValue(30)
    user.IsActive.SetValue(true)

    // Save the user
    if err := user.Save(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## License

MIT
