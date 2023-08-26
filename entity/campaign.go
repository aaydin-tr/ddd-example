package entity

import (
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/google/uuid"
)

type Campaign struct {
	ID                     uuid.UUID
	Name                   valueobject.Name
	Product                *Product
	Duration               valueobject.Duration
	PriceManipulationLimit valueobject.PriceManipulationLimit
	TargetSalesCount       valueobject.TargetSalesCount
	Status                 valueobject.Status
	TotalSales             valueobject.Quantity
	AverageItemPrice       valueobject.Price
}

func (c *Campaign) IncreaseTotalSales(amount int) error {
	newTotalSales, err := valueobject.NewQuantity(c.TotalSales.Value() + amount)
	if err != nil {
		return err
	}

	c.TotalSales = newTotalSales
	return nil
}

func (c *Campaign) Close() error {
	closeStatus, err := valueobject.NewStatus(valueobject.Ended)
	if err != nil {
		return err
	}

	c.Status = closeStatus
	return nil
}

func (c *Campaign) IsActive() bool {
	return c.Status.Value() == valueobject.Active
}

func (c *Campaign) DecreaseDuration(duration int) error {
	decreaseDuration := c.Duration.Value() - duration
	if decreaseDuration < 0 {
		decreaseDuration = 0
	}

	newDuration, err := valueobject.NewDuration(decreaseDuration)
	if err != nil {
		return err
	}

	c.Duration = newDuration
	return nil
}

func (c *Campaign) RemainingTargetSalesCount(quantity int) int {
	return c.TargetSalesCount.Value() - (quantity + c.TotalSales.Value())
}

func (c *Campaign) UpdateAverageItemPrice(orderPrice float64, orderQuantity int) error {
	if c.TotalSales.Value() == 0 {
		return nil
	}

	orderRevenue := orderPrice * float64(orderQuantity)
	oldTotalRevenue := c.AverageItemPrice.Value() * float64(c.TotalSales.Value()-orderQuantity)

	newTotalRevenue := oldTotalRevenue + orderRevenue

	newAverageItemPrice, err := valueobject.NewPrice(newTotalRevenue / float64(c.TotalSales.Value()))
	if err != nil {
		return err
	}
	c.AverageItemPrice = newAverageItemPrice
	return nil
}
