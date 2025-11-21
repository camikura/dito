#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

echo "Creating tables in Oracle NoSQL Database..."

# Check if container is running
if ! docker ps | grep -q nosql-database; then
    echo "Error: nosql-database container is not running. Start it with 'mise run db-start'"
    exit 1
fi

# Execute Go program to create tables via HTTP Proxy
cd "$PROJECT_ROOT"
go run docker/scripts/init-db.go
