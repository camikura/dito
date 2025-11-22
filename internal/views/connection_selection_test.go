package views

import (
	"strings"
	"testing"
)

func TestRenderConnectionSelection(t *testing.T) {
	tests := []struct {
		name     string
		model    ConnectionSelectionModel
		contains []string // strings that should be in the output
	}{
		{
			name: "first item selected",
			model: ConnectionSelectionModel{
				Choices: []string{"On-Premise", "Cloud"},
				Cursor:  0,
			},
			contains: []string{"Select Connection", "On-Premise", "Cloud"},
		},
		{
			name: "second item selected",
			model: ConnectionSelectionModel{
				Choices: []string{"On-Premise", "Cloud"},
				Cursor:  1,
			},
			contains: []string{"Select Connection", "On-Premise", "Cloud"},
		},
		{
			name: "single choice",
			model: ConnectionSelectionModel{
				Choices: []string{"On-Premise"},
				Cursor:  0,
			},
			contains: []string{"Select Connection", "On-Premise"},
		},
		{
			name: "empty choices",
			model: ConnectionSelectionModel{
				Choices: []string{},
				Cursor:  0,
			},
			contains: []string{"Select Connection"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderConnectionSelection(tt.model)

			// Check that all expected strings are present
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("RenderConnectionSelection() = %q, should contain %q", result, substr)
				}
			}

			// Check that result is not empty
			if result == "" {
				t.Error("RenderConnectionSelection() should not return empty string")
			}
		})
	}
}
