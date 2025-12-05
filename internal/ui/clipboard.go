package ui

import (
	"bytes"
	"encoding/json"

	"golang.design/x/clipboard"
)

var clipboardInitialized bool

// InitClipboard initializes the clipboard. Should be called once at startup.
func InitClipboard() error {
	if clipboardInitialized {
		return nil
	}
	err := clipboard.Init()
	if err != nil {
		return err
	}
	clipboardInitialized = true
	return nil
}

// CopyRowToClipboard copies a row to clipboard as JSON with columns in specified order.
// If columnOrder is nil or empty, falls back to default JSON marshaling (alphabetical).
func CopyRowToClipboard(row map[string]interface{}, columnOrder []string) error {
	if !clipboardInitialized {
		if err := InitClipboard(); err != nil {
			return err
		}
	}

	var jsonBytes []byte
	var err error

	if len(columnOrder) > 0 {
		jsonBytes, err = marshalOrderedJSON(row, columnOrder)
	} else {
		jsonBytes, err = json.MarshalIndent(row, "", "  ")
	}

	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, jsonBytes)
	return nil
}

// marshalOrderedJSON marshals a map to JSON with keys in the specified order.
func marshalOrderedJSON(row map[string]interface{}, columnOrder []string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{\n")

	first := true
	for _, col := range columnOrder {
		val, exists := row[col]
		if !exists {
			continue
		}

		if !first {
			buf.WriteString(",\n")
		}
		first = false

		// Marshal the key
		keyBytes, err := json.Marshal(col)
		if err != nil {
			return nil, err
		}

		// Marshal the value
		valBytes, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}

		buf.WriteString("  ")
		buf.Write(keyBytes)
		buf.WriteString(": ")
		buf.Write(valBytes)
	}

	buf.WriteString("\n}")
	return buf.Bytes(), nil
}

// CopyTextToClipboard copies plain text to clipboard.
func CopyTextToClipboard(text string) error {
	if !clipboardInitialized {
		if err := InitClipboard(); err != nil {
			return err
		}
	}

	clipboard.Write(clipboard.FmtText, []byte(text))
	return nil
}
