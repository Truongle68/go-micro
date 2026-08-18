package req

type CreateProduct struct {
	Name            string               `json:"name" validate:"required,min=3,max=200"`
	CategoryID      string               `json:"category_id" validate:"required"`
	Description     string               `json:"description" validate:"required"`
	DescriptionHTML string               `json:"description_html,omitempty"`
	Highlights      []string             `json:"highlights,omitempty"`
	Tags            []string             `json:"tags,omitempty"`
	OptionTypes     []OptionTypeInput    `json:"option_types,omitempty"`
	Variants        []CreateVariantInput `json:"variants" validate:"required,min=1,dive"`
	Specifications  []SpecGroupInput     `json:"specifications,omitempty"`
	Status          string               `json:"status" validate:"required,oneof=draft active"`
}

type OptionTypeInput struct {
	Name   string   `json:"name" validate:"required"`
	Values []string `json:"values" validate:"required,min=1"`
}

type CreateVariantInput struct {
	ID         string            `json:"id,omitempty"`
	SKU        string            `json:"sku" validate:"required"`
	Attributes map[string]string `json:"attributes"`
	Price      PriceInput        `json:"price" validate:"required"`
	Image      string            `json:"image,omitempty"`
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

type UpdateProduct struct {
	Version         int                  `json:"version" validate:"required,min=1"`
	Name            *string              `json:"name,omitempty" validate:"omitempty,min=3,max=200"`
	NameTranslation map[string]string    `json:"name_translation,omitempty"`
	CategoryID      *string              `json:"category_id,omitempty"`
	Description     *string              `json:"description,omitempty"`
	DescriptionHTML *string              `json:"description_html,omitempty"`
	Highlights      []string             `json:"highlights,omitempty"`
	Tags            []string             `json:"tags,omitempty"`
	OptionTypes     []OptionTypeInput    `json:"option_types,omitempty"`
	Variants        []CreateVariantInput `json:"variants,omitempty" validate:"omitempty,dive"`
	Specifications  []SpecGroupInput     `json:"specifications,omitempty"`
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
