package v1

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/usecase"
	"context"
)

type ProductGRPCUsecase interface {
	GetVariantsBySKUs(ctx context.Context, skus []string) ([]domain.Variant, error)
}

var _ ProductGRPCUsecase = (*usecase.ProductUC)(nil)
