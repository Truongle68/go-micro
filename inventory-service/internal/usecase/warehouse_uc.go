package usecase

import (
	"context"
	"inventory-service/internal/domain"
)

type WarehouseUC struct {
	repo WarehouseRepository
}

func NewWarehouseUC(r WarehouseRepository) *WarehouseUC {
	return &WarehouseUC{
		repo: r,
	}
}

func (u *WarehouseUC) List(ctx context.Context) ([]domain.Warehouse, error) {
	return u.repo.Find(ctx)
}
