package v1

import (
	"catalog-service/internal/delivery/http/v1/req"
	"catalog-service/internal/domain"
	"catalog-service/internal/usecase"
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TruongLe68/go-micro/pkg/httpbind"
	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/TruongLe68/go-micro/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (r *V1) createProduct(c *gin.Context) {
	request, ok := httpbind.BindAndValidate[req.CreateProduct](c, r.v, r.l, "createProduct")
	if !ok {
		return
	}

	product, err := r.product.Create(c.Request.Context(), toCreateProductInput(request))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "create product success", toProductRead(product))
}

func (r *V1) getProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := r.product.GetByID(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "get product success", toDetailedProductRead(product))
}

func (r *V1) responseProductsByCategory(c *gin.Context, categoryID string, p pagination.Params, msg string) {
	listResult, err := r.product.GetByCategory(c.Request.Context(), categoryID, p)
	if err != nil {
		r.handleError(c, err)
		return
	}

	result := pagination.NewResult(toProductList(listResult.Products), p, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, msg, result)
}

func (r *V1) getProductsByCategory(c *gin.Context) {
	r.responseProductsByCategory(c, c.Param("id"), pagination.FromQuery(c), "get products by category success")
}

func (r *V1) updateProduct(c *gin.Context) {
	id := c.Param("id")
	request, ok := httpbind.BindAndValidate[req.UpdateProduct](c, r.v, r.l, "updateProduct")
	if !ok {
		return
	}

	product, err := r.product.Update(c.Request.Context(), toUpdateProductInput(request, id))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "update product success", toProductRead(product))
}

func (r *V1) deleteProduct(c *gin.Context) {
	id := c.Param("id")
	err := r.product.Delete(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "delete product success", nil)
}

func (r *V1) listProducts(
	c *gin.Context,
	fetchFunc func(ctx context.Context, p pagination.Params) (*usecase.ProductList, error),
	defaultSuccessMsg string) {
	p := pagination.FromQuery(c)
	categoryID := c.Query("category_id")

	if categoryID != "" {
		r.responseProductsByCategory(c, categoryID, p, "list products by category success")
		return
	}

	listResult, err := fetchFunc(c.Request.Context(), p)
	if err != nil {
		r.handleError(c, err)
		return
	}

	result := pagination.NewResult(toProductList(listResult.Products), p, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, defaultSuccessMsg, result)
}
func (r *V1) listPublicProducts(c *gin.Context) {
	minPrice, ok := parseOptionalInt(c, "min_price")
	if !ok {
		return
	}
	maxPrice, ok := parseOptionalInt(c, "max_price")
	if !ok {
		return
	}

	p := pagination.FromQuery(c)

	in := usecase.PublicListInput{
		CategoryID: c.Query("category_id"),
		Query:      c.Query("q"),
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Page:       p,
		Sort:       parseSortParams(c),
	}

	listResult, err := r.product.ListPublic(c.Request.Context(), in)
	if err != nil {
		r.l.Error("listPublicProducts: %v", err)
		response.InternalServerError(c)
		return
	}

	result := pagination.NewResult(toProductList(listResult.Products), p, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, "list products success", result)
}

// parseAdminListInput parses all shared admin product query parameters from the
// request and returns a fully populated AdminListInput. It writes an error
// response and returns ok=false if any parameter is invalid.
func parseAdminListInput(c *gin.Context) (usecase.AdminListInput, bool) {
	var statuses []domain.ProductStatus
	if raw := c.Query("status"); raw != "" {
		for _, s := range strings.Split(raw, ",") {
			if s = strings.TrimSpace(s); s != "" {
				statuses = append(statuses, domain.ProductStatus(s))
			}
		}
	}

	minPrice, ok := parseOptionalInt(c, "min_price")
	if !ok {
		return usecase.AdminListInput{}, false
	}
	maxPrice, ok := parseOptionalInt(c, "max_price")
	if !ok {
		return usecase.AdminListInput{}, false
	}

	var createdFrom *time.Time
	if v := c.Query("created_from"); v != "" {
		t, err := utils.ParseQueryTime(v)
		if err != nil {
			response.Error(c, http.StatusBadRequest, response.CodeValidationError, "invalid created_from")
			return usecase.AdminListInput{}, false
		}
		createdFrom = t
	}

	var createdTo *time.Time
	if v := c.Query("created_to"); v != "" {
		t, err := utils.ParseQueryTime(v)
		if err != nil {
			response.Error(c, http.StatusBadRequest, response.CodeValidationError, "invalid created_to")
			return usecase.AdminListInput{}, false
		}
		if len(v) == len("2026-01-01") {
			endDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC)
			createdTo = &endDay
		} else {
			createdTo = t
		}
	}

	return usecase.AdminListInput{
		Statuses:    statuses,
		CategoryID:  c.Query("category_id"),
		Query:       c.Query("q"),
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		SKU:         c.Query("sku"),
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		Sort:        parseSortParams(c),
		Page:        pagination.FromQuery(c),
	}, true
}

// parseOptionalInt parses an optional integer query parameter.
// Returns (nil, true) when the param is absent, (*int, true) on success,
// and (nil, false) after writing a 400 response on parse failure.
func parseOptionalInt(c *gin.Context, param string) (*int, bool) {
	v := c.Query(param)
	if v == "" {
		return nil, true
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "invalid "+param)
		return nil, false
	}
	return &n, true
}

func (r *V1) listAdminProducts(c *gin.Context) {
	in, ok := parseAdminListInput(c)
	if !ok {
		return
	}

	listResult, err := r.product.ListAdmin(c.Request.Context(), in)
	if err != nil {
		r.l.Error("listAdminProducts: %v", err)
		response.InternalServerError(c)
		return
	}

	result := pagination.NewResult(toProductList(listResult.Products), in.Page, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, "list products success", result)
}

func (r *V1) searchProducts(c *gin.Context) {
	var query domain.SearchProductsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.InvalidQueryParams(c)
		return
	}
	p := pagination.FromQuery(c)

	listResult, err := r.product.Search(c.Request.Context(), query.ToDomainParams(), p)
	if err != nil {
		r.handleError(c, err)
		return
	}

	result := pagination.NewResult(toProductList(listResult.Products), p, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, "search products success", result)
}

func (r *V1) countProductsPerStatus(c *gin.Context) {
	in, ok := parseAdminListInput(c)
	if !ok {
		return
	}

	counts, err := r.product.CountPerStatus(c.Request.Context(), in)
	if err != nil {
		r.l.Error("countProductsPerStatus: %v", err)
		response.InternalServerError(c)
		return
	}

	response.Success(c, http.StatusOK, "get product status counts success", counts)
}

func parseSortParams(c *gin.Context) []domain.SortKey {
	var raw []string
	raw = append(raw, c.QueryArray("sort")...)
	if len(raw) == 1 && strings.Contains(raw[0], ",") { // ?sort=name_asc,price_asc
		raw = strings.Split(raw[0], ",")
	}

	out := make([]domain.SortKey, 0, len(raw))
	for _, s := range raw {
		s := strings.TrimSpace(strings.ToLower(s))
		if s == "" {
			continue
		}

		parts := strings.Split(s, "_")
		if len(parts) != 2 {
			continue
		}

		field := domain.SortField(parts[0])
		dir := domain.SortDir(parts[1])

		if !field.IsValid() || !dir.IsValid() {
			continue
		}
		out = append(out, domain.SortKey{
			Field: field,
			Dir:   dir,
		})
	}
	return out
}
