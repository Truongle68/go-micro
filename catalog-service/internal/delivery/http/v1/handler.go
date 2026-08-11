package v1

import (
	"catalog-service/internal/domain"
	"errors"
	"net/http"

	"github.com/TruongLe68/go-micro/pkg/ginmw"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *V1) handleError(c *gin.Context, err error) {
	appErr := domain.ToAppError(err)
	if appErr != nil {
		var status int
		switch {
		case errors.Is(err, domain.ErrProductNotFound),
			errors.Is(err, domain.ErrCategoryNotFound):
			status = http.StatusNotFound
		case errors.Is(err, domain.ErrDuplicateSKU),
			errors.Is(err, domain.ErrDuplicateSlug),
			errors.Is(err, domain.ErrDuplicateField),
			errors.Is(err, domain.ErrConcurrentUpdate):
			status = http.StatusConflict
		default:
			status = http.StatusBadRequest
		}
		response.Error(c, status, string(appErr.Code), appErr.Message)
		return
	}

	r.l.Error(err, "catalog handler unexpected error")
	response.InternalServerError(c)
}

func (r *V1) getUserID(c *gin.Context) (string, bool) {
	return ginmw.GetUserID(c)
}
