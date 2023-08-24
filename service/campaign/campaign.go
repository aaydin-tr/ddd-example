package campaign

import (
	"errors"

	"github.com/aaydin-tr/e-commerce/domain/campaign"
	"github.com/aaydin-tr/e-commerce/domain/product"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/google/uuid"
)

var (
	ErrNoCampaign                          = errors.New("No campaign found")
	ErrTargetSalesCountMustBeLessThanStock = errors.New("Target sales count must be less than stock")
	ErrCampaignDoesNotHaveProduct          = errors.New("Campaign does not have product")
)

type CampaignService struct {
	campaignRepository campaign.CampaignRepository
	productRepository  product.ProductRepository
}

func NewCampaignService(campaignRepository campaign.CampaignRepository,
	productRepository product.ProductRepository,
) *CampaignService {
	return &CampaignService{
		campaignRepository: campaignRepository,
		productRepository:  productRepository,
	}
}

func (c *CampaignService) CreateCampaign(campaignName string, productCode string, campaignDuration int, campaignPriceManipulationLimit int, campaignTargetSalesCount int) error {
	name, err := valueobject.NewName(campaignName)
	if err != nil {
		return err
	}

	if c.campaignRepository.Exist(name) {
		return campaign.ErrCampaignAlreadyExist
	}

	code, err := valueobject.NewCode(productCode)
	if err != nil {
		return err
	}

	p, err := c.productRepository.Get(code)
	if err != nil {
		return err
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

	if p.Stock.Value() < targetSalesCount.Value() {
		return ErrTargetSalesCountMustBeLessThanStock
	}

	status, err := valueobject.NewStatus(valueobject.Active)
	if err != nil {
		return err
	}

	newCampaign := &entity.Campaign{
		ID:                     uuid.New(),
		Name:                   name,
		ProductCode:            code,
		Duration:               duration,
		PriceManipulationLimit: priceManipulationLimit,
		TargetSalesCount:       targetSalesCount,
		Status:                 status,
	}

	err = c.campaignRepository.Create(newCampaign)

	if err != nil {
		return err
	}

	p.Campaign = newCampaign

	return nil

}

func (c *CampaignService) GetCampaignInfo(campaignName string) (*entity.Campaign, error) {
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

func (c *CampaignService) IncreaseTime(duration int) error {
	campaigns := c.campaignRepository.GetAll()
	if len(campaigns) == 0 {
		return ErrNoCampaign
	}

	for _, campaign := range campaigns {
		campaignProduct, err := c.productRepository.Get(campaign.ProductCode)
		if err != nil {
			return ErrCampaignDoesNotHaveProduct
		}

		campaign.DecreaseDuration(duration)
		if campaign.Duration.Value() <= 0 {
			campaign.Close()
			campaignProduct.RemoveCampaign()
			continue
		}

		campaignProduct.Discount(campaign.PriceManipulationLimit)
	}

	return nil
}
