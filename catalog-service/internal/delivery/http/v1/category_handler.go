package v1

import (
	"catalog-service/internal/delivery/http/v1/req"
	"catalog-service/internal/delivery/http/v1/res"
	"net/http"

	"github.com/TruongLe68/go-micro/pkg/httpbind"
	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) createCategory(c *gin.Context) {
	request, ok := httpbind.BindAndValidate[req.CreateCategory](c, r.v, r.l, "CreateCategory")
	if !ok {
		return
	}

	category, err := r.category.Create(c.Request.Context(), request.ToCreateCategoryInput())
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "create category success", res.ToCategoryRead(category))
}

func (r *V1) getCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := r.category.GetByID(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "get category success", res.ToCategoryRead(category))
}

func (r *V1) getCategoryChildren(c *gin.Context) {
	id := c.Param("id")
	p := pagination.FromQuery(c)
	listResult, err := r.category.GetChildren(c.Request.Context(), id, p)
	if err != nil {
		r.handleError(c, err)
		return
	}

	result := pagination.NewResult(res.ToCategoryList(listResult.Categories), p, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, "get category children success", result)
}

func (r *V1) updateCategory(c *gin.Context) {
	id := c.Param("id")
	request, ok := httpbind.BindAndValidate[req.UpdateCategory](c, r.v, r.l, "UpdateCategory")
	if !ok {
		return
	}

	category, err := r.category.Update(c.Request.Context(), request.ToUpdateCategoryInput(id))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "update category success", res.ToCategoryRead(category))
}

func (r *V1) deleteCategory(c *gin.Context) {
	id := c.Param("id")
	err := r.category.Delete(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "delete category success", nil)
}

func (r *V1) listCategories(c *gin.Context) {
	p := pagination.FromQuery(c)
	listResult, err := r.category.List(c.Request.Context(), p)
	if err != nil {
		r.handleError(c, err)
		return
	}

	result := pagination.NewResult(res.ToCategoryList(listResult.Categories), p, listResult.TotalCount)
	response.SuccessPaginated(c, http.StatusOK, "list categories success", result)
}
