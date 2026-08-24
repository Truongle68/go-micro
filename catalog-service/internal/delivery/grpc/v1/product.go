package v1

import (
	"catalog-service/internal/domain"
	"catalog-service/internal/usecase"
	"catalog-service/pkg/sliceutil"
	"context"
	"errors"

	catalogv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/catalog/v1"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServer struct {
	catalogv1.UnimplementedProductServiceServer
	uc ProductGRPCUsecase
	l  logger.Interface
}

func NewProductServer(uc ProductGRPCUsecase, l logger.Interface) *ProductServer {
	return &ProductServer{
		uc: uc,
		l:  l,
	}
}

func (s *ProductServer) GetVariantsBySKUs(ctx context.Context, req *catalogv1.GetVariantsBySKURequest) (*catalogv1.GetVariantsBySKUResponse, error) {
	if len(req.GetSku()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "skus are required")
	}

	uniqueSKUs := sliceutil.DedupeString(req.GetSku(), func(sku string) string {
		return sku
	})

	variants, err := s.uc.GetVariantsBySKUs(ctx, uniqueSKUs)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSKUNotFound):
			return nil, status.Errorf(codes.NotFound, "%v", err)
		case errors.Is(err, domain.ErrInactiveVariant):
			return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
		case errors.Is(err, domain.ErrEmptySKU):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			s.l.Error("grpc.GetVariantsBySKUs: %v", err)
			return nil, status.Error(codes.Internal, "failed to get variants by skus")
		}
	}

	return &catalogv1.GetVariantsBySKUResponse{
		Variants: mapVariantsToProto(variants),
	}, nil
}

func mapPriceToProto(p domain.Price) *catalogv1.Price {
	return &catalogv1.Price{
		Currency: p.Currency,
		Amount:   int32(p.Amount),
	}
}

func mapVariantToProto(v usecase.VariantView) *catalogv1.Variant {
	return &catalogv1.Variant{
		Id:          v.ID,
		ProductId:   v.ProductID,
		ProductName: v.ProductName,
		Sku:         v.SKU,
		Attributes:  v.Attributes,
		Price:       mapPriceToProto(v.Price),
		Image:       v.Image,
		IsActive:    v.IsActive,
	}
}

func mapVariantsToProto(variants []usecase.VariantView) []*catalogv1.Variant {
	protoVariants := make([]*catalogv1.Variant, 0, len(variants))
	for _, v := range variants {
		protoVariants = append(protoVariants, mapVariantToProto(v))
	}
	return protoVariants
}
