@echo off
echo ===================================
echo   PostgreSQL Database Backup Tool
echo ===================================
cd /d "%~dp0.."
mkdir database 2>nul

docker ps -q -f name=inventory_postgres >nul 2>nul
if %errorlevel% equ 0 (
    echo Detected Docker container: inventory_postgres
    echo Creating backup via Docker pg_dump...
    docker exec -t inventory_postgres pg_dump -U postgres -d inventory_db > database\backup.sql
) else (
    echo Docker container not running. Attempting local pg_dump...
    echo Please make sure pg_dump is in your PATH.
    pg_dump -U postgres -d inventory_db > database\backup.sql
)

if %errorlevel% equ 0 (
    echo Backup completed successfully: database/backup.sql
) else (
    echo Backup failed. Please verify database connection settings.
)
pause
