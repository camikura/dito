package ui

import (
	"testing"
)

func TestCopyRowToClipboard(t *testing.T) {
	// Test with a simple row
	row := map[string]interface{}{
		"id":   1,
		"name": "test",
	}

	// This test may fail in CI environments without display
	// We just ensure it doesn't panic
	err := CopyRowToClipboard(row, []string{"id", "name"})
	if err != nil {
		t.Skipf("Clipboard not available in this environment: %v", err)
	}
}

func TestCopyRowToClipboardNoOrder(t *testing.T) {
	// Test with nil column order (falls back to alphabetical)
	row := map[string]interface{}{
		"id":   1,
		"name": "test",
	}

	err := CopyRowToClipboard(row, nil)
	if err != nil {
		t.Skipf("Clipboard not available in this environment: %v", err)
	}
}

func TestCopyTextToClipboard(t *testing.T) {
	// This test may fail in CI environments without display
	// We just ensure it doesn't panic
	err := CopyTextToClipboard("test text")
	if err != nil {
		t.Skipf("Clipboard not available in this environment: %v", err)
	}
}

func TestMarshalOrderedJSON(t *testing.T) {
	row := map[string]interface{}{
		"id":    1,
		"name":  "test",
		"email": "test@example.com",
	}

	// Test that order is preserved
	result, err := marshalOrderedJSON(row, []string{"name", "id", "email"})
	if err != nil {
		t.Fatalf("marshalOrderedJSON failed: %v", err)
	}

	expected := `{
  "name": "test",
  "id": 1,
  "email": "test@example.com"
}`
	if string(result) != expected {
		t.Errorf("marshalOrderedJSON() = %s, want %s", string(result), expected)
	}
}

func TestMarshalOrderedJSONMissingColumn(t *testing.T) {
	row := map[string]interface{}{
		"id":   1,
		"name": "test",
	}

	// Column "missing" doesn't exist in row - should be skipped
	result, err := marshalOrderedJSON(row, []string{"name", "missing", "id"})
	if err != nil {
		t.Fatalf("marshalOrderedJSON failed: %v", err)
	}

	expected := `{
  "name": "test",
  "id": 1
}`
	if string(result) != expected {
		t.Errorf("marshalOrderedJSON() = %s, want %s", string(result), expected)
	}
}
