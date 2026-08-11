package usecase_test

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/repo"
	"catalog-service/internal/usecase"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCategoryRepoForUC struct {
	repo.CategoryRepository
	categories map[string]*domain.Category
	createFn   func(ctx context.Context, c *domain.Category) error
	updateFn   func(ctx context.Context, c *domain.Category) (*domain.Category, error)
}

func (m *mockCategoryRepoForUC) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	c, ok := m.categories[id]
	if !ok {
		return nil, domain.ErrCategoryNotFound
	}
	return c, nil
}

func (m *mockCategoryRepoForUC) Create(ctx context.Context, c *domain.Category) error {
	if m.createFn != nil {
		return m.createFn(ctx, c)
	}
	c.ID = "cat-child"
	m.categories[c.ID] = c
	return nil
}

func (m *mockCategoryRepoForUC) Update(ctx context.Context, c *domain.Category) (*domain.Category, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, c)
	}
	m.categories[c.ID] = c
	return c, nil
}

func TestCategoryUC_CreateWithAncestors(t *testing.T) {
	parent := &domain.Category{
		ID:   "cat-parent",
		Name: "Electronics",
		Ancestors: []domain.CategoryRef{},
	}

	r := &mockCategoryRepoForUC{
		categories: map[string]*domain.Category{
			"cat-parent": parent,
		},
	}

	uc := usecase.NewCategoryUC(r)
	parentID := "cat-parent"
	dto, err := uc.Create(context.Background(), usecase.CreateCategoryInput{
		ParentID: &parentID,
		Name:     "Laptops",
		Slug:     "laptops",
	})

	assert.NoError(t, err)
	assert.NotNil(t, dto)
	assert.Len(t, dto.Ancestors, 1)
	assert.Equal(t, "cat-parent", dto.Ancestors[0].ID)
	assert.Equal(t, "Electronics", dto.Ancestors[0].Name)
}

func TestCategoryUC_UpdateCircularParentRejection(t *testing.T) {
	child := &domain.Category{
		ID:        "cat-child",
		Name:      "Laptops",
		Ancestors: []domain.CategoryRef{{ID: "cat-parent", Name: "Electronics"}},
	}
	parent := &domain.Category{
		ID:        "cat-parent",
		Name:      "Electronics",
		Ancestors: []domain.CategoryRef{},
	}

	r := &mockCategoryRepoForUC{
		categories: map[string]*domain.Category{
			"cat-parent": parent,
			"cat-child":  child,
		},
	}

	uc := usecase.NewCategoryUC(r)

	// Attempt to set parent of cat-parent to cat-child (circular loop)
	childID := "cat-child"
	_, err := uc.Update(context.Background(), usecase.UpdateCategoryInput{
		ID:       "cat-parent",
		ParentID: &childID,
	})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrCircularCategoryParent))
}

func TestCategoryUC_ValidationErrors(t *testing.T) {
	r := &mockCategoryRepoForUC{categories: map[string]*domain.Category{}}
	uc := usecase.NewCategoryUC(r)

	_, err := uc.Create(context.Background(), usecase.CreateCategoryInput{
		Name: "",
		Slug: "slug",
	})
	assert.ErrorIs(t, err, domain.ErrEmptyName)

	_, err = uc.Create(context.Background(), usecase.CreateCategoryInput{
		Name: "Name",
		Slug: "",
	})
	assert.ErrorIs(t, err, domain.ErrEmptySlug)
}
