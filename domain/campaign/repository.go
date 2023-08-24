package campaign

import (
	"errors"

	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
)

var (
	ErrCampaignAlreadyExist = errors.New("Campaign already exist")
	ErrCampaignNotFound     = errors.New("Campaign not found")
)

//go:generate mockgen -destination=../../mock/repository/campaign/campaign.go -package=repository github.com/aaydin-tr/e-commerce/domain/campaign CampaignRepository
type CampaignRepository interface {
	Create(campaign *entity.Campaign) error
	Get(name valueobject.Name) (*entity.Campaign, error)
	GetAll() []*entity.Campaign
	Exist(name valueobject.Name) bool
}
