//go:build integration
// +build integration

package app

import (
	"strings"
	"testing"

	"github.com/aaydin-tr/e-commerce/service/order"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/stretchr/testify/assert"
)

type integationTestCase struct {
	name                     string
	productCode              string
	campaignCode             string
	commands                 []commandTestCase
	isCampaign               bool
	expectedLastPrice        float64
	expectedLastStock        int
	expectedSalesCount       int
	expectedTurnover         float64
	expectedAverageItemPrice float64
	expectedCampaignStatus   string
}

type commandTestCase struct {
	args string
	msg  string
	err  error
}

func TestAppIntegration(t *testing.T) {
	t.Parallel()

	testCases := []integationTestCase{
		{
			name:        "Create product without campaign and sell half of the stock",
			isCampaign:  false,
			productCode: "P1",
			commands: []commandTestCase{
				{args: "create_product P1 100 100", msg: "Product created; code P1, price 100.0, stock 100"},
				{args: "create_order P1 50", msg: "Order created; product P1, quantity 50"},
				{args: "get_product_info P1", msg: "Product P1 info; price 100.0, stock 50"},
			},
			expectedLastPrice: 100,
			expectedLastStock: 50,
		},
		{
			name:        "Create product without campaign and sell all of the stock",
			isCampaign:  false,
			productCode: "P1",
			commands: []commandTestCase{
				{args: "create_product P1 100 100", msg: "Product created; code P1, price 100.0, stock 100"},
				{args: "create_order P1 100", msg: "Order created; product P1, quantity 100"},
				{args: "get_product_info P1", msg: "Product P1 info; price 100.0, stock 0"},
			},
			expectedLastPrice: 100,
			expectedLastStock: 0,
		},
		{
			name:        "Create product without campaign and sell more than the stock",
			isCampaign:  false,
			productCode: "P1",
			commands: []commandTestCase{
				{args: "create_product P1 100 100", msg: "Product created; code P1, price 100.0, stock 100"},
				{args: "create_order P1 150", msg: "", err: order.ErrInsufficientStock},
				{args: "get_product_info P1", msg: "Product P1 info; price 100.0, stock 100"},
			},
			expectedLastPrice: 100,
			expectedLastStock: 100,
		},
		{
			name:         "Create product with campaign and sell half of the stock without price change",
			isCampaign:   true,
			campaignCode: "C1",
			productCode:  "P1",
			commands: []commandTestCase{
				{args: "create_product P1 100 100", msg: "Product created; code P1, price 100.0, stock 100"},
				{args: "create_campaign C1 P1 10 20 100", msg: "Campaign created; name C1, product P1, duration 10, limit 20, target sales count 100"},
				{args: "create_order P1 50", msg: "Order created; product P1, quantity 50"},
				{args: "get_product_info P1", msg: "Product P1 info; price 100.0, stock 50"},
				{args: "get_campaign_info C1", msg: "Campaign C1 info; Status Active, Target Sales 100, Total Sales 50, Turnover 5000.0, Average Item Price 100.0"},
			},
			expectedLastPrice:        100,
			expectedLastStock:        50,
			expectedSalesCount:       50,
			expectedTurnover:         5000,
			expectedAverageItemPrice: 100,
			expectedCampaignStatus:   valueobject.Active,
		},
		{
			name:         "Create product with campaign and sell all of the stock without price change",
			isCampaign:   true,
			productCode:  "P1",
			campaignCode: "C1",
			commands: []commandTestCase{
				{args: "create_product P1 100 100", msg: "Product created; code P1, price 100.0, stock 100"},
				{args: "create_campaign C1 P1 10 20 100", msg: "Campaign created; name C1, product P1, duration 10, limit 20, target sales count 100"},
				{args: "create_order P1 100", msg: "Order created; product P1, quantity 100"},
				{args: "get_product_info P1", msg: "Product P1 info; price 100.0, stock 0"},
				{args: "get_campaign_info C1", msg: "Campaign C1 info; Status Ended, Target Sales 100, Total Sales 100, Turnover 10000.0, Average Item Price 100.0"},
			},
			expectedLastPrice:        100,
			expectedLastStock:        0,
			expectedSalesCount:       100,
			expectedTurnover:         10000,
			expectedAverageItemPrice: 100,
			expectedCampaignStatus:   valueobject.Ended,
		},
		{
			name:         "Create product with campaign and sell more than the stock without price change",
			isCampaign:   true,
			campaignCode: "C1",
			productCode:  "P1",
			commands: []commandTestCase{
				{args: "create_product P1 100 200", msg: "Product created; code P1, price 100.0, stock 200"},
				{args: "create_campaign C1 P1 10 20 100", msg: "Campaign created; name C1, product P1, duration 10, limit 20, target sales count 100"},
				{args: "create_order P1 150", msg: "Order created; product P1, quantity 150"},
				{args: "get_product_info P1", msg: "Product P1 info; price 100.0, stock 50"},
				{args: "get_campaign_info C1", msg: "Campaign C1 info; Status Ended, Target Sales 100, Total Sales 100, Turnover 10000.0, Average Item Price 100.0"},
			},
			expectedLastPrice:        100,
			expectedLastStock:        50,
			expectedSalesCount:       100,
			expectedTurnover:         10000,
			expectedAverageItemPrice: 100,
			expectedCampaignStatus:   valueobject.Ended,
		},
		{
			name:         "Create product with campaign price increase",
			isCampaign:   true,
			productCode:  "P1",
			campaignCode: "C1",
			commands: []commandTestCase{
				{args: "create_product P1 100 100", msg: "Product created; code P1, price 100.0, stock 100"},
				{args: "create_campaign C1 P1 10 20 100", msg: "Campaign created; name C1, product P1, duration 10, limit 20, target sales count 100"},
				{args: "create_order P1 50", msg: "Order created; product P1, quantity 50"},
				{args: "increase_time 1", msg: "Time is 01:00"},
				{args: "get_product_info P1", msg: "Product P1 info; price 120.0, stock 50"},
				{args: "get_campaign_info C1", msg: "Campaign C1 info; Status Active, Target Sales 100, Total Sales 50, Turnover 5000.0, Average Item Price 100.0"},
			},
			expectedLastPrice:        120,
			expectedLastStock:        50,
			expectedSalesCount:       50,
			expectedTurnover:         5000,
			expectedAverageItemPrice: 100,
			expectedCampaignStatus:   valueobject.Active,
		},
		{
			name:         "Create product with campaign price decrease",
			isCampaign:   true,
			productCode:  "P1",
			campaignCode: "C1",
			commands: []commandTestCase{
				{args: "create_product P1 100 100", msg: "Product created; code P1, price 100.0, stock 100"},
				{args: "create_campaign C1 P1 10 20 100", msg: "Campaign created; name C1, product P1, duration 10, limit 20, target sales count 100"},
				{args: "get_product_info P1", msg: "Product P1 info; price 100.0, stock 100"},
				{args: "increase_time 1", msg: "Time is 01:00"},
				{args: "get_product_info P1", msg: "Product P1 info; price 80.0, stock 100"},
				{args: "get_campaign_info C1", msg: "Campaign C1 info; Status Active, Target Sales 100, Total Sales 0, Turnover 0.0, Average Item Price 0.0"},
			},
			expectedLastPrice:        80,
			expectedLastStock:        100,
			expectedSalesCount:       0,
			expectedTurnover:         0,
			expectedAverageItemPrice: 0,
			expectedCampaignStatus:   valueobject.Active,
		},

		{
			name:         "Campaign should be ended when the duration is over",
			isCampaign:   true,
			productCode:  "P1",
			campaignCode: "C1",
			commands: []commandTestCase{
				{args: "create_product P1 100 100", msg: "Product created; code P1, price 100.0, stock 100"},
				{args: "create_campaign C1 P1 1 20 100", msg: "Campaign created; name C1, product P1, duration 1, limit 20, target sales count 100"},
				{args: "increase_time 1", msg: "Time is 01:00"},
				{args: "get_campaign_info C1", msg: "Campaign C1 info; Status Ended, Target Sales 100, Total Sales 0, Turnover 0.0, Average Item Price 0.0"},
			},
			expectedLastPrice:        100,
			expectedLastStock:        100,
			expectedSalesCount:       0,
			expectedTurnover:         0,
			expectedAverageItemPrice: 0,
			expectedCampaignStatus:   valueobject.Ended,
		},
		{
			name:         "Campaign should be ended when the target sales count is reached",
			isCampaign:   true,
			productCode:  "P1",
			campaignCode: "C1",
			commands: []commandTestCase{
				{args: "create_product P1 100 200", msg: "Product created; code P1, price 100.0, stock 200"},
				{args: "create_campaign C1 P1 10 20 100", msg: "Campaign created; name C1, product P1, duration 10, limit 20, target sales count 100"},
				{args: "create_order P1 100", msg: "Order created; product P1, quantity 100"},
				{args: "increase_time 1", msg: "Time is 01:00"},
				{args: "get_campaign_info C1", msg: "Campaign C1 info; Status Ended, Target Sales 100, Total Sales 100, Turnover 10000.0, Average Item Price 100.0"},
			},
			expectedLastPrice:        100,
			expectedLastStock:        100,
			expectedSalesCount:       100,
			expectedTurnover:         10000,
			expectedAverageItemPrice: 100,
			expectedCampaignStatus:   valueobject.Ended,
		},
		{
			name:         "Product price should be initial price when the campaign is ended",
			isCampaign:   true,
			productCode:  "P1",
			campaignCode: "C1",
			commands: []commandTestCase{
				{args: "create_product P1 100 200", msg: "Product created; code P1, price 100.0, stock 200"},
				{args: "create_campaign C1 P1 10 20 100", msg: "Campaign created; name C1, product P1, duration 10, limit 20, target sales count 100"},
				{args: "create_order P1 100", msg: "Order created; product P1, quantity 100"},
				{args: "increase_time 1", msg: "Time is 01:00"},
				{args: "get_product_info P1", msg: "Product P1 info; price 100.0, stock 100"},
			},
			expectedLastPrice:        100,
			expectedLastStock:        100,
			expectedSalesCount:       100,
			expectedTurnover:         10000,
			expectedAverageItemPrice: 100,
			expectedCampaignStatus:   valueobject.Ended,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			app := setup(t)

			for _, command := range testCase.commands {
				msg, err := app.Run(strings.Split(command.args, " "))
				assert.Equal(t, command.err, err)
				assert.Equal(t, command.msg, msg)
			}

			product, err := app.productService.Get(testCase.productCode)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedLastPrice, product.Price.Value())
			assert.Equal(t, testCase.expectedLastStock, product.Stock.Value())

			if testCase.isCampaign {
				campaign, err := app.campaignSerivce.Get(testCase.campaignCode)
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedSalesCount, campaign.TotalSales.Value())
				assert.Equal(t, testCase.expectedTurnover, (float64(campaign.TotalSales.Value()) * campaign.AverageItemPrice.Value()))
				assert.Equal(t, testCase.expectedAverageItemPrice, campaign.AverageItemPrice.Value())
				assert.Equal(t, testCase.expectedCampaignStatus, campaign.Status.Value())
			}
		})
	}

}
