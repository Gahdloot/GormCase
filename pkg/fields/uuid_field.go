package fields

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// UUIDField represents a UUID field
type UUIDField struct {
	*BaseField
	autoGenerate bool
}

// NewUUIDField creates a new UUIDField
func NewUUIDField(name string) *UUIDField {
	return &UUIDField{
		BaseField: NewBaseField(name, "UUID"),
	}
}

// Validate validates the field value
func (f *UUIDField) Validate(value interface{}) error {
	if err := f.BaseField.Validate(value); err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case uuid.UUID:
		return nil
	case string:
		_, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("invalid UUID string: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("invalid type for UUIDField: %T", value)
	}
}

// Value implements driver.Valuer
func (f *UUIDField) Value() (driver.Value, error) {
	if f.value == nil {
		if f.autoGenerate {
			newUUID := uuid.New()
			f.value = newUUID
			return newUUID.String(), nil
		}
		return nil, nil
	}

	switch v := f.value.(type) {
	case uuid.UUID:
		return v.String(), nil
	case string:
		return v, nil
	default:
		return nil, fmt.Errorf("invalid UUID value type: %T", f.value)
	}
}

// Scan implements sql.Scanner
func (f *UUIDField) Scan(value interface{}) error {
	if value == nil {
		f.value = nil
		return nil
	}

	switch v := value.(type) {
	case uuid.UUID:
		f.value = v
	case string:
		// Remove any surrounding quotes
		v = strings.Trim(v, `"`)
		parsedUUID, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("cannot parse UUID string: %v", err)
		}
		f.value = parsedUUID
	case []byte:
		parsedUUID, err := uuid.ParseBytes(v)
		if err != nil {
			return fmt.Errorf("cannot parse UUID bytes: %v", err)
		}
		f.value = parsedUUID
	default:
		return fmt.Errorf("cannot scan %T into UUIDField", value)
	}

	return nil
}

// SetAutoGenerate sets the field to automatically generate a UUID on save
func (f *UUIDField) SetAutoGenerate() *UUIDField {
	f.autoGenerate = true
	return f
}

// SetValue sets the field value
func (f *UUIDField) SetValue(value interface{}) error {
	if err := f.Validate(value); err != nil {
		return err
	}
	f.value = value
	return nil
}

// GetValue returns the field value
func (f *UUIDField) GetValue() interface{} {
	if f.value == nil && f.autoGenerate {
		return uuid.New()
	}
	return f.value
}

// GetUUID returns the value as a uuid.UUID
func (f *UUIDField) GetUUID() (uuid.UUID, error) {
	if f.value == nil {
		if f.autoGenerate {
			return uuid.New(), nil
		}
		return uuid.Nil, fmt.Errorf("UUID value is nil")
	}

	switch v := f.value.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		return uuid.Parse(v)
	default:
		return uuid.Nil, fmt.Errorf("invalid UUID value type: %T", f.value)
	}
}
