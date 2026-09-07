package storage

import "sync"

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func New() *Store {
	return &Store{data: make(map[string]string)}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) SetNX(key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; ok {
		return false
	}
	s.data[key] = value
	return true
}

func (s *Store) Del(keys ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, k := range keys {
		if _, ok := s.data[k]; ok {
			delete(s.data, k)
			n++
		}
	}
	return n
}

func (s *Store) Exists(keys ...string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for _, k := range keys {
		if _, ok := s.data[k]; ok {
			n++
		}
	}
	return n
}
