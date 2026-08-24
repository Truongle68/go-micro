package usecase_test

import (
	"context"
	"errors"
	"testing"

	"cart-service/internal/client"
	"cart-service/internal/domain"
	"cart-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/logger"
)

type mockCartRepo struct {
	carts map[string]*domain.Cart
}

func newMockCartRepo() *mockCartRepo {
	return &mockCartRepo{
		carts: make(map[string]*domain.Cart),
	}
}

type mockCatalogClient struct {
	variants map[string]client.VariantDTO
	err      error
}

func newMockCatalogClient(variants ...client.VariantDTO) *mockCatalogClient {
	m := &mockCatalogClient{
		variants: make(map[string]client.VariantDTO),
	}
	for _, v := range variants {
		m.variants[v.SKU] = v
	}
	return m
}

func (m *mockCatalogClient) GetVariantsBySKUs(ctx context.Context, skus []string) ([]client.VariantDTO, error) {
	if m.err != nil {
		return nil, m.err
	}
	out := make([]client.VariantDTO, 0, len(skus))
	for _, sku := range skus {
		if v, ok := m.variants[sku]; ok {
			out = append(out, v)
		}
	}
	return out, nil
}

func (m *mockCartRepo) GetByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	c, ok := m.carts[userID]
	if !ok {
		return nil, domain.ErrCartNotFound
	}
	// return copy
	clone := *c
	clone.Items = append([]domain.CartItem(nil), c.Items...)
	return &clone, nil
}

func (m *mockCartRepo) Save(ctx context.Context, cart *domain.Cart) error {
	clone := *cart
	clone.Items = append([]domain.CartItem(nil), cart.Items...)
	m.carts[cart.UserID] = &clone
	return nil
}

func (m *mockCartRepo) Delete(ctx context.Context, userID string) error {
	delete(m.carts, userID)
	return nil
}

func TestCartUseCaseFlow(t *testing.T) {
	repo := newMockCartRepo()
	client := newMockCatalogClient()
	l := logger.New("error")
	uc := usecase.NewCartUC(repo, client, l)

	ctx := context.Background()
	userID := "user_test_1"

	// 1. GetCart initially -> returns empty cart
	cart, err := uc.GetCart(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error getting cart: %v", err)
	}
	if cart.UserID != userID || len(cart.Items) != 0 {
		t.Fatalf("expected empty cart for user, got %+v", cart)
	}

	// 2. AddItem
	cart, err = uc.AddItem(ctx, userID, "SKU-A", 2)
	if err != nil {
		t.Fatalf("unexpected error adding item: %v", err)
	}
	if len(cart.Items) != 1 || cart.Items[0].Quantity != 2 {
		t.Fatalf("expected 1 item with qty 2, got %+v", cart)
	}

	// 3. UpdateItemQuantity
	cart, err = uc.UpdateItemQuantity(ctx, userID, "SKU-A", 5)
	if err != nil {
		t.Fatalf("unexpected error updating item qty: %v", err)
	}
	if cart.Items[0].Quantity != 5 {
		t.Fatalf("expected qty 5, got %d", cart.Items[0].Quantity)
	}

	// 4. RemoveItem
	cart, err = uc.RemoveItem(ctx, userID, "SKU-A")
	if err != nil {
		t.Fatalf("unexpected error removing item: %v", err)
	}
	if len(cart.Items) != 0 {
		t.Fatalf("expected 0 items after removal, got %d", len(cart.Items))
	}

	// 5. Add two items and test RemoveItems
	_, _ = uc.AddItem(ctx, userID, "SKU-A", 2)
	_, _ = uc.AddItem(ctx, userID, "SKU-B", 3)
	cart, err = uc.RemoveItems(ctx, userID, []string{"SKU-A"})
	if err != nil {
		t.Fatalf("unexpected error removing items: %v", err)
	}
	if len(cart.Items) != 1 || cart.Items[0].SKU != "SKU-B" {
		t.Fatalf("expected only SKU-B remaining, got %+v", cart.Items)
	}

	// 6. ClearCart
	err = uc.ClearCart(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error clearing cart: %v", err)
	}

	// 7. UpdateItemQuantity on missing cart -> ErrCartNotFound
	_, err = uc.UpdateItemQuantity(ctx, userID, "SKU-B", 1)
	if !errors.Is(err, domain.ErrCartNotFound) {
		t.Fatalf("expected ErrCartNotFound, got %v", err)
	}
}
