package v1

import (
	"errors"
	"net/http"

	"inventory-service/internal/delivery/http/v1/req"
	"inventory-service/internal/domain"
	"inventory-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *V1) createPurchaseOrder(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	var r req.CreatePurchaseOrderReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	lines := make([]domain.NewPurchaseOrderLineInput, len(r.Lines))
	for i, l := range r.Lines {
		lines[i] = domain.NewPurchaseOrderLineInput{
			SKU:             l.SKU,
			ProductName:     "",
			QuantityOrdered: l.QuantityOrdered,
			UnitCost:        l.UnitCost,
		}
	}

	input := usecase.CreatePurchaseOrderInput{
		Code:        r.Code,
		SupplierID:  r.SupplierID,
		WarehouseID: r.WarehouseID,
		CreatedBy:   userID,
		Lines:       lines,
	}

	po, err := h.po.CreatePurchaseOrder(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrSuppNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeSuppNotFound), "supplier not found")
			return
		}
		if errors.Is(err, domain.ErrWhNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeWhNotFound), "warehouse not found")
			return
		}
		if errors.Is(err, domain.ErrSKUNotFound) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeSKUNotFound), err.Error())
			return
		}
		if errors.Is(err, domain.ErrDuplicatePOCode) {
			response.Error(c, http.StatusConflict, string(domain.CodeDuplicatePOCode), "purchase order with this code already exists")
			return
		}
		if appErr := domain.ToAppError(err); appErr != nil {
			response.Error(c, http.StatusBadRequest, string(appErr.Code), appErr.Message)
			return
		}
		h.l.Error("h.createPurchaseOrder: %v", err)
		response.InternalServerError(c, "failed to create purchase order")
		return
	}

	response.Success(c, http.StatusCreated, "purchase order created successfully", po)
}

func (h *V1) getPurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "purchase order id is required")
		return
	}

	po, err := h.po.GetPurchaseOrder(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPONotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodePONotFound), err.Error())
			return
		}
		h.l.Error("h.getPurchaseOrder: %v", err)
		response.InternalServerError(c, "failed to get purchase order")
		return
	}

	response.Success(c, http.StatusOK, "purchase order retrieved successfully", po)
}

func (h *V1) confirmPurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "purchase order id is required")
		return
	}

	po, err := h.po.ConfirmPurchaseOrder(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPONotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodePONotFound), err.Error())
			return
		}
		if errors.Is(err, domain.ErrInvalidPOTransition) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeInvalidPOTransition), err.Error())
			return
		}
		h.l.Error("h.confirmPurchaseOrder: %v", err)
		response.InternalServerError(c, "failed to confirm purchase order")
		return
	}

	response.Success(c, http.StatusOK, "purchase order confirmed successfully", po)
}

func (h *V1) receiveLine(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "purchase order id is required")
		return
	}

	var r req.ReceiveLineReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	po, err := h.po.ReceiveLine(c.Request.Context(), id, r.SKU, r.Quantity)
	if err != nil {
		if errors.Is(err, domain.ErrPONotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodePONotFound), err.Error())
			return
		}
		if appErr := domain.ToAppError(err); appErr != nil {
			response.Error(c, http.StatusBadRequest, string(appErr.Code), appErr.Message)
			return
		}
		h.l.Error("h.receiveLine: %v", err)
		response.InternalServerError(c, "failed to receive line")
		return
	}

	response.Success(c, http.StatusOK, "line received successfully", po)
}

func (h *V1) cancelPurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "purchase order id is required")
		return
	}

	po, err := h.po.CancelPurchaseOrder(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPONotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodePONotFound), err.Error())
			return
		}
		if errors.Is(err, domain.ErrInvalidPOTransition) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeInvalidPOTransition), err.Error())
			return
		}
		h.l.Error("h.cancelPurchaseOrder: %v", err)
		response.InternalServerError(c, "failed to cancel purchase order")
		return
	}

	response.Success(c, http.StatusOK, "purchase order cancelled successfully", po)
}

func (h *V1) listPurchaseOrders(c *gin.Context) {
	filter := domain.PurchaseOrderFilter{
		SupplierID:  c.Query("supplier_id"),
		WarehouseID: c.Query("warehouse_id"),
		Status:      c.Query("status"),
	}

	params := pagination.FromQuery(c)

	orders, err := h.po.ListPurchaseOrders(c.Request.Context(), filter, params)
	if err != nil {
		h.l.Error("h.listPurchaseOrders: %v", err)
		response.InternalServerError(c, "failed to list purchase orders")
		return
	}

	response.Success(c, http.StatusOK, "purchase orders retrieved successfully", orders)
}
