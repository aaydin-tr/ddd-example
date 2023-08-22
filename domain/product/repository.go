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

type ProductRepository interface {
	Get(code valueobject.Code) (*entity.Product, error)
	Save(product *entity.Product) error
}
