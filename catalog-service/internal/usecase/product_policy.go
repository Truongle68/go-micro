package usecase

import "catalog-service/internal/domain"

func applyPublicListPolicy(in PublicListInput) domain.ProductListFilter {
	return domain.ProductListFilter{
		Statuses:   []domain.ProductStatus{domain.ProductStatusActive},
		CategoryID: in.CategoryID,
		Query:      in.Query,
		MinPrice:   in.MinPrice,
		MaxPrice:   in.MaxPrice,
		Sort:       in.Sort,
	}
}

func applyAdminListPolicy(in AdminListInput) domain.ProductListFilter {
	return domain.ProductListFilter{
		Statuses:    in.Statuses,
		CategoryID:  in.CategoryID,
		Query:       in.Query,
		MinPrice:    in.MinPrice,
		MaxPrice:    in.MaxPrice,
		SKU:         in.SKU,
		CreatedFrom: in.CreatedFrom,
		CreatedTo:   in.CreatedTo,
		Sort:        in.Sort,
	}
}
