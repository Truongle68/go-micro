package v1

import (
	"catalog-service/internal/delivery/http/v1/req"
	"catalog-service/internal/delivery/http/v1/res"
	"catalog-service/internal/domain"
	"net/http"
	"strconv"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) CreateProduct(c *gin.Context) {
	var request req.CreateProduct
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - CreateProduct - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - CreateProduct - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := r.product.Create(c.Request.Context(), request.ToCreateProductInput())
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "create product success", res.ToProductResponse(product))
}

func (r *V1) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := r.product.GetByID(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "get product success", res.ToProductResponse(product))
}

func (r *V1) GetProductsByCategory(c *gin.Context) {
	id := c.Param("id")
	p := pagination.FromQuery(c)
	listResult, err := r.product.GetByCategory(c.Request.Context(), id, p)
	if err != nil {
		r.handleError(c, err)
		return
	}

	result := pagination.NewResult(res.ToProductListResponse(listResult.Products), p, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, "get products by category success", result)
}

func (r *V1) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var request req.UpdateProduct
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - UpdateProduct - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - UpdateProduct - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := r.product.Update(c.Request.Context(), request.ToUpdateProductInput(id))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "update product success", res.ToProductResponse(product))
}

func (r *V1) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	err := r.product.Delete(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "delete product success", nil)
}

func (r *V1) ListProducts(c *gin.Context) {
	p := pagination.FromQuery(c)
	categoryID := c.Query("category_id")

	if categoryID != "" {
		listResult, err := r.product.GetByCategory(c.Request.Context(), categoryID, p)
		if err != nil {
			r.handleError(c, err)
			return
		}
		result := pagination.NewResult(res.ToProductListResponse(listResult.Products), p, listResult.TotalCount)
		response.SuccessPaginated(c, http.StatusOK, "list products by category success", result)
		return
	}

	listResult, err := r.product.List(c.Request.Context(), p)
	if err != nil {
		r.handleError(c, err)
		return
	}

	result := pagination.NewResult(res.ToProductListResponse(listResult.Products), p, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, "list products success", result)
}

func (r *V1) SearchProducts(c *gin.Context) {
	var params domain.SearchProductParams
	params.Query = c.Query("q")
	if categoryID := c.Query("category_id"); categoryID != "" {
		params.CategoryID = categoryID
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			params.MinPrice = &val
		}
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			params.MaxPrice = &val
		}
	}
	if activeStr := c.Query("is_active"); activeStr != "" {
		if val, err := strconv.ParseBool(activeStr); err == nil {
			params.IsActive = &val
		}
	}
	p := pagination.FromQuery(c)

	listResult, err := r.product.Search(c.Request.Context(), params, p)
	if err != nil {
		r.handleError(c, err)
		return
	}

	result := pagination.NewResult(res.ToProductListResponse(listResult.Products), p, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, "search products success", result)
}
