package memory

import (
	"testing"

	"github.com/aaydin-tr/e-commerce/domain/campaign"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name         string
	campaignName string
	expectedErr  error
}

func TestMemoryCreateCampaign(t *testing.T) {
	mockRepo := NewCampaignRepository()

	t.Run("Create campaign", func(t *testing.T) {
		c := testCase{
			name:         "Create campaign new campaign",
			campaignName: "C1",
			expectedErr:  nil,
		}

		name, err := valueobject.NewName(c.campaignName)
		assert.NoError(t, err)

		campaign := &entity.Campaign{Name: name}

		err = mockRepo.Create(campaign)
		assert.ErrorIs(t, c.expectedErr, err)
	})

	t.Run("Create campaign which already exist", func(t *testing.T) {
		c := testCase{
			name:         "Create campaign new campaign",
			campaignName: "C1",
			expectedErr:  campaign.ErrCampaignAlreadyExist,
		}

		name, err := valueobject.NewName(c.campaignName)
		assert.NoError(t, err)

		campaign := &entity.Campaign{Name: name}

		err = mockRepo.Create(campaign)
		assert.ErrorIs(t, c.expectedErr, err)
	})

}

func TestMemoryGetCampaign(t *testing.T) {
	mockRepo := NewCampaignRepository()
	name, _ := valueobject.NewName("C1")
	mockRepo.campaigns["C1"] = &entity.Campaign{Name: name}

	t.Run("Get campaign", func(t *testing.T) {
		name, err := valueobject.NewName("C1")
		assert.NoError(t, err)

		campaign, err := mockRepo.Get(name)

		assert.NoError(t, err)
		assert.Equal(t, name, campaign.Name)
	})

	t.Run("Get campaign which not exist", func(t *testing.T) {
		name, err := valueobject.NewName("C2")
		assert.NoError(t, err)

		c, err := mockRepo.Get(name)
		assert.ErrorIs(t, err, campaign.ErrCampaignNotFound)
		assert.Nil(t, c)
	})
}
