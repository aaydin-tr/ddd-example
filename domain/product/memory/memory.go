package memory

import (
	"sync"

	"github.com/aaydin-tr/e-commerce/domain/product"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
)

type ProductRepository struct {
	products map[string]*entity.Product
	sync.RWMutex
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{products: make(map[string]*entity.Product)}
}

func (r *ProductRepository) Get(code valueobject.Code) (*entity.Product, error) {
	r.RLock()
	defer r.RUnlock()

	if r.products[code.Value()] == nil {
		return nil, product.ErrNotFound
	}

	return r.products[code.Value()], nil
}

func (r *ProductRepository) Create(newProduct *entity.Product) error {
	r.Lock()
	defer r.Unlock()

	if r.products[newProduct.Code.Value()] != nil {
		return product.ErrAlreadyExist
	}

	r.products[newProduct.Code.Value()] = newProduct
	return nil
}
