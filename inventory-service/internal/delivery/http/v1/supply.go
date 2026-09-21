package v1

import (
	"errors"
	"inventory-service/internal/delivery/http/v1/req"
	"inventory-service/internal/delivery/http/v1/res"
	"inventory-service/internal/domain"
	"net/http"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *V1) assignProducts(c *gin.Context) {
	var r req.AssignProductsReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	suppInputs := make([]domain.SupplyInput, 0, len(r.Supplies))
	for _, input := range r.Supplies {
		suppInputs = append(suppInputs, input.ToDomain())
	}

	results, err := h.supply.AssignProducts(c.Request.Context(), r.SupplierID, suppInputs)
	if err != nil {
		if errors.Is(err, domain.ErrSuppNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeSuppNotFound), "supplier with this id not found")
			return
		}
		if appErr := domain.ToAppError(err); appErr != nil {
			response.Error(c, http.StatusBadRequest, string(appErr.Code), appErr.Message)
			return
		}
		h.l.Error("h.assignProducts: %v", err)
		response.InternalServerError(c, "failed to assign products")
		return
	}

	response.Success(c, http.StatusCreated, "products assigned successfully", res.ToProductAssignedRes(results))
}

func (h *V1) getSupply(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, domain.ErrEmptySupplyID.Error())
		return
	}

	supply, err := h.supply.GetSupply(c, id)
	if err != nil {
		if appErr := domain.ToAppError(err); appErr != nil {
			response.Error(c, http.StatusBadRequest, string(appErr.Code), appErr.Message)
			return
		}
		h.l.Error("h.getSupply: %v", err)
		response.InternalServerError(c, "failed to get supply")
		return
	}

	response.Success(c, http.StatusOK, "get supply successfully", supply)
}

func (h *V1) listSupplierSupplies(c *gin.Context) {
	supplierID := c.Param("supplierId")
	if supplierID == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, domain.ErrEmptySuppID.Error())
		return
	}

	p := pagination.FromQuery(c)

	supply, total, err := h.supply.ListSupplierSupplies(c, supplierID, p)
	if err != nil {
		if appErr := domain.ToAppError(err); appErr != nil {
			response.Error(c, http.StatusBadRequest, string(appErr.Code), appErr.Message)
			return
		}
		h.l.Error("h.listSupplierSupplies: %v", err)
		response.InternalServerError(c, "failed to list supplier supplies")
		return
	}

	response.SuccessPaginated(c, http.StatusOK, "get supply successfully", pagination.NewResult(supply, p, total))
}
