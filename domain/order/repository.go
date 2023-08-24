package order

import (
	"errors"

	"github.com/aaydin-tr/e-commerce/entity"
)

var (
	ErrOrderAlreadyExist = errors.New("Order already exist")
)

//go:generate mockgen -destination=../../mock/repository/order/order.go -package=repository github.com/aaydin-tr/e-commerce/domain/order OrderRepository
type OrderRepository interface {
	Create(order *entity.Order) error
}
