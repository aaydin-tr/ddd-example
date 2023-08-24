package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	campaignRepo "github.com/aaydin-tr/e-commerce/domain/campaign/memory"
	orderRepo "github.com/aaydin-tr/e-commerce/domain/order/memory"
	productRepo "github.com/aaydin-tr/e-commerce/domain/product/memory"

	"github.com/aaydin-tr/e-commerce/app"
	"github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/pkg/storage"
	"github.com/aaydin-tr/e-commerce/service/campaign"
	"github.com/aaydin-tr/e-commerce/service/order"
	"github.com/aaydin-tr/e-commerce/service/product"
)

func main() {
	productRepository := productRepo.NewProductRepository(storage.New[*entity.Product]())
	orderRepository := orderRepo.NewOrderRepository(storage.New[*entity.Order]())
	campaignRepository := campaignRepo.NewCampaignRepository(storage.New[*entity.Campaign]())

	productService := product.NewProductService(productRepository)
	orderService := order.NewOrderService(productRepository, orderRepository)
	campaignService := campaign.NewCampaignService(campaignRepository, orderRepository, productRepository)

	app := app.NewApp(productService, orderService, campaignService)
	fmt.Println("Please enter command")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if !scanner.Scan() {
			fmt.Printf("Error while reading input: %s\n", scanner.Err())
			return
		}
		input := scanner.Text()
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		msg, err := app.Run(args)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}

		fmt.Println(msg)
	}

}
