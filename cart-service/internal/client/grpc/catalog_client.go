package grpc

import (
	"cart-service/internal/client"
	"context"
	"fmt"

	catalogv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/catalog/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type catalogGRPCClient struct {
	client catalogv1.ProductServiceClient
}

func NewCatalogGRPCClient(target string) (client.CatalogClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product gRPC service at %s: %w", target, err)
	}
	return &catalogGRPCClient{
		client: catalogv1.NewProductServiceClient(conn),
	}, nil
}

func (c *catalogGRPCClient) GetVariantsBySKUs(ctx context.Context, skus []string) ([]client.VariantDTO, error) {
	res, err := c.client.GetVariantsBySKUs(ctx, &catalogv1.GetVariantsBySKURequest{
		Sku: skus,
	})
	if err != nil {
		return []client.VariantDTO{}, fmt.Errorf("gRPC get variants by skus failed: %w", err)
	}
	if res.GetVariants() == nil {
		return []client.VariantDTO{}, nil
	}

	variants := make([]client.VariantDTO, len(res.GetVariants()))
	for i, v := range res.GetVariants() {
		variants[i] = client.VariantDTO{
			SKU:        v.GetSku(),
			Attributes: v.GetAttributes(),
			Price: client.Price{
				Amount:   int(v.GetPrice().GetAmount()),
				Currency: v.GetPrice().GetCurrency(),
			},
			Image:    v.GetImage(),
			IsActive: v.GetIsActive(),
		}
	}
	return variants, nil
}
