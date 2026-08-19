package domain_test

import (
	"testing"

	"order-service/internal/domain"
)

func TestOrderStateMachine(t *testing.T) {
	address := domain.AddressSnapshot{
		FullName: "John Doe",
		Phone:    "0901234567",
		Street:   "123 Main St",
		City:     "Ho Chi Minh",
	}

	items := []domain.OrderItem{
		{
			SKU:         "SKU-A",
			ProductName: "Test Product A",
			UnitPrice:   100,
			Quantity:    2,
		},
	}

	// 1. Create New Order -> status pending_payment
	order, history, err := domain.NewOrder("user_123", items, address, 15, "cod")
	if err != nil {
		t.Fatalf("unexpected error creating order: %v", err)
	}

	if order.Status != domain.OrderStatusPendingPayment {
		t.Fatalf("expected initial status pending_payment, got %s", order.Status)
	}
	if order.Subtotal != 200 || order.Total != 215 {
		t.Fatalf("expected subtotal 200 and total 215, got subtotal %d, total %d", order.Subtotal, order.Total)
	}
	if history.ToStatus != domain.OrderStatusPendingPayment {
		t.Fatalf("expected history toStatus pending_payment, got %s", history.ToStatus)
	}

	// 2. MarkConfirmed
	histConfirm, err := order.MarkConfirmed("PAY-123")
	if err != nil {
		t.Fatalf("unexpected error confirming order: %v", err)
	}
	if order.Status != domain.OrderStatusConfirmed || histConfirm.FromStatus != domain.OrderStatusPendingPayment {
		t.Fatalf("expected confirmed status, got %s", order.Status)
	}

	// 3. Prepare
	histPrep, err := order.Prepare()
	if err != nil {
		t.Fatalf("unexpected error preparing order: %v", err)
	}
	if order.Status != domain.OrderStatusPreparing || histPrep.FromStatus != domain.OrderStatusConfirmed {
		t.Fatalf("expected preparing status, got %s", order.Status)
	}

	// 4. Ship
	histShip, err := order.Ship("TRACK-999")
	if err != nil {
		t.Fatalf("unexpected error shipping order: %v", err)
	}
	if order.Status != domain.OrderStatusShipped || order.TrackingCode != "TRACK-999" || histShip.ToStatus != domain.OrderStatusShipped {
		t.Fatalf("expected shipped status with tracking code, got status %s, tracking %s", order.Status, order.TrackingCode)
	}

	// 5. Deliver
	histDeliv, err := order.Deliver()
	if err != nil {
		t.Fatalf("unexpected error delivering order: %v", err)
	}
	if order.Status != domain.OrderStatusDelivered || histDeliv.ToStatus != domain.OrderStatusDelivered {
		t.Fatalf("expected delivered status, got %s", order.Status)
	}

	// 6. Attempt invalid transition (e.g. deliver again or ship delivered order)
	_, err = order.Ship("TRACK-NEW")
	if err != domain.ErrInvalidOrderTransition {
		t.Fatalf("expected ErrInvalidOrderTransition, got %v", err)
	}
}

func TestOrderInvalidCreation(t *testing.T) {
	address := domain.AddressSnapshot{
		FullName: "John Doe",
		Phone:    "0901234567",
		Street:   "123 Main St",
		City:     "Ho Chi Minh",
	}

	items := []domain.OrderItem{
		{SKU: "SKU-A", UnitPrice: 100, Quantity: 1},
	}

	// Empty userID
	_, _, err := domain.NewOrder("", items, address, 0, "cod")
	if err != domain.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}

	// Empty items
	_, _, err = domain.NewOrder("user_1", nil, address, 0, "cod")
	if err != domain.ErrEmptyOrderItems {
		t.Fatalf("expected ErrEmptyOrderItems, got %v", err)
	}

	// Missing address
	_, _, err = domain.NewOrder("user_1", items, domain.AddressSnapshot{}, 0, "cod")
	if err != domain.ErrInvalidAddress {
		t.Fatalf("expected ErrInvalidAddress, got %v", err)
	}
}
