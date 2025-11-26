package ui

import (
	"strings"
	"testing"
)

func TestRecordDetail_NewRecordDetail(t *testing.T) {
	config := RecordDetailConfig{
		Row: map[string]interface{}{
			"id":   1,
			"name": "Alice",
		},
		Columns: []string{"id", "name"},
		Width:   40,
		Height:  10,
	}

	rd := NewRecordDetail(config)

	if rd == nil {
		t.Fatal("NewRecordDetail returned nil")
	}

	if rd.config.Title != " Record Details " {
		t.Errorf("Default title not set, got %q", rd.config.Title)
	}

	if rd.config.BorderColor != "#00D9FF" {
		t.Errorf("Default border color not set, got %q", rd.config.BorderColor)
	}
}

func TestRecordDetail_TotalLines(t *testing.T) {
	config := RecordDetailConfig{
		Row: map[string]interface{}{
			"id":    1,
			"name":  "Alice",
			"email": "alice@example.com",
		},
		Columns: []string{"id", "name", "email"},
		Width:   40,
		Height:  10,
	}

	rd := NewRecordDetail(config)

	// Should have 3 lines (one per column)
	if rd.TotalLines() != 3 {
		t.Errorf("TotalLines() = %d, want 3", rd.TotalLines())
	}
}

func TestRecordDetail_MaxScroll(t *testing.T) {
	tests := []struct {
		name      string
		rows      int
		height    int
		expectMax int
	}{
		{
			name:      "content fits",
			rows:      5,
			height:   10,
			expectMax: 0,
		},
		{
			name:      "content overflows",
			rows:      15,
			height:   10,
			expectMax: 7, // 15 - (10-2) = 7
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := make(map[string]interface{})
			columns := make([]string, tt.rows)
			for i := 0; i < tt.rows; i++ {
				col := "col" + string(rune('a'+i))
				row[col] = i
				columns[i] = col
			}

			config := RecordDetailConfig{
				Row:     row,
				Columns: columns,
				Width:   40,
				Height:  tt.height,
			}

			rd := NewRecordDetail(config)

			if rd.MaxScroll() != tt.expectMax {
				t.Errorf("MaxScroll() = %d, want %d", rd.MaxScroll(), tt.expectMax)
			}
		})
	}
}

func TestRecordDetail_Render(t *testing.T) {
	config := RecordDetailConfig{
		Row: map[string]interface{}{
			"id":   1,
			"name": "Alice",
		},
		Columns: []string{"id", "name"},
		Width:   30,
		Height:  6,
	}

	rd := NewRecordDetail(config)
	result := rd.Render()

	// Should contain borders
	if !strings.Contains(result, "╭") {
		t.Error("Result should contain top-left corner")
	}
	if !strings.Contains(result, "╮") {
		t.Error("Result should contain top-right corner")
	}
	if !strings.Contains(result, "╰") {
		t.Error("Result should contain bottom-left corner")
	}
	if !strings.Contains(result, "╯") {
		t.Error("Result should contain bottom-right corner")
	}

	// Should contain title
	if !strings.Contains(result, "Record Details") {
		t.Error("Result should contain title 'Record Details'")
	}
}

func TestRecordDetail_RenderCentered(t *testing.T) {
	config := RecordDetailConfig{
		Row: map[string]interface{}{
			"id": 1,
		},
		Columns: []string{"id"},
		Width:   20,
		Height:  5,
	}

	rd := NewRecordDetail(config)
	result := rd.RenderCentered(80, 24)

	// Result should be larger than the dialog itself (due to centering)
	lines := strings.Split(result, "\n")
	if len(lines) < 5 {
		t.Errorf("Centered result should have at least 5 lines, got %d", len(lines))
	}
}

func TestRecordDetail_ScrollOffset(t *testing.T) {
	// Create config with many columns to force scrolling
	row := make(map[string]interface{})
	columns := make([]string, 20)
	for i := 0; i < 20; i++ {
		col := "column" + string(rune('a'+i))
		row[col] = i
		columns[i] = col
	}

	config := RecordDetailConfig{
		Row:          row,
		Columns:      columns,
		Width:        40,
		Height:       10,
		ScrollOffset: 5,
	}

	rd := NewRecordDetail(config)

	// Max scroll should be 20 - 8 = 12
	if rd.MaxScroll() != 12 {
		t.Errorf("MaxScroll() = %d, want 12", rd.MaxScroll())
	}
}
