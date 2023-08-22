package order

import (
	"errors"

	"github.com/aaydin-tr/e-commerce/entity"
)

var (
	ErrOrderAlreadyExist = errors.New("Order already exist")
)

type OrderRepository interface {
	Create(order *entity.Order) error
}
