package ui

import (
	"regexp"
	"strings"
)

// ExtractTableNameFromSQL extracts the table name from a SQL query.
// Supports SELECT ... FROM table and SELECT ... FROM namespace.table
func ExtractTableNameFromSQL(sql string) string {
	// Case-insensitive regex to find FROM clause
	re := regexp.MustCompile(`(?i)\bFROM\s+([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)?)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// InsertAt inserts a string at the specified rune position in text.
// Returns the new text.
func InsertAt(text string, pos int, insert string) string {
	runes := []rune(text)
	if pos < 0 {
		pos = 0
	}
	if pos > len(runes) {
		pos = len(runes)
	}
	result := make([]rune, 0, len(runes)+len([]rune(insert)))
	result = append(result, runes[:pos]...)
	result = append(result, []rune(insert)...)
	result = append(result, runes[pos:]...)
	return string(result)
}

// DeleteAt deletes the character at the specified rune position (forward delete).
// Returns the new text.
func DeleteAt(text string, pos int) string {
	runes := []rune(text)
	if pos < 0 || pos >= len(runes) {
		return text
	}
	return string(append(runes[:pos], runes[pos+1:]...))
}

// Backspace deletes the character before the cursor position (rune-based).
// Returns the new text and new cursor position.
func Backspace(text string, pos int) (string, int) {
	runes := []rune(text)
	if pos <= 0 || pos > len(runes) {
		return text, pos
	}
	return string(append(runes[:pos-1], runes[pos:]...)), pos - 1
}

// InsertWithCursor inserts text at cursor position and moves cursor forward.
// Returns the new text and new cursor position (rune-based).
func InsertWithCursor(text string, pos int, insert string) (string, int) {
	newText := InsertAt(text, pos, insert)
	return newText, pos + len([]rune(insert))
}

// RuneLen returns the number of runes in a string.
func RuneLen(text string) int {
	return len([]rune(text))
}
