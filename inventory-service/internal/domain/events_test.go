package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"inventory-service/internal/domain"
)

func TestNewOutboxEvent(t *testing.T) {
	t.Run("valid payload struct", func(t *testing.T) {
		evt := domain.GoodsReceivedEvent{
			EventType:       domain.EventGoodsReceived,
			PurchaseOrderID: "po-123",
			POCode:          "PO-001",
			SupplierID:      "supp-1",
			SupplierCode:    "SUP-01",
			WarehouseID:     "wh-1",
			WarehouseCode:   "WH-01",
			ReceivedAt:      time.Now().UTC(),
			Lines: []domain.GoodsReceivedLine{
				{
					SKU:              "SKU-A",
					ProductName:      "Product A",
					QuantityOrdered:  10,
					QuantityReceived: 5,
					UnitCost:         1000,
				},
			},
		}

		outboxEvt, err := domain.NewOutboxEvent(
			domain.AggregatePurchaseOrder,
			evt.PurchaseOrderID,
			domain.EventGoodsReceived,
			evt,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if outboxEvt.AggregateType != domain.AggregatePurchaseOrder {
			t.Errorf("expected aggregate type %s, got %s", domain.AggregatePurchaseOrder, outboxEvt.AggregateType)
		}
		if outboxEvt.AggregateID != "po-123" {
			t.Errorf("expected aggregate id po-123, got %s", outboxEvt.AggregateID)
		}
		if outboxEvt.EventType != domain.EventGoodsReceived {
			t.Errorf("expected event type %s, got %s", domain.EventGoodsReceived, outboxEvt.EventType)
		}
		if len(outboxEvt.Payload) == 0 {
			t.Error("expected non-empty payload")
		}

		var unmarshaled domain.GoodsReceivedEvent
		if err := json.Unmarshal(outboxEvt.Payload, &unmarshaled); err != nil {
			t.Fatalf("failed to unmarshal payload: %v", err)
		}
		if unmarshaled.POCode != "PO-001" || len(unmarshaled.Lines) != 1 {
			t.Errorf("unexpected unmarshaled event: %+v", unmarshaled)
		}
	})

	t.Run("empty aggregate type", func(t *testing.T) {
		_, err := domain.NewOutboxEvent("", "123", "GoodsReceived", map[string]string{"foo": "bar"})
		if err == nil {
			t.Error("expected error for empty aggregate type")
		}
	})

	t.Run("empty aggregate id", func(t *testing.T) {
		_, err := domain.NewOutboxEvent("purchase_order", "", "GoodsReceived", map[string]string{"foo": "bar"})
		if err == nil {
			t.Error("expected error for empty aggregate id")
		}
	})

	t.Run("empty event type", func(t *testing.T) {
		_, err := domain.NewOutboxEvent("purchase_order", "123", "", map[string]string{"foo": "bar"})
		if err == nil {
			t.Error("expected error for empty event type")
		}
	})

	t.Run("nil payload", func(t *testing.T) {
		_, err := domain.NewOutboxEvent("purchase_order", "123", "GoodsReceived", nil)
		if err == nil {
			t.Error("expected error for nil payload")
		}
	})
}
