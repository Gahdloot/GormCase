package fields

import (
	"database/sql/driver"
	"fmt"
)

// IntegerField represents an integer field
type IntegerField struct {
	*BaseField
	maxValue *int64
	minValue *int64
}

// NewIntegerField creates a new IntegerField
func NewIntegerField(name string) *IntegerField {
	return &IntegerField{
		BaseField: NewBaseField(name, "INTEGER"),
	}
}

// Validate validates the field value
func (f *IntegerField) Validate(value interface{}) error {
	if err := f.BaseField.Validate(value); err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	var intVal int64
	switch v := value.(type) {
	case int:
		intVal = int64(v)
	case int32:
		intVal = int64(v)
	case int64:
		intVal = v
	default:
		return fmt.Errorf("invalid type for IntegerField: %T", value)
	}

	if f.minValue != nil && intVal < *f.minValue {
		return fmt.Errorf("value %d is less than minimum %d", intVal, *f.minValue)
	}

	if f.maxValue != nil && intVal > *f.maxValue {
		return fmt.Errorf("value %d is greater than maximum %d", intVal, *f.maxValue)
	}

	return nil
}

// Value implements driver.Valuer
func (f *IntegerField) Value() (driver.Value, error) {
	if f.value == nil {
		return nil, nil
	}
	return f.value, nil
}

// Scan implements sql.Scanner
func (f *IntegerField) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into IntegerField", value)
	}

	f.value = intVal
	return nil
}

// SetMaxValue sets the maximum allowed value
func (f *IntegerField) SetMaxValue(max int64) *IntegerField {
	f.maxValue = &max
	return f
}

// SetMinValue sets the minimum allowed value
func (f *IntegerField) SetMinValue(min int64) *IntegerField {
	f.minValue = &min
	return f
}

// SetValue sets the field value
func (f *IntegerField) SetValue(value interface{}) error {
	if err := f.Validate(value); err != nil {
		return err
	}
	f.value = value
	return nil
}

// GetValue returns the field value
func (f *IntegerField) GetValue() interface{} {
	return f.value
}
