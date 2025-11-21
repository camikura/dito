package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/oracle/nosql-go-sdk/nosqldb"
)

func main() {
	// Read SQL file
	sqlBytes, err := os.ReadFile("docker/schema/init.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	// Split SQL statements (assuming each statement ends with semicolon)
	sqlContent := string(sqlBytes)
	statements := splitSQLStatements(sqlContent)

	// Configure client for on-premise database
	cfg := nosqldb.Config{
		Endpoint: "http://localhost:8080",
		Mode:     "onprem",
	}

	client, err := nosqldb.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create NoSQL client: %v", err)
	}
	defer client.Close()

	// Execute each SQL statement
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		fmt.Printf("Executing statement %d: %s...\n", i+1, truncate(stmt, 60))

		// Check if it's a DDL statement (CREATE, DROP, ALTER) or DML (INSERT, UPDATE, DELETE, SELECT)
		isDDL := strings.HasPrefix(strings.ToUpper(stmt), "CREATE") ||
			strings.HasPrefix(strings.ToUpper(stmt), "DROP") ||
			strings.HasPrefix(strings.ToUpper(stmt), "ALTER")

		if isDDL {
			// Use SystemRequest for DDL
			req := &nosqldb.SystemRequest{
				Statement: stmt,
			}

			result, err := client.DoSystemRequest(req)
			if err != nil {
				log.Printf("Error executing DDL: %v\nStatement: %s", err, stmt)
				continue
			}

			// Wait for operation to complete (timeout: 30s, poll interval: 1s)
			result.WaitForCompletion(client, 30*time.Second, 1*time.Second)
			fmt.Printf("  ✓ Completed (State: %v)\n", result.State)
		} else {
			// Use Query for DML
			req := &nosqldb.QueryRequest{
				Statement: stmt,
			}

			_, err := client.Query(req)
			if err != nil {
				log.Printf("Error executing DML: %v\nStatement: %s", err, stmt)
				continue
			}

			fmt.Printf("  ✓ Completed\n")
		}
	}

	fmt.Println("\nAll tables created successfully!")
}

func splitSQLStatements(sql string) []string {
	// Remove comments and split by semicolon
	var statements []string
	lines := strings.Split(sql, "\n")
	var currentStmt strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		currentStmt.WriteString(line)
		currentStmt.WriteString(" ")

		// Check if line ends with semicolon
		if strings.HasSuffix(line, ";") {
			stmt := currentStmt.String()
			stmt = strings.TrimSuffix(strings.TrimSpace(stmt), ";")
			if stmt != "" {
				statements = append(statements, stmt)
			}
			currentStmt.Reset()
		}
	}

	// Add any remaining statement
	if currentStmt.Len() > 0 {
		stmt := strings.TrimSuffix(strings.TrimSpace(currentStmt.String()), ";")
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
