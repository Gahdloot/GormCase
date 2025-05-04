package fields

import (
	"database/sql/driver"
	"fmt"
)

// BooleanField represents a boolean field
type BooleanField struct {
	*BaseField
	defaultValue bool
}

// NewBooleanField creates a new BooleanField
func NewBooleanField(name string) *BooleanField {
	return &BooleanField{
		BaseField: NewBaseField(name, "BOOLEAN"),
	}
}

// Validate validates the field value
func (f *BooleanField) Validate(value interface{}) error {
	if err := f.BaseField.Validate(value); err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	switch value.(type) {
	case bool:
		return nil
	default:
		return fmt.Errorf("invalid type for BooleanField: %T", value)
	}
}

// Value implements driver.Valuer
func (f *BooleanField) Value() (driver.Value, error) {
	if f.value == nil {
		return f.defaultValue, nil
	}
	return f.value, nil
}

// Scan implements sql.Scanner
func (f *BooleanField) Scan(value interface{}) error {
	if value == nil {
		f.value = f.defaultValue
		return nil
	}

	switch v := value.(type) {
	case bool:
		f.value = v
	case []byte:
		f.value = string(v) == "true"
	case string:
		f.value = v == "true"
	default:
		return fmt.Errorf("cannot scan %T into BooleanField", value)
	}

	return nil
}

// SetDefault sets the default value
func (f *BooleanField) SetDefault(value bool) *BooleanField {
	f.defaultValue = value
	return f
}

// SetValue sets the field value
func (f *BooleanField) SetValue(value interface{}) error {
	if err := f.Validate(value); err != nil {
		return err
	}
	f.value = value
	return nil
}

// GetValue returns the field value
func (f *BooleanField) GetValue() interface{} {
	if f.value == nil {
		return f.defaultValue
	}
	return f.value
}
