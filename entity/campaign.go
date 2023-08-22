package entity

import (
	"time"

	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/google/uuid"
)

type Campaign struct {
	ID                     uuid.UUID
	Name                   valueobject.Name
	ProductCode            valueobject.Code
	Duration               valueobject.Duration
	PriceManipulationLimit valueobject.PriceManipulationLimit
	TargetSalesCount       valueobject.TargetSalesCount
	StartTime              time.Time
	TotalSales             int
	Status                 valueobject.Status
	// AverageItemPrice       valueobject.AverageItemPrice
}
