package ui

// InsertAt inserts a string at the specified position in text.
// Returns the new text.
func InsertAt(text string, pos int, insert string) string {
	if pos < 0 {
		pos = 0
	}
	if pos > len(text) {
		pos = len(text)
	}
	return text[:pos] + insert + text[pos:]
}

// DeleteAt deletes the character at the specified position (forward delete).
// Returns the new text.
func DeleteAt(text string, pos int) string {
	if pos < 0 || pos >= len(text) {
		return text
	}
	return text[:pos] + text[pos+1:]
}

// Backspace deletes the character before the cursor position.
// Returns the new text and new cursor position.
func Backspace(text string, pos int) (string, int) {
	if pos <= 0 || pos > len(text) {
		return text, pos
	}
	return text[:pos-1] + text[pos:], pos - 1
}

// InsertWithCursor inserts text at cursor position and moves cursor forward.
// Returns the new text and new cursor position.
func InsertWithCursor(text string, pos int, insert string) (string, int) {
	newText := InsertAt(text, pos, insert)
	return newText, pos + len(insert)
}
