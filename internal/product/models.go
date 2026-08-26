package product

import (
	"errors"
	"time"
)

// CreateProductRequest defines the payload for creating a product.
type CreateProductRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         int64  `json:"price"` // in cents
	StockQuantity int32  `json:"stock_quantity"`
}

// Validate validates the CreateProductRequest.
func (r *CreateProductRequest) Validate() error {
	if r.Name == "" {
		return errors.New("product name is required")
	}
	if r.Price <= 0 {
		return errors.New("product price must be greater than zero")
	}
	if r.StockQuantity < 0 {
		return errors.New("product stock quantity cannot be negative")
	}
	return nil
}

// UpdateProductRequest defines the payload for updating a product.
type UpdateProductRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         int64  `json:"price"`
	StockQuantity int32  `json:"stock_quantity"`
}

// Validate validates the UpdateProductRequest.
func (r *UpdateProductRequest) Validate() error {
	if r.Name == "" {
		return errors.New("product name is required")
	}
	if r.Price <= 0 {
		return errors.New("product price must be greater than zero")
	}
	if r.StockQuantity < 0 {
		return errors.New("product stock quantity cannot be negative")
	}
	return nil
}

// ProductResponse defines the JSON structure of a product in API responses.
type ProductResponse struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         int64     `json:"price"`
	StockQuantity int32     `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
