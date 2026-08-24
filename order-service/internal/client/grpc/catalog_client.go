package grpc

import (
	"context"
	"fmt"
	"order-service/internal/client"
	"order-service/internal/domain"

	catalogv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/catalog/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, fmt.Errorf("%w: %s", domain.ErrSKUNotFound, st.Message())
			case codes.FailedPrecondition:
				return nil, fmt.Errorf("%w: %s", domain.ErrInactiveVariant, st.Message())
			case codes.InvalidArgument:
				return nil, fmt.Errorf("%w: %s", domain.ErrEmptySKU, st.Message())
			}
		}
		return nil, fmt.Errorf("catalog GetVariantsBySKUs: %w", err)
	}
	if res.GetVariants() == nil {
		return []client.VariantDTO{}, nil
	}

	variants := make([]client.VariantDTO, len(res.GetVariants()))
	for i, v := range res.GetVariants() {
		variants[i] = client.VariantDTO{
			ID:          v.GetId(),
			ProductID:   v.GetProductId(),
			ProductName: v.GetProductName(),
			SKU:         v.GetSku(),
			Attributes:  v.GetAttributes(),
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
