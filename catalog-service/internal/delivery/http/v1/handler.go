package v1

import (
	"catalog-service/internal/delivery/http/middleware"
	"catalog-service/internal/domain"
	"errors"
	"net/http"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound),
		errors.Is(err, domain.ErrCategoryNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrEmptyProductID),
		errors.Is(err, domain.ErrEmptyCategoryID),
		errors.Is(err, domain.ErrEmptyName),
		errors.Is(err, domain.ErrEmptySku),
		errors.Is(err, domain.ErrInvalidPrice),
		errors.Is(err, domain.ErrInvalidPriceRange),
		errors.Is(err, domain.ErrInvalidCategoryID),
		errors.Is(err, domain.ErrInvalidProductID),
		errors.Is(err, domain.ErrNoFieldsToUpdate):
		response.Error(c, http.StatusBadRequest, err.Error())
	default:
		r.l.Error(err, "catalog handler unexpected error")
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}

func (r *V1) getUserID(c *gin.Context) (string, bool) {
	return middleware.GetUserID(c)
}
