package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cart-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/logger"
)

type cartUC struct {
	repo   domain.CartRepository
	logger logger.Interface
}

func NewCartUC(repo domain.CartRepository, l logger.Interface) Cart {
	return &cartUC{
		repo:   repo,
		logger: l,
	}
}

func (uc *cartUC) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidUserID
	}

	cart, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			return domain.NewCart(userID), nil
		}
		return nil, fmt.Errorf("cartUC.GetCart - repo.GetByUserID: %w", err)
	}

	return cart, nil
}

func (uc *cartUC) AddItem(ctx context.Context, userID string, sku string, quantity int) (*domain.Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidUserID
	}

	cart, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			cart = domain.NewCart(userID)
		} else {
			return nil, fmt.Errorf("cartUC.AddItem - repo.GetByUserID: %w", err)
		}
	}

	if err := cart.AddItem(sku, quantity); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, cart); err != nil {
		return nil, fmt.Errorf("cartUC.AddItem - repo.Save: %w", err)
	}

	return cart, nil
}

func (uc *cartUC) UpdateItemQuantity(ctx context.Context, userID string, sku string, quantity int) (*domain.Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidUserID
	}

	cart, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			return nil, domain.ErrCartNotFound
		}
		return nil, fmt.Errorf("cartUC.UpdateItemQuantity - repo.GetByUserID: %w", err)
	}

	if err := cart.UpdateItem(sku, quantity); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, cart); err != nil {
		return nil, fmt.Errorf("cartUC.UpdateItemQuantity - repo.Save: %w", err)
	}

	return cart, nil
}

func (uc *cartUC) RemoveItem(ctx context.Context, userID string, sku string) (*domain.Cart, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidUserID
	}

	cart, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			return nil, domain.ErrCartNotFound
		}
		return nil, fmt.Errorf("cartUC.RemoveItem - repo.GetByUserID: %w", err)
	}

	if err := cart.RemoveItem(sku); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, cart); err != nil {
		return nil, fmt.Errorf("cartUC.RemoveItem - repo.Save: %w", err)
	}

	return cart, nil
}

func (uc *cartUC) ClearCart(ctx context.Context, userID string) error {
	if strings.TrimSpace(userID) == "" {
		return domain.ErrInvalidUserID
	}

	if err := uc.repo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("cartUC.ClearCart - repo.Delete: %w", err)
	}

	return nil
}
