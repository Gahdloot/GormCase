package fields

import (
	"database/sql/driver"
	"fmt"
)

// TextField represents a text field for long content
type TextField struct {
	*BaseField
}

// NewTextField creates a new TextField
func NewTextField(name string) *TextField {
	return &TextField{
		BaseField: NewBaseField(name, "TEXT"),
	}
}

// Validate validates the field value
func (f *TextField) Validate(value interface{}) error {
	if err := f.BaseField.Validate(value); err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	switch value.(type) {
	case string:
		return nil
	default:
		return fmt.Errorf("invalid type for TextField: %T", value)
	}
}

// Value implements driver.Valuer
func (f *TextField) Value() (driver.Value, error) {
	if f.value == nil {
		return nil, nil
	}
	return f.value, nil
}

// Scan implements sql.Scanner
func (f *TextField) Scan(value interface{}) error {
	if value == nil {
		f.value = nil
		return nil
	}

	switch v := value.(type) {
	case string:
		f.value = v
	case []byte:
		f.value = string(v)
	default:
		return fmt.Errorf("cannot scan %T into TextField", value)
	}

	return nil
}

// SetValue sets the field value
func (f *TextField) SetValue(value interface{}) error {
	if err := f.Validate(value); err != nil {
		return err
	}
	f.value = value
	return nil
}

// GetValue returns the field value
func (f *TextField) GetValue() interface{} {
	return f.value
}
