package types

type Storage[T any] interface {
	Set(key string, value T)
	Get(key string) (T, bool)
	Delete(key string)
	Len() int
	Keys() []string
	Values() []T
}
