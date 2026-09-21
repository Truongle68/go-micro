package usecase

import (
	"context"
	"fmt"
	"inventory-service/internal/domain"
	"strings"

	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type SupplierRepo interface {
	FindByID(ctx context.Context, id string) (*domain.Supplier, error)
}

type SupplyUC struct {
	supplierRepo  SupplierRepo
	repo          SupplyRepository
	catalogClient CatalogClient
}

func NewSupplyUC(supplierRepo SupplierRepo, repo SupplyRepository, catalogClient CatalogClient) *SupplyUC {
	return &SupplyUC{
		supplierRepo:  supplierRepo,
		repo:          repo,
		catalogClient: catalogClient,
	}
}

func (uc *SupplyUC) AssignProducts(ctx context.Context, supplierID string, inputs []domain.SupplyInput) ([]domain.SupplyAssignmentResult, error) {
	supplier, err := uc.supplierRepo.FindByID(ctx, supplierID)
	if err != nil {
		return nil, err
	}
	if !supplier.IsActive {
		return nil, domain.ErrInactiveSupp
	}

	skuSlice := make([]string, 0, len(inputs))
	for _, in := range inputs {
		skuSlice = append(skuSlice, in.SKU)
	}

	existingSKUs, err := uc.catalogClient.FindExistingSKUs(ctx, skuSlice)
	if err != nil {
		return nil, fmt.Errorf("SupplyUC.AssignProducts: %w", err)
	}

	valids := make([]*domain.Supply, 0, len(inputs))
	results := make([]domain.SupplyAssignmentResult, 0, len(inputs))
	seen := make(map[string]bool)

	for _, in := range inputs {
		if seen[in.SKU] {
			results = append(results, domain.SupplyAssignmentResult{
				SKU:     in.SKU,
				Success: false,
				Err:     domain.ErrDuplicateSKU,
			})
			continue
		}
		seen[in.SKU] = true

		if _, exists := existingSKUs[in.SKU]; !exists {
			results = append(results, domain.SupplyAssignmentResult{
				SKU:     in.SKU,
				Success: false,
				Err:     domain.ErrSKUNotFound,
			})
			continue
		}

		s, err := domain.NewSupply(supplier.ID, in.SKU, in.Price, in.LeadTime, in.MinOrderQty)
		if err != nil {
			results = append(results, domain.SupplyAssignmentResult{
				SKU:     in.SKU,
				Success: false,
				Err:     err,
			})
			continue
		}

		valids = append(valids, s)
		results = append(results, domain.SupplyAssignmentResult{
			SKU:     in.SKU,
			Success: true,
			Supply:  s,
		})
	}

	if len(valids) > 0 {
		// create batch supplies
		if err := uc.repo.CreateBatch(ctx, valids); err != nil {
			for i := range results {
				if results[i].Success {
					results[i].Success = false
					results[i].Err = fmt.Errorf("batch write failed: %w", err)
				}
			}
			return results, fmt.Errorf("create supplies batch: %w", err)
		}
	}
	return results, nil
}

func (uc *SupplyUC) GetSupply(ctx context.Context, id string) (*domain.Supply, error) {
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrEmptySupplyID
	}
	s, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("SupplyUC.GetSupply - error: %w", err)
	}
	return s, nil
}

func (uc *SupplyUC) ListSupplierSupplies(ctx context.Context, supplierID string, page pagination.Params) ([]domain.Supply, int64, error) {
	if strings.TrimSpace(supplierID) == "" {
		return nil, 0, domain.ErrEmptySuppID
	}
	list, total, err := uc.repo.FindBySupplierID(ctx, supplierID, page)
	if err != nil {
		return nil, 0, fmt.Errorf("SupplyUC.ListSupplierSupplies - error: %w", err)
	}
	return list, total, nil
}
