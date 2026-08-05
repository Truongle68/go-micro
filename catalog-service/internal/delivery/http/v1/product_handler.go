package v1

import (
	"catalog-service/internal/delivery/http/v1/req"
	"catalog-service/internal/domain"
	"net/http"

	"github.com/TruongLe68/go-micro/pkg/httpbind"
	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/response"
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

	response.Success(c, http.StatusCreated, "create product success", toProductResponse(product))
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

func (r *V1) listProducts(c *gin.Context) {
	p := pagination.FromQuery(c)
	categoryID := c.Query("category_id")

	if categoryID != "" {
		r.responseProductsByCategory(c, categoryID, p, "list products by category success")
		return
	}

	listResult, err := r.product.List(c.Request.Context(), p)
	if err != nil {
		r.handleError(c, err)
		return
	}

	result := pagination.NewResult(toProductList(listResult.Products), p, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, "list products success", result)
}

func (r *V1) searchProducts(c *gin.Context) {
	var query domain.SearchProductsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid query parameters")
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
