package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"product-inventory/internal/database/db"
)

var (
	// ErrProductNotFound occurs when a product does not exist.
	ErrProductNotFound = errors.New("product not found")
)

// Service coordinates product business operations.
type Service struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create validates and creates a new product.
func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*ProductResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	p, err := s.repo.CreateProduct(ctx, db.CreateProductParams{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return s.mapToResponse(p), nil
}

// Get fetches a product by ID.
func (s *Service) Get(ctx context.Context, id int64) (*ProductResponse, error) {
	p, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}
	return s.mapToResponse(p), nil
}

// List returns a paginated list of products.
func (s *Service) List(ctx context.Context, page, limit int32) ([]ProductResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	products, err := s.repo.ListProducts(ctx, db.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	res := make([]ProductResponse, len(products))
	for i, p := range products {
		res[i] = *s.mapToResponse(p)
	}
	return res, nil
}

// Update validates and updates a product's details.
func (s *Service) Update(ctx context.Context, id int64, req UpdateProductRequest) (*ProductResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	_, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to fetch product for update: %w", err)
	}

	p, err := s.repo.UpdateProduct(ctx, db.UpdateProductParams{
		ID:            id,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return s.mapToResponse(p), nil
}

// Delete deletes a product by ID.
func (s *Service) Delete(ctx context.Context, id int64) error {
	_, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return fmt.Errorf("failed to fetch product for deletion: %w", err)
	}

	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (s *Service) mapToResponse(p db.Products) *ProductResponse {
	return &ProductResponse{
		ID:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		CreatedAt:     p.CreatedAt.Time,
		UpdatedAt:     p.UpdatedAt.Time,
	}
}
