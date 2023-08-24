package order

import (
	"testing"

	"github.com/aaydin-tr/e-commerce/domain/product"
	"github.com/aaydin-tr/e-commerce/entity"
	mockOrder "github.com/aaydin-tr/e-commerce/mock/repository/order"
	mockProduct "github.com/aaydin-tr/e-commerce/mock/repository/product"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var mockOrderRepo *mockOrder.MockOrderRepository
var mockProductRepo *mockProduct.MockProductRepository

func setup(t *testing.T) (*OrderService, func()) {
	ct := gomock.NewController(t)

	mockOrderRepo = mockOrder.NewMockOrderRepository(ct)
	mockProductRepo = mockProduct.NewMockProductRepository(ct)

	orderService := NewOrderService(mockProductRepo, mockOrderRepo)

	return orderService, func() {
		ct.Finish()
		mockOrderRepo = nil
		mockProductRepo = nil
	}
}

func TestNewOrderService(t *testing.T) {
	ct := gomock.NewController(t)

	mockOrderRepo = mockOrder.NewMockOrderRepository(ct)
	mockProductRepo = mockProduct.NewMockProductRepository(ct)

	orderService := NewOrderService(mockProductRepo, mockOrderRepo)

	assert.Equal(t, orderService.orderRepository, mockOrderRepo)
	assert.Equal(t, orderService.productRepository, mockProductRepo)

	ct.Finish()
}

func TestOrderService_Create(t *testing.T) {
	orderService, teardown := setup(t)
	defer teardown()

	t.Run("should return error when product code is invalid", func(t *testing.T) {
		err := orderService.Create("", 1)
		assert.Error(t, err)
	})

	t.Run("should return error when order quantity is invalid", func(t *testing.T) {
		err := orderService.Create("P1", 0)
		assert.Error(t, err)
	})

	t.Run("should return error when product is not found", func(t *testing.T) {
		code, _ := valueobject.NewCode("P1")
		mockProductRepo.EXPECT().Get(code).Return(nil, product.ErrNotFound)

		err := orderService.Create("P1", 1)
		assert.ErrorIs(t, err, product.ErrNotFound)
	})

	t.Run("should return error when product stock is insufficient", func(t *testing.T) {
		code, _ := valueobject.NewCode("P1")
		stock, _ := valueobject.NewStock(1)
		mockProductRepo.EXPECT().Get(code).Return(&entity.Product{Stock: stock}, nil)

		err := orderService.Create("P1", 2)
		assert.ErrorIs(t, err, ErrInsufficientStock)
	})

	t.Run("should return error when order repository returns error", func(t *testing.T) {
		code, _ := valueobject.NewCode("P1")
		stock, _ := valueobject.NewStock(1)
		mockProductRepo.EXPECT().Get(code).Return(&entity.Product{Stock: stock}, nil)
		mockOrderRepo.EXPECT().Create(gomock.Any()).Return(product.ErrNotFound)

		err := orderService.Create("P1", 1)
		assert.ErrorIs(t, err, product.ErrNotFound)
	})

	t.Run("success without campaign", func(t *testing.T) {
		code, _ := valueobject.NewCode("P1")
		stock, _ := valueobject.NewStock(1)
		mockProductRepo.EXPECT().Get(code).Return(&entity.Product{Stock: stock}, nil)
		mockOrderRepo.EXPECT().Create(gomock.Any()).Return(nil)

		err := orderService.Create("P1", 1)
		assert.NoError(t, err)
	})

	t.Run("success with campaign", func(t *testing.T) {
		code, _ := valueobject.NewCode("P1")
		stock, _ := valueobject.NewStock(10)
		campaignName, _ := valueobject.NewName("C1")
		targetSalesCount, _ := valueobject.NewTargetSalesCount(10)
		status, _ := valueobject.NewStatus(valueobject.Active)
		price, _ := valueobject.NewPrice(10)

		campaign := &entity.Campaign{Name: campaignName, ProductCode: code, TargetSalesCount: targetSalesCount, Status: status}
		product := &entity.Product{Stock: stock, Campaign: campaign, Price: price}
		mockProductRepo.EXPECT().Get(code).Return(product, nil)
		mockOrderRepo.EXPECT().Create(gomock.Any()).Return(nil)

		err := orderService.Create("P1", 1)
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

		campaign := &entity.Campaign{Name: campaignName, ProductCode: code, TargetSalesCount: targetSalesCount, Status: status}
		product := &entity.Product{Stock: stock, Campaign: campaign, Price: price}
		mockProductRepo.EXPECT().Get(code).Return(product, nil)
		mockOrderRepo.EXPECT().Create(gomock.Any()).Return(nil)

		err := orderService.Create("P1", 15)
		assert.NoError(t, err)

		assert.Equal(t, 10, campaign.TotalSales.Value())
		assert.Equal(t, float64(10), campaign.AverageItemPrice.Value())
		assert.Equal(t, 0, product.Stock.Value())
		assert.Equal(t, 15, product.TotalDemandCount.Value())
		assert.Equal(t, valueobject.Ended, campaign.Status.Value())
	})
}
