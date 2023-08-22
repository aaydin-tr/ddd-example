package entity

import (
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/google/uuid"
)

type Product struct {
	ID    uuid.UUID
	Code  valueobject.Code
	Price valueobject.Price
	Stock valueobject.Stock
}
