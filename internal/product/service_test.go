package product

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"product-inventory/internal/database/db"
)

type mockRepository struct {
	products map[int64]db.Products
	nextID   int64
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		products: make(map[int64]db.Products),
		nextID:   1,
	}
}

func (m *mockRepository) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Products, error) {
	p := db.Products{
		ID:            m.nextID,
		Name:          arg.Name,
		Description:   arg.Description,
		Price:         arg.Price,
		StockQuantity: arg.StockQuantity,
	}
	m.products[m.nextID] = p
	m.nextID++
	return p, nil
}

func (m *mockRepository) GetProduct(ctx context.Context, id int64) (db.Products, error) {
	p, exists := m.products[id]
	if !exists {
		return db.Products{}, pgx.ErrNoRows
	}
	return p, nil
}

func (m *mockRepository) ListProducts(ctx context.Context, arg db.ListProductsParams) ([]db.Products, error) {
	var list []db.Products
	// Simple mock pagination
	start := int(arg.Offset)
	count := 0
	for i := int64(1); i < m.nextID; i++ {
		p, exists := m.products[i]
		if exists {
			if count >= start && count < start+int(arg.Limit) {
				list = append(list, p)
			}
			count++
		}
	}
	return list, nil
}

func (m *mockRepository) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Products, error) {
	p, exists := m.products[arg.ID]
	if !exists {
		return db.Products{}, pgx.ErrNoRows
	}
	p.Name = arg.Name
	p.Description = arg.Description
	p.Price = arg.Price
	p.StockQuantity = arg.StockQuantity
	m.products[arg.ID] = p
	return p, nil
}

func (m *mockRepository) DeleteProduct(ctx context.Context, id int64) error {
	delete(m.products, id)
	return nil
}

func TestProductService_Create(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		req := CreateProductRequest{
			Name:          "Laptop",
			Description:   "M3 Pro",
			Price:         199900,
			StockQuantity: 10,
		}
		res, err := service.Create(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, int64(1), res.ID)
		assert.Equal(t, "Laptop", res.Name)
		assert.Equal(t, int32(10), res.StockQuantity)
	})

	t.Run("invalid name", func(t *testing.T) {
		req := CreateProductRequest{
			Name:          "",
			Description:   "M3 Pro",
			Price:         199900,
			StockQuantity: 10,
		}
		res, err := service.Create(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "product name is required")
	})

	t.Run("invalid price", func(t *testing.T) {
		req := CreateProductRequest{
			Name:          "Laptop",
			Price:         -10,
			StockQuantity: 10,
		}
		res, err := service.Create(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "product price must be greater than zero")
	})

	t.Run("invalid stock quantity", func(t *testing.T) {
		req := CreateProductRequest{
			Name:          "Laptop",
			Price:         1000,
			StockQuantity: -5,
		}
		res, err := service.Create(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "product stock quantity cannot be negative")
	})
}

func TestProductService_Get(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo)
	ctx := context.Background()

	// Seed product
	p, _ := repo.CreateProduct(ctx, db.CreateProductParams{
		Name:          "Tablet",
		Price:         49900,
		StockQuantity: 5,
	})

	t.Run("get existing", func(t *testing.T) {
		res, err := service.Get(ctx, p.ID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, p.ID, res.ID)
		assert.Equal(t, "Tablet", res.Name)
	})

	t.Run("get non-existent", func(t *testing.T) {
		res, err := service.Get(ctx, 999)
		assert.ErrorIs(t, err, ErrProductNotFound)
		assert.Nil(t, res)
	})
}

func TestProductService_Update(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo)
	ctx := context.Background()

	p, _ := repo.CreateProduct(ctx, db.CreateProductParams{
		Name:          "Phone",
		Price:         79900,
		StockQuantity: 15,
	})

	t.Run("successful update", func(t *testing.T) {
		req := UpdateProductRequest{
			Name:          "Phone Pro",
			Description:   "Updated description",
			Price:         89900,
			StockQuantity: 20,
		}
		res, err := service.Update(ctx, p.ID, req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Phone Pro", res.Name)
		assert.Equal(t, int32(20), res.StockQuantity)
		assert.Equal(t, int64(89900), res.Price)
	})

	t.Run("update non-existent", func(t *testing.T) {
		req := UpdateProductRequest{
			Name:          "Unknown",
			Price:         1000,
			StockQuantity: 5,
		}
		res, err := service.Update(ctx, 999, req)
		assert.ErrorIs(t, err, ErrProductNotFound)
		assert.Nil(t, res)
	})
}

func TestProductService_Delete(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo)
	ctx := context.Background()

	p, _ := repo.CreateProduct(ctx, db.CreateProductParams{
		Name:          "Keyboard",
		Price:         9900,
		StockQuantity: 3,
	})

	t.Run("successful delete", func(t *testing.T) {
		err := service.Delete(ctx, p.ID)
		assert.NoError(t, err)

		// Assert it's deleted
		_, err = service.Get(ctx, p.ID)
		assert.ErrorIs(t, err, ErrProductNotFound)
	})

	t.Run("delete non-existent", func(t *testing.T) {
		err := service.Delete(ctx, 999)
		assert.ErrorIs(t, err, ErrProductNotFound)
	})
}
