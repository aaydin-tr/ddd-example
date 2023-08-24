package memory

import (
	"testing"

	"github.com/aaydin-tr/e-commerce/domain/campaign"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/pkg/storage"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	expectedErr error
}

func TestMemoryCreateCampaign(t *testing.T) {
	mockRepo := NewCampaignRepository(storage.New[*entity.Campaign]())

	cname, _ := valueobject.NewName("C1")

	t.Run("Create campaign", func(t *testing.T) {
		c := testCase{
			name:        "Create campaign new campaign",
			expectedErr: nil,
		}

		campaign := &entity.Campaign{Name: cname}

		err := mockRepo.Create(campaign)
		assert.ErrorIs(t, c.expectedErr, err)
	})

	t.Run("Create campaign which already exist", func(t *testing.T) {
		c := testCase{
			name:        "Create campaign new campaign",
			expectedErr: campaign.ErrCampaignAlreadyExist,
		}

		name, err := valueobject.NewName(cname.Value())
		assert.NoError(t, err)

		campaign := &entity.Campaign{Name: name}

		err = mockRepo.Create(campaign)
		assert.ErrorIs(t, c.expectedErr, err)
	})

}

func TestMemoryGetCampaign(t *testing.T) {
	mockRepo := NewCampaignRepository(storage.New[*entity.Campaign]())
	name, _ := valueobject.NewName("C1")
	mockRepo.storage.Set(name.Value(), &entity.Campaign{Name: name})

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

func TestMemoryGetAll(t *testing.T) {
	mockRepo := NewCampaignRepository(storage.New[*entity.Campaign]())
	first, _ := valueobject.NewName("C1")
	second, _ := valueobject.NewName("C2")
	mockRepo.storage.Set(first.Value(), &entity.Campaign{Name: first})
	mockRepo.storage.Set(second.Value(), &entity.Campaign{Name: second})

	t.Run("Get all campaign", func(t *testing.T) {
		campaigns := mockRepo.GetAll()
		assert.Len(t, campaigns, 2)
	})

}

func TestMemoryExist(t *testing.T) {
	mockRepo := NewCampaignRepository(storage.New[*entity.Campaign]())
	name, _ := valueobject.NewName("C1")
	mockRepo.storage.Set(name.Value(), &entity.Campaign{Name: name})

	t.Run("Campaign exist", func(t *testing.T) {
		name, err := valueobject.NewName("C1")
		assert.NoError(t, err)

		ok := mockRepo.Exist(name)
		assert.True(t, ok)
	})

	t.Run("Campaign not exist", func(t *testing.T) {
		name, err := valueobject.NewName("C2")
		assert.NoError(t, err)

		ok := mockRepo.Exist(name)
		assert.False(t, ok)
	})
}
