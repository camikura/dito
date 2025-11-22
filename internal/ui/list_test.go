package ui

import (
	"strings"
	"testing"
)

func TestSelectableList_Render(t *testing.T) {
	tests := []struct {
		name     string
		list     SelectableList
		contains []string
	}{
		{
			name: "simple list with title",
			list: SelectableList{
				Title:         "Tables (3)",
				Items:         []string{"users", "orders", "products"},
				SelectedIndex: 0,
				Focused:       true,
			},
			contains: []string{"Tables (3)", "users", "orders", "products", ">"},
		},
		{
			name: "list without title",
			list: SelectableList{
				Title:         "",
				Items:         []string{"item1", "item2"},
				SelectedIndex: 1,
				Focused:       true,
			},
			contains: []string{"item1", "item2", ">"},
		},
		{
			name: "empty list",
			list: SelectableList{
				Title:         "Empty",
				Items:         []string{},
				SelectedIndex: 0,
				Focused:       true,
			},
			contains: []string{"Empty"},
		},
		{
			name: "unfocused list",
			list: SelectableList{
				Title:         "Unfocused",
				Items:         []string{"item1", "item2"},
				SelectedIndex: 0,
				Focused:       false,
			},
			contains: []string{"Unfocused", "item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.list.Render()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Render() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestSelectableList_Selection(t *testing.T) {
	list := SelectableList{
		Title:         "Test",
		Items:         []string{"item1", "item2", "item3"},
		SelectedIndex: 1,
		Focused:       true,
	}

	result := list.Render()

	// Should have exactly one selection indicator
	count := strings.Count(result, ">")
	if count != 1 {
		t.Errorf("Should have exactly 1 selection indicator, got %d", count)
	}

	// Should contain the selected item
	if !strings.Contains(result, "item2") {
		t.Errorf("Result should contain selected item 'item2'")
	}
}

func TestSelectableList_FirstItemSelected(t *testing.T) {
	list := SelectableList{
		Items:         []string{"first", "second", "third"},
		SelectedIndex: 0,
		Focused:       true,
	}

	result := list.Render()

	// Should select first item
	if !strings.Contains(result, "> first") {
		t.Errorf("First item should be selected")
	}
}

func TestSelectableList_LastItemSelected(t *testing.T) {
	list := SelectableList{
		Items:         []string{"first", "second", "third"},
		SelectedIndex: 2,
		Focused:       true,
	}

	result := list.Render()

	// Should select last item
	if !strings.Contains(result, "> third") {
		t.Errorf("Last item should be selected")
	}
}

func TestSelectableList_FocusedVsUnfocused(t *testing.T) {
	items := []string{"item1", "item2"}

	// Focused list
	focusedList := SelectableList{
		Items:         items,
		SelectedIndex: 0,
		Focused:       true,
	}
	focusedResult := focusedList.Render()

	// Unfocused list
	unfocusedList := SelectableList{
		Items:         items,
		SelectedIndex: 0,
		Focused:       false,
	}
	unfocusedResult := unfocusedList.Render()

	// Both should contain items
	if !strings.Contains(focusedResult, "item1") {
		t.Errorf("Focused list should contain items")
	}
	if !strings.Contains(unfocusedResult, "item1") {
		t.Errorf("Unfocused list should contain items")
	}

	// Both should have selection indicator
	if !strings.Contains(focusedResult, ">") {
		t.Errorf("Focused list should have selection indicator")
	}
	if !strings.Contains(unfocusedResult, ">") {
		t.Errorf("Unfocused list should have selection indicator")
	}
}

func TestSelectableList_WithTitle(t *testing.T) {
	list := SelectableList{
		Title:         "My List (5)",
		Items:         []string{"a", "b", "c", "d", "e"},
		SelectedIndex: 2,
		Focused:       true,
	}

	result := list.Render()

	// Should contain title
	if !strings.Contains(result, "My List (5)") {
		t.Errorf("Result should contain title")
	}

	// Should contain all items
	for _, item := range list.Items {
		if !strings.Contains(result, item) {
			t.Errorf("Result should contain item %q", item)
		}
	}
}

func TestSelectableList_NoTrailingNewline(t *testing.T) {
	list := SelectableList{
		Items:         []string{"item1"},
		SelectedIndex: 0,
		Focused:       true,
	}

	result := list.Render()

	if strings.HasSuffix(result, "\n") {
		t.Errorf("Result should not have trailing newline")
	}
}
