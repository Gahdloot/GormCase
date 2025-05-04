package fields

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// ManyToMany represents a many-to-many relationship field
type ManyToMany struct {
	*BaseField
	toModel     interface{}
	through     interface{}
	relatedName string
}

// NewManyToMany creates a new ManyToMany
func NewManyToMany(name string, toModel interface{}) *ManyToMany {
	return &ManyToMany{
		BaseField: NewBaseField(name, "MANY_TO_MANY"),
		toModel:   toModel,
	}
}

// Validate validates the field value
func (m *ManyToMany) Validate(value interface{}) error {
	if err := m.BaseField.Validate(value); err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	// Check if value is a slice of integers
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("invalid type for ManyToMany: %T", value)
	}

	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		var intVal int64
		switch v := elem.(type) {
		case int:
			intVal = int64(v)
		case int32:
			intVal = int64(v)
		case int64:
			intVal = v
		default:
			return fmt.Errorf("invalid element type in ManyToMany: %T", elem)
		}

		if intVal <= 0 {
			return fmt.Errorf("many-to-many value must be positive")
		}
	}

	return nil
}

// Value implements driver.Valuer
func (m *ManyToMany) Value() (driver.Value, error) {
	if m.value == nil {
		return nil, nil
	}
	return m.value, nil
}

// Scan implements sql.Scanner
func (m *ManyToMany) Scan(value interface{}) error {
	if value == nil {
		m.value = nil
		return nil
	}

	// Handle string representation of array
	if _, ok := value.(string); ok {
		// TODO: Implement parsing of string array
		return fmt.Errorf("string array parsing not implemented yet")
	}

	// Handle byte array representation
	if _, ok := value.([]byte); ok {
		// TODO: Implement parsing of byte array
		return fmt.Errorf("byte array parsing not implemented yet")
	}

	return fmt.Errorf("cannot scan %T into ManyToMany", value)
}

// SetThrough sets the through model for the relationship
func (m *ManyToMany) SetThrough(through interface{}) *ManyToMany {
	m.through = through
	return m
}

// SetRelatedName sets the related name for reverse lookup
func (m *ManyToMany) SetRelatedName(name string) *ManyToMany {
	m.relatedName = name
	return m
}

// GetToModel returns the target model
func (m *ManyToMany) GetToModel() interface{} {
	return m.toModel
}

// GetToModelName returns the name of the target model
func (m *ManyToMany) GetToModelName() string {
	return reflect.TypeOf(m.toModel).Name()
}

// GetThrough returns the through model
func (m *ManyToMany) GetThrough() interface{} {
	return m.through
}

// SetValue sets the field value
func (m *ManyToMany) SetValue(value interface{}) error {
	if err := m.Validate(value); err != nil {
		return err
	}
	m.value = value
	return nil
}

// GetValue returns the field value
func (m *ManyToMany) GetValue() interface{} {
	return m.value
}
