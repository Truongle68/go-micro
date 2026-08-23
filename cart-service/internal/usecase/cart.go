package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cart-service/internal/client"
	"cart-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/logger"
)

type CartUC struct {
	repo          CartRepository
	catalogClient client.CatalogClient
	logger        logger.Interface
}

func NewCartUC(repo CartRepository, client client.CatalogClient, l logger.Interface) *CartUC {
	return &CartUC{
		repo:          repo,
		catalogClient: client,
		logger:        l,
	}
}

func (uc *CartUC) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidUserID
	}

	cart, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			return domain.NewCart(userID), nil
		}
		return nil, fmt.Errorf("CartUC.GetCart - repo.GetByUserID: %w", err)
	}

	return cart, nil
}

func (uc *CartUC) AddItem(ctx context.Context, userID string, sku string, quantity int) (*domain.Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidUserID
	}

	cart, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			cart = domain.NewCart(userID)
		} else {
			return nil, fmt.Errorf("CartUC.AddItem - repo.GetByUserID: %w", err)
		}
	}

	if err := cart.AddItem(sku, quantity); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, cart); err != nil {
		return nil, fmt.Errorf("CartUC.AddItem - repo.Save: %w", err)
	}

	return cart, nil
}

func (uc *CartUC) UpdateItemQuantity(ctx context.Context, userID string, sku string, quantity int) (*domain.Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidUserID
	}

	cart, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			return nil, domain.ErrCartNotFound
		}
		return nil, fmt.Errorf("CartUC.UpdateItemQuantity - repo.GetByUserID: %w", err)
	}

	if err := cart.UpdateItem(sku, quantity); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, cart); err != nil {
		return nil, fmt.Errorf("CartUC.UpdateItemQuantity - repo.Save: %w", err)
	}

	return cart, nil
}

func (uc *CartUC) RemoveItem(ctx context.Context, userID string, sku string) (*domain.Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidUserID
	}

	cart, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			return nil, domain.ErrCartNotFound
		}
		return nil, fmt.Errorf("CartUC.RemoveItem - repo.GetByUserID: %w", err)
	}

	if err := cart.RemoveItem(sku); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, cart); err != nil {
		return nil, fmt.Errorf("CartUC.RemoveItem - repo.Save: %w", err)
	}

	return cart, nil
}

func (uc *CartUC) ClearCart(ctx context.Context, userID string) error {
	if strings.TrimSpace(userID) == "" {
		return domain.ErrInvalidUserID
	}

	if err := uc.repo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("CartUC.ClearCart - repo.Delete: %w", err)
	}

	return nil
}
