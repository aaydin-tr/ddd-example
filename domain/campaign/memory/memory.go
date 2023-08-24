package memory

import (
	"github.com/aaydin-tr/e-commerce/domain/campaign"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/types"
	"github.com/aaydin-tr/e-commerce/valueobject"
)

type CampaignRepository struct {
	storage types.Storage[*entity.Campaign]
}

func NewCampaignRepository(storage types.Storage[*entity.Campaign]) *CampaignRepository {
	return &CampaignRepository{storage: storage}
}

func (r *CampaignRepository) Create(newCampaign *entity.Campaign) error {
	_, ok := r.storage.Get(newCampaign.Name.Value())
	if ok {
		return campaign.ErrCampaignAlreadyExist
	}

	r.storage.Set(newCampaign.Name.Value(), newCampaign)
	return nil
}

func (r *CampaignRepository) Get(name valueobject.Name) (*entity.Campaign, error) {
	result, ok := r.storage.Get(name.Value())
	if !ok {
		return nil, campaign.ErrCampaignNotFound
	}

	return result, nil
}

func (r *CampaignRepository) Exist(name valueobject.Name) bool {
	_, ok := r.storage.Get(name.Value())
	return ok
}

func (r *CampaignRepository) GetAll() []*entity.Campaign {
	var result []*entity.Campaign
	for _, item := range r.storage.Values() {
		result = append(result, item)
	}

	return result
}
