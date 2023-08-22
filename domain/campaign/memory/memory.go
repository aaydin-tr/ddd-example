package memory

import (
	"sync"

	"github.com/aaydin-tr/e-commerce/domain/campaign"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
)

type CampaignRepository struct {
	campaigns map[string]*entity.Campaign
	mu        sync.RWMutex
}

func NewCampaignRepository() *CampaignRepository {
	return &CampaignRepository{campaigns: make(map[string]*entity.Campaign)}
}

func (r *CampaignRepository) Create(newCampaign *entity.Campaign) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.campaigns[newCampaign.Name.Value()] != nil {
		return campaign.ErrCampaignAlreadyExist
	}

	r.campaigns[newCampaign.Name.Value()] = newCampaign
	return nil
}

func (r *CampaignRepository) Get(name valueobject.Name) (*entity.Campaign, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := r.campaigns[name.Value()]
	if result == nil {
		return nil, campaign.ErrCampaignNotFound
	}

	return result, nil
}
