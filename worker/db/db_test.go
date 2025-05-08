package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestJSONBMap_Value(t *testing.T) {
	tests := []struct {
		name     string
		input    JSONBMap
		expected driver.Value
		err      error
	}{
		{
			name:     "valid map",
			input:    JSONBMap{"key": "value", "num": 42},
			expected: []byte(`{"key":"value","num":42}`),
			err:      nil,
		},
		{
			name:     "empty map",
			input:    JSONBMap{},
			expected: []byte(`{}`),
			err:      nil,
		},
		{
			name:     "nil map",
			input:    nil,
			expected: nil,
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Value()
			if !errors.Is(err, tt.err) {
				t.Errorf("Value() error = %v, want %v", err, tt.err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Value() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestJSONBMap_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected JSONBMap
		err      error
	}{
		{
			name:     "valid JSON",
			input:    []byte(`{"key":"value","num":42}`),
			expected: JSONBMap{"key": "value", "num": float64(42)}, // JSON numbers are float64
			err:      nil,
		},
		{
			name:     "empty JSON",
			input:    []byte(`{}`),
			expected: JSONBMap{},
			err:      nil,
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			err:      nil,
		},
		{
			name:     "empty byte slice",
			input:    []byte{},
			expected: JSONBMap{},
			err:      nil,
		},
		{
			name:     "invalid JSON",
			input:    []byte(`{invalid}`),
			expected: nil,
			err:      &json.SyntaxError{},
		},
		{
			name:     "non-byte input",
			input:    "not a byte slice",
			expected: nil,
			err:      errors.New("scan source is not []byte"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var j JSONBMap
			err := j.Scan(tt.input)
			if tt.err != nil && (err == nil || !reflect.TypeOf(err).AssignableTo(reflect.TypeOf(tt.err))) {
				t.Errorf("Scan() error = %v, want type %T", err, tt.err)
			} else if tt.err == nil && !errors.Is(err, tt.err) {
				t.Errorf("Scan() error = %v, want %v", err, tt.err)
			}
			if !reflect.DeepEqual(j, tt.expected) {
				t.Errorf("Scan() result = %v, want %v", j, tt.expected)
			}
		})
	}
}
