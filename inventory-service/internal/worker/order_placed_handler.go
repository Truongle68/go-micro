package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq"
)

type OrderPlacedHandler struct {
	stockUC stockUseCase
	logger  logger.Interface
}

func NewOrderPlacedHandler(stockUC stockUseCase, l logger.Interface) *OrderPlacedHandler {
	return &OrderPlacedHandler{
		stockUC: stockUC,
		logger:  l,
	}
}

func (h *OrderPlacedHandler) Handle(ctx context.Context, event rabbitmq.Event) error {
	if event.Type != "" && event.Type != domain.EventGoodsReceived {
		return fmt.Errorf("unexpected event type=%s", event.Type)
	}

	var payload domain.GoodsReceivedEvent
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("OrderPlacedHandler unmarshal: %w", err)
	}

	if err := h.stockUC.ApplyGoodsReceived(ctx, payload); err != nil {
		return fmt.Errorf("OrderPlacedHandler apply: %w", err)
	}
	h.logger.Info("handled goods received po=%s, lines=%d", payload.PurchaseOrderID, len(payload.Lines))
	return nil
}
