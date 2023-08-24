package campaign

import (
	"errors"
	"testing"

	"github.com/aaydin-tr/e-commerce/domain/campaign"
	"github.com/aaydin-tr/e-commerce/domain/product"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/google/uuid"

	mockCampaign "github.com/aaydin-tr/e-commerce/mock/repository/campaign"
	mockOrder "github.com/aaydin-tr/e-commerce/mock/repository/order"
	mockProduct "github.com/aaydin-tr/e-commerce/mock/repository/product"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var mockCampaignRepo *mockCampaign.MockCampaignRepository
var mockOrderRepo *mockOrder.MockOrderRepository
var mockProductRepo *mockProduct.MockProductRepository

func setup(t *testing.T) (*CampaignService, func()) {
	ct := gomock.NewController(t)

	mockCampaignRepo = mockCampaign.NewMockCampaignRepository(ct)
	mockOrderRepo = mockOrder.NewMockOrderRepository(ct)
	mockProductRepo = mockProduct.NewMockProductRepository(ct)

	campaignService := NewCampaignService(mockCampaignRepo, mockOrderRepo, mockProductRepo)

	return campaignService, func() {
		ct.Finish()
		mockCampaignRepo = nil
		mockOrderRepo = nil
		mockProductRepo = nil
	}
}

func TestNewCampaignService(t *testing.T) {
	ct := gomock.NewController(t)

	mockCampaignRepo = mockCampaign.NewMockCampaignRepository(ct)
	mockOrderRepo = mockOrder.NewMockOrderRepository(ct)
	mockProductRepo = mockProduct.NewMockProductRepository(ct)

	campaignService := NewCampaignService(mockCampaignRepo, mockOrderRepo, mockProductRepo)

	assert.Equal(t, campaignService.campaignRepository, mockCampaignRepo)
	assert.Equal(t, campaignService.orderRepository, mockOrderRepo)
	assert.Equal(t, campaignService.productRepository, mockProductRepo)

	ct.Finish()
}

func TestCampaignServiceCreateCampaign(t *testing.T) {
	campaignService, teardown := setup(t)
	defer teardown()

	t.Run("should return error when campaign name is invalid", func(t *testing.T) {
		err := campaignService.CreateCampaign("", "P1", 10, 20, 100)
		assert.NotNil(t, err)
	})

	t.Run("should return error when campaign name is already exist", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(true)

		err := campaignService.CreateCampaign("C1", "P1", 10, 20, 100)
		assert.ErrorIs(t, err, campaign.ErrCampaignAlreadyExist)
	})

	t.Run("should return error when product code is invalid", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)

		err := campaignService.CreateCampaign("C1", "", 10, 20, 100)
		assert.NotNil(t, err)
	})

	t.Run("should return error when product is not found", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(nil, product.ErrNotFound)

		err := campaignService.CreateCampaign("C1", "P1", 10, 20, 100)
		assert.ErrorIs(t, err, product.ErrNotFound)
	})

	t.Run("should return error when campaign duration is invalid", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(nil, nil)

		err := campaignService.CreateCampaign("C1", "P1", 0, 20, 100)
		assert.NotNil(t, err)
	})

	t.Run("should return error when campaign price manipulation limit is invalid", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(nil, nil)

		err := campaignService.CreateCampaign("C1", "P1", 10, 0, 100)
		assert.NotNil(t, err)
	})

	t.Run("should return error when campaign target sales count is invalid", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(nil, nil)

		err := campaignService.CreateCampaign("C1", "P1", 10, 20, 0)
		assert.NotNil(t, err)
	})

	t.Run("should return error when campaign target sales count is greater than product stock", func(t *testing.T) {
		stock, _ := valueobject.NewStock(100)
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(&entity.Product{Stock: stock}, nil)

		err := campaignService.CreateCampaign("C1", "P1", 10, 20, 200)
		assert.ErrorIs(t, err, ErrTargetSalesCountMustBeLessThanStock)
	})

	t.Run("should return error when campaign repo create returns error", func(t *testing.T) {
		returnErr := errors.New("error")
		stock, _ := valueobject.NewStock(100)
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(&entity.Product{Stock: stock}, nil)
		mockCampaignRepo.EXPECT().Create(gomock.Any()).Return(returnErr)

		err := campaignService.CreateCampaign("C1", "P1", 10, 20, 50)
		assert.ErrorIs(t, err, returnErr)

	})

	t.Run("success", func(t *testing.T) {
		stock, _ := valueobject.NewStock(100)
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(&entity.Product{Stock: stock}, nil)
		mockCampaignRepo.EXPECT().Create(gomock.Any()).Return(nil)

		err := campaignService.CreateCampaign("C1", "P1", 10, 20, 50)
		assert.Nil(t, err)
	})

}

