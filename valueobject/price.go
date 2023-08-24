package valueobject

import "errors"

var (
	ErrPriceMustBePositive = errors.New("Price must be positive")
)

type Price struct {
	value float64
}

func NewPrice(value float64) (Price, error) {
	if value <= 0 {
		return Price{}, ErrPriceMustBePositive
	}

	return Price{value: value}, nil
}

func (p Price) Value() float64 {
	return p.value
}

func (p Price) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	price, ok := value.(Price)
	if !ok {
		return false
	}

	return p.value == price.value
}
