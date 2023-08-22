package valueobject

import (
	"errors"
)

var (
	ErrDurationLessThanZero = errors.New("Duration can not be less than zero or equal to zero")
)

type Duration struct {
	value int
}

func NewDuration(value int) (Duration, error) {
	if value <= 0 {
		return Duration{}, ErrDurationLessThanZero
	}

	return Duration{value: value}, nil
}

func (d Duration) Value() int {
	return d.value
}

func (d Duration) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	duration, ok := value.(Duration)
	if !ok {
		return false
	}

	return d.value == duration.value
}
