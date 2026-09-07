package v1

import (
	"errors"
	"net/http"
	"strings"

	"inventory-service/internal/delivery/http/v1/req"
	"inventory-service/internal/domain"
	"inventory-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

// Public endpoints

// checkAvailability returns the available stock for multiple SKUs.
// Supports ?skus=SKU1,SKU2 or multiple ?skus=SKU1&skus=SKU2.
func (h *V1) checkAvailability(c *gin.Context) {
	rawSKUs := c.QueryArray("skus")
	var skus []string

	for _, s := range rawSKUs {
		for _, part := range strings.Split(s, ",") {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				skus = append(skus, trimmed)
			}
		}
	}

	if len(skus) == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "at least one sku must be provided")
		return
	}

	avail, err := h.stock.CheckStock(c.Request.Context(), skus)
	if err != nil {
		h.l.Error("h.checkAvailability: %v", err)
		response.InternalServerError(c, "failed to check stock availability")
		return
	}

	response.Success(c, http.StatusOK, "stock availability retrieved successfully", gin.H{
		"availability": avail,
	})
}

// getSKUAvailability returns the available stock for a single SKU.
func (h *V1) getSKUAvailability(c *gin.Context) {
	sku := strings.TrimSpace(c.Param("sku"))
	if sku == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "sku is required")
		return
	}

	avail, err := h.stock.GetSKUAvailability(c.Request.Context(), sku)
	if err != nil {
		if errors.Is(err, domain.ErrEmptySKU) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeEmptySKU), err.Error())
			return
		}
		h.l.Error("h.getSKUAvailability: %v", err)
		response.InternalServerError(c, "failed to check sku availability")
		return
	}

	response.Success(c, http.StatusOK, "sku availability retrieved successfully", gin.H{
		"sku":       sku,
		"available": avail,
		"in_stock":  avail > 0,
	})
}

// Protected Admin endpoints

// listStockLevels returns paginated detailed stock levels with joined warehouse and catalog info.
func (h *V1) listStockLevels(c *gin.Context) {
	filter := domain.StockLevelFilter{
		WarehouseID: c.Query("warehouse_id"),
		SKU:         c.Query("sku"),
		LowStock:    c.Query("low_stock") == "true",
	}

	params := pagination.FromQuery(c)

	levels, total, err := h.stock.ListStockLevels(c.Request.Context(), filter, params)
	if err != nil {
		h.l.Error("h.listStockLevels: %v", err)
		response.InternalServerError(c, "failed to list stock levels")
		return
	}

	response.SuccessPaginated(c, http.StatusOK, "stock levels retrieved successfully", pagination.NewResult(levels, params, total))
}

// getStockLevel returns a single detailed stock level by ID.
func (h *V1) getStockLevel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "stock level id is required")
		return
	}

	level, err := h.stock.GetStockLevel(c.Request.Context(), id)
	if err != nil {
		h.l.Error("h.getStockLevel: %v", err)
		response.InternalServerError(c, "failed to get stock level")
		return
	}
	if level == nil {
		response.Error(c, http.StatusNotFound, string(domain.CodeStockLevelNotFound), "stock level not found")
		return
	}

	response.Success(c, http.StatusOK, "stock level retrieved successfully", level)
}

// adjustStock handles manual stock adjustments and logs stock movements.
func (h *V1) adjustStock(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	var r req.AdjustStockReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	level, err := h.stock.AdjustStock(c.Request.Context(), usecase.AdjustStockInput{
		SKU:           r.SKU,
		WarehouseID:   r.WarehouseID,
		QuantityDelta: r.QuantityDelta,
		Reason:        r.Reason,
		Note:          r.Note,
		CreatedBy:     userID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrEmptySKU) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeEmptySKU), err.Error())
			return
		}
		if errors.Is(err, domain.ErrEmptyWhID) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeEmptyWarehouseID), err.Error())
			return
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeInsufficientStock), err.Error())
			return
		}
		if appErr := domain.ToAppError(err); appErr != nil {
			response.Error(c, http.StatusBadRequest, string(appErr.Code), appErr.Message)
			return
		}
		h.l.Error("h.adjustStock: %v", err)
		response.InternalServerError(c, "failed to adjust stock")
		return
	}

	response.Success(c, http.StatusOK, "stock adjusted successfully", level)
}

// updateThresholds updates reorder threshold and quantity for a stock level.
func (h *V1) updateThresholds(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "stock level id is required")
		return
	}

	var r req.UpdateThresholdsReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	level, err := h.stock.UpdateThresholds(c.Request.Context(), id, r.ReorderThreshold, r.ReorderQuantity)
	if err != nil {
		if errors.Is(err, domain.ErrStockLevelNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeStockLevelNotFound), err.Error())
			return
		}
		if errors.Is(err, domain.ErrNegativeQuantity) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeNegativeQuantity), err.Error())
			return
		}
		h.l.Error("h.updateThresholds: %v", err)
		response.InternalServerError(c, "failed to update stock thresholds")
		return
	}

	response.Success(c, http.StatusOK, "stock thresholds updated successfully", level)
}

// getStockSummary returns KPI cards for inventory management.
func (h *V1) getStockSummary(c *gin.Context) {
	warehouseID := c.Query("warehouse_id")

	summary, err := h.stock.GetStockSummary(c.Request.Context(), warehouseID)
	if err != nil {
		h.l.Error("h.getStockSummary: %v", err)
		response.InternalServerError(c, "failed to get stock summary")
		return
	}

	response.Success(c, http.StatusOK, "stock summary retrieved successfully", summary)
}

// listStockMovements returns paginated audit history of stock movements.
func (h *V1) listStockMovements(c *gin.Context) {
	filter := domain.StockMovementFilter{
		SKU:         c.Query("sku"),
		WarehouseID: c.Query("warehouse_id"),
		Type:        domain.MovementType(c.Query("type")),
	}

	params := pagination.FromQuery(c)

	movements, total, err := h.stock.ListMovements(c.Request.Context(), filter, params)
	if err != nil {
		h.l.Error("h.listStockMovements: %v", err)
		response.InternalServerError(c, "failed to list stock movements")
		return
	}

	response.SuccessPaginated(c, http.StatusOK, "stock movements retrieved successfully", pagination.NewResult(movements, params, total))
}
