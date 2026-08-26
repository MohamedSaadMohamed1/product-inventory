package order

import (
	"errors"
	"time"
)

// OrderStatus represents the current state of an order.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// CreateOrderItemRequest defines the payload for an individual item in a create order request.
type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

// CreateOrderRequest defines the payload for creating an order.
type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}

// Validate validates the CreateOrderRequest.
func (r *CreateOrderRequest) Validate() error {
	if len(r.Items) == 0 {
		return errors.New("order must contain at least one item")
	}

	seenProducts := make(map[int64]bool)
	for _, item := range r.Items {
		if item.ProductID <= 0 {
			return errors.New("invalid product ID in order items")
		}
		if item.Quantity <= 0 {
			return errors.New("quantity must be greater than zero")
		}
		if seenProducts[item.ProductID] {
			return errors.New("duplicate product ID in order; consolidate quantities instead")
		}
		seenProducts[item.ProductID] = true
	}
	return nil
}

// OrderItemResponse represents an item inside an order response.
type OrderItemResponse struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int32     `json:"quantity"`
	UnitPrice int64     `json:"unit_price"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderResponse represents the response details of an order.
type OrderResponse struct {
	ID          int64               `json:"id"`
	Status      OrderStatus         `json:"status"`
	TotalAmount int64               `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	CancelledAt *time.Time          `json:"cancelled_at,omitempty"`
	Items       []OrderItemResponse `json:"items"`
}
