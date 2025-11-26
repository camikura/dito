package views

import (
	"fmt"
	"strings"

	"github.com/camikura/dito/internal/ui"
	"github.com/oracle/nosql-go-sdk/nosqldb"
)

// ColumnInfo is an alias for ui.ColumnInfo for backward compatibility.
type ColumnInfo = ui.ColumnInfo

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
// Deprecated: Use ui.ParsePrimaryKeysFromDDL instead.
func ParsePrimaryKeysFromDDL(ddl string) []string {
	return ui.ParsePrimaryKeysFromDDL(ddl)
}

// ParseColumnsFromDDL extracts column information from DDL string.
// Deprecated: Use ui.ParseColumnsFromDDL instead.
// Note: This function adds " (Primary Key)" suffix to type for backward compatibility.
func ParseColumnsFromDDL(ddl string, primaryKeys []string) []ColumnInfo {
	cols := ui.ParseColumnsFromDDL(ddl, primaryKeys)

	// Add " (Primary Key)" suffix to type for backward compatibility
	for i := range cols {
		if cols[i].IsPrimaryKey {
			cols[i].Type = cols[i].Type + " (Primary Key)"
		}
	}

	return cols
}
