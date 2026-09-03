package grpc

import (
	"context"
	"fmt"

	"order-service/internal/client"
	"order-service/internal/domain"

	inventoryv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type inventoryGRPCClient struct {
	client inventoryv1.InventoryServiceClient
}

// NewInventoryGRPCClient creates a new gRPC-backed InventoryClient.
func NewInventoryGRPCClient(target string) (client.InventoryClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory gRPC service at %s: %w", target, err)
	}

	return &inventoryGRPCClient{
		client: inventoryv1.NewInventoryServiceClient(conn),
	}, nil
}

func (c *inventoryGRPCClient) CheckStock(ctx context.Context, items []client.SKUQty) (map[string]int, error) {
	skus := make([]string, len(items))
	for i, item := range items {
		skus[i] = item.SKU
	}

	res, err := c.client.CheckStock(ctx, &inventoryv1.CheckStockRequest{
		Skus: skus,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc CheckStock: %w", err)
	}

	result := make(map[string]int, len(res.GetItems()))
	for _, item := range res.GetItems() {
		result[item.GetSku()] = int(item.GetAvailable())
	}

	return result, nil
}

func (c *inventoryGRPCClient) ReserveStock(ctx context.Context, orderID string, items []client.SKUQty) error {
	protoItems := make([]*inventoryv1.SKUQuantity, len(items))
	for i, item := range items {
		protoItems[i] = &inventoryv1.SKUQuantity{
			Sku:      item.SKU,
			Quantity: int32(item.Quantity),
		}
	}

	_, err := c.client.ReserveStock(ctx, &inventoryv1.ReserveStockRequest{
		OrderId: orderID,
		Items:   protoItems,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.FailedPrecondition {
			return fmt.Errorf("%w: %s", domain.ErrInsufficientStock, st.Message())
		}
		return fmt.Errorf("grpc ReserveStock: %w", err)
	}

	return nil
}

func (c *inventoryGRPCClient) ConfirmReservation(ctx context.Context, orderID string) error {
	_, err := c.client.ConfirmReservation(ctx, &inventoryv1.ConfirmReservationRequest{
		OrderId: orderID,
	})
	if err != nil {
		return fmt.Errorf("grpc ConfirmReservation: %w", err)
	}

	return nil
}

func (c *inventoryGRPCClient) ReleaseReservation(ctx context.Context, orderID string) error {
	_, err := c.client.ReleaseReservation(ctx, &inventoryv1.ReleaseReservationRequest{
		OrderId: orderID,
	})
	if err != nil {
		return fmt.Errorf("grpc ReleaseReservation: %w", err)
	}

	return nil
}
