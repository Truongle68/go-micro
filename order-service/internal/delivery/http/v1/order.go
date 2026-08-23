package v1

import (
	"errors"
	"net/http"
	"strings"

	"order-service/internal/delivery/http/v1/req"
	"order-service/internal/domain"
	"order-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/pagination"
	"github.com/TruongLe68/go-micro/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *V1) checkout(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	var r req.CheckoutReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	token := extractBearerToken(c)

	itemsInput := make([]usecase.CheckoutItemInput, len(r.Items))
	for i, item := range r.Items {
		itemsInput[i] = usecase.CheckoutItemInput{
			ProductID:    item.ProductID,
			VariantID:    item.VariantID,
			SKU:          item.SKU,
			ProductName:  item.ProductName,
			Image:        item.Image,
			VariantAttrs: item.VariantAttrs,
			UnitPrice:    item.UnitPrice,
			Quantity:     item.Quantity,
		}
	}

	input := usecase.CheckoutInput{
		Items:           itemsInput,
		ShippingAddress: r.ShippingAddress,
		ShippingFee:     r.ShippingFee,
		PaymentMethod:   r.PaymentMethod,
	}

	order, err := h.o.Checkout(c.Request.Context(), userID, input, token)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyOrderItems) || errors.Is(err, domain.ErrInvalidAddress) || errors.Is(err, domain.ErrCartEmpty) {
			response.Error(c, http.StatusBadRequest, response.CodeValidationError, err.Error())
			return
		}
		h.l.Error("h.checkout - h.o.Checkout: %v", err)
		response.InternalServerError(c, "checkout failed")
		return
	}

	response.Success(c, http.StatusCreated, "order created successfully", order)
}

func (h *V1) listOrders(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	params := pagination.FromQuery(c)

	result, err := h.o.ListOrdersByUser(c.Request.Context(), userID, params)
	if err != nil {
		h.l.Error("h.listOrders - h.o.ListOrdersByUser: %v", err)
		response.InternalServerError(c, "failed to list orders")
		return
	}

	response.SuccessPaginated(c, http.StatusOK, "orders retrieved successfully", *result)
}

func (h *V1) getOrder(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "order id is required")
		return
	}

	order, err := h.o.GetOrder(c.Request.Context(), orderID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeOrderNotFound), err.Error())
			return
		}
		h.l.Error("h.getOrder - h.o.GetOrder: %v", err)
		response.InternalServerError(c, "failed to get order")
		return
	}

	response.Success(c, http.StatusOK, "order retrieved successfully", order)
}

func (h *V1) getTracking(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "order id is required")
		return
	}

	timeline, err := h.o.GetTrackingTimeline(c.Request.Context(), orderID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeOrderNotFound), err.Error())
			return
		}
		h.l.Error("h.getTracking - h.o.GetTrackingTimeline: %v", err)
		response.InternalServerError(c, "failed to get tracking timeline")
		return
	}

	response.Success(c, http.StatusOK, "tracking timeline retrieved successfully", timeline)
}

func (h *V1) shipOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "order id is required")
		return
	}

	var r req.ShipOrderReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.InvalidRequestBody(c, err.Error())
		return
	}

	order, err := h.o.ShipOrder(c.Request.Context(), orderID, r.TrackingCode)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeOrderNotFound), err.Error())
			return
		}
		if errors.Is(err, domain.ErrInvalidOrderTransition) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeInvalidOrderTransition), err.Error())
			return
		}
		h.l.Error("h.shipOrder - h.o.ShipOrder: %v", err)
		response.InternalServerError(c, "failed to ship order")
		return
	}

	response.Success(c, http.StatusOK, "order shipped successfully", order)
}

func (h *V1) deliverOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "order id is required")
		return
	}

	order, err := h.o.DeliverOrder(c.Request.Context(), orderID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeOrderNotFound), err.Error())
			return
		}
		if errors.Is(err, domain.ErrInvalidOrderTransition) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeInvalidOrderTransition), err.Error())
			return
		}
		h.l.Error("h.deliverOrder - h.o.DeliverOrder: %v", err)
		response.InternalServerError(c, "failed to deliver order")
		return
	}

	response.Success(c, http.StatusOK, "order delivered successfully", order)
}

func (h *V1) cancelOrder(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok || userID == "" {
		response.Unauthorized(c)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "order id is required")
		return
	}

	var r req.CancelOrderReq
	_ = c.ShouldBindJSON(&r)

	order, err := h.o.CancelOrder(c.Request.Context(), orderID, userID, r.Reason)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, string(domain.CodeOrderNotFound), err.Error())
			return
		}
		if errors.Is(err, domain.ErrOrderAlreadyCancelled) || errors.Is(err, domain.ErrInvalidOrderTransition) {
			response.Error(c, http.StatusBadRequest, string(domain.CodeOrderAlreadyCancelled), err.Error())
			return
		}
		h.l.Error("h.cancelOrder - h.o.CancelOrder: %v", err)
		response.InternalServerError(c, "failed to cancel order")
		return
	}

	response.Success(c, http.StatusOK, "order cancelled successfully", order)
}

func extractBearerToken(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	const prefix = "Bearer "
	if strings.HasPrefix(h, prefix) {
		return strings.TrimSpace(h[len(prefix):])
	}
	cookie, err := c.Cookie("access_token")
	if err == nil {
		return cookie
	}
	return ""
}
