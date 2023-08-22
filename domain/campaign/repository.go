package campaign

import (
	"errors"

	"github.com/aaydin-tr/e-commerce/entity"
)

var (
	ErrCampaignAlreadyExist = errors.New("Campaign already exist")
	ErrCampaignNotFound     = errors.New("Campaign not found")
)

type CampaignRepository interface {
	Create(campaign *entity.Campaign) error
	Get(name string) (*entity.Campaign, error)
}
