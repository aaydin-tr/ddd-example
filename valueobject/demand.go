package valueobject

import "errors"

var (
	ErrDemandCannotBeNegative = errors.New("Demand cannot be negative")
)

type Demand struct {
	value int
}

func NewDemand(value int) (Demand, error) {
	if value < 0 {
		return Demand{}, ErrDemandCannotBeNegative
	}

	return Demand{value: value}, nil
}

func (d Demand) Value() int {
	return d.value
}

func (n Demand) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	demand, ok := value.(Demand)
	if !ok {
		return false
	}

	return n.value == demand.value
}
