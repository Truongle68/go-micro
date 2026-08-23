package usecase

import (
	"context"
	"fmt"

	"order-service/internal/client"
	"order-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type OrderUC struct {
	repo       OrderRepository
	cartClient client.CartClient
	logger     logger.Interface
}

func NewOrderUC(repo OrderRepository, cartClient client.CartClient, l logger.Interface) *OrderUC {
	return &OrderUC{
		repo:       repo,
		cartClient: cartClient,
		logger:     l,
	}
}

func (uc *OrderUC) Checkout(ctx context.Context, userID string, input CheckoutInput, token string) (*domain.Order, error) {
	var domainItems []domain.OrderItem

	if len(input.Items) > 0 {
		domainItems = make([]domain.OrderItem, len(input.Items))
		for i, item := range input.Items {
			productName := item.ProductName
			if productName == "" {
				productName = item.SKU
			}
			domainItems[i] = domain.OrderItem{
				ProductID:    item.ProductID,
				VariantID:    item.VariantID,
				SKU:          item.SKU,
				ProductName:  productName,
				Image:        item.Image,
				VariantAttrs: item.VariantAttrs,
				UnitPrice:    item.UnitPrice,
				Quantity:     item.Quantity,
			}
		}
	} else if uc.cartClient != nil {
		cart, err := uc.cartClient.GetCart(ctx, userID, token)
		if err != nil {
			return nil, fmt.Errorf("OrderUC.Checkout - cartClient.GetCart: %w", err)
		}
		if len(cart.Items) == 0 {
			return nil, domain.ErrCartEmpty
		}

		domainItems = make([]domain.OrderItem, len(cart.Items))
		for i, item := range cart.Items {
			domainItems[i] = domain.OrderItem{
				SKU:         item.SKU,
				ProductName: item.SKU,
				UnitPrice:   0,
				Quantity:    item.Quantity,
			}
		}
	} else {
		return nil, domain.ErrEmptyOrderItems
	}

	order, history, err := domain.NewOrder(userID, domainItems, input.ShippingAddress, input.ShippingFee, input.PaymentMethod)
	if err != nil {
		return nil, fmt.Errorf("OrderUC.Checkout - domain.NewOrder: %w", err)
	}

	if err := uc.repo.Create(ctx, order, history); err != nil {
		return nil, fmt.Errorf("OrderUC.Checkout - repo.Create: %w", err)
	}

	// COD Payment path (stub): transition pending_payment -> confirmed
	confirmHist, err := order.MarkConfirmed("COD-" + order.ID)
	if err != nil {
		return nil, fmt.Errorf("OrderUC.Checkout - MarkConfirmed: %w", err)
	}

	if err := uc.repo.UpdateStatus(ctx, order, confirmHist); err != nil {
		return nil, fmt.Errorf("OrderUC.Checkout - repo.UpdateStatus: %w", err)
	}

	// Best-effort clear cart
	if uc.cartClient != nil {
		go func() {
			_ = uc.cartClient.ClearCart(context.Background(), userID, token)
		}()
	}

	return order, nil
}

func (uc *OrderUC) GetOrder(ctx context.Context, orderID string, userID string) (*domain.Order, error) {
	order, err := uc.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if userID != "" && order.UserID != userID {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}

func (uc *OrderUC) ListOrdersByUser(ctx context.Context, userID string, page pagination.Params) (*pagination.Result[domain.Order], error) {
	normParams := page.Normalize()

	orders, total, err := uc.repo.FindByUserID(ctx, userID, normParams.Limit, normParams.Skip())
	if err != nil {
		return nil, fmt.Errorf("OrderUC.ListOrdersByUser: %w", err)
	}

	result := pagination.NewResult(orders, normParams, total)
	return &result, nil
}

func (uc *OrderUC) GetTrackingTimeline(ctx context.Context, orderID string, userID string) ([]domain.OrderStatusHistory, error) {
	order, err := uc.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if userID != "" && order.UserID != userID {
		return nil, domain.ErrOrderNotFound
	}

	return uc.repo.GetTrackingHistory(ctx, orderID)
}

func (uc *OrderUC) ShipOrder(ctx context.Context, orderID string, trackingCode string) (*domain.Order, error) {
	order, err := uc.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	shipHist, err := order.Ship(trackingCode)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.UpdateStatus(ctx, order, shipHist); err != nil {
		return nil, fmt.Errorf("OrderUC.ShipOrder - repo.UpdateStatus: %w", err)
	}

	return order, nil
}

func (uc *OrderUC) DeliverOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := uc.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	delivHist, err := order.Deliver()
	if err != nil {
		return nil, err
	}

	if err := uc.repo.UpdateStatus(ctx, order, delivHist); err != nil {
		return nil, fmt.Errorf("OrderUC.DeliverOrder - repo.UpdateStatus: %w", err)
	}

	return order, nil
}

func (uc *OrderUC) CancelOrder(ctx context.Context, orderID string, userID string, reason string) (*domain.Order, error) {
	order, err := uc.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if userID != "" && order.UserID != userID {
		return nil, domain.ErrOrderNotFound
	}

	cancelHist, err := order.Cancel(reason)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.UpdateStatus(ctx, order, cancelHist); err != nil {
		return nil, fmt.Errorf("OrderUC.CancelOrder - repo.UpdateStatus: %w", err)
	}

	return order, nil
}
