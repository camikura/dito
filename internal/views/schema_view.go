package views

import (
	"fmt"
	"strings"

	"github.com/oracle/nosql-go-sdk/nosqldb"
)

// ColumnInfo represents a column with its name and type.
type ColumnInfo struct {
	Name string
	Type string
}

// SchemaViewModel represents the data needed to render the schema view.
type SchemaViewModel struct {
	TableName      string
	AllTables      []string                // All tables for finding children
	TableSchema    *nosqldb.TableResult    // Schema information (DDL, etc.)
	Indexes        []nosqldb.IndexInfo     // Index information
	LoadingDetails bool                    // Whether details are being loaded
}

// RenderSchemaView renders the schema view for a selected table.
// This is a pure rendering function that takes model data and returns a string.
func RenderSchemaView(m SchemaViewModel) string {
	var content strings.Builder

	// Header
	content.WriteString(fmt.Sprintf("Table:    %s\n", m.TableName))

	// Parent/Child relationship
	if strings.Contains(m.TableName, ".") {
		// Child table
		parts := strings.Split(m.TableName, ".")
		content.WriteString(fmt.Sprintf("Parent:   %s\n", parts[0]))
		content.WriteString("Children: -\n")
	} else {
		// Parent table - count child tables
		content.WriteString("Parent:   -\n")
		childCount := 0
		var childNames []string
		prefix := m.TableName + "."
		for _, t := range m.AllTables {
			if strings.HasPrefix(t, prefix) {
				childCount++
				childNames = append(childNames, strings.TrimPrefix(t, prefix))
			}
		}
		if childCount > 0 {
			content.WriteString(fmt.Sprintf("Children: %s\n", strings.Join(childNames, ", ")))
		} else {
			content.WriteString("Children: -\n")
		}
	}

	// Column information and index information
	if m.TableSchema != nil {
		// Column information (extracted from DDL string)
		content.WriteString("\nColumns:\n")
		if m.TableSchema.DDL != "" {
			// Extract column information from DDL
			primaryKeys := ParsePrimaryKeysFromDDL(m.TableSchema.DDL)
			columns := ParseColumnsFromDDL(m.TableSchema.DDL, primaryKeys)
			if len(columns) > 0 {
				for _, col := range columns {
					content.WriteString(fmt.Sprintf("  %-20s %s\n", col.Name, col.Type))
				}
			} else {
				content.WriteString("  (No column information available)\n")
			}
		} else if m.TableSchema.Schema != "" {
			content.WriteString("  " + m.TableSchema.Schema + "\n")
		} else {
			content.WriteString("  (No column information available)\n")
		}

		// Index information
		content.WriteString("\nIndexes:\n")
		if len(m.Indexes) > 0 {
			for _, index := range m.Indexes {
				fields := strings.Join(index.FieldNames, ", ")
				content.WriteString(fmt.Sprintf("  %-20s (%s)\n", index.IndexName, fields))
			}
		} else {
			content.WriteString("  (none)\n")
		}
	} else if m.LoadingDetails {
		content.WriteString("\nColumns:\n")
		content.WriteString("  Loading...\n")
		content.WriteString("\nIndexes:\n")
		content.WriteString("  Loading...\n")
	} else {
		content.WriteString("\nColumns:\n")
		content.WriteString("  (Schema information will be displayed here)\n")
		content.WriteString("\nIndexes:\n")
		content.WriteString("  (Index information will be displayed here)\n")
	}

	return content.String()
}

// ParsePrimaryKeysFromDDL extracts primary key column names from DDL string.
func ParsePrimaryKeysFromDDL(ddl string) []string {
	var primaryKeys []string

	// Find PRIMARY KEY(col1, col2, ...) part
	upperDDL := strings.ToUpper(ddl)
	pkIndex := strings.Index(upperDDL, "PRIMARY KEY")
	if pkIndex == -1 {
		return primaryKeys
	}

	// Get content inside parentheses after PRIMARY KEY
	pkPart := ddl[pkIndex:]
	start := strings.Index(pkPart, "(")
	end := strings.LastIndex(pkPart, ")") // Get last parenthesis
	if start == -1 || end == -1 || start >= end {
		return primaryKeys
	}

	// Extract column name list
	keysPart := pkPart[start+1 : end]

	// Handle SHARD() syntax
	// Support format like PRIMARY KEY(SHARD(id), name)
	keysPart = strings.ReplaceAll(keysPart, "SHARD(", "")
	keysPart = strings.ReplaceAll(keysPart, "shard(", "")
	keysPart = strings.ReplaceAll(keysPart, ")", "")

	keys := strings.Split(keysPart, ",")
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key != "" {
			primaryKeys = append(primaryKeys, key)
		}
	}

	return primaryKeys
}

// ParseColumnsFromDDL extracts column information from DDL string.
func ParseColumnsFromDDL(ddl string, primaryKeys []string) []ColumnInfo {
	var columns []ColumnInfo

	// Create PRIMARY KEY map for fast lookup
	pkMap := make(map[string]bool)
	for _, pk := range primaryKeys {
		pkMap[pk] = true
	}

	// Extract column definition part from CREATE TABLE ... ( ... )
	start := strings.Index(ddl, "(")
	end := strings.LastIndex(ddl, ")")
	if start == -1 || end == -1 || start >= end {
		return columns
	}

	columnDefs := ddl[start+1 : end]

	// Exclude PRIMARY KEY definition
	lines := strings.Split(columnDefs, ",")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip PRIMARY KEY line
		if strings.HasPrefix(strings.ToUpper(line), "PRIMARY KEY") {
			continue
		}

		// Extract column name and type
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			typ := parts[1]

			// Check if it's a PRIMARY KEY
			if pkMap[name] {
				typ += " (Primary Key)"
			}

			columns = append(columns, ColumnInfo{Name: name, Type: typ})
		}
	}

	return columns
}
