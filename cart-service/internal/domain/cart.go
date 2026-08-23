package domain

import (
	"strings"
	"time"
)

type Cart struct {
	UserID    string     `json:"user_id"`
	Items     []CartItem `json:"items"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

func NewCart(userID string) *Cart {
	return &Cart{
		UserID:    userID,
		Items:     make([]CartItem, 0),
		UpdatedAt: time.Now().UTC(),
	}
}

func (c *Cart) AddItem(sku string, quantity int) error {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return ErrInvalidSKU
	}
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	for i := range c.Items {
		if c.Items[i].SKU == sku {
			c.Items[i].Quantity += quantity
			c.UpdatedAt = time.Now().UTC()
			return nil
		}
	}

	c.Items = append(c.Items, CartItem{
		SKU:      sku,
		Quantity: quantity,
	})
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Cart) UpdateItem(sku string, quantity int) error {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return ErrInvalidSKU
	}

	if quantity <= 0 {
		return c.RemoveItem(sku)
	}

	found := false
	for i := range c.Items {
		if c.Items[i].SKU == sku {
			c.Items[i].Quantity = quantity
			found = true
			break
		}
	}

	if !found {
		return ErrItemNotFound
	}

	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Cart) RemoveItem(sku string) error {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return ErrInvalidSKU
	}

	foundIdx := -1
	for i := range c.Items {
		if c.Items[i].SKU == sku {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		return ErrItemNotFound
	}

	c.Items = append(c.Items[:foundIdx], c.Items[foundIdx+1:]...)
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Cart) Clear() {
	c.Items = make([]CartItem, 0)
	c.UpdatedAt = time.Now().UTC()
}
