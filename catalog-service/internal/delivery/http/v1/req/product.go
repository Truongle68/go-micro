package req

type CreateProduct struct {
	Name            string               `json:"name" validate:"required,min=3,max=200"`
	CategoryID      string               `json:"category_id" validate:"required"`
	Description     string               `json:"description" validate:"required"`
	DescriptionHTML string               `json:"description_html,omitempty"`
	Highlights      []string             `json:"highlights,omitempty"`
	Tags            []string             `json:"tags,omitempty"`
	Images          []ImageInput         `json:"images" validate:"required,min=1,dive"`
	OptionTypes     []OptionTypeInput    `json:"option_types,omitempty"`
	Variants        []CreateVariantInput `json:"variants" validate:"required,min=1,dive"`
	Specifications  []SpecGroupInput     `json:"specifications,omitempty"`
	Shipping        ShippingInput        `json:"shipping" validate:"required"`
}

type ImageInput struct {
	URL       string `json:"url" validate:"required,url"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
	AltText   string `json:"alt_text,omitempty"`
}

type OptionTypeInput struct {
	Name   string   `json:"name" validate:"required"`
	Values []string `json:"values" validate:"required,min=1"`
}

type CreateVariantInput struct {
	ID          string            `json:"id,omitempty"`
	SKU         string            `json:"sku" validate:"required"`
	Attributes  map[string]string `json:"attributes"`
	Price       PriceInput        `json:"price" validate:"required"`
	Stock       int               `json:"stock" validate:"gte=0"`
	WeightGrams int               `json:"weight_grams,omitempty"`
	Images      []ImageInput      `json:"images,omitempty"`
}

type PriceInput struct {
	Amount   int64  `json:"amount" validate:"required,min=1"`
	Currency string `json:"currency" validate:"required,len=3"`
}

type SpecGroupInput struct {
	Group string          `json:"group" validate:"required"`
	Items []SpecItemInput `json:"items" validate:"required,min=1"`
}

type SpecItemInput struct {
	Label string `json:"label" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type ShippingInput struct {
	IsFreeShipping bool   `json:"is_free_shipping"`
	Fragile        bool   `json:"fragile"`
	ShippingClass  string `json:"shipping_class,omitempty"`
}

type UpdateProduct struct {
	Version         int                  `json:"version" validate:"required,min=1"`
	Name            *string              `json:"name,omitempty" validate:"omitempty,min=3,max=200"`
	NameTranslation map[string]string    `json:"name_translation,omitempty"`
	CategoryID      *string              `json:"category_id,omitempty"`
	Description     *string              `json:"description,omitempty"`
	DescriptionHTML *string              `json:"description_html,omitempty"`
	Highlights      []string             `json:"highlights,omitempty"`
	Tags            []string             `json:"tags,omitempty"`
	Images          []ImageInput         `json:"images,omitempty" validate:"omitempty,dive"`
	OptionTypes     []OptionTypeInput    `json:"option_types,omitempty"`
	Variants        []CreateVariantInput `json:"variants,omitempty" validate:"omitempty,dive"`
	Specifications  []SpecGroupInput     `json:"specifications,omitempty"`
	Shipping        *ShippingInput       `json:"shipping,omitempty"`
	Status          *string              `json:"status,omitempty" validate:"omitempty,oneof=draft active inactive archived"`
}

type SearchProduct struct {
	Query      string  `json:"query"`
	CategoryID string  `json:"category_id"`
	MinPrice   *int64  `json:"min_price"`
	MaxPrice   *int64  `json:"max_price"`
	Status     *string `json:"status"`
	Page       int64   `json:"page"`
	Limit      int64   `json:"limit"`
}
