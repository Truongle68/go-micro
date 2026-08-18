package v1

import (
	"errors"
	"net/http"

	"cart-service/internal/delivery/http/v1/req"
	"cart-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *V1) getCart(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	cart, err := h.c.GetCart(c.Request.Context(), userID)
	if err != nil {
		h.l.Error("h.getCart - h.c.GetCart: %v", err)
		response.InternalServerError(c, domain.ErrFailedToRetrieveCart.Error())
		return
	}

	response.Success(c, http.StatusOK, "cart retrieved successfully", cart)
}

func (h *V1) addItem(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	var r req.AddItemReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	cart, err := h.c.AddItem(c.Request.Context(), userID, r.SKU, r.Quantity)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidSKU) || errors.Is(err, domain.ErrInvalidQuantity) {
			response.Error(c, http.StatusBadRequest, response.CodeValidationError, err.Error())
			return
		}
		h.l.Error("h.addItem - h.c.AddItem: %v", err)
		response.InternalServerError(c, domain.ErrFailedToAddToCart.Error())
		return
	}

	response.Success(c, http.StatusOK, "item added to cart", cart)
}

func (h *V1) updateItemQuantity(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	sku := c.Param("sku")
	if sku == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, domain.ErrSKURequired.Error())
		return
	}

	var r req.UpdateItemReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	cart, err := h.c.UpdateItemQuantity(c.Request.Context(), userID, sku, r.Quantity)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) || errors.Is(err, domain.ErrCartNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeNotFound), err.Error())
			return
		}
		if errors.Is(err, domain.ErrInvalidSKU) || errors.Is(err, domain.ErrInvalidQuantity) {
			response.Error(c, http.StatusBadRequest, response.CodeValidationError, err.Error())
			return
		}
		h.l.Error("h.updateItemQuantity - h.c.UpdateItemQuantity: %v", err)
		response.InternalServerError(c, domain.ErrFailedToUpdateItem.Error())
		return
	}

	response.Success(c, http.StatusOK, "item quantity updated", cart)
}

func (h *V1) removeItem(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	sku := c.Param("sku")
	if sku == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, domain.ErrSKURequired.Error())
		return
	}

	cart, err := h.c.RemoveItem(c.Request.Context(), userID, sku)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) || errors.Is(err, domain.ErrCartNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeNotFound), err.Error())
			return
		}
		h.l.Error("h.removeItem - h.c.RemoveItem: %v", err)
		response.InternalServerError(c, domain.ErrFailedToRemoveItem.Error())
		return
	}

	response.Success(c, http.StatusOK, "item removed from cart", cart)
}

func (h *V1) clearCart(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	if err := h.c.ClearCart(c.Request.Context(), userID); err != nil {
		h.l.Error("h.clearCart - h.c.ClearCart: %v", err)
		response.InternalServerError(c, domain.ErrFailedToClearCart.Error())
		return
	}

	response.Success(c, http.StatusOK, "cart cleared successfully", nil)
}
