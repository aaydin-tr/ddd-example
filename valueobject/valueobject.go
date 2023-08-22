package valueobject

type ValueObject interface {
	Equals(value ValueObject) bool
}
