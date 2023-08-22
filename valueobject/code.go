package valueobject

import "errors"

type Code struct {
	value string
}

var (
	ErrCodeIsRequired = errors.New("code is required")
)

func NewCode(value string) (Code, error) {
	if value == "" {
		return Code{}, ErrCodeIsRequired
	}

	return Code{value: value}, nil
}

func (c Code) Value() string {
	return c.value
}

func (c Code) Equals(value ValueObject) bool {
	if value == nil {
		return false
	}

	code, ok := value.(Code)
	if !ok {
		return false
	}

	return c.value == code.value
}
