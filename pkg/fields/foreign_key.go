package fields

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// ForeignKey represents a foreign key field
type ForeignKey struct {
	*BaseField
	toModel     interface{}
	onDelete    string
	onUpdate    string
	relatedName string
}

// NewForeignKey creates a new ForeignKey
func NewForeignKey(name string, toModel interface{}) *ForeignKey {
	return &ForeignKey{
		BaseField: NewBaseField(name, "INTEGER"),
		toModel:   toModel,
	}
}

// Validate validates the field value
func (f *ForeignKey) Validate(value interface{}) error {
	if err := f.BaseField.Validate(value); err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	// Check if value is an integer
	var intVal int64
	switch v := value.(type) {
	case int:
		intVal = int64(v)
	case int32:
		intVal = int64(v)
	case int64:
		intVal = v
	default:
		return fmt.Errorf("invalid type for ForeignKey: %T", value)
	}

	// Check if value is positive
	if intVal <= 0 {
		return fmt.Errorf("foreign key value must be positive")
	}

	return nil
}

// Value implements driver.Valuer
func (f *ForeignKey) Value() (driver.Value, error) {
	if f.value == nil {
		return nil, nil
	}
	return f.value, nil
}

// Scan implements sql.Scanner
func (f *ForeignKey) Scan(value interface{}) error {
	if value == nil {
		f.value = nil
		return nil
	}

	var intVal int64
	switch v := value.(type) {
	case int64:
		intVal = v
	case int32:
		intVal = int64(v)
	case int:
		intVal = int64(v)
	default:
		return fmt.Errorf("cannot scan %T into ForeignKey", value)
	}

	f.value = intVal
	return nil
}

// SetOnDelete sets the ON DELETE behavior
func (f *ForeignKey) SetOnDelete(action string) *ForeignKey {
	switch action {
	case "CASCADE", "SET NULL", "SET DEFAULT", "RESTRICT", "NO ACTION":
		f.onDelete = action
	default:
		panic(fmt.Sprintf("invalid ON DELETE action: %s", action))
	}
	return f
}

// SetOnUpdate sets the ON UPDATE behavior
func (f *ForeignKey) SetOnUpdate(action string) *ForeignKey {
	switch action {
	case "CASCADE", "SET NULL", "SET DEFAULT", "RESTRICT", "NO ACTION":
		f.onUpdate = action
	default:
		panic(fmt.Sprintf("invalid ON UPDATE action: %s", action))
	}
	return f
}

// SetRelatedName sets the related name for reverse lookup
func (f *ForeignKey) SetRelatedName(name string) *ForeignKey {
	f.relatedName = name
	return f
}

// GetToModel returns the target model
func (f *ForeignKey) GetToModel() interface{} {
	return f.toModel
}

// GetToModelName returns the name of the target model
func (f *ForeignKey) GetToModelName() string {
	return reflect.TypeOf(f.toModel).Name()
}

// SetValue sets the field value
func (f *ForeignKey) SetValue(value interface{}) error {
	if err := f.Validate(value); err != nil {
		return err
	}
	f.value = value
	return nil
}

// GetValue returns the field value
func (f *ForeignKey) GetValue() interface{} {
	return f.value
}
