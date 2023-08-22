package memory

import (
	"github.com/aaydin-tr/e-commerce/domain/order"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/types"
)

type OrderRepository struct {
	storage types.Storage[*entity.Order]
}

func NewOrderRepository(storage types.Storage[*entity.Order]) *OrderRepository {
	return &OrderRepository{storage: storage}
}

func (r *OrderRepository) Create(newOrder *entity.Order) error {
	_, ok := r.storage.Get(newOrder.ID.String())
	if ok {
		return order.ErrOrderAlreadyExist
	}

	r.storage.Set(newOrder.ID.String(), newOrder)
	return nil
}
