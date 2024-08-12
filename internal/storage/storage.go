package storage

// Storage представляет собой интерфейс для различных хранилищ
type Storage interface {
	Get(key string) (string, bool)
	Set(key, value string) error
	Close() error
}
