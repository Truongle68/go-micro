package v1

import (
	"catalog-service/internal/delivery/http/v1/req"
	"catalog-service/internal/delivery/http/v1/res"
	"catalog-service/internal/domain"
	"catalog-service/internal/usecase"
)

func toCreateProductInput(req req.CreateProduct) usecase.CreateProductInput {
	variants := make([]usecase.CreateVariantInput, len(req.Variants))
	for i, vr := range req.Variants {
		variants[i] = usecase.CreateVariantInput{
			ID:         vr.ID,
			SKU:        vr.SKU,
			Attributes: vr.Attributes,
			Price: usecase.PriceInput{
				Amount:   vr.Price.Amount,
				Currency: vr.Price.Currency,
			},
			Stock:       vr.Stock,
			WeightGrams: vr.WeightGrams,
			Images:      toImageInputs(vr.Images),
		}
	}

	optionTypes := make([]usecase.OptionTypeInput, len(req.OptionTypes))
	for i, ot := range req.OptionTypes {
		optionTypes[i] = usecase.OptionTypeInput{Name: ot.Name, Values: ot.Values}
	}

	specs := make([]usecase.SpecGroupInput, len(req.Specifications))
	for i, sg := range req.Specifications {
		items := make([]usecase.SpecItemInput, len(sg.Items))
		for j, item := range sg.Items {
			items[j] = usecase.SpecItemInput{Label: item.Label, Value: item.Value}
		}
		specs[i] = usecase.SpecGroupInput{Group: sg.Group, Items: items}
	}

	return usecase.CreateProductInput{
		Name:            req.Name,
		CategoryID:      req.CategoryID,
		Description:     req.Description,
		DescriptionHTML: req.DescriptionHTML,
		Highlights:      req.Highlights,
		Tags:            req.Tags,
		Images:          toImageInputs(req.Images),
		OptionTypes:     optionTypes,
		Variants:        variants,
		Specifications:  specs,
		Shipping: usecase.ShippingInput{
			IsFreeShipping: req.Shipping.IsFreeShipping,
			Fragile:        req.Shipping.Fragile,
			ShippingClass:  req.Shipping.ShippingClass,
		},
	}
}

func toImageInputs(in []req.ImageInput) []usecase.ImageInput {
	out := make([]usecase.ImageInput, len(in))
	for i, img := range in {
		out[i] = usecase.ImageInput{URL: img.URL, IsPrimary: img.IsPrimary, SortOrder: img.SortOrder, AltText: img.AltText}
	}
	return out
}

func toUpdateProductInput(req req.UpdateProduct, id string) usecase.UpdateProductInput {
	var status *domain.ProductStatus
	if req.Status != nil {
		s := domain.ProductStatus(*req.Status)
		status = &s
	}

	var images []usecase.ImageInput
	if req.Images != nil {
		images = make([]usecase.ImageInput, len(req.Images))
		for i, img := range req.Images {
			images[i] = usecase.ImageInput{
				URL:       img.URL,
				IsPrimary: img.IsPrimary,
				SortOrder: img.SortOrder,
				AltText:   img.AltText,
			}
		}
	}

	var optionTypes []usecase.OptionTypeInput
	if req.OptionTypes != nil {
		optionTypes = make([]usecase.OptionTypeInput, len(req.OptionTypes))
		for i, ot := range req.OptionTypes {
			optionTypes[i] = usecase.OptionTypeInput{
				Name:   ot.Name,
				Values: ot.Values,
			}
		}
	}

	var variants []usecase.CreateVariantInput
	if req.Variants != nil {
		variants = make([]usecase.CreateVariantInput, len(req.Variants))
		for i, v := range req.Variants {
			var vImages []usecase.ImageInput
			if v.Images != nil {
				vImages = make([]usecase.ImageInput, len(v.Images))
				for j, img := range v.Images {
					vImages[j] = usecase.ImageInput{
						URL:       img.URL,
						IsPrimary: img.IsPrimary,
						SortOrder: img.SortOrder,
						AltText:   img.AltText,
					}
				}
			}
			variants[i] = usecase.CreateVariantInput{
				ID:         v.ID,
				SKU:        v.SKU,
				Attributes: v.Attributes,
				Price: usecase.PriceInput{
					Amount:   v.Price.Amount,
					Currency: v.Price.Currency,
				},
				Stock:       v.Stock,
				WeightGrams: v.WeightGrams,
				Images:      vImages,
			}
		}
	}

	var specs []usecase.SpecGroupInput
	if req.Specifications != nil {
		specs = make([]usecase.SpecGroupInput, len(req.Specifications))
		for i, sg := range req.Specifications {
			items := make([]usecase.SpecItemInput, len(sg.Items))
			for j, item := range sg.Items {
				items[j] = usecase.SpecItemInput{
					Label: item.Label,
					Value: item.Value,
				}
			}
			specs[i] = usecase.SpecGroupInput{
				Group: sg.Group,
				Items: items,
			}
		}
	}

	var shipping *usecase.ShippingInput
	if req.Shipping != nil {
		shipping = &usecase.ShippingInput{
			IsFreeShipping: req.Shipping.IsFreeShipping,
			Fragile:        req.Shipping.Fragile,
			ShippingClass:  req.Shipping.ShippingClass,
		}
	}

	return usecase.UpdateProductInput{
		ID:              id,
		Version:         req.Version,
		Name:            req.Name,
		NameTranslation: req.NameTranslation,
		CategoryID:      req.CategoryID,
		Description:     req.Description,
		DescriptionHTML: req.DescriptionHTML,
		Highlights:      req.Highlights,
		Tags:            req.Tags,
		Images:          images,
		OptionTypes:     optionTypes,
		Variants:        variants,
		Specifications:  specs,
		Shipping:        shipping,
		Status:          status,
	}
}

