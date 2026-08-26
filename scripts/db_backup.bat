@echo off
echo ===================================
echo   PostgreSQL Database Backup Tool
echo ===================================
cd /d "%~dp0.."
mkdir database 2>nul

set "PGPASSWORD=postgres"

docker ps -q -f name=inventory_postgres >nul 2>nul
if %errorlevel% equ 0 (
    echo Detected Docker container: inventory_postgres
    echo Creating backup via Docker pg_dump...
    docker exec -t inventory_postgres pg_dump -U postgres -d inventory_db > database\backup.sql
) else (
    echo Docker container not running. Attempting local pg_dump...
    
    :: Check if pg_dump is in PATH
    where pg_dump >nul 2>nul
    if %errorlevel% neq 0 (
        echo pg_dump not found in system PATH. Attempting automatic detection...
        if exist "C:\Program Files\PostgreSQL\16\bin\pg_dump.exe" (
            set "PATH=C:\Program Files\PostgreSQL\16\bin;%PATH%"
            echo Found PostgreSQL 16 bin directory
        ) else if exist "C:\Program Files\PostgreSQL\15\bin\pg_dump.exe" (
            set "PATH=C:\Program Files\PostgreSQL\15\bin;%PATH%"
            echo Found PostgreSQL 15 bin directory
        ) else if exist "C:\Program Files\PostgreSQL\14\bin\pg_dump.exe" (
            set "PATH=C:\Program Files\PostgreSQL\14\bin;%PATH%"
            echo Found PostgreSQL 14 bin directory
        )
    )
    
    pg_dump -U postgres -h localhost -p 5432 -d inventory_db > database\backup.sql
)

if %errorlevel% equ 0 (
    echo Backup completed successfully: database/backup.sql
) else (
    echo Backup failed. Please verify database connection settings.
)
pause
