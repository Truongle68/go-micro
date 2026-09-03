package domain

import (
	"strings"
	"time"
)

type PurchaseOrder struct {
	ID           string              `json:"id"`
	Version      int                 `json:"version"`
	Code         string              `json:"code"`
	SupplierID   string              `json:"supplier_id"`
	SupplierName string              `json:"supplier_name"`
	WarehouseID  string              `json:"warehouse_id"`
	Lines        []PurchaseOrderLine `json:"lines"`
	Status       PurchaseOrderStatus `json:"status"`
	CreatedBy    string              `json:"created_by"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	ExpectedAt   *time.Time          `json:"expected_at,omitempty"`
	ReceivedAt   *time.Time          `json:"received_at,omitempty"`
}

type PurchaseOrderLine struct {
	SKU              string `json:"sku"`
	ProductName      string `json:"product_name"`
	QuantityOrdered  int    `json:"quantity_ordered"`
	QuantityReceived int    `json:"quantity_received"`
	UnitCost         int64  `json:"unit_cost"`
}

type PurchaseOrderStatus string

const (
	POStatusDraft             PurchaseOrderStatus = "draft"
	POStatusOrdered           PurchaseOrderStatus = "ordered"
	POStatusPartiallyReceived PurchaseOrderStatus = "partially_received"
	POStatusReceived          PurchaseOrderStatus = "received"
	POStatusCancelled         PurchaseOrderStatus = "cancelled"
)

type NewPurchaseOrderLineInput struct {
	SKU             string
	ProductName     string
	QuantityOrdered int
	UnitCost        int64
}

func NewPurchaseOrder(
	code, supplierID, supplierName, warehouseID, createdBy string,
	lineInputs []NewPurchaseOrderLineInput,
) (*PurchaseOrder, error) {
	if strings.TrimSpace(supplierID) == "" {
		return nil, ErrEmptySuppID
	}
	if strings.TrimSpace(warehouseID) == "" {
		return nil, ErrEmptyWhID
	}
	if len(lineInputs) == 0 {
		return nil, ErrEmptyPurchaseOrderLines
	}

	lines := make([]PurchaseOrderLine, 0, len(lineInputs))
	seenSKUs := make(map[string]bool, len(lineInputs))
	for _, in := range lineInputs {
		if strings.TrimSpace(in.SKU) == "" {
			return nil, ErrEmptySKU
		}
		if in.QuantityOrdered <= 0 {
			return nil, ErrNonPositiveQuantity
		}
		if in.UnitCost < 0 {
			return nil, ErrInvalidPrice
		}
		if seenSKUs[in.SKU] {
			return nil, ErrPODuplicateSKU
		}
		seenSKUs[in.SKU] = true

		lines = append(lines, PurchaseOrderLine{
			SKU:              in.SKU,
			ProductName:      in.ProductName,
			QuantityOrdered:  in.QuantityOrdered,
			QuantityReceived: 0,
			UnitCost:         in.UnitCost,
		})
	}

	now := time.Now().UTC()
	return &PurchaseOrder{
		Code: code, Version: 1, SupplierID: supplierID, SupplierName: supplierName,
		WarehouseID: warehouseID, Lines: lines,
		Status: POStatusDraft, CreatedBy: createdBy,
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (po PurchaseOrder) TotalCost() int64 {
	var total int64
	for _, l := range po.Lines {
		total += l.UnitCost * int64(l.QuantityOrdered)
	}
	return total
}

func (po *PurchaseOrder) Confirm() error {
	if po.Status != POStatusDraft {
		return ErrInvalidPOTransition
	}
	po.Status = POStatusOrdered
	po.Version++
	po.UpdatedAt = time.Now().UTC()
	return nil
}

// ReceiveLine records that `qty` units of a specific SKU physically
// arrived. This does NOT touch StockLevel/StockBatch directly — that's
// the usecase layer's job (same separation as Order.Cancel not calling
// inventory-service itself). This method only updates the PO's OWN
// bookkeeping and decides whether the whole PO's status should flip.
func (po *PurchaseOrder) ReceiveLine(sku string, qty int) error {
	if po.Status != POStatusOrdered && po.Status != POStatusPartiallyReceived {
		return ErrInvalidPOTransition
	}
	if qty <= 0 {
		return ErrNonPositiveQuantity
	}

	found := false
	for i := range po.Lines {
		if po.Lines[i].SKU != sku {
			continue
		}
		found = true
		newReceived := po.Lines[i].QuantityReceived + qty
		if newReceived > po.Lines[i].QuantityOrdered {
			return ErrReceivedExceedsOrdered
		}
		po.Lines[i].QuantityReceived = newReceived
		break
	}
	if !found {
		return ErrEmptySKU
	}

	po.Status = po.computeReceivingStatus()
	po.Version++
	po.UpdatedAt = time.Now().UTC()
	if po.Status == POStatusReceived {
		now := time.Now().UTC()
		po.ReceivedAt = &now
	}
	return nil
}

func (po PurchaseOrder) computeReceivingStatus() PurchaseOrderStatus {
	allReceived, anyReceived := true, false
	for _, l := range po.Lines {
		if l.QuantityReceived > 0 {
			anyReceived = true
		}
		if l.QuantityReceived < l.QuantityOrdered {
			allReceived = false
		}
	}
	switch {
	case allReceived:
		return POStatusReceived
	case anyReceived:
		return POStatusPartiallyReceived
	default:
		return po.Status // no lines received yet — stays wherever it was (Ordered)
	}
}

func (po *PurchaseOrder) Cancel() error {
	if po.Status == POStatusReceived || po.Status == POStatusCancelled {
		return ErrInvalidPOTransition
	}
	po.Status = POStatusCancelled
	po.Version++
	po.UpdatedAt = time.Now().UTC()
	return nil
}

type PurchaseOrderFilter struct {
	SupplierID  string
	WarehouseID string
	Status      string
}
