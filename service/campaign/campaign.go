package campaign

import (
	"errors"

	"github.com/aaydin-tr/e-commerce/domain/campaign"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/google/uuid"
)

var (
	ErrNoCampaign                          = errors.New("No campaign found")
	ErrTargetSalesCountMustBeLessThanStock = errors.New("Target sales count must be less than stock")
)

type CampaignServiceInterface interface {
	Create(campaignName string, product *entity.Product, campaignDuration int, campaignPriceManipulationLimit int, campaignTargetSalesCount int) error
	Get(campaignName string) (*entity.Campaign, error)
	GetAll() ([]*entity.Campaign, error)
}

type CampaignService struct {
	campaignRepository campaign.CampaignRepository
}

func NewCampaignService(campaignRepository campaign.CampaignRepository) *CampaignService {
	return &CampaignService{
		campaignRepository: campaignRepository,
	}
}

func (c *CampaignService) Create(campaignName string, product *entity.Product, campaignDuration int, campaignPriceManipulationLimit int, campaignTargetSalesCount int) error {
	name, err := valueobject.NewName(campaignName)
	if err != nil {
		return err
	}

	if c.campaignRepository.Exist(name) {
		return campaign.ErrCampaignAlreadyExist
	}

	duration, err := valueobject.NewDuration(campaignDuration)
	if err != nil {
		return err
	}

	priceManipulationLimit, err := valueobject.NewPriceManipulationLimit(campaignPriceManipulationLimit)
	if err != nil {
		return err
	}

	targetSalesCount, err := valueobject.NewTargetSalesCount(campaignTargetSalesCount)
	if err != nil {
		return err
	}

	if product.Stock.Value() < targetSalesCount.Value() {
		return ErrTargetSalesCountMustBeLessThanStock
	}

	status, err := valueobject.NewStatus(valueobject.Active)
	if err != nil {
		return err
	}

	newCampaign := &entity.Campaign{
		ID:                     uuid.New(),
		Name:                   name,
		Product:                product,
		Duration:               duration,
		PriceManipulationLimit: priceManipulationLimit,
		TargetSalesCount:       targetSalesCount,
		Status:                 status,
	}

	err = c.campaignRepository.Create(newCampaign)

	if err != nil {
		return err
	}

	product.Campaign = newCampaign

	return nil

}

func (c *CampaignService) Get(campaignName string) (*entity.Campaign, error) {
	name, err := valueobject.NewName(campaignName)
	if err != nil {
		return nil, err
	}

	campaign, err := c.campaignRepository.Get(name)
	if err != nil {
		return nil, err
	}

	return campaign, nil
}

func (c *CampaignService) GetAll() ([]*entity.Campaign, error) {
	campaigns := c.campaignRepository.GetAll()
	if len(campaigns) == 0 {
		return nil, ErrNoCampaign
	}

	return campaigns, nil
}
