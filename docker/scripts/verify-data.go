package main

import (
	"fmt"
	"log"

	"github.com/oracle/nosql-go-sdk/nosqldb"
)

func main() {
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

	// Get list of tables
	fmt.Println("=== Tables ===")
	listReq := &nosqldb.ListTablesRequest{}
	listResult, err := client.ListTables(listReq)
	if err != nil {
		log.Fatalf("Failed to list tables: %v", err)
	}

	for _, tableName := range listResult.Tables {
		fmt.Printf("- %s\n", tableName)
	}

	fmt.Println("\n✓ Tables are available!")
	fmt.Println("You can now implement the table list display in the TUI.")
}
