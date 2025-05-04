package fields

import (
	"database/sql/driver"
	"errors"
)

// Field is the base interface that all fields must implement
type Field interface {
	// Name returns the name of the field
	Name() string

	// SetName sets the name of the field
	SetName(name string)

	// Type returns the SQL type of the field
	Type() string

	// IsNullable returns whether the field can be null
	IsNullable() bool

	// IsPrimaryKey returns whether the field is a primary key
	IsPrimaryKey() bool

	// IsUnique returns whether the field must be unique
	IsUnique() bool

	// Default returns the default value of the field
	Default() interface{}

	// Validate validates the field value
	Validate(value interface{}) error

	// Value returns the field value for database storage
	Value() (driver.Value, error)

	// Scan scans the field value from database storage
	Scan(value interface{}) error
}

// BaseField provides a default implementation of the Field interface
type BaseField struct {
	name         string
	fieldType    string
	nullable     bool
	primaryKey   bool
	unique       bool
	defaultValue interface{}
	value        interface{}
}

// NewBaseField creates a new BaseField
func NewBaseField(name string, fieldType string) *BaseField {
	return &BaseField{
		name:      name,
		fieldType: fieldType,
	}
}

// Name implements Field.Name
func (f *BaseField) Name() string {
	return f.name
}

// SetName implements Field.SetName
func (f *BaseField) SetName(name string) {
	f.name = name
}

// Type implements Field.Type
func (f *BaseField) Type() string {
	return f.fieldType
}

// IsNullable implements Field.IsNullable
func (f *BaseField) IsNullable() bool {
	return f.nullable
}

// IsPrimaryKey implements Field.IsPrimaryKey
func (f *BaseField) IsPrimaryKey() bool {
	return f.primaryKey
}

// IsUnique implements Field.IsUnique
func (f *BaseField) IsUnique() bool {
	return f.unique
}

// Default implements Field.Default
func (f *BaseField) Default() interface{} {
	return f.defaultValue
}

// Validate implements Field.Validate
func (f *BaseField) Validate(value interface{}) error {
	if !f.nullable && value == nil {
		return errors.New("field cannot be null")
	}
	return nil
}

// Value implements Field.Value
func (f *BaseField) Value() (driver.Value, error) {
	if f.value == nil {
		return nil, nil
	}
	return f.value, nil
}

// Scan implements Field.Scan
func (f *BaseField) Scan(value interface{}) error {
	f.value = value
	return nil
}

// SetValue sets the field value
func (f *BaseField) SetValue(value interface{}) error {
	if err := f.Validate(value); err != nil {
		return err
	}
	f.value = value
	return nil
}

// GetValue returns the field value
func (f *BaseField) GetValue() interface{} {
	return f.value
}
