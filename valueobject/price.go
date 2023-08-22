package valueobject

import "errors"

var (
	ErrPriceMustBePositive = errors.New("Price must be positive")
)

type Price struct {
	value float32
}

func NewPrice(value float32) (Price, error) {
	if value < 0 {
		return Price{}, ErrPriceMustBePositive
	}

	return Price{value: value}, nil
}

func (p Price) Value() float32 {
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
