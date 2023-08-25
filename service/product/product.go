package product

import (
	"github.com/aaydin-tr/e-commerce/domain/product"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/google/uuid"
)

type ProductServiceInterface interface {
	Create(productCode string, productPrice float64, productStock int) error
	Get(productCode string) (*entity.Product, error)
}

type ProductService struct {
	productRepository product.ProductRepository
}

func NewProductService(productRepository product.ProductRepository) *ProductService {
	return &ProductService{productRepository: productRepository}
}

func (s *ProductService) Create(productCode string, productPrice float64, productStock int) error {
	code, err := valueobject.NewCode(productCode)
	if err != nil {
		return err
	}

	price, err := valueobject.NewPrice(productPrice)
	if err != nil {
		return err
	}

	stock, err := valueobject.NewStock(productStock)
	if err != nil {
		return err
	}

	product := &entity.Product{
		ID:            uuid.New(),
		Code:          code,
		Price:         price,
		Stock:         stock,
		InititalStock: stock,
		InititalPrice: price,
	}

	return s.productRepository.Create(product)
}

func (s *ProductService) Get(productCode string) (*entity.Product, error) {
	code, err := valueobject.NewCode(productCode)
	if err != nil {
		return nil, err
	}

	result, err := s.productRepository.Get(code)
	if err != nil {
		return nil, err
	}

	return result, nil
}