func TestCampaignServiceGetCampaignInfo(t *testing.T) {
	campaignService, teardown := setup(t)
	defer teardown()

	t.Run("should return error when campaign name is invalid", func(t *testing.T) {
		c, err := campaignService.GetCampaignInfo("")
		assert.NotNil(t, err)
		assert.Nil(t, c)
	})

	t.Run("should return error when campaign is not found", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Get(gomock.Any()).Return(nil, campaign.ErrCampaignNotFound)

		c, err := campaignService.GetCampaignInfo("C1")
		assert.ErrorIs(t, err, campaign.ErrCampaignNotFound)
		assert.Nil(t, c)
	})

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		name, _ := valueobject.NewName("C1")
		code, _ := valueobject.NewCode("P1")
		duration, _ := valueobject.NewDuration(10)
		priceManipulationLimit, _ := valueobject.NewPriceManipulationLimit(20)
		targetSalesCount, _ := valueobject.NewTargetSalesCount(100)

		mockCampaignData := &entity.Campaign{
			ID:                     id,
			Name:                   name,
			ProductCode:            code,
			Duration:               duration,
			PriceManipulationLimit: priceManipulationLimit,
			TargetSalesCount:       targetSalesCount,
		}

		mockCampaignRepo.EXPECT().Get(gomock.Any()).Return(mockCampaignData, nil)

		c, err := campaignService.GetCampaignInfo("C1")
		assert.Nil(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, c.ID, id)
		assert.Equal(t, c.Name, name)
		assert.Equal(t, c.ProductCode, code)
		assert.Equal(t, c.Duration, duration)
		assert.Equal(t, c.PriceManipulationLimit, priceManipulationLimit)
		assert.Equal(t, c.TargetSalesCount, targetSalesCount)
	})
}

func TestCampaignServiceIncreaseTime(t *testing.T) {
	campaignService, teardown := setup(t)
	defer teardown()

	t.Run("should return No campaign found error when campaign is not found", func(t *testing.T) {
		mockCampaignRepo.EXPECT().GetAll().Return(nil)

		err := campaignService.IncreaseTime(1)
		assert.ErrorIs(t, err, ErrNoCampaign)
	})

	t.Run("should return error when campaign does not have valid product", func(t *testing.T) {
		code, _ := valueobject.NewCode("P1")
		mockCampaignRepo.EXPECT().GetAll().Return([]*entity.Campaign{{ProductCode: code}})
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(nil, product.ErrNotFound)

		err := campaignService.IncreaseTime(1)
		assert.ErrorIs(t, err, ErrCampaignDoesNotHaveProduct)
	})

	t.Run("should close campaign when campaign duration is over", func(t *testing.T) {
		id := uuid.New()
		name, _ := valueobject.NewName("C1")
		pcode, _ := valueobject.NewCode("P1")

		mockCampaignData := &entity.Campaign{
			ID:          id,
			Name:        name,
			ProductCode: pcode,
		}
		mockProductData := &entity.Product{Code: pcode, Campaign: mockCampaignData}

		mockCampaignRepo.EXPECT().GetAll().Return([]*entity.Campaign{mockCampaignData})
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(mockProductData, nil)

		err := campaignService.IncreaseTime(1)
		assert.Nil(t, err)
		assert.Equal(t, mockCampaignData.Status.Value(), valueobject.Ended)
		assert.Nil(t, mockProductData.Campaign)

	})

	t.Run("success", func(t *testing.T) {
		code, _ := valueobject.NewCode("P1")
		stock, _ := valueobject.NewStock(100)
		price, _ := valueobject.NewPrice(100)
		inititalStock, _ := valueobject.NewStock(100)
		inititalPrice, _ := valueobject.NewPrice(100)
		totalDemandCount, _ := valueobject.NewDemand(0)
		duration, _ := valueobject.NewDuration(10)

		mockCampaignRepo.EXPECT().GetAll().Return([]*entity.Campaign{{ProductCode: code, Duration: duration}})
		mockProductRepo.EXPECT().Get(gomock.Any()).Return(&entity.Product{
			Stock:            stock,
			Price:            price,
			InititalStock:    inititalStock,
			InititalPrice:    inititalPrice,
			TotalDemandCount: totalDemandCount,
		}, nil)

		err := campaignService.IncreaseTime(1)
		assert.Nil(t, err)
	})

}
