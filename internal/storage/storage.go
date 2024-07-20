package storage

import (
	"github.com/ChebuRashkaRF/urlshortener/internal/logger"
	"github.com/ChebuRashkaRF/urlshortener/internal/models"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

type URLStorage struct {
	URLMap   map[string]string
	URLMapMx sync.RWMutex
	UUID     int
	producer *Producer
	consumer *Consumer
}

func NewURLStorage(filePath string) (*URLStorage, error) {
	// Создаем producer для записи в файл
	producer, err := NewProducer(filePath)
	if err != nil {
		return nil, err
	}

	// Создаем consumer для чтения из файла
	consumer, err := NewConsumer(filePath)
	if err != nil {
		return nil, err
	}

	// Создаем URLStorage
	storage := &URLStorage{
		URLMap:   make(map[string]string),
		UUID:     1,
		producer: producer,
		consumer: consumer,
	}

	// Загружаем данные из файла
	err = storage.LoadFromFile()
	if err != nil {
		return nil, err
	}

	return storage, nil
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
	s.saveToFile(key, value)
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

// saveToFile сохранение данных в файл
func (s *URLStorage) saveToFile(key, value string) {
	record := models.URLRecord{
		UUID:        strconv.Itoa(s.UUID),
		ShortURL:    key,
		OriginalURL: value,
	}
	err := s.producer.WriteRecord(&record)
	if err != nil {
		logger.Log.Error("Error save file", zap.Error(err))
		return
	}
	s.UUID++
}

// LoadFromFile загрузка данных из файла
func (s *URLStorage) LoadFromFile() error {
	maxUUID := 0
	for {
		record, err := s.consumer.ReadRecord()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
		s.URLMap[record.ShortURL] = record.OriginalURL

		currentUUID, err := strconv.Atoi(record.UUID)
		if err == nil && currentUUID > maxUUID {
			maxUUID = currentUUID
		}
	}

	s.UUID = maxUUID + 1

	return nil
}

func (s *URLStorage) Close() error {

	if err := s.producer.Close(); err != nil {
		logger.Log.Error("Error closing Producer", zap.Error(err))
		return err
	}

	if err := s.consumer.Close(); err != nil {
		logger.Log.Error("Error closing Consumer", zap.Error(err))
		return err
	}

	return nil
}
