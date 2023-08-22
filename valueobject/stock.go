package valueobject

import "errors"

var (
	ErrStockMustBePositive = errors.New("Price must be positive")
)

type Stock struct {
	value int
}

func NewStock(value int) (Stock, error) {
	if value < 0 {
		return Stock{}, ErrStockMustBePositive
	}

	return Stock{value: value}, nil
}

func (s Stock) Value() int {
	return s.value
}

func (s Stock) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	stock, ok := value.(Stock)
	if !ok {
		return false
	}

	return s.value == stock.value
}
