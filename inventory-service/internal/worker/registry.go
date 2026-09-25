package worker

import (
	"context"
	"inventory-service/internal/domain"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscription struct {
	Name     string
	Exchange string
	Queue    string
	Handler  rabbitmq.HandlerFunc

	routingKey string
	fanout     bool
}

func NewSubscription(name, exchange, queue, routingKey string, handler rabbitmq.HandlerFunc) Subscription {
	return Subscription{Name: name, Exchange: exchange, Queue: queue, Handler: handler, routingKey: routingKey}
}

func NewFanoutSubscription(name, exchange, queue string, handler rabbitmq.HandlerFunc) Subscription {
	return Subscription{Name: name, Exchange: exchange, Queue: queue, Handler: handler, fanout: true}
}

// Register opens this subscription's consumer on conn, dispatching to the
// matching rabbitmq constructor.
func (s Subscription) Register(conn *amqp.Connection) (*rabbitmq.Consumer, error) {
	if s.fanout {
		return rabbitmq.NewFanoutConsumer(conn, s.Exchange, s.Queue, s.Handler)
	}
	return rabbitmq.NewConsumer(conn, s.Exchange, s.Queue, s.routingKey, s.Handler)
}

type stockUseCase interface {
	ApplyGoodsReceived(ctx context.Context, goods domain.GoodsReceivedEvent) error
}

func BuildSubscriptions(stockUC stockUseCase, l logger.Interface) []Subscription {
	goodsReceived := NewGoodsReceivedHandler(stockUC, l)
	orderPlaced := NewOrderPlacedHandler(stockUC, l)

	return []Subscription{
		NewSubscription("goods-received", rabbitmq.ExchangeInventory, "inventory.goods-received.q", domain.EventGoodsReceived, goodsReceived.Handle),
		NewFanoutSubscription("order-placed", rabbitmq.ExchangeOrderPlaced, "inventory.order-placed.q", orderPlaced.Handle),
	}
}
