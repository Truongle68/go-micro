package domain_test

import (
	"catalog-service/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplyUpdate(t *testing.T) {
	initialTime := time.Now().Add(-1 * time.Hour)
	p := &domain.Product{
		ID:          "p-1",
		Slug:        "test-product",
		Name:        "Old Name",
		CategoryID:  "cat-1",
		Description: "Old Description",
		Status:      domain.ProductStatusDraft,
		CreatedAt:   initialTime,
		UpdatedAt:   initialTime,
	}

	newName := "New Product Name"
	newCatID := "cat-2"
	newDesc := "New Description"
	newStatus := domain.ProductStatusActive

	params := domain.UpdateProductParams{
		Name:        &newName,
		CategoryID:  &newCatID,
		Description: &newDesc,
		Status:      &newStatus,
		Variants: []domain.Variant{
			{
				SKU: "SKU-001",
				Price: domain.Price{
					Amount:   1000,
					Currency: "USD",
				},
			},
		},
	}

	err := p.ApplyUpdate(params)
	assert.NoError(t, err)

	assert.Equal(t, "New Product Name", p.Name)
	assert.Equal(t, "cat-2", p.CategoryID)
	assert.Equal(t, "New Description", p.Description)
	assert.Equal(t, domain.ProductStatusActive, p.Status)
	assert.Len(t, p.Variants, 1)
	assert.Equal(t, "SKU-001", p.Variants[0].SKU)
	assert.True(t, p.UpdatedAt.After(initialTime))
}

func TestApplyUpdate_ValidationErrors(t *testing.T) {
	p := &domain.Product{
		ID:         "p-1",
		Name:       "Valid Name",
		CategoryID: "cat-1",
		Status:     domain.ProductStatusDraft,
	}

	emptyName := ""
	err := p.ApplyUpdate(domain.UpdateProductParams{Name: &emptyName})
	assert.ErrorIs(t, err, domain.ErrEmptyName)

	emptyCatID := ""
	err = p.ApplyUpdate(domain.UpdateProductParams{CategoryID: &emptyCatID})
	assert.ErrorIs(t, err, domain.ErrEmptyCategoryID)

	invalidStatus := domain.ProductStatus("invalid_status")
	err = p.ApplyUpdate(domain.UpdateProductParams{Status: &invalidStatus})
	assert.ErrorIs(t, err, domain.ErrInvalidProductStatus)
}
