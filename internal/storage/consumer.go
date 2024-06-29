package storage

import (
	"encoding/json"
	"os"

	"github.com/ChebuRashkaRF/urlshortener/internal/models"
)

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadRecord() (*models.URLRecord, error) {
	record := &models.URLRecord{}

	if err := c.decoder.Decode(&record); err != nil {
		return nil, err
	}

	return record, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
