package order

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"product-inventory/internal/database/db"
)

// Repository defines data access operations for orders, items, and inventory transactions.
type Repository interface {
	CreateOrder(ctx context.Context, status string, totalAmount int64, userID int64) (db.Orders, error)
	GetOrder(ctx context.Context, id int64) (db.Orders, error)
	GetOrderForUpdate(ctx context.Context, id int64) (db.Orders, error)
	UpdateOrderStatus(ctx context.Context, arg db.UpdateOrderStatusParams) (db.Orders, error)
	CreateOrderItem(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItems, error)
	ListOrderItems(ctx context.Context, orderID int64) ([]db.OrderItems, error)
	ListOrdersByUserID(ctx context.Context, userID int64, limit, offset int32) ([]db.Orders, error)
	GetProductsForUpdate(ctx context.Context, ids []int64) ([]db.Products, error)
	UpdateProductStock(ctx context.Context, arg db.UpdateProductStockParams) (db.Products, error)
}

// SQLRepository implements Repository using database/db generated code.
type SQLRepository struct {
	queries *db.Queries
}

// NewRepository creates a new SQLRepository.
func NewRepository(queries *db.Queries) Repository {
	return &SQLRepository{queries: queries}
}

// CreateOrder inserts an order.
func (r *SQLRepository) CreateOrder(ctx context.Context, status string, totalAmount int64, userID int64) (db.Orders, error) {
	return r.queries.CreateOrder(ctx, db.CreateOrderParams{
		Status:      status,
		TotalAmount: totalAmount,
		UserID:      pgtype.Int8{Int64: userID, Valid: true},
	})
}

// GetOrder fetches an order.
func (r *SQLRepository) GetOrder(ctx context.Context, id int64) (db.Orders, error) {
	return r.queries.GetOrder(ctx, id)
}

// GetOrderForUpdate fetches an order with an exclusive row lock.
func (r *SQLRepository) GetOrderForUpdate(ctx context.Context, id int64) (db.Orders, error) {
	return r.queries.GetOrderForUpdate(ctx, id)
}

// UpdateOrderStatus updates the status of an order.
func (r *SQLRepository) UpdateOrderStatus(ctx context.Context, arg db.UpdateOrderStatusParams) (db.Orders, error) {
	return r.queries.UpdateOrderStatus(ctx, arg)
}

// CreateOrderItem inserts an item in an order.
func (r *SQLRepository) CreateOrderItem(ctx context.Context, arg db.CreateOrderItemParams) (db.OrderItems, error) {
	return r.queries.CreateOrderItem(ctx, arg)
}

// ListOrderItems lists items under an order.
func (r *SQLRepository) ListOrderItems(ctx context.Context, orderID int64) ([]db.OrderItems, error) {
	return r.queries.ListOrderItemsByOrderId(ctx, orderID)
}

// GetProductsForUpdate fetches multiple products in order by ID with an exclusive row lock.
func (r *SQLRepository) GetProductsForUpdate(ctx context.Context, ids []int64) ([]db.Products, error) {
	return r.queries.GetProductsForUpdate(ctx, ids)
}

// UpdateProductStock adjusts the stock quantity of a product (positive to add, negative to subtract).
func (r *SQLRepository) UpdateProductStock(ctx context.Context, arg db.UpdateProductStockParams) (db.Products, error) {
	return r.queries.UpdateProductStock(ctx, arg)
}

// ListOrdersByUserID lists paginated orders for a user.
func (r *SQLRepository) ListOrdersByUserID(ctx context.Context, userID int64, limit, offset int32) ([]db.Orders, error) {
	return r.queries.ListOrdersByUserID(ctx, db.ListOrdersByUserIDParams{
		UserID: pgtype.Int8{Int64: userID, Valid: true},
		Limit:  limit,
		Offset: offset,
	})
}
