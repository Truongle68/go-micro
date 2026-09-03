package v1

import (
	"errors"
	"net/http"
	"strconv"

	"inventory-service/internal/delivery/http/v1/req"
	"inventory-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *V1) createSupplier(c *gin.Context) {
	var r req.CreateSupplierReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	address := domain.SupplierAddress{
		Line1:    r.Address.Line1,
		Ward:     r.Address.Ward,
		District: r.Address.District,
		City:     r.Address.City,
	}

	supplier, err := h.supplier.CreateSupplier(c.Request.Context(), r.Code, r.Name, r.Phone, address)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateSuppCode) {
			response.Error(c, http.StatusConflict, string(domain.CodeDuplicateSuppCode), "supplier with this code already exists")
			return
		}
		if appErr := domain.ToAppError(err); appErr != nil {
			response.Error(c, http.StatusBadRequest, string(appErr.Code), appErr.Message)
			return
		}
		h.l.Error("h.createSupplier: %v", err)
		response.InternalServerError(c, "failed to create supplier")
		return
	}

	response.Success(c, http.StatusCreated, "supplier created successfully", supplier)
}

func (h *V1) getSupplier(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "supplier id is required")
		return
	}

	supplier, err := h.supplier.GetSupplier(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSuppNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeSuppNotFound), err.Error())
			return
		}
		h.l.Error("h.getSupplier: %v", err)
		response.InternalServerError(c, "failed to get supplier")
		return
	}

	response.Success(c, http.StatusOK, "supplier retrieved successfully", supplier)
}

func (h *V1) updateSupplier(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "supplier id is required")
		return
	}

	var r req.UpdateSupplierReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	address := domain.SupplierAddress{
		Line1:    r.Address.Line1,
		Ward:     r.Address.Ward,
		District: r.Address.District,
		City:     r.Address.City,
	}

	supplier, err := h.supplier.UpdateSupplier(c.Request.Context(), id, r.Name, r.Email, r.Phone, address)
	if err != nil {
		if errors.Is(err, domain.ErrSuppNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeSuppNotFound), err.Error())
			return
		}
		if appErr := domain.ToAppError(err); appErr != nil {
			response.Error(c, http.StatusBadRequest, string(appErr.Code), appErr.Message)
			return
		}
		h.l.Error("h.updateSupplier: %v", err)
		response.InternalServerError(c, "failed to update supplier")
		return
	}

	response.Success(c, http.StatusOK, "supplier updated successfully", supplier)
}

func (h *V1) deactivateSupplier(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "supplier id is required")
		return
	}

	supplier, err := h.supplier.DeactivateSupplier(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSuppNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeSuppNotFound), err.Error())
			return
		}
		if errors.Is(err, domain.ErrSuppAlreadyInactive) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeSuppAlreadyInactive), err.Error())
			return
		}
		h.l.Error("h.deactivateSupplier: %v", err)
		response.InternalServerError(c, "failed to deactivate supplier")
		return
	}

	response.Success(c, http.StatusOK, "supplier deactivated successfully", supplier)
}

func (h *V1) reactivateSupplier(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "supplier id is required")
		return
	}

	supplier, err := h.supplier.ReactivateSupplier(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSuppNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeSuppNotFound), err.Error())
			return
		}
		if errors.Is(err, domain.ErrSuppAlreadyActive) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeSuppAlreadyActive), err.Error())
			return
		}
		h.l.Error("h.reactivateSupplier: %v", err)
		response.InternalServerError(c, "failed to reactivate supplier")
		return
	}

	response.Success(c, http.StatusOK, "supplier reactivated successfully", supplier)
}

func (h *V1) listSuppliers(c *gin.Context) {
	activeOnlyStr := c.DefaultQuery("active_only", "false")
	activeOnly, _ := strconv.ParseBool(activeOnlyStr)

	params := pagination.FromQuery(c)

	suppliers, err := h.supplier.ListSuppliers(c.Request.Context(), activeOnly, params)
	if err != nil {
		h.l.Error("h.listSuppliers: %v", err)
		response.InternalServerError(c, "failed to list suppliers")
		return
	}

	response.Success(c, http.StatusOK, "suppliers retrieved successfully", suppliers)
}
