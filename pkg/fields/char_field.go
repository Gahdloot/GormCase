package fields

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// CharField represents a character field in the database
type CharField struct {
	*BaseField
	maxLength int
}

// NewCharField creates a new CharField
func NewCharField(name string, maxLength int) *CharField {
	return &CharField{
		BaseField: NewBaseField(name, fmt.Sprintf("VARCHAR(%d)", maxLength)),
		maxLength: maxLength,
	}
}

// Validate implements Field.Validate
func (f *CharField) Validate(value interface{}) error {
	if err := f.BaseField.Validate(value); err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}

	if len(str) > f.maxLength {
		return fmt.Errorf("string length %d exceeds maximum length %d", len(str), f.maxLength)
	}

	return nil
}

// Value implements Field.Value
func (f *CharField) Value() (driver.Value, error) {
	if f.value == nil {
		return nil, nil
	}

	str, ok := f.value.(string)
	if !ok {
		return nil, errors.New("value must be a string")
	}

	return str, nil
}

// Scan implements Field.Scan
func (f *CharField) Scan(value interface{}) error {
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
		return fmt.Errorf("cannot scan %T into CharField", value)
	}

	return nil
}

// SetMaxLength sets the maximum length of the field
func (f *CharField) SetMaxLength(maxLength int) {
	f.maxLength = maxLength
	f.fieldType = fmt.Sprintf("VARCHAR(%d)", maxLength)
}

// GetMaxLength returns the maximum length of the field
func (f *CharField) GetMaxLength() int {
	return f.maxLength
}

// SetNullable sets whether the field can be null
func (f *CharField) SetNullable(nullable bool) {
	f.nullable = nullable
}

// SetPrimaryKey sets whether the field is a primary key
func (f *CharField) SetPrimaryKey(primaryKey bool) {
	f.primaryKey = primaryKey
}

// SetUnique sets whether the field must be unique
func (f *CharField) SetUnique(unique bool) {
	f.unique = unique
}

// SetDefault sets the default value of the field
func (f *CharField) SetDefault(defaultValue interface{}) error {
	if err := f.Validate(defaultValue); err != nil {
		return err
	}
	f.defaultValue = defaultValue
	return nil
}
