package campaign

import (
	"errors"
	"testing"

	"github.com/aaydin-tr/e-commerce/domain/campaign"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/google/uuid"

	mockCampaign "github.com/aaydin-tr/e-commerce/mock/repository/campaign"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var mockCampaignRepo *mockCampaign.MockCampaignRepository

func setup(t *testing.T) (*CampaignService, func()) {
	ct := gomock.NewController(t)

	mockCampaignRepo = mockCampaign.NewMockCampaignRepository(ct)

	campaignService := NewCampaignService(mockCampaignRepo)

	return campaignService, func() {
		ct.Finish()
		mockCampaignRepo = nil
	}
}

func TestNewCampaignService(t *testing.T) {
	ct := gomock.NewController(t)

	mockCampaignRepo = mockCampaign.NewMockCampaignRepository(ct)

	campaignService := NewCampaignService(mockCampaignRepo)

	assert.Equal(t, campaignService.campaignRepository, mockCampaignRepo)

	ct.Finish()
}

func TestCampaignServiceCreateCampaign(t *testing.T) {
	campaignService, teardown := setup(t)
	defer teardown()

	code, _ := valueobject.NewCode("P1")
	stokc, _ := valueobject.NewStock(100)
	price, _ := valueobject.NewPrice(100)

	mockProduct := &entity.Product{Code: code, Stock: stokc, Price: price}

	t.Run("should return error when campaign name is invalid", func(t *testing.T) {
		err := campaignService.Create("", mockProduct, 10, 20, 100)
		assert.NotNil(t, err)
	})

	t.Run("should return error when campaign name is already exist", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(true)

		err := campaignService.Create("C1", mockProduct, 10, 20, 100)
		assert.ErrorIs(t, err, campaign.ErrCampaignAlreadyExist)
	})

	t.Run("should return error when campaign duration is invalid", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)

		err := campaignService.Create("C1", mockProduct, 0, 20, 100)
		assert.NotNil(t, err)
	})

	t.Run("should return error when campaign price manipulation limit is invalid", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)

		err := campaignService.Create("C1", mockProduct, 10, 0, 100)
		assert.NotNil(t, err)
	})

	t.Run("should return error when campaign target sales count is invalid", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)

		err := campaignService.Create("C1", mockProduct, 10, 20, 0)
		assert.NotNil(t, err)
	})

	t.Run("should return error when campaign target sales count is greater than product stock", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)

		err := campaignService.Create("C1", mockProduct, 10, 20, 200)
		assert.ErrorIs(t, err, ErrTargetSalesCountMustBeLessThanStock)
	})

	t.Run("should return error when campaign repo create returns error", func(t *testing.T) {
		returnErr := errors.New("error")
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)
		mockCampaignRepo.EXPECT().Create(gomock.Any()).Return(returnErr)

		err := campaignService.Create("C1", mockProduct, 10, 20, 50)
		assert.ErrorIs(t, err, returnErr)

	})

	t.Run("success", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Exist(gomock.Any()).Return(false)
		mockCampaignRepo.EXPECT().Create(gomock.Any()).Return(nil)

		err := campaignService.Create("C1", mockProduct, 10, 20, 50)
		assert.Nil(t, err)
	})

}

func TestCampaignServiceGetCampaignInfo(t *testing.T) {
	campaignService, teardown := setup(t)
	defer teardown()

	t.Run("should return error when campaign name is invalid", func(t *testing.T) {
		c, err := campaignService.Get("")
		assert.NotNil(t, err)
		assert.Nil(t, c)
	})

	t.Run("should return error when campaign is not found", func(t *testing.T) {
		mockCampaignRepo.EXPECT().Get(gomock.Any()).Return(nil, campaign.ErrCampaignNotFound)

		c, err := campaignService.Get("C1")
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
		mockProduct := &entity.Product{Code: code}
		mockCampaignData := &entity.Campaign{
			ID:                     id,
			Name:                   name,
			Product:                mockProduct,
			Duration:               duration,
			PriceManipulationLimit: priceManipulationLimit,
			TargetSalesCount:       targetSalesCount,
		}

		mockCampaignRepo.EXPECT().Get(gomock.Any()).Return(mockCampaignData, nil)

		c, err := campaignService.Get("C1")
		assert.Nil(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, c.ID, id)
		assert.Equal(t, c.Name, name)
		assert.Equal(t, c.Product.Code, code)
		assert.Equal(t, c.Duration, duration)
		assert.Equal(t, c.PriceManipulationLimit, priceManipulationLimit)
		assert.Equal(t, c.TargetSalesCount, targetSalesCount)
	})
}

func TestCampaignServiceGetAllCampaigns(t *testing.T) {
	campaignService, teardown := setup(t)
	defer teardown()

	t.Run("should return error when campaign repo get all empty array", func(t *testing.T) {
		mockCampaignRepo.EXPECT().GetAll().Return([]*entity.Campaign{})

		c, err := campaignService.GetAll()
		assert.ErrorIs(t, err, ErrNoCampaign)
		assert.Nil(t, c)
	})

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		name, _ := valueobject.NewName("C1")
		code, _ := valueobject.NewCode("P1")
		duration, _ := valueobject.NewDuration(10)
		priceManipulationLimit, _ := valueobject.NewPriceManipulationLimit(20)
		targetSalesCount, _ := valueobject.NewTargetSalesCount(100)
		mockProduct := &entity.Product{Code: code}
		mockCampaignData := &entity.Campaign{
			ID:                     id,
			Name:                   name,
			Product:                mockProduct,
			Duration:               duration,
			PriceManipulationLimit: priceManipulationLimit,
			TargetSalesCount:       targetSalesCount,
		}

		mockCampaignRepo.EXPECT().GetAll().Return([]*entity.Campaign{mockCampaignData})

		c, err := campaignService.GetAll()
		assert.Nil(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, c[0].ID, id)
		assert.Equal(t, c[0].Name, name)
		assert.Equal(t, c[0].Product.Code, code)
		assert.Equal(t, c[0].Duration, duration)
		assert.Equal(t, c[0].PriceManipulationLimit, priceManipulationLimit)
		assert.Equal(t, c[0].TargetSalesCount, targetSalesCount)
	})
}
