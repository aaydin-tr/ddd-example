package valueobject

import "errors"

var (
	ErrNameCannotBeEmpty = errors.New("Name cannot be empty")
)

type Name struct {
	value string
}

func NewName(value string) (Name, error) {
	if value == "" {
		return Name{}, ErrNameCannotBeEmpty
	}

	return Name{value: value}, nil
}

func (n Name) Value() string {
	return n.value
}

func (n Name) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	name, ok := value.(Name)
	if !ok {
		return false
	}

	return n.value == name.value
}
