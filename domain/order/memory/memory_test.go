package memory

import (
	"testing"

	"github.com/aaydin-tr/e-commerce/domain/order"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	expectedErr error
}

func TestMemoryCreateOrder(t *testing.T) {
	mockRepo := NewOrderRepository()
	uuid := uuid.New()

	t.Run("Create order", func(t *testing.T) {
		c := testCase{
			name:        "Create order new order",
			expectedErr: nil,
		}

		err := mockRepo.Create(&entity.Order{ID: uuid})
		assert.ErrorIs(t, c.expectedErr, err)
	})

	t.Run("Create order which already exist", func(t *testing.T) {
		c := testCase{
			name:        "Create order new order",
			expectedErr: order.ErrOrderAlreadyExist,
		}

		err := mockRepo.Create(&entity.Order{ID: uuid})
		assert.ErrorIs(t, c.expectedErr, err)
	})

}
