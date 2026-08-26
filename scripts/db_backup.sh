#!/bin/bash
echo "==================================="
echo "  PostgreSQL Database Backup Tool  "
echo "==================================="
# Navigate to project root
cd "$(dirname "$0")/.."
mkdir -p database

if [ "$(docker ps -q -f name=inventory_postgres)" ]; then
    echo "Detected Docker container: inventory_postgres"
    echo "Creating backup via Docker pg_dump..."
    docker exec -t inventory_postgres pg_dump -U postgres -d inventory_db > database/backup.sql
else
    echo "Docker container not running. Attempting local pg_dump..."
    pg_dump -U postgres -d inventory_db > database/backup.sql
fi

if [ $? -eq 0 ]; then
    echo "Backup completed successfully: database/backup.sql"
else
    echo "Backup failed. Please verify database connection settings."
fi
