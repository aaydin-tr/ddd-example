package memory

import (
	"testing"

	"github.com/aaydin-tr/e-commerce/domain/product"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	expectedErr error
}

func TestMemoryCreateProduct(t *testing.T) {
	mockRepo := NewProductRepository()
	code, _ := valueobject.NewCode("P1")

	t.Run("Create product", func(t *testing.T) {
		c := testCase{
			name:        "Create product new product",
			expectedErr: nil,
		}

		err := mockRepo.Create(&entity.Product{Code: code})
		assert.ErrorIs(t, c.expectedErr, err)
	})

	t.Run("Create product which already exist", func(t *testing.T) {
		c := testCase{
			name:        "Create product new product",
			expectedErr: product.ErrAlreadyExist,
		}

		err := mockRepo.Create(&entity.Product{Code: code})
		assert.ErrorIs(t, c.expectedErr, err)
	})

}

func TestMemoryGetProduct(t *testing.T) {
	mockRepo := NewProductRepository()
	code, _ := valueobject.NewCode("P1")
	mockRepo.products["P1"] = &entity.Product{Code: code}

	t.Run("Get product", func(t *testing.T) {
		name, err := valueobject.NewCode(code.Value())
		assert.NoError(t, err)

		product, err := mockRepo.Get(name)
		assert.NoError(t, err)
		assert.Equal(t, code.Value(), product.Code.Value())
	})

	t.Run("Get product which not exist", func(t *testing.T) {
		name, err := valueobject.NewCode("P2")
		assert.NoError(t, err)

		_, err = mockRepo.Get(name)
		assert.ErrorIs(t, product.ErrNotFound, err)
	})
}
