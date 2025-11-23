package ui

import (
	"strings"
	"testing"
)

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "nil value",
			value:    nil,
			expected: "(null)",
		},
		{
			name:     "string value",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "integer value",
			value:    42,
			expected: "42",
		},
		{
			name:     "boolean value",
			value:    true,
			expected: "true",
		},
		{
			name:     "empty JSON object",
			value:    map[string]interface{}{},
			expected: "{}",
		},
		{
			name: "simple JSON object",
			value: map[string]interface{}{
				"theme": "dark",
			},
			expected: `{"theme":"dark"}`,
		},
		{
			name: "JSON object with multiple keys",
			value: map[string]interface{}{
				"theme":         "dark",
				"notifications": true,
				"language":      "en",
			},
			// Keys are sorted alphabetically by json.Marshal
			expected: `{"language":"en","notifications":true,"theme":"dark"}`,
		},
		{
			name:     "empty JSON array",
			value:    []interface{}{},
			expected: "[]",
		},
		{
			name:     "JSON array with strings",
			value:    []interface{}{"developer", "senior"},
			expected: `["developer","senior"]`,
		},
		{
			name:     "JSON array with mixed types",
			value:    []interface{}{"tag1", 123, true},
			expected: `["tag1",123,true]`,
		},
		{
			name: "nested JSON object",
			value: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  30,
				},
			},
			expected: `{"user":{"age":30,"name":"Alice"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatValue(tt.value)
			if result != tt.expected {
				t.Errorf("FormatValue() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatValuePretty(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "nil value",
			value:    nil,
			expected: "(null)",
		},
		{
			name:     "string value",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "empty JSON object",
			value:    map[string]interface{}{},
			expected: "{}",
		},
		{
			name: "simple JSON object",
			value: map[string]interface{}{
				"theme": "dark",
			},
			expected: `{
  "theme": "dark"
}`,
		},
		{
			name: "JSON object with multiple keys",
			value: map[string]interface{}{
				"theme":         "dark",
				"notifications": true,
				"language":      "en",
			},
			expected: `{
  "language": "en",
  "notifications": true,
  "theme": "dark"
}`,
		},
		{
			name:     "empty JSON array",
			value:    []interface{}{},
			expected: "[]",
		},
		{
			name:     "JSON array with strings",
			value:    []interface{}{"developer", "senior"},
			expected: `[
  "developer",
  "senior"
]`,
		},
		{
			name: "nested JSON object",
			value: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Alice",
					"age":  30,
				},
			},
			expected: `{
  "user": {
    "age": 30,
    "name": "Alice"
  }
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatValuePretty(tt.value)
			if result != tt.expected {
				t.Errorf("FormatValuePretty() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatJSONMinimal(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: "{}",
		},
		{
			name: "single key",
			input: map[string]interface{}{
				"key1": "value1",
			},
			expected: `{"key1":"value1"}`,
		},
		{
			name: "multiple keys - keys sorted alphabetically",
			input: map[string]interface{}{
				"zebra":  "z",
				"apple":  "a",
				"banana": "b",
			},
			expected: `{"apple":"a","banana":"b","zebra":"z"}`,
		},
		{
			name: "keys with different types",
			input: map[string]interface{}{
				"string": "value",
				"number": 123,
				"bool":   true,
			},
			expected: `{"bool":true,"number":123,"string":"value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatJSONMinimal(tt.input)
			if result != tt.expected {
				t.Errorf("formatJSONMinimal() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatJSONArrayMinimal(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected string
	}{
		{
			name:     "empty array",
			input:    []interface{}{},
			expected: "[]",
		},
		{
			name:     "single element",
			input:    []interface{}{"item1"},
			expected: `["item1"]`,
		},
		{
			name:     "multiple elements",
			input:    []interface{}{"a", "b", "c"},
			expected: `["a","b","c"]`,
		},
		{
			name:     "mixed types",
			input:    []interface{}{"string", 123, true, nil},
			expected: `["string",123,true,null]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatJSONArrayMinimal(tt.input)
			if result != tt.expected {
				t.Errorf("formatJSONArrayMinimal() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatJSONPretty(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    string
		checkLines  bool
		lineCount   int
		containsAll []string
	}{
		{
			name:  "simple object",
			input: map[string]interface{}{"key": "value"},
			expected: `{
  "key": "value"
}`,
		},
		{
			name: "nested object",
			input: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "value",
				},
			},
			checkLines: true,
			lineCount:  5,
			containsAll: []string{
				`"outer"`,
				`"inner"`,
				`"value"`,
			},
		},
		{
			name: "array",
			input: []interface{}{
				"item1",
				"item2",
			},
			checkLines: true,
			lineCount:  4,
			containsAll: []string{
				`"item1"`,
				`"item2"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatJSONPretty(tt.input)
			if tt.expected != "" && result != tt.expected {
				t.Errorf("formatJSONPretty() = %q, want %q", result, tt.expected)
			}
			if tt.checkLines {
				lines := strings.Split(result, "\n")
				if len(lines) != tt.lineCount {
					t.Errorf("formatJSONPretty() line count = %d, want %d", len(lines), tt.lineCount)
				}
			}
			if tt.containsAll != nil {
				for _, needle := range tt.containsAll {
					if !strings.Contains(result, needle) {
						t.Errorf("formatJSONPretty() does not contain %q", needle)
					}
				}
			}
		})
	}
}

func TestIsJSONType(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
		{
			name:     "string",
			value:    "hello",
			expected: false,
		},
		{
			name:     "integer",
			value:    42,
			expected: false,
		},
		{
			name:     "boolean",
			value:    true,
			expected: false,
		},
		{
			name:     "map[string]interface{}",
			value:    map[string]interface{}{"key": "value"},
			expected: true,
		},
		{
			name:     "empty map",
			value:    map[string]interface{}{},
			expected: true,
		},
		{
			name:     "[]interface{}",
			value:    []interface{}{"item1", "item2"},
			expected: true,
		},
		{
			name:     "empty slice",
			value:    []interface{}{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsJSONType(tt.value)
			if result != tt.expected {
				t.Errorf("IsJSONType() = %v, want %v", result, tt.expected)
			}
		})
	}
}
