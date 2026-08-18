package usecase_test

import (
	"context"
	"errors"
	"testing"

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
	l := logger.New("error")
	uc := usecase.NewCartUC(repo, l)

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

	// 5. ClearCart
	err = uc.ClearCart(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error clearing cart: %v", err)
	}

	// 6. UpdateItemQuantity on missing cart -> ErrCartNotFound
	_, err = uc.UpdateItemQuantity(ctx, userID, "SKU-B", 1)
	if !errors.Is(err, domain.ErrCartNotFound) {
		t.Fatalf("expected ErrCartNotFound, got %v", err)
	}
}
