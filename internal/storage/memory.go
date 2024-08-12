package storage

import (
	"sync"
)

type MemoryStorage struct {
	URLMap   map[string]string
	URLMapMx sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		URLMap: make(map[string]string),
	}
}

func (s *MemoryStorage) Get(key string) (string, bool) {
	s.URLMapMx.RLock()
	defer s.URLMapMx.RUnlock()
	url, ok := s.URLMap[key]
	return url, ok
}

func (s *MemoryStorage) Set(key, value string) error {
	s.URLMapMx.Lock()
	defer s.URLMapMx.Unlock()
	s.URLMap[key] = value
	return nil
}

func (s *MemoryStorage) GetURLMap() map[string]string {
	s.URLMapMx.RLock()
	defer s.URLMapMx.RUnlock()

	mapCopy := make(map[string]string, len(s.URLMap))
	for key, val := range s.URLMap {
		mapCopy[key] = val
	}

	return mapCopy
}

func (s *MemoryStorage) Close() error {
	return nil
}
