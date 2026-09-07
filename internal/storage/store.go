package storage

import (
	"errors"
	"strconv"
	"sync"
)

var errNotAnInteger = errors.New("value is not an integer or out of range")

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

func (s *Store) Incr(key string) (int64, error) {
	return s.incrBy(key, 1)
}

func (s *Store) Decr(key string) (int64, error) {
	return s.incrBy(key, -1)
}

func (s *Store) incrBy(key string, delta int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur := int64(0)
	if v, ok := s.data[key]; ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errNotAnInteger
		}
		cur = n
	}
	cur += delta
	s.data[key] = strconv.FormatInt(cur, 10)
	return cur, nil
}

func (s *Store) Append(key, value string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] += value
	return len(s.data[key])
}

func (s *Store) StrLen(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data[key])
}

func (s *Store) MGet(keys ...string) []*string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*string, len(keys))
	for i, k := range keys {
		if v, ok := s.data[k]; ok {
			vc := v
			out[i] = &vc
		}
	}
	return out
}

func (s *Store) MSet(pairs map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range pairs {
		s.data[k] = v
	}
}
