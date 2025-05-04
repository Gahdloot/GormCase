package fields

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONField represents a JSON field
type JSONField struct {
	*BaseField
}

// NewJSONField creates a new JSONField
func NewJSONField(name string) *JSONField {
	return &JSONField{
		BaseField: NewBaseField(name, "JSONB"),
	}
}

// Validate validates the field value
func (f *JSONField) Validate(value interface{}) error {
	if err := f.BaseField.Validate(value); err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	// Try to marshal the value to JSON to validate it
	_, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("invalid JSON value: %v", err)
	}

	return nil
}

// Value implements driver.Valuer
func (f *JSONField) Value() (driver.Value, error) {
	if f.value == nil {
		return nil, nil
	}

	// Marshal the value to JSON
	jsonBytes, err := json.Marshal(f.value)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	return string(jsonBytes), nil
}

// Scan implements sql.Scanner
func (f *JSONField) Scan(value interface{}) error {
	if value == nil {
		f.value = nil
		return nil
	}

	var jsonBytes []byte
	switch v := value.(type) {
	case string:
		jsonBytes = []byte(v)
	case []byte:
		jsonBytes = v
	default:
		return fmt.Errorf("cannot scan %T into JSONField", value)
	}

	// Try to unmarshal into a generic interface{}
	var result interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	f.value = result
	return nil
}

// SetValue sets the field value
func (f *JSONField) SetValue(value interface{}) error {
	if err := f.Validate(value); err != nil {
		return err
	}
	f.value = value
	return nil
}

// GetValue returns the field value
func (f *JSONField) GetValue() interface{} {
	return f.value
}

// GetMap returns the value as a map[string]interface{}
func (f *JSONField) GetMap() (map[string]interface{}, error) {
	if f.value == nil {
		return nil, nil
	}

	jsonBytes, err := json.Marshal(f.value)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return result, nil
}

// GetSlice returns the value as a []interface{}
func (f *JSONField) GetSlice() ([]interface{}, error) {
	if f.value == nil {
		return nil, nil
	}

	jsonBytes, err := json.Marshal(f.value)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	var result []interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return result, nil
}
