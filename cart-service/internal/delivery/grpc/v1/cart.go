package v1

import (
	"context"

	"cart-service/internal/domain"

	cartv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/cart/v1"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartServer struct {
	cartv1.UnimplementedCartServiceServer
	uc CartUseCase
	l  logger.Interface
}

type CartUseCase interface {
	GetCart(ctx context.Context, userID string) (*domain.Cart, error)
	AddItem(ctx context.Context, userID string, sku string, quantity int) (*domain.Cart, error)
	UpdateItemQuantity(ctx context.Context, userID string, sku string, quantity int) (*domain.Cart, error)
	RemoveItem(ctx context.Context, userID string, sku string) (*domain.Cart, error)
	RemoveItems(ctx context.Context, userID string, skus []string) (*domain.Cart, error)
	ClearCart(ctx context.Context, userID string) error
}

func NewCartServer(uc CartUseCase, l logger.Interface) *CartServer {
	return &CartServer{
		uc: uc,
		l:  l,
	}
}

func (s *CartServer) GetCart(ctx context.Context, req *cartv1.GetCartRequest) (*cartv1.GetCartResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	cart, err := s.uc.GetCart(ctx, req.GetUserId())
	if err != nil {
		s.l.Error("grpc.GetCart - uc.GetCart: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get cart: %v", err)
	}

	return &cartv1.GetCartResponse{
		Cart: mapCartToProto(cart),
	}, nil
}

func (s *CartServer) AddItem(ctx context.Context, req *cartv1.AddItemRequest) (*cartv1.AddItemResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetSku() == "" {
		return nil, status.Error(codes.InvalidArgument, "sku is required")
	}

	cart, err := s.uc.AddItem(ctx, req.GetUserId(), req.GetSku(), int(req.GetQuantity()))
	if err != nil {
		s.l.Error("grpc.AddItem - uc.AddItem: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to add item: %v", err)
	}

	return &cartv1.AddItemResponse{
		Cart: mapCartToProto(cart),
	}, nil
}

func (s *CartServer) UpdateItemQuantity(ctx context.Context, req *cartv1.UpdateItemQuantityRequest) (*cartv1.UpdateItemQuantityResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	cart, err := s.uc.UpdateItemQuantity(ctx, req.GetUserId(), req.GetSku(), int(req.GetQuantity()))
	if err != nil {
		s.l.Error("grpc.UpdateItemQuantity - uc.UpdateItemQuantity: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update item quantity: %v", err)
	}

	return &cartv1.UpdateItemQuantityResponse{
		Cart: mapCartToProto(cart),
	}, nil
}

func (s *CartServer) RemoveItem(ctx context.Context, req *cartv1.RemoveItemRequest) (*cartv1.RemoveItemResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	cart, err := s.uc.RemoveItem(ctx, req.GetUserId(), req.GetSku())
	if err != nil {
		s.l.Error("grpc.RemoveItem - uc.RemoveItem: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to remove item: %v", err)
	}

	return &cartv1.RemoveItemResponse{
		Cart: mapCartToProto(cart),
	}, nil
}

func (s *CartServer) RemoveItems(ctx context.Context, req *cartv1.RemoveItemsRequest) (*cartv1.RemoveItemsResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	cart, err := s.uc.RemoveItems(ctx, req.GetUserId(), req.GetSkus())
	if err != nil {
		s.l.Error("grpc.RemoveItems - uc.RemoveItems: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to remove items: %v", err)
	}

	return &cartv1.RemoveItemsResponse{
		Cart: mapCartToProto(cart),
	}, nil
}

func (s *CartServer) ClearCart(ctx context.Context, req *cartv1.ClearCartRequest) (*cartv1.ClearCartResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if err := s.uc.ClearCart(ctx, req.GetUserId()); err != nil {
		s.l.Error("grpc.ClearCart - uc.ClearCart: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to clear cart: %v", err)
	}

	return &cartv1.ClearCartResponse{
		Success: true,
	}, nil
}

func mapCartToProto(c *domain.Cart) *cartv1.Cart {
	if c == nil {
		return nil
	}
	items := make([]*cartv1.CartItem, len(c.Items))
	for i, item := range c.Items {
		items[i] = &cartv1.CartItem{
			Sku:      item.SKU,
			Quantity: int32(item.Quantity),
		}
	}

	return &cartv1.Cart{
		UserId:    c.UserID,
		Items:     items,
		UpdatedAt: c.UpdatedAt.String(),
	}
}
