package cache

import "time"

// Cache interface contains all behaviors for cache adapter.
type Cache interface {
	Put(key string, val interface{}, ttl time.Duration) error
	Get(key string, val interface{}) error
	Clear() error
}
