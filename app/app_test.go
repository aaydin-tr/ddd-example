package app

import (
	"testing"

	campaignRepo "github.com/aaydin-tr/e-commerce/domain/campaign/memory"
	orderRepo "github.com/aaydin-tr/e-commerce/domain/order/memory"
	productRepo "github.com/aaydin-tr/e-commerce/domain/product/memory"
	"github.com/aaydin-tr/e-commerce/service/campaign"
	"github.com/aaydin-tr/e-commerce/service/order"
	"github.com/aaydin-tr/e-commerce/service/product"
	"github.com/stretchr/testify/assert"

	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/pkg/storage"
)

func setup(t *testing.T) *App {
	mockProductRepository := productRepo.NewProductRepository(storage.New[*entity.Product]())
	mockOrderRepository := orderRepo.NewOrderRepository(storage.New[*entity.Order]())
	mockCampaignRepository := campaignRepo.NewCampaignRepository(storage.New[*entity.Campaign]())

	mockProductService := product.NewProductService(mockProductRepository)
	mockOrderService := order.NewOrderService(mockProductRepository, mockOrderRepository)
	mockCampaignService := campaign.NewCampaignService(mockCampaignRepository, mockOrderRepository, mockProductRepository)

	return NewApp(mockProductService, mockOrderService, mockCampaignService)
}

func TestNewApp(t *testing.T) {
	app := setup(t)
	assert.NotNil(t, app)
}

func TestAppRun(t *testing.T) {
	app := setup(t)

	t.Parallel()
	t.Run("valid command", func(t *testing.T) {
		msg, err := app.Run([]string{"create_product", "P1", "100", "1000"})
		assert.NoError(t, err)
		assert.Equal(t, "Product created; code P1, price 100.0, stock 1000", msg)
	})

	t.Run("invalid command", func(t *testing.T) {
		msg, err := app.Run([]string{"invalid_command"})
		assert.ErrorIs(t, err, ErrCommandNotFound)
		assert.Equal(t, "", msg)
	})

}

func TestAppCreateProduct(t *testing.T) {
	app := setup(t)
	app.productService.Create("P1", 100, 1000)

	t.Parallel()
	t.Run("invalid parameters", func(t *testing.T) {
		msg, err := app.createProduct([]string{"P2", "100"})
		assert.ErrorIs(t, err, ErrInvalidParameters)
		assert.Equal(t, "", msg)
	})
	t.Run("invalid price", func(t *testing.T) {
		msg, err := app.createProduct([]string{"P3", "invalid_price", "1000"})
		assert.ErrorIs(t, err, ErrPriceMustBeFloat)
		assert.Equal(t, "", msg)
	})
	t.Run("invalid stock", func(t *testing.T) {
		msg, err := app.createProduct([]string{"P4", "100", "invalid_stock"})
		assert.ErrorIs(t, err, ErrStockMustBeInt)
		assert.Equal(t, "", msg)
	})

	t.Run("invalid stock", func(t *testing.T) {
		msg, err := app.createProduct([]string{"P5", "100", "invalid_stock"})
		assert.ErrorIs(t, err, ErrStockMustBeInt)
		assert.Equal(t, "", msg)
	})

	t.Run("product already exist", func(t *testing.T) {
		msg, err := app.createProduct([]string{"P1", "100", "1000"})
		assert.NotNil(t, err)
		assert.Equal(t, "", msg)
	})

	t.Run("valid parameters", func(t *testing.T) {
		msg, err := app.createProduct([]string{"P6", "100", "1000"})
		assert.NoError(t, err)
		assert.Equal(t, "Product created; code P6, price 100.0, stock 1000", msg)

		p, err := app.productService.Get("P6")
		assert.NoError(t, err)
		assert.Equal(t, "P6", p.Code.Value())
		assert.Equal(t, 100.0, p.Price.Value())
		assert.Equal(t, 1000, p.Stock.Value())
	})

}

func TestAppGetProductInfo(t *testing.T) {
	app := setup(t)
	app.productService.Create("P1", 100, 1000)

	t.Parallel()
	t.Run("invalid parameters", func(t *testing.T) {
		msg, err := app.getProductInfo([]string{})
		assert.ErrorIs(t, err, ErrInvalidParameters)
		assert.Equal(t, "", msg)
	})

	t.Run("product not found", func(t *testing.T) {
		msg, err := app.getProductInfo([]string{"P3"})
		assert.NotNil(t, err)
		assert.Equal(t, "", msg)
	})

	t.Run("valid parameters", func(t *testing.T) {
		msg, err := app.getProductInfo([]string{"P1"})
		assert.NoError(t, err)
		assert.Equal(t, "Product P1 info; price 100.0, stock 1000", msg)
	})

}

