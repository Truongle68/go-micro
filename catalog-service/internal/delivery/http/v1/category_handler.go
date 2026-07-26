package v1

import (
	"catalog-service/internal/delivery/http/v1/req"
	"catalog-service/internal/delivery/http/v1/res"
	"net/http"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) CreateCategory(c *gin.Context) {
	var request req.CreateCategory
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - CreateCategory - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - CreateCategory - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := r.category.Create(c.Request.Context(), request.ToCreateCategoryInput())
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "create category success", res.ToCategoryResponse(category))
}

func (r *V1) GetCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := r.category.GetByID(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "get category success", res.ToCategoryResponse(category))
}

func (r *V1) GetCategoryChildren(c *gin.Context) {
	id := c.Param("id")
	children, err := r.category.GetChildren(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "get category children success", res.ToCategoryListResponse(children))
}

func (r *V1) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var request req.UpdateCategory
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Warn("restapi - v1 - UpdateCategory - ShouldBindJSON: %v", err)
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := r.v.Struct(request); err != nil {
		r.l.Warn("restapi - v1 - UpdateCategory - validate: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := r.category.Update(c.Request.Context(), request.ToUpdateCategoryInput(id))
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "update category success", res.ToCategoryResponse(category))
}

func (r *V1) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	err := r.category.Delete(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "delete category success", nil)
}

func (r *V1) ListCategories(c *gin.Context) {
	p := pagination.FromQuery(c)
	categories, err := r.category.List(c.Request.Context(), p)
	if err != nil {
		r.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "list categories success", res.ToCategoryListResponse(categories))
}
