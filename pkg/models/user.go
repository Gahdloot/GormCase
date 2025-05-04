package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gahdloot/GormCase/pkg/fields"
	"github.com/pkg/errors"
)

// User represents a user in the system
type User struct {
	*BaseModel
	Username   *fields.CharField
	Email      *fields.CharField
	Password   *fields.CharField
	FirstName  *fields.CharField
	LastName   *fields.CharField
	IsActive   bool
	LastLogin  time.Time
	DateJoined time.Time
}

// NewUser creates a new User instance
func NewUser() *User {
	return &User{
		BaseModel:  &BaseModel{},
		Username:   fields.NewCharField("username", 150),
		Email:      fields.NewCharField("email", 254),
		Password:   fields.NewCharField("password", 128),
		FirstName:  fields.NewCharField("first_name", 150),
		LastName:   fields.NewCharField("last_name", 150),
		IsActive:   true,
		DateJoined: time.Now(),
	}
}

// TableName implements Model.TableName
func (u *User) TableName() string {
	return "users"
}

// Validate implements Model.Validate
func (u *User) Validate() error {
	// Validate required fields
	if err := u.Username.Validate(u.Username.GetValue()); err != nil {
		return err
	}
	if err := u.Email.Validate(u.Email.GetValue()); err != nil {
		return err
	}
	if err := u.Password.Validate(u.Password.GetValue()); err != nil {
		return err
	}

	// Add custom validation logic here
	return nil
}

// Save implements Model.Save
func (u *User) Save(ctx context.Context) error {
	if err := u.Validate(); err != nil {
		return err
	}

	// Get database connection
	db := GetDBManager()

	// Prepare fields and values for insert/update
	fields := []string{
		"username",
		"email",
		"password",
		"first_name",
		"last_name",
		"is_active",
		"last_login",
		"date_joined",
	}
	values := []interface{}{
		u.Username.GetValue(),
		u.Email.GetValue(),
		u.Password.GetValue(),
		u.FirstName.GetValue(),
		u.LastName.GetValue(),
		u.IsActive,
		u.LastLogin,
		u.DateJoined,
	}

	// If ID exists, update; otherwise insert
	if u.GetID() != nil {
		// Update
		query := "UPDATE users SET "
		placeholders := make([]string, len(fields))
		for i, field := range fields {
			placeholders[i] = fmt.Sprintf("%s = $%d", field, i+1)
		}
		query += strings.Join(placeholders, ", ")
		query += " WHERE id = $" + fmt.Sprintf("%d", len(fields)+1)

		values = append(values, u.GetID())
		_, err := db.Exec(ctx, query, values...)
		return err
	} else {
		// Insert
		placeholders := make([]string, len(fields))
		for i := range fields {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}
		query := fmt.Sprintf(
			"INSERT INTO users (%s) VALUES (%s) RETURNING id",
			strings.Join(fields, ", "),
			strings.Join(placeholders, ", "),
		)

		var id int
		err := db.QueryRow(ctx, query, values...).Scan(&id)
		if err != nil {
			return err
		}
		u.SetID(id)
		return nil
	}
}

// Delete implements Model.Delete
func (u *User) Delete(ctx context.Context) error {
	if u.GetID() == nil {
		return errors.New("cannot delete user without ID")
	}

	db := GetDBManager()
	_, err := db.Exec(ctx, "DELETE FROM users WHERE id = $1", u.GetID())
	return err
}

// FindByID finds a user by ID
func (u *User) FindByID(ctx context.Context, id int) error {
	db := GetDBManager()
	query := "SELECT id, username, email, password, first_name, last_name, is_active, last_login, date_joined FROM users WHERE id = $1"
	return db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.IsActive,
		&u.LastLogin,
		&u.DateJoined,
	)
}

// FindByUsername finds a user by username
func (u *User) FindByUsername(ctx context.Context, username string) error {
	db := GetDBManager()
	query := "SELECT id, username, email, password, first_name, last_name, is_active, last_login, date_joined FROM users WHERE username = $1"
	return db.QueryRow(ctx, query, username).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.IsActive,
		&u.LastLogin,
		&u.DateJoined,
	)
}
