package usecase

import (
	"context"
	"fmt"
	"time"

	"inventory-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type SupplierUC struct {
	repo   SupplierRepository
	logger logger.Interface
}

func NewSupplierUC(repo SupplierRepository, l logger.Interface) *SupplierUC {
	return &SupplierUC{
		repo:   repo,
		logger: l,
	}
}

func (uc *SupplierUC) CreateSupplier(ctx context.Context, code, name, phone string, address domain.SupplierAddress) (*domain.Supplier, error) {
	s, err := domain.NewSupplier(code, name, phone, address)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, s); err != nil {
		return nil, fmt.Errorf("SupplierUC.CreateSupplier - repo.Create: %w", err)
	}
	return s, nil
}

func (uc *SupplierUC) GetSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	s, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *SupplierUC) UpdateSupplier(ctx context.Context, id, name, email, phone string, address domain.SupplierAddress) (*domain.Supplier, error) {
	s, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.Name = name
	s.Email = email
	s.Phone = phone
	s.Address = address
	s.Version++
	s.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Update(ctx, s); err != nil {
		return nil, fmt.Errorf("SupplierUC.UpdateSupplier - repo.Update: %w", err)
	}
	return s, nil
}

func (uc *SupplierUC) DeactivateSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	s, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.Deactivate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, s); err != nil {
		return nil, fmt.Errorf("SupplierUC.DeactivateSupplier - repo.Update: %w", err)
	}
	return s, nil
}

func (uc *SupplierUC) ReactivateSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	s, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.Reactivate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, s); err != nil {
		return nil, fmt.Errorf("SupplierUC.ReactivateSupplier - repo.Update: %w", err)
	}
	return s, nil
}

func (uc *SupplierUC) ListSuppliers(ctx context.Context, activeOnly bool, page pagination.Params) ([]domain.Supplier, error) {
	suppliers, err := uc.repo.List(ctx, activeOnly, page)
	if err != nil {
		return nil, fmt.Errorf("SupplierUC.ListSuppliers - repo.List: %w", err)
	}
	return suppliers, nil
}
