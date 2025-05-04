package models

import (
	"context"
	"time"
)

// Model is the base interface that all models must implement
type Model interface {
	// TableName returns the name of the table for this model
	TableName() string

	// GetID returns the primary key value of the model
	GetID() interface{}

	// SetID sets the primary key value of the model
	SetID(id interface{})

	// CreatedAt returns the creation timestamp
	CreatedAt() time.Time

	// UpdatedAt returns the last update timestamp
	UpdatedAt() time.Time

	// Save persists the model to the database
	Save(ctx context.Context) error

	// Delete removes the model from the database
	Delete(ctx context.Context) error

	// Validate performs model validation
	Validate() error
}

// BaseModel provides a default implementation of the Model interface
type BaseModel struct {
	ID            interface{} `json:"id"`
	CreatedAtTime time.Time   `json:"created_at"`
	UpdatedAtTime time.Time   `json:"updated_at"`
}

// GetID implements Model.GetID
func (m *BaseModel) GetID() interface{} {
	return m.ID
}

// SetID implements Model.SetID
func (m *BaseModel) SetID(id interface{}) {
	m.ID = id
}

// CreatedAt implements Model.CreatedAt
func (m *BaseModel) CreatedAt() time.Time {
	return m.CreatedAtTime
}

// UpdatedAt implements Model.UpdatedAt
func (m *BaseModel) UpdatedAt() time.Time {
	return m.UpdatedAtTime
}

// Save implements Model.Save
func (m *BaseModel) Save(ctx context.Context) error {
	// To be implemented
	return nil
}

// Delete implements Model.Delete
func (m *BaseModel) Delete(ctx context.Context) error {
	// To be implemented
	return nil
}

// Validate implements Model.Validate
func (m *BaseModel) Validate() error {
	// To be implemented
	return nil
}
