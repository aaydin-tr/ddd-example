package valueobject

import "errors"

var (
	ErrTargetSalesCountLessThanZero = errors.New("TargetSalesCount can not be less than zero or equal to zero")
)

type TargetSalesCount struct {
	value int
}

func NewTargetSalesCount(value int) (TargetSalesCount, error) {
	if value <= 0 {
		return TargetSalesCount{}, ErrTargetSalesCountLessThanZero
	}

	return TargetSalesCount{value: value}, nil
}

func (t TargetSalesCount) Value() int {
	return t.value
}

func (t TargetSalesCount) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	targetSalesCount, ok := value.(TargetSalesCount)
	if !ok {
		return false
	}

	return t.value == targetSalesCount.value
}
