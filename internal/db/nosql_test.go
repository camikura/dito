package db

import (
	"reflect"
	"testing"
)

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name:  "simple string",
			input: "hello",
			want:  "'hello'",
		},
		{
			name:  "string with single quote",
			input: "it's",
			want:  "'it''s'",
		},
		{
			name:  "string with multiple single quotes",
			input: "it's a 'test'",
			want:  "'it''s a ''test'''",
		},
		{
			name:  "integer",
			input: 42,
			want:  "42",
		},
		{
			name:  "int32",
			input: int32(123),
			want:  "123",
		},
		{
			name:  "int64",
			input: int64(9876543210),
			want:  "9876543210",
		},
		{
			name:  "float32",
			input: float32(3.14),
			want:  "3.14",
		},
		{
			name:  "float64",
			input: float64(2.718281828),
			want:  "2.718281828",
		},
		{
			name:  "boolean",
			input: true,
			want:  "'true'",
		},
		{
			name:  "nil",
			input: nil,
			want:  "'<nil>'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatValue(tt.input)
			if got != tt.want {
				t.Errorf("formatValue(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestConvertRowValues(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		checkFn  func(map[string]interface{}) bool
		errMsg   string
	}{
		{
			name:  "nil input",
			input: nil,
			checkFn: func(got map[string]interface{}) bool {
				return len(got) == 0
			},
			errMsg: "should return empty map for nil input",
		},
		{
			name:  "empty map",
			input: map[string]interface{}{},
			checkFn: func(got map[string]interface{}) bool {
				return len(got) == 0
			},
			errMsg: "should return empty map for empty input",
		},
		{
			name: "simple values",
			input: map[string]interface{}{
				"id":   1,
				"name": "Alice",
			},
			checkFn: func(got map[string]interface{}) bool {
				return len(got) == 2 && got["name"] == "Alice"
			},
			errMsg: "should preserve simple values",
		},
		{
			name: "nil value",
			input: map[string]interface{}{
				"id":   1,
				"data": nil,
			},
			checkFn: func(got map[string]interface{}) bool {
				_, exists := got["data"]
				return len(got) == 2 && exists && got["data"] == nil
			},
			errMsg: "should preserve nil values",
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"id": 1,
				"data": map[string]interface{}{
					"key": "value",
				},
			},
			checkFn: func(got map[string]interface{}) bool {
				if len(got) != 2 {
					return false
				}
				data, ok := got["data"].(map[string]interface{})
				return ok && data["key"] == "value"
			},
			errMsg: "should preserve nested maps",
		},
		{
			name: "slice value",
			input: map[string]interface{}{
				"id":   1,
				"tags": []interface{}{"tag1", "tag2"},
			},
			checkFn: func(got map[string]interface{}) bool {
				if len(got) != 2 {
					return false
				}
				tags, ok := got["tags"].([]interface{})
				return ok && len(tags) == 2 && tags[0] == "tag1" && tags[1] == "tag2"
			},
			errMsg: "should preserve slice values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertRowValues(tt.input)
			if !tt.checkFn(got) {
				t.Errorf("convertRowValues() failed: %s, got = %v", tt.errMsg, got)
			}
		})
	}
}

func TestConvertValue(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		checkFn func(interface{}) bool
		errMsg  string
	}{
		{
			name:  "nil value",
			input: nil,
			checkFn: func(got interface{}) bool {
				return got == nil
			},
			errMsg: "should return nil",
		},
		{
			name:  "string value",
			input: "hello",
			checkFn: func(got interface{}) bool {
				return got == "hello"
			},
			errMsg: "should preserve string value",
		},
		{
			name:  "integer value",
			input: 42,
			checkFn: func(got interface{}) bool {
				// After JSON conversion, integers may become float64
				return got == 42 || got == float64(42)
			},
			errMsg: "should preserve integer value (possibly as float64)",
		},
		{
			name:  "boolean value",
			input: true,
			checkFn: func(got interface{}) bool {
				return got == true
			},
			errMsg: "should preserve boolean value",
		},
		{
			name: "map value",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
			checkFn: func(got interface{}) bool {
				m, ok := got.(map[string]interface{})
				return ok && len(m) == 2 && m["key1"] == "value1"
			},
			errMsg: "should preserve map structure",
		},
		{
			name:  "slice value",
			input: []interface{}{"a", "b", "c"},
			checkFn: func(got interface{}) bool {
				s, ok := got.([]interface{})
				return ok && len(s) == 3 && s[0] == "a" && s[1] == "b" && s[2] == "c"
			},
			errMsg: "should preserve slice values",
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "value",
				},
			},
			checkFn: func(got interface{}) bool {
				m, ok := got.(map[string]interface{})
				if !ok {
					return false
				}
				inner, ok := m["outer"].(map[string]interface{})
				return ok && inner["inner"] == "value"
			},
			errMsg: "should preserve nested map structure",
		},
		{
			name:  "slice with nested maps",
			input: []interface{}{map[string]interface{}{"key": "value"}},
			checkFn: func(got interface{}) bool {
				s, ok := got.([]interface{})
				if !ok || len(s) != 1 {
					return false
				}
				m, ok := s[0].(map[string]interface{})
				return ok && m["key"] == "value"
			},
			errMsg: "should preserve slice with nested maps",
		},
		{
			name: "pointer to string",
			input: func() interface{} {
				s := "test"
				return &s
			}(),
			checkFn: func(got interface{}) bool {
				return got == "test"
			},
			errMsg: "should dereference pointer to string",
		},
		{
			name: "pointer to int",
			input: func() interface{} {
				i := 42
				return &i
			}(),
			checkFn: func(got interface{}) bool {
				// After JSON conversion, integers may become float64
				return got == 42 || got == float64(42)
			},
			errMsg: "should dereference pointer to int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertValue(tt.input)
			if !tt.checkFn(got) {
				t.Errorf("convertValue() failed: %s, got = %v (type %T)", tt.errMsg, got, got)
			}
		})
	}
}

func TestConvertMapValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{} // Use interface{} to allow nil
		expected map[string]interface{}
	}{
		{
			name:     "nil MapValue",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertMapValue(nil)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("convertMapValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConvertValueWithPointers(t *testing.T) {
	// Test conversion of pointers to primitive types
	t.Run("pointer to map", func(t *testing.T) {
		m := map[string]interface{}{"key": "value"}
		ptr := &m
		result := convertValue(ptr)
		expected := map[string]interface{}{"key": "value"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("convertValue() with pointer to map = %v, want %v", result, expected)
		}
	})

	t.Run("pointer to slice", func(t *testing.T) {
		s := []interface{}{"a", "b"}
		ptr := &s
		result := convertValue(ptr)
		expected := []interface{}{"a", "b"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("convertValue() with pointer to slice = %v, want %v", result, expected)
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		var ptr *string
		result := convertValue(ptr)
		if result != nil {
			t.Errorf("convertValue() with nil pointer = %v, want nil", result)
		}
	})
}

func TestConvertValueWithComplexStructures(t *testing.T) {
	t.Run("deeply nested structure", func(t *testing.T) {
		input := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"level3": []interface{}{
						map[string]interface{}{
							"key": "value",
						},
					},
				},
			},
		}
		result := convertValue(input)

		// Check structure is preserved
		m, ok := result.(map[string]interface{})
		if !ok {
			t.Errorf("result is not a map")
			return
		}
		l1, ok := m["level1"].(map[string]interface{})
		if !ok {
			t.Errorf("level1 is not a map")
			return
		}
		l2, ok := l1["level2"].(map[string]interface{})
		if !ok {
			t.Errorf("level2 is not a map")
			return
		}
		l3, ok := l2["level3"].([]interface{})
		if !ok || len(l3) != 1 {
			t.Errorf("level3 is not a slice with 1 element")
			return
		}
		l3m, ok := l3[0].(map[string]interface{})
		if !ok || l3m["key"] != "value" {
			t.Errorf("nested structure not preserved correctly")
		}
	})

	t.Run("mixed types in slice", func(t *testing.T) {
		input := []interface{}{
			"string",
			123,
			true,
			nil,
			map[string]interface{}{"key": "value"},
			[]interface{}{"nested"},
		}
		result := convertValue(input)

		s, ok := result.([]interface{})
		if !ok {
			t.Errorf("result is not a slice")
			return
		}
		if len(s) != 6 {
			t.Errorf("slice length = %d, want 6", len(s))
			return
		}
		if s[0] != "string" {
			t.Errorf("s[0] = %v, want 'string'", s[0])
		}
		// s[1] could be 123 or float64(123) after JSON conversion
		if s[1] != 123 && s[1] != float64(123) {
			t.Errorf("s[1] = %v, want 123", s[1])
		}
		if s[2] != true {
			t.Errorf("s[2] = %v, want true", s[2])
		}
		if s[3] != nil {
			t.Errorf("s[3] = %v, want nil", s[3])
		}
		m, ok := s[4].(map[string]interface{})
		if !ok || m["key"] != "value" {
			t.Errorf("s[4] is not the expected map")
		}
		nested, ok := s[5].([]interface{})
		if !ok || len(nested) != 1 || nested[0] != "nested" {
			t.Errorf("s[5] is not the expected nested slice")
		}
	})
}
