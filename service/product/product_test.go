package product

import (
	"testing"

	"github.com/aaydin-tr/e-commerce/domain/product"
	"github.com/aaydin-tr/e-commerce/entity"
	mockProduct "github.com/aaydin-tr/e-commerce/mock/repository/product"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var mockProductRepo *mockProduct.MockProductRepository

func setup(t *testing.T) (*ProductService, func()) {
	ct := gomock.NewController(t)

	mockProductRepo = mockProduct.NewMockProductRepository(ct)

	productService := NewProductService(mockProductRepo)

	return productService, func() {
		ct.Finish()
		mockProductRepo = nil
	}
}

func TestNewProductService(t *testing.T) {
	ct := gomock.NewController(t)

	mockProductRepo = mockProduct.NewMockProductRepository(ct)

	productService := NewProductService(mockProductRepo)

	assert.Equal(t, productService.productRepository, mockProductRepo)

	ct.Finish()
}

func TestProductServiceCreate(t *testing.T) {
	productService, teardown := setup(t)
	defer teardown()

	t.Run("should return error when product code is invalid", func(t *testing.T) {
		err := productService.Create("", 100, 100)
		assert.ErrorIs(t, err, valueobject.ErrCodeIsRequired)
	})

	t.Run("should return error when product price is invalid", func(t *testing.T) {
		err := productService.Create("P1", -1, 100)
		assert.ErrorIs(t, err, valueobject.ErrPriceMustBePositive)
	})

	t.Run("should return error when product stock is invalid", func(t *testing.T) {
		err := productService.Create("P1", 100, -1)
		assert.ErrorIs(t, err, valueobject.ErrStockMustBePositive)
	})

	t.Run("should return error when product already exist", func(t *testing.T) {
		mockProductRepo.EXPECT().Create(gomock.Any()).Return(product.ErrAlreadyExist)
		err := productService.Create("P1", 100, 100)
		assert.ErrorIs(t, err, product.ErrAlreadyExist)
	})

	t.Run("success", func(t *testing.T) {
		mockProductRepo.EXPECT().Create(gomock.Any()).Return(nil)
		err := productService.Create("P1", 100, 100)
		assert.Nil(t, err)
	})
}

func TestProductServiceGet(t *testing.T) {
	productService, teardown := setup(t)
	defer teardown()

	t.Run("should return error when product code is invalid", func(t *testing.T) {
		_, err := productService.Get("")
		assert.ErrorIs(t, err, valueobject.ErrCodeIsRequired)
	})

	t.Run("should return error when product is not found", func(t *testing.T) {
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(nil, product.ErrNotFound)
		_, err := productService.Get("P1")
		assert.ErrorIs(t, err, product.ErrNotFound)
	})

	t.Run("success", func(t *testing.T) {
		code, _ := valueobject.NewCode("P1")
		mockProductData := &entity.Product{Code: code}
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(mockProductData, nil)
		p, err := productService.Get(mockProductData.Code.Value())
		assert.Nil(t, err)
		assert.Equal(t, p.Code.Value(), mockProductData.Code.Value())
	})
}
