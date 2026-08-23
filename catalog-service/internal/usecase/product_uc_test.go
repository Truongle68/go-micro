package usecase_test

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/usecase"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockProductRepo struct {
	usecase.ProductRepository
	createFn    func(ctx context.Context, p *domain.Product) error
	existSlugFn func(ctx context.Context, name string) (bool, error)
	findByIDFn  func(ctx context.Context, id string) (*domain.Product, error)
	updateFn    func(ctx context.Context, p *domain.Product, expectedVersion int) (*domain.Product, error)
}

func (m *mockProductRepo) EnsureIndexes(ctx context.Context) error {
	return nil
}

func (m *mockProductRepo) Create(ctx context.Context, p *domain.Product) error {
	if m.createFn != nil {
		return m.createFn(ctx, p)
	}
	p.ID = "507f1f77bcf86cd799439011"
	return nil
}

func (m *mockProductRepo) ExistSlug(ctx context.Context, name string) (bool, error) {
	if m.existSlugFn != nil {
		return m.existSlugFn(ctx, name)
	}
	return false, nil
}

func (m *mockProductRepo) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return &domain.Product{
		ID:          id,
		Version:     1,
		Name:        "Existing Product",
		CategoryID:  "507f1f77bcf86cd799439011",
		OptionTypes: []domain.OptionType{{Name: "Color", Values: []string{"Red", "Blue"}}},
		Variants: []domain.Variant{
			{SKU: "SKU-RED", Attributes: map[string]string{"Color": "Red"}},
		},
	}, nil
}

func (m *mockProductRepo) Update(ctx context.Context, p *domain.Product, expectedVersion int) (*domain.Product, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, p, expectedVersion)
	}
	return p, nil
}

type mockCategoryRepo struct {
	usecase.CategoryRepository
}

func (m *mockCategoryRepo) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	return &domain.Category{ID: id, Name: "Electronics"}, nil
}

func (m *mockCategoryRepo) BuildBreadcrumb(ctx context.Context, id string) ([]domain.CategoryRef, error) {
	return []domain.CategoryRef{{ID: id, Name: "Electronics"}}, nil
}

func TestCreateProduct_ImageExtraction(t *testing.T) {
	pRepo := &mockProductRepo{}
	cRepo := &mockCategoryRepo{}
	uc := usecase.NewProductUC(pRepo, cRepo)

	input := usecase.CreateProductInput{
		Name:       "Test Laptop",
		CategoryID: "507f1f77bcf86cd799439011",
		Variants: []usecase.CreateVariantInput{
			{SKU: "SKU-001", Price: usecase.PriceInput{Amount: 1000, Currency: "USD"}, Image: "https://example.com/img1.png"},
		},
		Status: string(domain.ProductStatusDraft),
	}

	product, err := uc.Create(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "https://example.com/img1.png", product.Images[0])
}

func TestCreateProduct_InvalidAttributeValue(t *testing.T) {
	pRepo := &mockProductRepo{}
	cRepo := &mockCategoryRepo{}
	uc := usecase.NewProductUC(pRepo, cRepo)

	input := usecase.CreateProductInput{
		Name:       "Test Shirt",
		CategoryID: "507f1f77bcf86cd799439011",
		OptionTypes: []usecase.OptionTypeInput{
			{Name: "Color", Values: []string{"Red", "Blue"}},
		},
		Variants: []usecase.CreateVariantInput{
			{
				SKU:        "SHIRT-GREEN",
				Attributes: map[string]string{"Color": "Green"},
				Price:      usecase.PriceInput{Amount: 2000, Currency: "USD"},
			},
		},
		Status: string(domain.ProductStatusDraft),
	}

	_, err := uc.Create(context.Background(), input)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidVariantAttribute))
}

func TestCreateProduct_SlugCollisionRetry(t *testing.T) {
	attempts := 0
	pRepo := &mockProductRepo{
		createFn: func(ctx context.Context, p *domain.Product) error {
			attempts++
			if attempts == 1 {
				return domain.ErrDuplicateSlug
			}
			p.ID = "507f1f77bcf86cd799439011"
			return nil
		},
	}
	cRepo := &mockCategoryRepo{}
	uc := usecase.NewProductUC(pRepo, cRepo)

	input := usecase.CreateProductInput{
		Name:       "Sample Item",
		CategoryID: "507f1f77bcf86cd799439011",
		Variants: []usecase.CreateVariantInput{
			{SKU: "ITEM-001", Price: usecase.PriceInput{Amount: 500, Currency: "USD"}},
		},
		Status: string(domain.ProductStatusDraft),
	}

	product, err := uc.Create(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, 2, attempts)
	assert.Equal(t, "sample-item-1", product.Slug)
}

func TestUpdateProduct_OptionTypesValidationOnExistingVariants(t *testing.T) {
	pRepo := &mockProductRepo{
		findByIDFn: func(ctx context.Context, id string) (*domain.Product, error) {
			return &domain.Product{
				ID:          id,
				Version:     1,
				Name:        "Shirt",
				CategoryID:  "507f1f77bcf86cd799439011",
				OptionTypes: []domain.OptionType{{Name: "Color", Values: []string{"Red", "Blue"}}},
				Variants: []domain.Variant{
					{SKU: "SHIRT-RED", Attributes: map[string]string{"Color": "Red"}},
				},
			}, nil
		},
	}
	cRepo := &mockCategoryRepo{}
	uc := usecase.NewProductUC(pRepo, cRepo)

	// Admin attempts to update OptionTypes to only "Yellow", leaving existing variant "SHIRT-RED" invalid
	input := usecase.UpdateProductInput{
		ID:      "507f1f77bcf86cd799439011",
		Version: 1,
		OptionTypes: []usecase.OptionTypeInput{
			{Name: "Color", Values: []string{"Yellow"}},
		},
	}

	_, err := uc.Update(context.Background(), input)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidVariantAttribute))
}
