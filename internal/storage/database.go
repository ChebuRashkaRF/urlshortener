package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type Database interface {
	Ping() error
	Close() error
}

type database struct {
	DB *sql.DB
}

func NewDatabase(dsn string) (Database, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &database{DB: db}, nil
}

func (d *database) Close() error {
	return d.DB.Close()
}

func (d *database) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return d.DB.PingContext(ctx)
}
