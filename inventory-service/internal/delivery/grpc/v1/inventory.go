package v1

import (
	"context"

	"inventory-service/internal/domain"
	"inventory-service/internal/usecase"

	inventoryv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/inventory/v1"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StockUseCase defines the methods the gRPC handler needs from StockUC.
type StockUseCase interface {
	CheckStock(ctx context.Context, skus []string) (map[string]int, error)
	ReserveStock(ctx context.Context, orderID string, items []usecase.SKUQty) error
	ConfirmReservation(ctx context.Context, orderID string) error
	ReleaseReservation(ctx context.Context, orderID string) error
}

var _ StockUseCase = (*usecase.StockUC)(nil)

// InventoryServer implements the gRPC InventoryServiceServer.
type InventoryServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	uc StockUseCase
	l  logger.Interface
}

// NewInventoryServer creates a new InventoryServer.
func NewInventoryServer(uc StockUseCase, l logger.Interface) *InventoryServer {
	return &InventoryServer{
		uc: uc,
		l:  l,
	}
}

// CheckStock returns available quantities for the requested SKUs.
func (s *InventoryServer) CheckStock(ctx context.Context, req *inventoryv1.CheckStockRequest) (*inventoryv1.CheckStockResponse, error) {
	if len(req.GetSkus()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "skus cannot be empty")
	}

	avail, err := s.uc.CheckStock(ctx, req.GetSkus())
	if err != nil {
		s.l.Error("grpc.CheckStock: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to check stock: %v", err)
	}

	items := make([]*inventoryv1.StockAvailability, 0, len(avail))
	for sku, qty := range avail {
		items = append(items, &inventoryv1.StockAvailability{
			Sku:       sku,
			Available: int32(qty),
		})
	}

	return &inventoryv1.CheckStockResponse{Items: items}, nil
}

// ReserveStock reserves inventory for an order.
func (s *InventoryServer) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	if req.GetOrderId() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if len(req.GetItems()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "items cannot be empty")
	}

	items := make([]usecase.SKUQty, len(req.GetItems()))
	for i, item := range req.GetItems() {
		items[i] = usecase.SKUQty{
			SKU:      item.GetSku(),
			Quantity: int(item.GetQuantity()),
		}
	}

	err := s.uc.ReserveStock(ctx, req.GetOrderId(), items)
	if err != nil {
		s.l.Error("grpc.ReserveStock order=%s: %v", req.GetOrderId(), err)
		return nil, mapDomainError(err)
	}

	return &inventoryv1.ReserveStockResponse{Success: true}, nil
}

// ConfirmReservation confirms all pending reservations for an order.
func (s *InventoryServer) ConfirmReservation(ctx context.Context, req *inventoryv1.ConfirmReservationRequest) (*inventoryv1.ConfirmReservationResponse, error) {
	if req.GetOrderId() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	err := s.uc.ConfirmReservation(ctx, req.GetOrderId())
	if err != nil {
		s.l.Error("grpc.ConfirmReservation order=%s: %v", req.GetOrderId(), err)
		return nil, status.Errorf(codes.Internal, "failed to confirm reservation: %v", err)
	}

	return &inventoryv1.ConfirmReservationResponse{Success: true}, nil
}

// ReleaseReservation releases all pending reservations for an order.
func (s *InventoryServer) ReleaseReservation(ctx context.Context, req *inventoryv1.ReleaseReservationRequest) (*inventoryv1.ReleaseReservationResponse, error) {
	if req.GetOrderId() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	err := s.uc.ReleaseReservation(ctx, req.GetOrderId())
	if err != nil {
		s.l.Error("grpc.ReleaseReservation order=%s: %v", req.GetOrderId(), err)
		return nil, status.Errorf(codes.Internal, "failed to release reservation: %v", err)
	}

	return &inventoryv1.ReleaseReservationResponse{Success: true}, nil
}

// mapDomainError converts domain errors to appropriate gRPC status codes.
func mapDomainError(err error) error {
	if err == nil {
		return nil
	}

	appErr := domain.ToAppError(err)
	if appErr == nil {
		return status.Errorf(codes.Internal, "internal error: %v", err)
	}

	switch appErr.Code {
	case domain.CodeInsufficientStock:
		return status.Errorf(codes.FailedPrecondition, "%s", appErr.Message)
	case domain.CodeNonPositiveQuantity:
		return status.Errorf(codes.InvalidArgument, "%s", appErr.Message)
	case domain.CodeEmptySKU, domain.CodeEmptyWarehouseID:
		return status.Errorf(codes.InvalidArgument, "%s", appErr.Message)
	default:
		return status.Errorf(codes.Internal, "%s", appErr.Message)
	}
}
