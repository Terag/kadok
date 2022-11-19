package cache

import "time"

type Cache interface {
	Get(key string) ([]byte, bool, error)
	Put(key string, value []byte, lifetime time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
}
