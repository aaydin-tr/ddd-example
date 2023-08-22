package memory

import (
	"sync"

	"github.com/aaydin-tr/e-commerce/domain/campaign"
	"github.com/aaydin-tr/e-commerce/entity"
)

type CampaignRepository struct {
	campaigns map[string]*entity.Campaign
	sync.RWMutex
}

func NewCampaignRepository() *CampaignRepository {
	return &CampaignRepository{campaigns: make(map[string]*entity.Campaign)}
}

func (r *CampaignRepository) Create(newCampaign *entity.Campaign) error {
	r.Lock()
	defer r.Unlock()

	if r.campaigns[newCampaign.Name.Value()] != nil {
		return campaign.ErrCampaignAlreadyExist
	}

	r.campaigns[newCampaign.Name.Value()] = newCampaign
	return nil
}

func (r *CampaignRepository) Get(name string) (*entity.Campaign, error) {
	r.RLock()
	defer r.RUnlock()

	result := r.campaigns[name]
	if result == nil {
		return nil, campaign.ErrCampaignNotFound
	}

	return result, nil
}
