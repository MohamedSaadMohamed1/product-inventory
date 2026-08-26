@echo off
SETLOCAL

:: Add standard Go installation path to PATH
set PATH=C:\Program Files\Go\bin;%PATH%

:: Check parameter
IF "%1"=="" GOTO usage
IF "%1"=="run" GOTO run
IF "%1"=="test" GOTO test
IF "%1"=="test-race" GOTO test-race
IF "%1"=="sqlc" GOTO sqlc
IF "%1"=="db-backup" GOTO db-backup
IF "%1"=="db-restore" GOTO db-restore
GOTO usage

:run
echo [make] Running API server...
go run cmd/api/main.go
GOTO end

:test
echo [make] Running tests...
go test -v ./...
GOTO end

:test-race
echo [make] Running tests with race detector...
set CGO_ENABLED=1
set PATH=%~dp0w64devkit-mini\w64devkit\bin;%PATH%
go test -v -race ./...
GOTO end

:sqlc
echo [make] Generating database layer...
go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate
GOTO end

:db-backup
echo [make] Backing up database...
call scripts\db_backup.bat
GOTO end

:db-restore
echo [make] Restoring database from backup.sql...
set PGPASSWORD=postgres
set PATH=C:\Program Files\PostgreSQL\15\bin;%PATH%
psql -U postgres -h localhost -p 5432 -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'inventory_db' AND pid <> pg_backend_pid();"
psql -U postgres -h localhost -p 5432 -d postgres -c "DROP DATABASE IF EXISTS inventory_db; CREATE DATABASE inventory_db;"
cmd.exe /c "psql -U postgres -h localhost -p 5432 -d inventory_db < database\backup.sql"
GOTO end

:usage
echo Usage: make [run ^| test ^| test-race ^| sqlc ^| db-backup ^| db-restore]
GOTO end

:end
ENDLOCAL
