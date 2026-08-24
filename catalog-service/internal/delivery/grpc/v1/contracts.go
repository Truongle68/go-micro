package v1

import (
	"catalog-service/internal/usecase"
	"context"
)

type ProductGRPCUsecase interface {
	GetVariantsBySKUs(ctx context.Context, skus []string) ([]usecase.VariantView, error)
}

var _ ProductGRPCUsecase = (*usecase.ProductUC)(nil)
