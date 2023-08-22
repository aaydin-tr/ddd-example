package valueobject

import (
	"errors"
)

const (
	Active = "Active"
	Ended  = "Ended"
)

var (
	ErrStatusCannotBeEmpty = errors.New("Status cannot be empty")
	ErrStatusMustBeOneOf   = errors.New("Status must be one of 'Active', 'Ended'")
)

type Status struct {
	value string
}

func NewStatus(value string) (Status, error) {
	if value == "" {
		return Status{}, ErrStatusCannotBeEmpty
	}

	if value != Active && value != Ended {
		return Status{}, ErrStatusMustBeOneOf
	}

	return Status{value: value}, nil
}

func (s Status) Value() string {
	return s.value
}

func (s Status) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	status, ok := value.(Status)
	if !ok {
		return false
	}

	return s.value == status.value
}
