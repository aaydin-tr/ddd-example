package memory

import (
	"sync"

	"github.com/aaydin-tr/e-commerce/domain/order"
	"github.com/aaydin-tr/e-commerce/entity"
)

type OrderRepository struct {
	orders map[string]*entity.Order
	sync.RWMutex
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{orders: make(map[string]*entity.Order)}
}

func (r *OrderRepository) Create(newOrder *entity.Order) error {
	r.Lock()
	defer r.Unlock()

	if r.orders[newOrder.ID.String()] != nil {
		return order.ErrOrderAlreadyExist
	}

	r.orders[newOrder.ID.String()] = newOrder
	return nil
}
