package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"product-inventory/internal/api"
	"product-inventory/internal/auth"
	"product-inventory/internal/database"
	"product-inventory/internal/database/db"
	"product-inventory/internal/order"
	"product-inventory/internal/product"
)

func main() {
	// Load environment variables from .env if present
	loadEnv()

	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://postgres:postgres@localhost:5432/inventory_db?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	// 1. Initialize Postgres connection pool
	pool, err := database.NewPostgresPool(ctx, dbSource)
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to PostgreSQL successfully")

	// 2. Run Database Migrations on startup
	if err := runMigrations(ctx, pool); err != nil {
		log.Printf("Migration warning: %v", err)
	} else {
		log.Println("Database migrations applied successfully")
	}

	// 3. Initialize layers
	txManager := database.NewTxManager(pool)
	queries := db.New(pool)

	jwtSecret := os.Getenv("JWT_SECRET")
	envMode := os.Getenv("ENV")
	if envMode == "production" || envMode == "prod" {
		if jwtSecret == "" || jwtSecret == "fallback-jwt-secret-key-for-local-testing" {
			log.Fatal("Fatal Error: Unsafe or missing JWT_SECRET configured in production. Server refused to start.")
		}
	} else {
		if jwtSecret == "" {
			jwtSecret = "fallback-jwt-secret-key-for-local-testing"
			log.Println("Warning: No JWT_SECRET configured. Falling back to local testing default key.")
		}
	}

	// Auth Module
	authService := auth.NewService(queries, jwtSecret)
	authHandler := auth.NewHandler(authService)

	// Product Module
	productRepo := product.NewRepository(queries)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	// Order Module
	orderService := order.NewService(pool, txManager)
	orderHandler := order.NewHandler(orderService)

	// 4. Initialize HTTP router
	r := api.NewRouter(productHandler, orderHandler, authHandler, jwtSecret)

	// 5. Start HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting API server on port %s in %s environment...", port, os.Getenv("ENV"))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Fatal: Server failed: %v", err)
	}
}

// loadEnv parses a local .env file line by line and sets environment variables if they are not already set.
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return // Ignore error, file may not exist (e.g. in Docker production)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}
}

// runMigrations reads all *.up.sql files and applies them in alphabetical order.
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	files, err := os.ReadDir("db/migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var upFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			upFiles = append(upFiles, f.Name())
		}
	}
	sort.Strings(upFiles)

	for _, fileName := range upFiles {
		filePath := filepath.Join("db/migrations", fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		log.Printf("Applying database migration: %s", fileName)
		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			log.Printf("Migration warning for %s (continuing): %v", fileName, err)
		}
	}
	return nil
}
