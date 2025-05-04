package fields

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// DateTimeField represents a datetime field
type DateTimeField struct {
	*BaseField
	autoNow     bool
	autoNowAdd  bool
	defaultTime *time.Time
}

// NewDateTimeField creates a new DateTimeField
func NewDateTimeField(name string) *DateTimeField {
	return &DateTimeField{
		BaseField: NewBaseField(name, "TIMESTAMP WITH TIME ZONE"),
	}
}

// Validate validates the field value
func (f *DateTimeField) Validate(value interface{}) error {
	if err := f.BaseField.Validate(value); err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	switch value.(type) {
	case time.Time:
		return nil
	default:
		return fmt.Errorf("invalid type for DateTimeField: %T", value)
	}
}

// Value implements driver.Valuer
func (f *DateTimeField) Value() (driver.Value, error) {
	if f.value == nil {
		if f.autoNow || f.autoNowAdd {
			return time.Now(), nil
		}
		if f.defaultTime != nil {
			return f.defaultTime, nil
		}
		return nil, nil
	}
	return f.value, nil
}

// Scan implements sql.Scanner
func (f *DateTimeField) Scan(value interface{}) error {
	if value == nil {
		f.value = nil
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		f.value = v
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return fmt.Errorf("cannot parse time string: %v", err)
		}
		f.value = t
	default:
		return fmt.Errorf("cannot scan %T into DateTimeField", value)
	}

	return nil
}

// SetAutoNow sets the field to automatically update to current time on save
func (f *DateTimeField) SetAutoNow() *DateTimeField {
	f.autoNow = true
	return f
}

// SetAutoNowAdd sets the field to automatically set to current time on creation
func (f *DateTimeField) SetAutoNowAdd() *DateTimeField {
	f.autoNowAdd = true
	return f
}

// SetDefault sets the default value
func (f *DateTimeField) SetDefault(t time.Time) *DateTimeField {
	f.defaultTime = &t
	return f
}

// SetValue sets the field value
func (f *DateTimeField) SetValue(value interface{}) error {
	if err := f.Validate(value); err != nil {
		return err
	}
	f.value = value
	return nil
}

// GetValue returns the field value
func (f *DateTimeField) GetValue() interface{} {
	if f.value == nil {
		if f.autoNow || f.autoNowAdd {
			return time.Now()
		}
		if f.defaultTime != nil {
			return f.defaultTime
		}
		return nil
	}
	return f.value
}
