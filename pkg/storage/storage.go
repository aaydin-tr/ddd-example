package storage

import "sync"

type Storage[T any] struct {
	datas map[string]T
	mu    sync.RWMutex
}

func New[T any]() *Storage[T] {
	return &Storage[T]{
		datas: make(map[string]T),
	}
}

func (s *Storage[T]) Set(key string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.datas[key] = value
}

func (s *Storage[T]) Get(key string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.datas[key]
	return value, ok
}

func (s *Storage[T]) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.datas, key)
}

func (s *Storage[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.datas)
}

func (s *Storage[T]) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.datas))
	for key := range s.datas {
		keys = append(keys, key)
	}
	return keys
}

func (s *Storage[T]) Values() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	values := make([]T, 0, len(s.datas))
	for _, value := range s.datas {
		values = append(values, value)
	}
	return values
}

func (s *Storage[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.datas = make(map[string]T)
}
