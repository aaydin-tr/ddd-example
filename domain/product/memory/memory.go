package memory

import (
	"github.com/aaydin-tr/e-commerce/domain/product"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/types"
	"github.com/aaydin-tr/e-commerce/valueobject"
)

type ProductRepository struct {
	storage types.Storage[*entity.Product]
}

func NewProductRepository(storage types.Storage[*entity.Product]) *ProductRepository {
	return &ProductRepository{storage: storage}
}

func (r *ProductRepository) Get(code valueobject.Code) (*entity.Product, error) {
	result, ok := r.storage.Get(code.Value())
	if !ok {
		return nil, product.ErrNotFound
	}

	return result, nil
}

func (r *ProductRepository) Create(newProduct *entity.Product) error {
	_, ok := r.storage.Get(newProduct.Code.Value())
	if ok {
		return product.ErrAlreadyExist
	}

	r.storage.Set(newProduct.Code.Value(), newProduct)
	return nil
}
