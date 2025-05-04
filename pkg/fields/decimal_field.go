package fields

import (
	"database/sql/driver"
	"fmt"
	"math"
)

// DecimalField represents a decimal field with precision and scale
type DecimalField struct {
	*BaseField
	maxDigits     int
	decimalPlaces int
	maxValue      *float64
	minValue      *float64
}

// NewDecimalField creates a new DecimalField
func NewDecimalField(name string, maxDigits, decimalPlaces int) *DecimalField {
	return &DecimalField{
		BaseField:     NewBaseField(name, fmt.Sprintf("DECIMAL(%d,%d)", maxDigits, decimalPlaces)),
		maxDigits:     maxDigits,
		decimalPlaces: decimalPlaces,
	}
}

// Validate validates the field value
func (f *DecimalField) Validate(value interface{}) error {
	if err := f.BaseField.Validate(value); err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	var floatVal float64
	switch v := value.(type) {
	case float64:
		floatVal = v
	case float32:
		floatVal = float64(v)
	case int:
		floatVal = float64(v)
	case int64:
		floatVal = float64(v)
	default:
		return fmt.Errorf("invalid type for DecimalField: %T", value)
	}

	// Check precision
	absVal := math.Abs(floatVal)
	if absVal >= math.Pow10(f.maxDigits-f.decimalPlaces) {
		return fmt.Errorf("value exceeds maximum precision of %d digits", f.maxDigits)
	}

	// Check scale
	scale := 0
	if floatVal != 0 {
		scale = int(math.Ceil(math.Log10(math.Abs(floatVal))))
	}
	if scale > f.decimalPlaces {
		return fmt.Errorf("value exceeds maximum scale of %d decimal places", f.decimalPlaces)
	}

	if f.minValue != nil && floatVal < *f.minValue {
		return fmt.Errorf("value %f is less than minimum %f", floatVal, *f.minValue)
	}

	if f.maxValue != nil && floatVal > *f.maxValue {
		return fmt.Errorf("value %f is greater than maximum %f", floatVal, *f.maxValue)
	}

	return nil
}

// Value implements driver.Valuer
func (f *DecimalField) Value() (driver.Value, error) {
	if f.value == nil {
		return nil, nil
	}
	return f.value, nil
}

// Scan implements sql.Scanner
func (f *DecimalField) Scan(value interface{}) error {
	if value == nil {
		f.value = nil
		return nil
	}

	var floatVal float64
	switch v := value.(type) {
	case float64:
		floatVal = v
	case float32:
		floatVal = float64(v)
	case int64:
		floatVal = float64(v)
	case int:
		floatVal = float64(v)
	case string:
		_, err := fmt.Sscanf(v, "%f", &floatVal)
		if err != nil {
			return fmt.Errorf("cannot parse string into DecimalField: %v", err)
		}
	default:
		return fmt.Errorf("cannot scan %T into DecimalField", value)
	}

	f.value = floatVal
	return nil
}

// SetMaxValue sets the maximum allowed value
func (f *DecimalField) SetMaxValue(max float64) *DecimalField {
	f.maxValue = &max
	return f
}

// SetMinValue sets the minimum allowed value
func (f *DecimalField) SetMinValue(min float64) *DecimalField {
	f.minValue = &min
	return f
}

// SetValue sets the field value
func (f *DecimalField) SetValue(value interface{}) error {
	if err := f.Validate(value); err != nil {
		return err
	}
	f.value = value
	return nil
}

// GetValue returns the field value
func (f *DecimalField) GetValue() interface{} {
	return f.value
}
