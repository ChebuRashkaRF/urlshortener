package storage

import (
	"context"
	"database/sql"
	"github.com/ChebuRashkaRF/urlshortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

type DatabaseStorage struct {
	DB *sql.DB
}

func NewDatabaseStorage(dsn string) (*DatabaseStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	storage := &DatabaseStorage{DB: db}

	err = storage.init()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (d *DatabaseStorage) init() error {
	_, err := d.DB.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_url TEXT UNIQUE NOT NULL,
			original_url TEXT NOT NULL
		)
	`)
	return err
}

func (d *DatabaseStorage) Get(key string) (string, bool) {
	var originalURL string
	err := d.DB.QueryRowContext(context.Background(), "SELECT original_url FROM urls WHERE short_url = $1", key).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Debug("Short URL not found in database", zap.String("short_url", key))
			return "", false
		}
		logger.Log.Error("Failed to fetch original URL from database", zap.Error(err), zap.String("short_url", key))
		return "", false
	}
	return originalURL, true
}

func (d *DatabaseStorage) Set(key, value string) error {
	_, err := d.DB.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING", key, value)
	if err != nil {
		logger.Log.Error("Failed to insert URL into database", zap.Error(err), zap.String("short_url", key), zap.String("original_url", value))
		return err
	}
	return nil
}

func (d *DatabaseStorage) Close() error {
	return d.DB.Close()
}

func (d *DatabaseStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return d.DB.PingContext(ctx)
}
