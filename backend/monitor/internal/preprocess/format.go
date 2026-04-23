package preprocess

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// FieldMapping describes how to rename/transform fields in raw data.
type FieldMapping struct {
	From      string
	To        string
	Transform func(any) any // optional transform function
}

// FormatMapper is a Stage that maps/renames fields.
type FormatMapper struct {
	name     string
	mappings []FieldMapping
}

// NewFormatMapper creates a FormatMapper stage.
func NewFormatMapper(name string, mappings []FieldMapping) *FormatMapper {
	return &FormatMapper{name: name, mappings: mappings}
}

func (f *FormatMapper) Name() string { return f.name }

func (f *FormatMapper) Process(_ context.Context, data any) (any, error) {
	m, ok := data.(map[string]any)
	if !ok {
		return data, nil
	}
	result := make(map[string]any, len(m))
	// copy all existing fields
	for k, v := range m {
		result[k] = v
	}
	// apply mappings
	for _, fm := range f.mappings {
		val, exists := m[fm.From]
		if !exists {
			continue
		}
		if fm.Transform != nil {
			val = fm.Transform(val)
		}
		result[fm.To] = val
		if fm.From != fm.To {
			delete(result, fm.From)
		}
	}
	return result, nil
}

// NormalizeTimestamp converts various timestamp formats to RFC3339.
func NormalizeTimestamp(v any) any {
	s, ok := v.(string)
	if !ok {
		return v
	}
	formats := []string{
		time.RFC3339,
		time.RFC1123,
		time.RFC1123Z,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"Mon, 02 Jan 2006 15:04:05 MST",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, strings.TrimSpace(s)); err == nil {
			return t.UTC().Format(time.RFC3339)
		}
	}
	return v
}

// UnitConvert returns a transform that multiplies by factor.
func UnitConvert(factor float64) func(any) any {
	return func(v any) any {
		switch n := v.(type) {
		case float64:
			return n * factor
		case int:
			return float64(n) * factor
		default:
			return v
		}
	}
}

// StringClean returns a transform that trims and lowercases.
func StringClean(v any) any {
	s, ok := v.(string)
	if !ok {
		return v
	}
	return strings.TrimSpace(strings.ToLower(s))
}

// ValidateRequired is a Stage that checks required fields exist.
type ValidateRequired struct {
	fields []string
}

// NewValidateRequired creates a validation stage.
func NewValidateRequired(fields ...string) *ValidateRequired {
	return &ValidateRequired{fields: fields}
}

func (v *ValidateRequired) Name() string { return "validate-required" }

func (v *ValidateRequired) Process(_ context.Context, data any) (any, error) {
	m, ok := data.(map[string]any)
	if !ok {
		return data, nil
	}
	for _, f := range v.fields {
		if _, exists := m[f]; !exists {
			return nil, fmt.Errorf("missing required field: %s", f)
		}
	}
	return data, nil
}
