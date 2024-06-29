package storage

import (
	"encoding/json"
	"os"

	"github.com/ChebuRashkaRF/urlshortener/internal/models"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteRecord(record *models.URLRecord) error {
	return p.encoder.Encode(record)
}

func (p *Producer) Close() error {
	return p.file.Close()
}
