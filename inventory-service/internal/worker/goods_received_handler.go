package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory-service/internal/domain"
	"inventory-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq"
)

type GoodsReceivedHandler struct {
	stockUC *usecase.StockUC
	logger  logger.Interface
}

func NewGoodsReceivedHandler(stockUC *usecase.StockUC, l logger.Interface) *GoodsReceivedHandler {
	return &GoodsReceivedHandler{
		stockUC: stockUC,
		logger:  l,
	}
}

func (h *GoodsReceivedHandler) Handle(ctx context.Context, event rabbitmq.Event) error {
	if event.Type != "" && event.Type != domain.EventGoodsReceived {
		return fmt.Errorf("unexpected event type=%s", event.Type)
	}

	var payload domain.GoodsReceivedEvent
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("GoodsReceivedHandler unmarshal: %w", err)
	}

	if err := h.stockUC.ApplyGoodsReceived(ctx, payload); err != nil {
		return fmt.Errorf("GoodsReceivedHandler apply: %w", err)
	}
	h.logger.Info("handled goods received po=%s, lines=%d", payload.PurchaseOrderID, len(payload.Lines))
	return nil
}
