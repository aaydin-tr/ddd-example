package app

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aaydin-tr/e-commerce/service/campaign"
	"github.com/aaydin-tr/e-commerce/service/order"
	"github.com/aaydin-tr/e-commerce/service/product"
)

var (
	ErrCommandNotFound            = errors.New("Command not found")
	ErrInvalidParameters          = errors.New("Invalid parameters")
	ErrPriceMustBeFloat           = errors.New("Price must be float")
	ErrStockMustBeInt             = errors.New("Stock must be integer")
	ErrQuantityMustBeInt          = errors.New("Quantity must be integer")
	ErrDurationMustBeInt          = errors.New("Duration must be integer")
	ErrLimitMustBeInt             = errors.New("Limit must be integer")
	ErrTargetSalesMustBeInt       = errors.New("Target sales must be integer")
	ErrHourMustBeInt              = errors.New("Hour must be integer")
	ErrCampaignDoesNotHaveProduct = errors.New("Campaign does not have product")
)

type App struct {
	systemTime      time.Time
	commands        map[string]func(params []string) (string, error)
	productService  *product.ProductService
	orderSerivce    *order.OrderService
	campaignSerivce *campaign.CampaignService
}

func NewApp(productService *product.ProductService, orderService *order.OrderService, campaignService *campaign.CampaignService) *App {

	app := &App{
		systemTime:      time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
		productService:  productService,
		orderSerivce:    orderService,
		campaignSerivce: campaignService,
	}

	commands := make(map[string]func(params []string) (string, error))
	commands["create_product"] = app.createProduct
	commands["get_product_info"] = app.getProductInfo
	commands["create_order"] = app.createOrder
	commands["create_campaign"] = app.createCampaign
	commands["get_campaign_info"] = app.getCampaignInfo
	commands["increase_time"] = app.increaseTime

	app.commands = commands
	return app
}

func (this *App) Run(args []string) (string, error) {
	if cmd, ok := this.commands[args[0]]; ok {
		return cmd(args[1:])
	}

	return "", ErrCommandNotFound
}

func (this *App) createProduct(params []string) (string, error) {
	if len(params) != 3 {
		return "", ErrInvalidParameters
	}

	code := params[0]
	price, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		return "", ErrPriceMustBeFloat
	}
	stock, err := strconv.Atoi(params[2])
	if err != nil {
		return "", ErrStockMustBeInt
	}

	err = this.productService.Create(code, price, stock)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Product created; code %s, price %.1f, stock %d", code, price, stock), nil
}

func (this *App) getProductInfo(params []string) (string, error) {
	if len(params) != 1 {
		return "", ErrInvalidParameters
	}

	result, err := this.productService.Get(params[0])
	if err != nil {
		return "", err
	}

	result.IncreaseDemand(1)

	return fmt.Sprintf("Product %s info; price %.1f, stock %d", result.Code.Value(), result.Price.Value(), result.Stock.Value()), nil
}

func (this *App) createOrder(params []string) (string, error) {
	if len(params) != 2 {
		return "", ErrInvalidParameters
	}

	code := params[0]
	quantity, err := strconv.Atoi(params[1])
	if err != nil {
		return "", ErrQuantityMustBeInt
	}

	product, err := this.productService.Get(code)
	if err != nil {
		return "", err
	}

	err = this.orderSerivce.Create(product, quantity)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Order created; product %s, quantity %d", code, quantity), nil
}

func (this *App) createCampaign(params []string) (string, error) {
	if len(params) != 5 {
		return "", ErrInvalidParameters
	}

	name := params[0]
	code := params[1]
	duration, err := strconv.Atoi(params[2])
	if err != nil {
		return "", ErrDurationMustBeInt
	}
	limit, err := strconv.Atoi(params[3])
	if err != nil {
		return "", ErrLimitMustBeInt
	}
	targetSalesCount, err := strconv.Atoi(params[4])
	if err != nil {
		return "", ErrTargetSalesMustBeInt
	}

	product, err := this.productService.Get(code)
	if err != nil {
		return "", err
	}

	err = this.campaignSerivce.Create(name, product, duration, limit, targetSalesCount)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Campaign created; name %s, product %s, duration %d, limit %d, target sales count %d", name, code, duration, limit, targetSalesCount), nil
}

func (this *App) getCampaignInfo(params []string) (string, error) {
	if len(params) != 1 {
		return "", ErrInvalidParameters
	}

	result, err := this.campaignSerivce.Get(params[0])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Campaign %s info; Status %s, Target Sales %d, Total Sales %d, Turnover %.1f, Average Item Price %.1f", result.Name.Value(), result.Status.Value(), result.TargetSalesCount.Value(), result.TotalSales.Value(), (float64(result.TotalSales.Value()) * result.AverageItemPrice.Value()), result.AverageItemPrice.Value()), nil
}

func (this *App) increaseTime(params []string) (string, error) {
	if len(params) != 1 {
		return "", ErrInvalidParameters
	}

	hour, err := strconv.Atoi(params[0])
	if err != nil {
		return "", ErrHourMustBeInt
	}

	this.systemTime = this.systemTime.Add(time.Duration(hour) * time.Hour)
	campaigns, err := this.campaignSerivce.GetAll()
	if err != nil {
		return "", err
	}

	for _, campaign := range campaigns {

		if campaign.Product == nil {
			return "", ErrCampaignDoesNotHaveProduct
		}

		campaign.DecreaseDuration(hour)
		if campaign.Duration.Value() <= 0 {
			campaign.Close()
			campaign.Product.RemoveCampaign()
			continue
		}

		campaign.Product.Discount(campaign.PriceManipulationLimit)
	}

	return fmt.Sprintf("Time is %s", this.systemTime.Format("15:00")), nil
}
