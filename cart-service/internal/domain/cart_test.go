package domain_test

import (
	"testing"

	"cart-service/internal/domain"
)

func TestCartOperations(t *testing.T) {
	c := domain.NewCart("user123")

	if len(c.Items) != 0 {
		t.Fatalf("expected empty cart, got %d items", len(c.Items))
	}

	// Add item
	err := c.AddItem("SKU-1", 2)
	if err != nil {
		t.Fatalf("unexpected error adding item: %v", err)
	}
	if len(c.Items) != 1 || c.Items[0].Quantity != 2 {
		t.Fatalf("expected 1 item with qty 2, got %+v", c.Items)
	}

	// Add same item (increment)
	err = c.AddItem("SKU-1", 3)
	if err != nil {
		t.Fatalf("unexpected error adding existing item: %v", err)
	}
	if len(c.Items) != 1 || c.Items[0].Quantity != 5 {
		t.Fatalf("expected 1 item with qty 5, got %+v", c.Items)
	}

	// Add second item
	err = c.AddItem("SKU-2", 1)
	if err != nil {
		t.Fatalf("unexpected error adding second item: %v", err)
	}
	if len(c.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(c.Items))
	}

	// Update item quantity
	err = c.UpdateItem("SKU-1", 10)
	if err != nil {
		t.Fatalf("unexpected error updating item: %v", err)
	}
	if c.Items[0].Quantity != 10 {
		t.Fatalf("expected qty 10, got %d", c.Items[0].Quantity)
	}

	// Update item quantity to 0 (should remove item)
	err = c.UpdateItem("SKU-1", 0)
	if err != nil {
		t.Fatalf("unexpected error updating item to 0: %v", err)
	}
	if len(c.Items) != 1 || c.Items[0].SKU != "SKU-2" {
		t.Fatalf("expected SKU-1 to be removed, got %+v", c.Items)
	}

	// Remove item
	err = c.RemoveItem("SKU-2")
	if err != nil {
		t.Fatalf("unexpected error removing item: %v", err)
	}
	if len(c.Items) != 0 {
		t.Fatalf("expected empty cart, got %d items", len(c.Items))
	}

	// Remove non-existent item
	err = c.RemoveItem("SKU-999")
	if err != domain.ErrItemNotFound {
		t.Fatalf("expected ErrItemNotFound, got %v", err)
	}
}
