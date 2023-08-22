package valueobject

import "errors"

var (
	ErrPriceManipulationLimitLessThanZero = errors.New("PriceManipulationLimit can not be less than zero or equal to zero")
)

type PriceManipulationLimit struct {
	value int
}

func NewPriceManipulationLimit(value int) (PriceManipulationLimit, error) {
	if value <= 0 {
		return PriceManipulationLimit{}, ErrPriceManipulationLimitLessThanZero
	}

	return PriceManipulationLimit{value: value}, nil
}

func (p PriceManipulationLimit) Value() int {
	return p.value
}

func (p PriceManipulationLimit) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	priceManipulationLimit, ok := value.(PriceManipulationLimit)
	if !ok {
		return false
	}

	return p.value == priceManipulationLimit.value
}
