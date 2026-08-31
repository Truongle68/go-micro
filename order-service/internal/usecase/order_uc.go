package usecase

import (
	"context"
	"fmt"
	"strings"

	"order-service/internal/client"
	"order-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/pagination"
)

type OrderUC struct {
	repo          OrderRepository
	cartClient    client.CartClient
	catalogClient client.CatalogClient
	transactor    Transactor
	logger        logger.Interface
}

func NewOrderUC(repo OrderRepository, cartClient client.CartClient, catalogClient client.CatalogClient, transactor Transactor, l logger.Interface) *OrderUC {
	return &OrderUC{
		repo:          repo,
		cartClient:    cartClient,
		catalogClient: catalogClient,
		transactor:    transactor,
		logger:        l,
	}
}

type rawLine struct {
	SKU      string
	Quantity int
}

func (uc *OrderUC) Checkout(ctx context.Context, userID string, input CheckoutInput, token string) (*domain.Order, error) {
	raws, err := uc.resolveCheckoutLines(ctx, userID, input, token)
	if err != nil {
		return nil, err
	}

	domainItems, err := uc.buildOrderItemsFromCatalog(ctx, raws)
	if err != nil {
		return nil, err
	}

	order, history, err := domain.NewOrder(userID, domainItems, input.ShippingAddress, input.ShippingFee, input.PaymentMethod)
	if err != nil {
		return nil, fmt.Errorf("OrderUC.Checkout - domain.NewOrder: %w", err)
	}

	err = uc.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		err := uc.repo.Create(txCtx, order, history)
		if err != nil {
			return fmt.Errorf("OrderUC.Checkout - Create: %w", err)
		}

		// COD Payment path (stub): transition pending_payment -> confirmed
		confirmHist, err := order.MarkConfirmed("COD-" + order.ID)
		if err != nil {
			return fmt.Errorf("OrderUC.Checkout - MarkConfirmed: %w", err)
		}

		if err := uc.repo.UpdateStatus(txCtx, order, confirmHist); err != nil {
			return fmt.Errorf("OrderUC.Checkout - repo.UpdateStatus: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("OrderUC.Checkout - CreateWithTransaction: %w", err)
	}
	// Best-effort remove checked out items from cart
	if uc.cartClient != nil && len(domainItems) > 0 {
		skus := make([]string, len(domainItems))
		for i, item := range domainItems {
			skus[i] = item.SKU
		}
		if err := uc.cartClient.RemoveItems(ctx, userID, skus, token); err != nil {
			uc.logger.Warn("checkout cart cleanup: %v", err)
		}
	}
	return order, nil
}

func (uc *OrderUC) resolveCheckoutLines(
	ctx context.Context,
	userID string,
	input CheckoutInput,
	token string,
) ([]rawLine, error) {
	items := input.Items
	if len(items) <= 0 {
		return nil, domain.ErrEmptyOrderItems
	}

	if uc.cartClient != nil {
		cart, err := uc.cartClient.GetCart(ctx, userID, token)
		if err != nil {
			return nil, fmt.Errorf("OrderUC.Checkout - cartClient.GetCart: %w", err)
		}

		cartItems := make(map[string]int, len(cart.Items))
		for _, item := range cart.Items {
			cartItems[item.SKU] = item.Quantity
		}

		for _, item := range items {
			cartQty, ok := cartItems[item.SKU]
			if !ok {
				return nil, domain.ErrCartItemNotFound
			}
			if item.Quantity > cartQty {
				return nil, domain.ErrCartItemQtyExceeded
			}
		}
	}

	return parseRawLine(
		items,
		func(item CheckoutItemInput) string { return item.SKU },
		func(item CheckoutItemInput) int { return item.Quantity },
	)
}

func (uc *OrderUC) buildOrderItemsFromCatalog(ctx context.Context, raws []rawLine) ([]domain.OrderItem, error) {
	if len(raws) == 0 {
		return nil, domain.ErrEmptyOrderItems
	}

	qtyBySKU := make(map[string]int, len(raws))
	skus := make([]string, 0, len(raws))
	for _, r := range raws {
		if _, ok := qtyBySKU[r.SKU]; !ok {
			skus = append(skus, r.SKU)
		}
		qtyBySKU[r.SKU] += r.Quantity
	}

	variants, err := uc.catalogClient.GetVariantsBySKUs(ctx, skus)
	if err != nil {
		return nil, fmt.Errorf("buildOrderItemsFromCatalog - GetVariantsBySKUs: %w", err)
	}
	bySKU := make(map[string]client.VariantDTO, len(variants))
	for _, v := range variants {
		bySKU[v.SKU] = v
	}

	domainItems := make([]domain.OrderItem, 0, len(skus))
	for _, sku := range skus {
		v, ok := bySKU[sku]
		if !ok {
			return nil, fmt.Errorf("%w: %s", domain.ErrEmptySKU, sku)
		}
		if !v.IsActive {
			return nil, fmt.Errorf("%w: %s", domain.ErrInactiveVariant, sku)
		}

		domainItems = append(domainItems, domain.OrderItem{
			ProductID:    v.ProductID,
			VariantID:    v.ID,
			SKU:          sku,
			ProductName:  v.ProductName,
			Image:        v.Image,
			VariantAttrs: v.Attributes,
			UnitPrice:    int64(v.Price.Amount),
			Quantity:     qtyBySKU[sku],
		})
	}
	return domainItems, nil
}

func parseRawLine[T any](items []T, getSKU func(T) string, getQuantity func(T) int) ([]rawLine, error) {
	raw := make([]rawLine, 0, len(items))
	for _, item := range items {
		sku := strings.TrimSpace(getSKU(item))
		if sku == "" {
			return nil, domain.ErrEmptySKU
		}
		qty := getQuantity(item)
		if qty <= 0 {
			return nil, domain.ErrInvalidQuantity
		}
		raw = append(raw, rawLine{SKU: sku, Quantity: qty})
	}
	return raw, nil
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
