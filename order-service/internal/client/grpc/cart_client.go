package grpc

import (
	"context"
	"fmt"

	"order-service/internal/client"

	cartv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/cart/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type cartGRPCClient struct {
	client cartv1.CartServiceClient
}

func NewCartGRPCClient(target string) (client.CartClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cart gRPC service at %s: %w", target, err)
	}

	return &cartGRPCClient{
		client: cartv1.NewCartServiceClient(conn),
	}, nil
}

func (c *cartGRPCClient) GetCart(ctx context.Context, userID string, token string) (*client.CartDTO, error) {
	res, err := c.client.GetCart(ctx, &cartv1.GetCartRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc GetCart failed for user %s: %w", userID, err)
	}

	if res.GetCart() == nil {
		return &client.CartDTO{UserID: userID, Items: []client.CartItemDTO{}}, nil
	}

	items := make([]client.CartItemDTO, len(res.GetCart().GetItems()))
	for i, item := range res.GetCart().GetItems() {
		items[i] = client.CartItemDTO{
			SKU:      item.GetSku(),
			Quantity: int(item.GetQuantity()),
		}
	}

	return &client.CartDTO{
		UserID: res.GetCart().GetUserId(),
		Items:  items,
	}, nil
}

func (c *cartGRPCClient) ClearCart(ctx context.Context, userID string, token string) error {
	res, err := c.client.ClearCart(ctx, &cartv1.ClearCartRequest{
		UserId: userID,
	})
	if err != nil {
		return fmt.Errorf("grpc ClearCart failed for user %s: %w", userID, err)
	}

	if !res.GetSuccess() {
		return fmt.Errorf("grpc ClearCart for user %s returned success=false", userID)
	}

	return nil
}
