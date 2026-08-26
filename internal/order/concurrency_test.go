package order

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"product-inventory/internal/database"
	"product-inventory/internal/database/db"
)

func TestOrderService_ConcurrencyAndRollback(t *testing.T) {
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://postgres:postgres@localhost:5432/inventory_db?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		t.Skipf("Skipping database integration tests: failed to connect to database: %v", err)
		return
	}
	defer pool.Close()

	// Verify database connection is alive
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("Skipping database integration tests: database ping failed: %v", err)
		return
	}

	// Drop existing tables to ensure migrations can run cleanly in the test context
	_, err = pool.Exec(ctx, "DROP TABLE IF EXISTS order_items, orders, products, users CASCADE;")
	assert.NoError(t, err)

	// Run migrations to ensure tables exist
	migrationContent1, err := os.ReadFile("../../db/migrations/001_init.up.sql")
	assert.NoError(t, err)
	_, err = pool.Exec(ctx, string(migrationContent1))
	assert.NoError(t, err)

	migrationContent2, err := os.ReadFile("../../db/migrations/002_auth.up.sql")
	assert.NoError(t, err)
	_, err = pool.Exec(ctx, string(migrationContent2))
	assert.NoError(t, err)

	migrationContent3, err := os.ReadFile("../../db/migrations/003_email_auth.up.sql")
	assert.NoError(t, err)
	_, err = pool.Exec(ctx, string(migrationContent3))
	assert.NoError(t, err)

	txManager := database.NewTxManager(pool)
	queries := db.New(pool)
	orderService := NewService(pool, txManager)

	// Seed a test user
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:        "testconcurrency@example.com",
		PasswordHash: "somehash",
		Role:         "customer",
	})
	assert.NoError(t, err)

	// 2. Insert test product with a stock of 10
	product, err := queries.CreateProduct(ctx, db.CreateProductParams{
		Name:          "Concurrent iPhone",
		Description:   "Stock of 10 items",
		Price:         99900,
		StockQuantity: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int32(10), product.StockQuantity)

	// 3. Spawn 100 concurrent ordering requests
	numRequests := 100
	var wg sync.WaitGroup
	wg.Add(numRequests)

	successCount := 0
	failureCount := 0
	var mu sync.Mutex

	var successfulOrders []int64

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()

			req := CreateOrderRequest{
				Items: []CreateOrderItemRequest{
					{
						ProductID: product.ID,
						Quantity:  1,
					},
				},
			}

			res, err := orderService.Create(ctx, user.ID, req)

			mu.Lock()
			defer mu.Unlock()

			if err == nil && res != nil {
				successCount++
				successfulOrders = append(successfulOrders, res.ID)
			} else {
				failureCount++
			}
		}()
	}

	wg.Wait()

	// 4. Assert correctness of concurrency limits
	assert.Equal(t, 10, successCount, "Exactly 10 concurrent requests must succeed")
	assert.Equal(t, 90, failureCount, "Exactly 90 concurrent requests must fail")

	// 5. Query final stock from DB
	updatedProduct, err := queries.GetProduct(ctx, product.ID)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), updatedProduct.StockQuantity, "Final stock quantity must be exactly 0 (no negative stock)")

	// 6. Test Double Cancellation and Concurrency
	// We will attempt to cancel the first successful order concurrently multiple times
	if len(successfulOrders) > 0 {
		orderToCancel := successfulOrders[0]

		numCancels := 10
		var cancelWg sync.WaitGroup
		cancelWg.Add(numCancels)

		cancelSuccess := 0
		cancelFailure := 0
		var cancelMu sync.Mutex

		for i := 0; i < numCancels; i++ {
			go func() {
				defer cancelWg.Done()
				err := orderService.Cancel(ctx, orderToCancel, user.ID, "customer")

				cancelMu.Lock()
				defer cancelMu.Unlock()

				if err == nil {
					cancelSuccess++
				} else {
					cancelFailure++
				}
			}()
		}

		cancelWg.Wait()

		// Assert that exactly one cancel request succeeded in changing the state and restoring the stock
		assert.Equal(t, 1, cancelSuccess, "Exactly 1 cancellation attempt must succeed")
		assert.Equal(t, numCancels-1, cancelFailure, "All other cancellation attempts must fail (preventing double release)")

		// Check database stock is restored by exactly 1 item
		restoredProduct, err := queries.GetProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Equal(t, int32(1), restoredProduct.StockQuantity, "Stock should be restored to exactly 1")
	}
}

func TestOrderService_ListOrders(t *testing.T) {
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://postgres:postgres@localhost:5432/inventory_db?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		t.Skipf("Skipping database integration tests: failed to connect to database: %v", err)
		return
	}
	defer pool.Close()

	txManager := database.NewTxManager(pool)
	queries := db.New(pool)
	orderService := NewService(pool, txManager)

	// Clean tables
	_, err = pool.Exec(ctx, "DELETE FROM order_items; DELETE FROM orders; DELETE FROM products; DELETE FROM users;")
	assert.NoError(t, err)

	// Seed user
	user1, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:        "listuser1@example.com",
		PasswordHash: "hash1",
		Role:         "customer",
	})
	assert.NoError(t, err)

	user2, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:        "listuser2@example.com",
		PasswordHash: "hash2",
		Role:         "customer",
	})
	assert.NoError(t, err)

	// Seed product
	product, err := queries.CreateProduct(ctx, db.CreateProductParams{
		Name:          "List Test Product",
		Description:   "Product details",
		Price:         1000,
		StockQuantity: 100,
	})
	assert.NoError(t, err)

	// Create 3 orders for user1 and 1 order for user2
	for i := 0; i < 3; i++ {
		_, err = orderService.Create(ctx, user1.ID, CreateOrderRequest{
			Items: []CreateOrderItemRequest{
				{ProductID: product.ID, Quantity: 1},
			},
		})
		assert.NoError(t, err)
	}

	_, err = orderService.Create(ctx, user2.ID, CreateOrderRequest{
		Items: []CreateOrderItemRequest{
			{ProductID: product.ID, Quantity: 1},
		},
	})
	assert.NoError(t, err)

	// List user1 orders (page 1, limit 2) - should return 2 orders
	orders1, err := orderService.List(ctx, user1.ID, 1, 2)
	assert.NoError(t, err)
	assert.Len(t, orders1, 2)

	// List user1 orders (page 2, limit 2) - should return 1 order
	orders2, err := orderService.List(ctx, user1.ID, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, orders2, 1)

	// List user2 orders (page 1, limit 5) - should return 1 order
	ordersUser2, err := orderService.List(ctx, user2.ID, 1, 5)
	assert.NoError(t, err)
	assert.Len(t, ordersUser2, 1)
}
