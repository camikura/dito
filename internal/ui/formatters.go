package ui

import (
	"encoding/json"
	"fmt"
	"sort"
)

// FormatValue formats a value for display.
// For JSON objects/arrays, it provides minimal representation.
func FormatValue(value interface{}) string {
	if value == nil {
		return "(null)"
	}

	// Check if it's a map (JSON object)
	if m, ok := value.(map[string]interface{}); ok {
		return formatJSONMinimal(m)
	}

	// Check if it's a slice (JSON array)
	if s, ok := value.([]interface{}); ok {
		return formatJSONArrayMinimal(s)
	}

	return fmt.Sprintf("%v", value)
}

// FormatValuePretty formats a value for detailed display.
// For JSON objects/arrays, it provides pretty-printed representation.
func FormatValuePretty(value interface{}) string {
	if value == nil {
		return "(null)"
	}

	// Check if it's a map (JSON object) or slice (JSON array)
	if m, ok := value.(map[string]interface{}); ok {
		return formatJSONPretty(m)
	}
	if s, ok := value.([]interface{}); ok {
		return formatJSONPretty(s)
	}

	return fmt.Sprintf("%v", value)
}

// formatJSONMinimal creates a minimal representation of a JSON object
// Generates compact JSON string - will be truncated by TruncateString if needed
func formatJSONMinimal(m map[string]interface{}) string {
	if len(m) == 0 {
		return "{}"
	}

	// Generate compact JSON (single line, no indentation)
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		// Fallback to showing first key
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return fmt.Sprintf(`{"%s":...}`, keys[0])
	}

	return string(jsonBytes)
}

// formatJSONArrayMinimal creates a minimal representation of a JSON array
// Generates compact JSON string - will be truncated by TruncateString if needed
func formatJSONArrayMinimal(s []interface{}) string {
	if len(s) == 0 {
		return "[]"
	}

	// Generate compact JSON (single line, no indentation)
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		// Fallback to showing length
		return fmt.Sprintf("[...%d]", len(s))
	}

	return string(jsonBytes)
}

// formatJSONPretty creates a pretty-printed JSON string
func formatJSONPretty(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(bytes)
}

// IsJSONType checks if a value is a JSON object or array
func IsJSONType(value interface{}) bool {
	switch value.(type) {
	case map[string]interface{}, []interface{}:
		return true
	default:
		return false
	}
}
