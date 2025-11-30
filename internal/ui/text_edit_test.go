package ui

import "testing"

func TestInsertAt(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		pos      int
		insert   string
		expected string
	}{
		{"insert at start", "hello", 0, "X", "Xhello"},
		{"insert at end", "hello", 5, "X", "helloX"},
		{"insert in middle", "hello", 2, "X", "heXllo"},
		{"insert empty", "hello", 2, "", "hello"},
		{"insert into empty", "", 0, "X", "X"},
		{"negative pos clamps to 0", "hello", -1, "X", "Xhello"},
		{"pos beyond length clamps", "hello", 100, "X", "helloX"},
		{"insert multi-char", "hello", 2, "XYZ", "heXYZllo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InsertAt(tt.text, tt.pos, tt.insert)
			if result != tt.expected {
				t.Errorf("InsertAt(%q, %d, %q) = %q, want %q", tt.text, tt.pos, tt.insert, result, tt.expected)
			}
		})
	}
}

func TestDeleteAt(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		pos      int
		expected string
	}{
		{"delete at start", "hello", 0, "ello"},
		{"delete at end", "hello", 4, "hell"},
		{"delete in middle", "hello", 2, "helo"},
		{"delete beyond length no-op", "hello", 5, "hello"},
		{"delete negative pos no-op", "hello", -1, "hello"},
		{"delete from empty no-op", "", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeleteAt(tt.text, tt.pos)
			if result != tt.expected {
				t.Errorf("DeleteAt(%q, %d) = %q, want %q", tt.text, tt.pos, result, tt.expected)
			}
		})
	}
}

func TestBackspace(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		pos         int
		expectedTxt string
		expectedPos int
	}{
		{"backspace in middle", "hello", 3, "helo", 2},
		{"backspace at end", "hello", 5, "hell", 4},
		{"backspace at start no-op", "hello", 0, "hello", 0},
		{"backspace negative pos no-op", "hello", -1, "hello", -1},
		{"backspace beyond length no-op", "hello", 10, "hello", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultTxt, resultPos := Backspace(tt.text, tt.pos)
			if resultTxt != tt.expectedTxt || resultPos != tt.expectedPos {
				t.Errorf("Backspace(%q, %d) = (%q, %d), want (%q, %d)",
					tt.text, tt.pos, resultTxt, resultPos, tt.expectedTxt, tt.expectedPos)
			}
		})
	}
}

func TestInsertWithCursor(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		pos         int
		insert      string
		expectedTxt string
		expectedPos int
	}{
		{"insert single char", "hello", 2, "X", "heXllo", 3},
		{"insert multi char", "hello", 2, "XYZ", "heXYZllo", 5},
		{"insert at start", "hello", 0, "X", "Xhello", 1},
		{"insert at end", "hello", 5, "X", "helloX", 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultTxt, resultPos := InsertWithCursor(tt.text, tt.pos, tt.insert)
			if resultTxt != tt.expectedTxt || resultPos != tt.expectedPos {
				t.Errorf("InsertWithCursor(%q, %d, %q) = (%q, %d), want (%q, %d)",
					tt.text, tt.pos, tt.insert, resultTxt, resultPos, tt.expectedTxt, tt.expectedPos)
			}
		})
	}
}

func TestExtractTableNameFromSQL(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		expected string
	}{
		{"simple select", "SELECT * FROM users", "users"},
		{"select with columns", "SELECT id, name FROM products", "products"},
		{"lowercase from", "select * from orders", "orders"},
		{"mixed case", "SELECT * From Users", "Users"},
		{"with namespace", "SELECT * FROM ns.tablename", "ns.tablename"},
		{"with where", "SELECT * FROM users WHERE id = 1", "users"},
		{"with join", "SELECT * FROM users JOIN orders", "users"},
		{"multiline", "SELECT *\nFROM users\nWHERE id = 1", "users"},
		{"extra spaces", "SELECT  *  FROM   users", "users"},
		{"no from clause", "SELECT 1", ""},
		{"empty string", "", ""},
		{"only select", "SELECT *", ""},
		{"from at end", "SELECT * FROM", ""},
		{"table with underscore", "SELECT * FROM user_accounts", "user_accounts"},
		{"table with numbers", "SELECT * FROM table123", "table123"},
		{"child table", "SELECT * FROM orders.items", "orders.items"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractTableNameFromSQL(tt.sql)
			if result != tt.expected {
				t.Errorf("ExtractTableNameFromSQL(%q) = %q, want %q", tt.sql, result, tt.expected)
			}
		})
	}
}

func TestRuneLen(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{"empty", "", 0},
		{"ascii", "hello", 5},
		{"japanese", "こんにちは", 5},
		{"mixed", "hello世界", 7},
		{"emoji", "👍", 1},
		{"multiple emoji", "👍👎", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RuneLen(tt.text)
			if result != tt.expected {
				t.Errorf("RuneLen(%q) = %d, want %d", tt.text, result, tt.expected)
			}
		})
	}
}

func TestMultibyteOperations(t *testing.T) {
	// Test InsertAt with Japanese characters
	t.Run("InsertAt with Japanese", func(t *testing.T) {
		result := InsertAt("こんにちは", 2, "X")
		expected := "こんXにちは"
		if result != expected {
			t.Errorf("InsertAt with Japanese = %q, want %q", result, expected)
		}
	})

	// Test DeleteAt with Japanese characters
	t.Run("DeleteAt with Japanese", func(t *testing.T) {
		result := DeleteAt("こんにちは", 2)
		expected := "こんちは"
		if result != expected {
			t.Errorf("DeleteAt with Japanese = %q, want %q", result, expected)
		}
	})

	// Test Backspace with Japanese characters
	t.Run("Backspace with Japanese", func(t *testing.T) {
		resultTxt, resultPos := Backspace("こんにちは", 3)
		expectedTxt := "こんちは"
		expectedPos := 2
		if resultTxt != expectedTxt || resultPos != expectedPos {
			t.Errorf("Backspace with Japanese = (%q, %d), want (%q, %d)",
				resultTxt, resultPos, expectedTxt, expectedPos)
		}
	})

	// Test InsertWithCursor with Japanese characters
	t.Run("InsertWithCursor with Japanese", func(t *testing.T) {
		resultTxt, resultPos := InsertWithCursor("こんにちは", 2, "世界")
		expectedTxt := "こん世界にちは"
		expectedPos := 4
		if resultTxt != expectedTxt || resultPos != expectedPos {
			t.Errorf("InsertWithCursor with Japanese = (%q, %d), want (%q, %d)",
				resultTxt, resultPos, expectedTxt, expectedPos)
		}
	})
}
