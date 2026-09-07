package v1

import (
	"net/http"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *V1) listWarehouses(c *gin.Context) {
	warehouse, err := h.warehouse.List(c.Request.Context())
	if err != nil {
		h.l.Error("h.listWarehouses: %v", err)
		response.InternalServerError(c, "failed to list warehouses")
		return
	}

	response.Success(c, http.StatusOK, "warehouses retrieved successfully", warehouse)
}
