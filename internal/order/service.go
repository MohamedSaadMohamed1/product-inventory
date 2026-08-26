package order

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"product-inventory/internal/database"
	"product-inventory/internal/database/db"
)

var (
	ErrOrderNotFound          = errors.New("order not found")
	ErrOrderAlreadyCancelled  = errors.New("order is already cancelled")
	ErrInvalidOrderState      = errors.New("order is not in PENDING state")
	ErrInsufficientStock      = errors.New("insufficient stock")
	ErrProductNotFound        = errors.New("product not found")
	ErrUnauthorized           = errors.New("unauthorized")
)

// Service coordinates order business logic and transaction boundaries.
type Service struct {
	pool      *pgxpool.Pool
	txManager *database.TxManager
}

// NewService creates a new Service instance.
func NewService(pool *pgxpool.Pool, txManager *database.TxManager) *Service {
	return &Service{
		pool:      pool,
		txManager: txManager,
	}
}

// Create places a new order, validating and reserving stock within a transaction.
func (s *Service) Create(ctx context.Context, userID int64, req CreateOrderRequest) (*OrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	var response *OrderResponse

	// Execute order creation inside a transaction
	err := s.txManager.ExecTx(ctx, func(q *db.Queries) error {
		txRepo := NewRepository(q)

		// 1. Extract product IDs and sort them to prevent deadlocks
		productIDs := make([]int64, len(req.Items))
		for i, item := range req.Items {
			productIDs[i] = item.ProductID
		}
		sort.Slice(productIDs, func(i, j int) bool {
			return productIDs[i] < productIDs[j]
		})

		// 2. Fetch and lock product rows
		products, err := txRepo.GetProductsForUpdate(ctx, productIDs)
		if err != nil {
			return fmt.Errorf("failed to lock product rows: %w", err)
		}

		// Ensure all requested products exist
		productMap := make(map[int64]db.Products)
		for _, p := range products {
			productMap[p.ID] = p
		}

		for _, item := range req.Items {
			if _, exists := productMap[item.ProductID]; !exists {
				return fmt.Errorf("%w: product ID %d does not exist", ErrProductNotFound, item.ProductID)
			}
		}

		// 3. Validate stock and calculate total amount
		var totalAmount int64
		for _, item := range req.Items {
			p := productMap[item.ProductID]
			if p.StockQuantity < item.Quantity {
				return fmt.Errorf("%w: product %d (%s) has only %d available, requested %d",
					ErrInsufficientStock, p.ID, p.Name, p.StockQuantity, item.Quantity)
			}
			totalAmount += p.Price * int64(item.Quantity)
		}

		// 4. Reserve stock by updating database stock quantities
		for _, item := range req.Items {
			_, err := txRepo.UpdateProductStock(ctx, db.UpdateProductStockParams{
				ID:            item.ProductID,
				StockQuantity: -item.Quantity, // Decrement stock
			})
			if err != nil {
				return fmt.Errorf("failed to reserve stock for product %d: %w", item.ProductID, err)
			}
		}

		// 5. Create Order row linked to user
		order, err := txRepo.CreateOrder(ctx, string(OrderStatusPending), totalAmount, userID)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// 6. Create Order Items
		itemResponses := make([]OrderItemResponse, len(req.Items))
		for i, item := range req.Items {
			p := productMap[item.ProductID]
			orderItem, err := txRepo.CreateOrderItem(ctx, db.CreateOrderItemParams{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: p.Price,
			})
			if err != nil {
				return fmt.Errorf("failed to insert order item for product %d: %w", item.ProductID, err)
			}

			itemResponses[i] = OrderItemResponse{
				ID:        orderItem.ID,
				OrderID:   orderItem.OrderID,
				ProductID: orderItem.ProductID,
				Quantity:  orderItem.Quantity,
				UnitPrice: orderItem.UnitPrice,
				CreatedAt: orderItem.CreatedAt.Time,
			}
		}

		response = &OrderResponse{
			ID:          order.ID,
			Status:      OrderStatus(order.Status),
			TotalAmount: order.TotalAmount,
			CreatedAt:   order.CreatedAt.Time,
			UpdatedAt:   order.UpdatedAt.Time,
			Items:       itemResponses,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// Get retrieves an order by its ID, including its associated items.
// Ownership check is performed: only admins or the owner of the order can view it.
func (s *Service) Get(ctx context.Context, id int64, userID int64, userRole string) (*OrderResponse, error) {
	queries := db.New(s.pool)
	repo := NewRepository(queries)

	order, err := repo.GetOrder(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Auth check: Admin can view any order, Customer can only view their own
	if userRole != "admin" && (!order.UserID.Valid || order.UserID.Int64 != userID) {
		return nil, ErrUnauthorized
	}

	items, err := repo.ListOrderItems(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list order items: %w", err)
	}

	itemResponses := make([]OrderItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = OrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			CreatedAt: item.CreatedAt.Time,
		}
	}

	var cancelledAt *time.Time
	if order.CancelledAt.Valid {
		cancelledAt = &order.CancelledAt.Time
	}

	return &OrderResponse{
		ID:          order.ID,
		Status:      OrderStatus(order.Status),
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt.Time,
		UpdatedAt:   order.UpdatedAt.Time,
		CancelledAt: cancelledAt,
		Items:       itemResponses,
	}, nil
}

// Cancel cancels a PENDING order, releasing the reserved stock back to inventory.
// Ownership check is performed: only admins, system, or the owner of the order can cancel it.
func (s *Service) Cancel(ctx context.Context, id int64, userID int64, userRole string) error {
	return s.txManager.ExecTx(ctx, func(q *db.Queries) error {
		txRepo := NewRepository(q)

		// 1. Lock the order row to ensure exclusive access and prevent duplicate releases
		order, err := txRepo.GetOrderForUpdate(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrOrderNotFound
			}
			return fmt.Errorf("failed to lock order row: %w", err)
		}

		// Auth check: Admin or System can cancel any order, Customer can only cancel their own
		if userRole != "admin" && userRole != "system" && (!order.UserID.Valid || order.UserID.Int64 != userID) {
			return ErrUnauthorized
		}

		// 2. Validate current state transitions
		if OrderStatus(order.Status) == OrderStatusCancelled {
			return ErrOrderAlreadyCancelled
		}
		if OrderStatus(order.Status) != OrderStatusPending {
			return fmt.Errorf("%w: expected PENDING but got %s", ErrInvalidOrderState, order.Status)
		}

		// 3. Fetch order items
		items, err := txRepo.ListOrderItems(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to retrieve order items: %w", err)
		}

		if len(items) > 0 {
			// 4. Extract product IDs and sort them to prevent deadlocks when locking products
			productIDs := make([]int64, len(items))
			for i, item := range items {
				productIDs[i] = item.ProductID
			}
			sort.Slice(productIDs, func(i, j int) bool {
				return productIDs[i] < productIDs[j]
			})

			// 5. Lock product rows
			_, err = txRepo.GetProductsForUpdate(ctx, productIDs)
			if err != nil {
				return fmt.Errorf("failed to lock product rows during cancellation: %w", err)
			}

			// 6. Restore stock quantities to products
			for _, item := range items {
				_, err = txRepo.UpdateProductStock(ctx, db.UpdateProductStockParams{
					ID:            item.ProductID,
					StockQuantity: item.Quantity, // Add back quantity to stock
				})
				if err != nil {
					return fmt.Errorf("failed to restore stock for product %d: %w", item.ProductID, err)
				}
			}
		}

		// 7. Update order status to CANCELLED
		_, err = txRepo.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
			ID:     id,
			Status: string(OrderStatusCancelled),
			CancelledAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return nil
	})
}

// List retrieves a paginated list of orders for a user.
func (s *Service) List(ctx context.Context, userID int64, page, limit int) ([]*OrderResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	queries := db.New(s.pool)
	repo := NewRepository(queries)

	orders, err := repo.ListOrdersByUserID(ctx, userID, int32(limit), int32(offset))
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	responses := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		items, err := repo.ListOrderItems(ctx, o.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list order items for order %d: %w", o.ID, err)
		}

		itemResponses := make([]OrderItemResponse, len(items))
		for j, item := range items {
			itemResponses[j] = OrderItemResponse{
				ID:        item.ID,
				OrderID:   item.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
				CreatedAt: item.CreatedAt.Time,
			}
		}

		var cancelledAt *time.Time
		if o.CancelledAt.Valid {
			cancelledAt = &o.CancelledAt.Time
		}

		responses[i] = &OrderResponse{
			ID:          o.ID,
			Status:      OrderStatus(o.Status),
			TotalAmount: o.TotalAmount,
			CreatedAt:   o.CreatedAt.Time,
			UpdatedAt:   o.UpdatedAt.Time,
			CancelledAt: cancelledAt,
			Items:       itemResponses,
		}
	}

	return responses, nil
}