func toProductResponse(p *domain.Product) res.ProductResponse {
	variants := make([]res.VariantResponse, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = res.VariantResponse{
			ID:    v.ID,
			SKU:   v.SKU,
			Price: v.Price.Amount,
			Stock: v.Inventory.TotalAvailable,
		}
	}
	return res.ProductResponse{
		ID:        p.ID,
		Slug:      p.Slug,
		Name:      p.Name,
		Status:    string(p.Status),
		Variants:  variants,
		CreatedAt: p.CreatedAt,
	}
}

func toProductCategoryRead(c domain.Category) *res.ProductCategoryRead {
	return &res.ProductCategoryRead{
		ID:              c.ID,
		Name:            c.Name,
		NameTranslation: c.NameTranslation,
		Slug:            c.Slug,
		Icon:            c.Icon,
		SortOrder:       c.SortOrder,
		IsActive:        c.IsActive,
	}
}

func toProductRead(p *domain.Product) res.ProductRead {
	if p == nil {
		return res.ProductRead{}
	}

	categoryPath := make([]res.CategoryRefRead, len(p.CategoryPath))
	for i, cp := range p.CategoryPath {
		categoryPath[i] = res.CategoryRefRead{
			ID:   cp.ID,
			Name: cp.Name,
		}
	}

	images := make([]res.ImageRead, len(p.Images))
	for i, img := range p.Images {
		images[i] = res.ImageRead{
			URL:       img.URL,
			IsPrimary: img.IsPrimary,
			SortOrder: img.SortOrder,
			AltText:   img.AltText,
		}
	}

	optionTypes := make([]res.OptionTypeRead, len(p.OptionTypes))
	for i, ot := range p.OptionTypes {
		optionTypes[i] = res.OptionTypeRead{
			Name:   ot.Name,
			Values: ot.Values,
		}
	}

	variants := make([]res.VariantRead, len(p.Variants))
	for i, v := range p.Variants {
		vImages := make([]res.ImageRead, len(v.Images))
		for j, img := range v.Images {
			vImages[j] = res.ImageRead{
				URL:       img.URL,
				IsPrimary: img.IsPrimary,
				SortOrder: img.SortOrder,
				AltText:   img.AltText,
			}
		}

		variants[i] = res.VariantRead{
			ID:         v.ID,
			SKU:        v.SKU,
			Attributes: v.Attributes,
			Price: res.PriceRead{
				Amount:   v.Price.Amount,
				Currency: v.Price.Currency,
			},
			Inventory: res.InventoryRead{
				TotalAvailable: v.Inventory.TotalAvailable,
				Reserved:       v.Inventory.Reserved,
			},
			WeightGrams: v.WeightGrams,
			Images:      vImages,
			IsActive:    v.IsActive,
			CreatedAt:   v.CreatedAt,
		}
	}

	specs := make([]res.SpecGroupRead, len(p.Specifications))
	for i, sg := range p.Specifications {
		items := make([]res.SpecItemRead, len(sg.Items))
		for j, item := range sg.Items {
			items[j] = res.SpecItemRead{
				Label: item.Label,
				Value: item.Value,
			}
		}
		specs[i] = res.SpecGroupRead{
			Group: sg.Group,
			Items: items,
		}
	}

	return res.ProductRead{
		ID:              p.ID,
		Version:         p.Version,
		Slug:            p.Slug,
		Name:            p.Name,
		NameTranslation: p.NameTranslation,
		CategoryID:      p.CategoryID,
		CategoryPath:    categoryPath,
		Description:     p.Description,
		DescriptionHTML: p.DescriptionHTML,
		Highlights:      p.Highlights,
		Tags:            p.Tags,
		Images:          images,
		OptionTypes:     optionTypes,
		Variants:        variants,
		Specifications:  specs,
		RatingSummary: res.RatingSummaryRead{
			Average:   p.RatingSummary.Average,
			Count:     p.RatingSummary.Count,
			Breakdown: p.RatingSummary.Breakdown,
		},
		SalesCount: p.SalesCount,
		Shipping: res.ShippingInfoRead{
			IsFreeShipping: p.Shipping.IsFreeShipping,
			ShipsFrom: res.ShipsFromRead{
				WarehouseID: p.Shipping.ShipsFrom.WarehouseID,
				Region:      p.Shipping.ShipsFrom.Region,
			},
			Fragile:       p.Shipping.Fragile,
			ShippingClass: p.Shipping.ShippingClass,
		},
		Status:    string(p.Status),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toDetailedProductRead(dp *usecase.DetailedProduct) res.ProductRead {
	if dp == nil {
		return res.ProductRead{}
	}

	p := toProductRead(&dp.Product)

	if dp.Category != nil {
		p.Category = toProductCategoryRead(*dp.Category)
	}

	return p
}

func toProductList(dps []usecase.DetailedProduct) []res.ProductRead {
	products := make([]res.ProductRead, len(dps))
	for i, dp := range dps {
		products[i] = toDetailedProductRead(&dp)
	}
	return products
}
