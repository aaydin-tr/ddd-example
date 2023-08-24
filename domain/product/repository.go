package product

import (
	"errors"

	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
)

var (
	ErrNotFound     = errors.New("Product not found")
	ErrAlreadyExist = errors.New("Product already exist")
)

//go:generate mockgen -destination=../../mock/repository/product/product.go -package=repository github.com/aaydin-tr/e-commerce/domain/product ProductRepository
type ProductRepository interface {
	Get(code valueobject.Code) (*entity.Product, error)
	Create(product *entity.Product) error
}
