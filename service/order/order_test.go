package order

import (
	"testing"

	"github.com/aaydin-tr/e-commerce/entity"
	mockOrder "github.com/aaydin-tr/e-commerce/mock/repository/order"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var mockOrderRepo *mockOrder.MockOrderRepository

func setup(t *testing.T) (*OrderService, func()) {
	ct := gomock.NewController(t)

	mockOrderRepo = mockOrder.NewMockOrderRepository(ct)

	orderService := NewOrderService(mockOrderRepo)

	return orderService, func() {
		ct.Finish()
		mockOrderRepo = nil
	}
}

func TestNewOrderService(t *testing.T) {
	ct := gomock.NewController(t)

	mockOrderRepo = mockOrder.NewMockOrderRepository(ct)

	orderService := NewOrderService(mockOrderRepo)

	assert.Equal(t, orderService.orderRepository, mockOrderRepo)

	ct.Finish()
}

func TestOrderService_Create(t *testing.T) {
	orderService, teardown := setup(t)
	defer teardown()
	code, _ := valueobject.NewCode("P1")
	stock, _ := valueobject.NewStock(10)
	price, _ := valueobject.NewPrice(10)

	mockProduct := &entity.Product{Code: code, Stock: stock, Price: price}

	t.Run("should return error when order quantity is invalid", func(t *testing.T) {
		err := orderService.Create(mockProduct, 0)
		assert.Error(t, err)
	})

	t.Run("should return error when product stock is insufficient", func(t *testing.T) {
		err := orderService.Create(mockProduct, 20)
		assert.ErrorIs(t, err, ErrInsufficientStock)
	})

	t.Run("success without campaign", func(t *testing.T) {
		mockOrderRepo.EXPECT().Create(gomock.Any()).Return(nil)

		err := orderService.Create(mockProduct, 1)
		assert.NoError(t, err)
	})

	t.Run("success with campaign", func(t *testing.T) {
		code, _ := valueobject.NewCode("P1")
		stock, _ := valueobject.NewStock(10)
		campaignName, _ := valueobject.NewName("C1")
		targetSalesCount, _ := valueobject.NewTargetSalesCount(10)
		status, _ := valueobject.NewStatus(valueobject.Active)
		price, _ := valueobject.NewPrice(10)

		product := &entity.Product{Stock: stock, Price: price, Code: code}
		campaign := &entity.Campaign{Name: campaignName, Product: product, TargetSalesCount: targetSalesCount, Status: status}
		product.Campaign = campaign

		mockOrderRepo.EXPECT().Create(gomock.Any()).Return(nil)

		err := orderService.Create(product, 1)
		assert.NoError(t, err)

		assert.Equal(t, 1, campaign.TotalSales.Value())
		assert.Equal(t, float64(10), campaign.AverageItemPrice.Value())
		assert.Equal(t, 9, product.Stock.Value())
		assert.Equal(t, 1, product.TotalDemandCount.Value())
		assert.Equal(t, valueobject.Active, campaign.Status.Value())
	})

	t.Run("success with campaign and order more than target sales count", func(t *testing.T) {
		code, _ := valueobject.NewCode("P1")
		stock, _ := valueobject.NewStock(15)
		campaignName, _ := valueobject.NewName("C1")
		targetSalesCount, _ := valueobject.NewTargetSalesCount(10)
		status, _ := valueobject.NewStatus(valueobject.Active)
		price, _ := valueobject.NewPrice(10)

		product := &entity.Product{Stock: stock, Price: price, Code: code}
		campaign := &entity.Campaign{Name: campaignName, Product: product, TargetSalesCount: targetSalesCount, Status: status}
		product.Campaign = campaign

		mockOrderRepo.EXPECT().Create(gomock.Any()).Return(nil)

		err := orderService.Create(product, 15)
		assert.NoError(t, err)

		assert.Equal(t, 10, campaign.TotalSales.Value())
		assert.Equal(t, float64(10), campaign.AverageItemPrice.Value())
		assert.Equal(t, 0, product.Stock.Value())
		assert.Equal(t, 15, product.TotalDemandCount.Value())
		assert.Equal(t, valueobject.Ended, campaign.Status.Value())
	})
}
