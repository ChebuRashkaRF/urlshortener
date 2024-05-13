package storage

import (
	"sync"
)

type URLStorage struct {
	URLMap   map[string]string
	URLMapMx sync.RWMutex
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		URLMap: make(map[string]string),
	}
}

func (s *URLStorage) Get(key string) (string, bool) {
	s.URLMapMx.RLock()
	defer s.URLMapMx.RUnlock()
	url, ok := s.URLMap[key]
	return url, ok
}

func (s *URLStorage) Set(key, value string) {
	s.URLMapMx.Lock()
	defer s.URLMapMx.Unlock()
	s.URLMap[key] = value
}

func (s *URLStorage) GetURLMap() map[string]string {
	s.URLMapMx.RLock()
	defer s.URLMapMx.RUnlock()

	mapCopy := make(map[string]string, len(s.URLMap))
	for key, val := range s.URLMap {
		mapCopy[key] = val
	}

	return mapCopy
}
