package db

import (
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
