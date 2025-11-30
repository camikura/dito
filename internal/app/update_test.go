package app

import (
	"reflect"
	"testing"
)

func TestSortTablesForTree(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single table",
			input:    []string{"users"},
			expected: []string{"users"},
		},
		{
			name:     "no child tables",
			input:    []string{"products", "users", "orders"},
			expected: []string{"orders", "products", "users"},
		},
		{
			name:     "parent before child",
			input:    []string{"users.phones", "users"},
			expected: []string{"users", "users.phones"},
		},
		{
			name:     "multiple children",
			input:    []string{"users.phones", "users", "users.addresses"},
			expected: []string{"users", "users.addresses", "users.phones"},
		},
		{
			name:     "mixed parents and children",
			input:    []string{"users.phones", "products", "users", "orders.items", "orders"},
			expected: []string{"orders", "orders.items", "products", "users", "users.phones"},
		},
		{
			name:     "already sorted",
			input:    []string{"orders", "users", "users.addresses"},
			expected: []string{"orders", "users", "users.addresses"},
		},
		{
			name:     "complex hierarchy",
			input:    []string{"a.b", "c", "a", "b.c", "b"},
			expected: []string{"a", "a.b", "b", "b.c", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the original
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			result := sortTablesForTree(input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("sortTablesForTree(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSortTablesForTree_DoesNotModifyInput(t *testing.T) {
	input := []string{"users.phones", "users", "products"}
	original := make([]string, len(input))
	copy(original, input)

	_ = sortTablesForTree(input)

	// The original input should not be modified in order
	// Note: sortTablesForTree copies the input, so this should pass
	if !reflect.DeepEqual(input, original) {
		t.Errorf("sortTablesForTree modified original slice: got %v, want %v", input, original)
	}
}