func TestAppCreateOrder(t *testing.T) {
	app := setup(t)
	app.productService.Create("P1", 100, 1000)

	t.Parallel()
	t.Run("invalid parameters", func(t *testing.T) {
		msg, err := app.createOrder([]string{"P2"})
		assert.ErrorIs(t, err, ErrInvalidParameters)
		assert.Equal(t, "", msg)
	})

	t.Run("product not found", func(t *testing.T) {
		msg, err := app.createOrder([]string{"P3", "10"})
		assert.NotNil(t, err)
		assert.Equal(t, "", msg)
	})

	t.Run("invalid quantity", func(t *testing.T) {
		msg, err := app.createOrder([]string{"P1", "invalid_quantity"})
		assert.ErrorIs(t, err, ErrQuantityMustBeInt)
		assert.Equal(t, "", msg)
	})

	t.Run("insufficient stock", func(t *testing.T) {
		msg, err := app.createOrder([]string{"P1", "1001"})
		assert.NotNil(t, err)
		assert.Equal(t, "", msg)
	})

	t.Run("valid parameters", func(t *testing.T) {
		msg, err := app.createOrder([]string{"P1", "10"})
		assert.NoError(t, err)
		assert.Equal(t, "Order created; product P1, quantity 10", msg)

		p, err := app.productService.Get("P1")
		assert.NoError(t, err)
		assert.Equal(t, 990, p.Stock.Value())
	})

}

func TestAppCreateCampaign(t *testing.T) {
	app := setup(t)
	app.productService.Create("P1", 100, 1000)

	t.Parallel()

	t.Run("invalid parameters", func(t *testing.T) {
		msg, err := app.createCampaign([]string{"C1", "P1"})
		assert.ErrorIs(t, err, ErrInvalidParameters)
		assert.Equal(t, "", msg)
	})

	t.Run("product not found", func(t *testing.T) {
		msg, err := app.createCampaign([]string{"C1", "P2", "10", "20", "3"})
		assert.NotNil(t, err)
		assert.Equal(t, "", msg)
	})

	t.Run("invalid duration", func(t *testing.T) {
		msg, err := app.createCampaign([]string{"C1", "P1", "invalid_duration", "20", "3"})
		assert.ErrorIs(t, err, ErrDurationMustBeInt)
		assert.Equal(t, "", msg)
	})

	t.Run("invalid limit", func(t *testing.T) {
		msg, err := app.createCampaign([]string{"C1", "P1", "10", "invalid_limit", "3"})
		assert.ErrorIs(t, err, ErrLimitMustBeInt)
		assert.Equal(t, "", msg)
	})

	t.Run("invalid target sales count", func(t *testing.T) {
		msg, err := app.createCampaign([]string{"C1", "P1", "10", "20", "invalid_target_sales_count"})
		assert.ErrorIs(t, err, ErrTargetSalesMustBeInt)
		assert.Equal(t, "", msg)
	})

	t.Run("valid parameters", func(t *testing.T) {
		msg, err := app.createCampaign([]string{"C1", "P1", "10", "20", "3"})
		assert.NoError(t, err)
		assert.Equal(t, "Campaign created; name C1, product P1, duration 10, limit 20, target sales count 3", msg)

		c, err := app.campaignSerivce.GetCampaignInfo("C1")
		assert.NoError(t, err)
		assert.Equal(t, "C1", c.Name.Value())
		assert.Equal(t, "P1", c.ProductCode.Value())
		assert.Equal(t, 10, c.Duration.Value())
		assert.Equal(t, 20, c.PriceManipulationLimit.Value())
		assert.Equal(t, 3, c.TargetSalesCount.Value())
	})

}

func TestAppGetCampaignInfo(t *testing.T) {
	app := setup(t)
	app.productService.Create("P1", 100, 1000)
	app.campaignSerivce.CreateCampaign("C1", "P1", 10, 20, 100)

	t.Parallel()

	t.Run("invalid parameters", func(t *testing.T) {
		msg, err := app.getCampaignInfo([]string{})
		assert.ErrorIs(t, err, ErrInvalidParameters)
		assert.Equal(t, "", msg)
	})

	t.Run("campaign not found", func(t *testing.T) {
		msg, err := app.getCampaignInfo([]string{"C2"})
		assert.NotNil(t, err)
		assert.Equal(t, "", msg)
	})

	t.Run("valid parameters", func(t *testing.T) {
		msg, err := app.getCampaignInfo([]string{"C1"})
		assert.NoError(t, err)
		assert.Equal(t, "Campaign C1 info; Status Active, Target Sales 100, Total Sales 0, Turnover 0.0, Average Item Price 0.0", msg)
	})

}

func TestAppIncreaseTime(t *testing.T) {
	app := setup(t)
	app.productService.Create("P1", 100, 1000)
	app.campaignSerivce.CreateCampaign("C1", "P1", 10, 20, 100)

	t.Parallel()

	t.Run("invalid parameters", func(t *testing.T) {
		msg, err := app.increaseTime([]string{})
		assert.NotNil(t, err)
		assert.Equal(t, "", msg)
	})

	t.Run("invalid time", func(t *testing.T) {
		msg, err := app.increaseTime([]string{"invalid_time"})
		assert.ErrorIs(t, err, ErrHourMustBeInt)
		assert.Equal(t, "", msg)
	})

	t.Run("no campaign", func(t *testing.T) {
		newApp := setup(t)
		msg, err := newApp.increaseTime([]string{"10"})
		assert.NotNil(t, err)
		assert.Equal(t, "", msg)
	})

	t.Run("valid parameters", func(t *testing.T) {
		msg, err := app.increaseTime([]string{"10"})
		assert.NoError(t, err)
		assert.Equal(t, "Time is 10:00", msg)
		assert.Equal(t, 10, app.systemTime.Hour())
	})

}
