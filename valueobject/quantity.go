package valueobject

import "errors"

var (
	ErrQuantityMustBePositive = errors.New("Quantity must be positive")
)

type Quantity struct {
	value int
}

func NewQuantity(value int) (Quantity, error) {
	if value <= 0 {
		return Quantity{}, ErrQuantityMustBePositive
	}

	return Quantity{value: value}, nil
}

func (q Quantity) Value() int {
	return q.value
}

func (q Quantity) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	quantity, ok := value.(Quantity)
	if !ok {
		return false
	}

	return q.value == quantity.value
}
