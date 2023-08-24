package entity

import (
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/google/uuid"
)

type Product struct {
	ID       uuid.UUID
	Code     valueobject.Code
	Price    valueobject.Price
	Stock    valueobject.Stock
	Campaign *Campaign

	InititalStock    valueobject.Stock
	InititalPrice    valueobject.Price
	TotalDemandCount valueobject.Demand
}

func (p *Product) DecreaseStock(amount int) error {
	newStock, err := valueobject.NewStock(p.Stock.Value() - amount)
	if err != nil {
		return err
	}

	p.Stock = newStock
	return nil
}

func (p *Product) UpdatePrice(price float64) error {
	newPrice, err := valueobject.NewPrice(price)
	if err != nil {
		return err
	}

	p.Price = newPrice
	return nil
}

func (p *Product) IncreaseDemand(amount int) error {
	newTotalDemand, err := valueobject.NewDemand(p.TotalDemandCount.Value() + amount)
	if err != nil {
		return err
	}

	p.TotalDemandCount = newTotalDemand
	return nil
}

func (p *Product) RemoveCampaign() {
	p.Campaign = nil
	p.UpdatePrice(p.InititalPrice.Value())
}

func (p *Product) Discount(pm valueobject.PriceManipulationLimit) {
	if p.Stock.Value() == 0 {
		return
	}

	sellCount := p.InititalStock.Value() - p.Stock.Value()
	salesRate := float64(sellCount) / float64(p.TotalDemandCount.Value()) * 100
	priceChange := (salesRate - 50) * float64(pm.Value()) / 50
	priceDiff := p.InititalPrice.Value() + priceChange

	p.UpdatePrice(priceDiff)
}
