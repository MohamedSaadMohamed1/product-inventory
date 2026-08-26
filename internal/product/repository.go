package product

import (
	"context"

	"product-inventory/internal/database/db"
)

// Repository defines the contract for product data access.
type Repository interface {
	CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Products, error)
	GetProduct(ctx context.Context, id int64) (db.Products, error)
	ListProducts(ctx context.Context, arg db.ListProductsParams) ([]db.Products, error)
	UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Products, error)
	DeleteProduct(ctx context.Context, id int64) error
}

// SQLRepository implements Repository using database/db generated code.
type SQLRepository struct {
	queries *db.Queries
}

// NewRepository creates a new SQLRepository.
func NewRepository(queries *db.Queries) Repository {
	return &SQLRepository{queries: queries}
}

// CreateProduct inserts a new product.
func (r *SQLRepository) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Products, error) {
	return r.queries.CreateProduct(ctx, arg)
}

// GetProduct retrieves a product by its ID.
func (r *SQLRepository) GetProduct(ctx context.Context, id int64) (db.Products, error) {
	return r.queries.GetProduct(ctx, id)
}

// ListProducts lists products with pagination.
func (r *SQLRepository) ListProducts(ctx context.Context, arg db.ListProductsParams) ([]db.Products, error) {
	return r.queries.ListProducts(ctx, arg)
}

// UpdateProduct updates all fields of a product.
func (r *SQLRepository) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Products, error) {
	return r.queries.UpdateProduct(ctx, arg)
}

// DeleteProduct deletes a product.
func (r *SQLRepository) DeleteProduct(ctx context.Context, id int64) error {
	return r.queries.DeleteProduct(ctx, id)
}
