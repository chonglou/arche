package redis

// https://www.iana.org/assignments/uri-schemes/prov/redis

import (
	"encoding/json"
	"time"

	"github.com/chonglou/arche/web/cache"
	"github.com/garyburd/redigo/redis"
)

// New open a redis cache
func New(pool *redis.Pool, prefix string) cache.Cache {
	return &Cache{
		prefix: prefix,
		pool:   pool,
	}
}

// Cache redis cache adapter.
type Cache struct {
	pool   *redis.Pool
	prefix string
}

// Put cached value with key and expire time.
func (p *Cache) Put(key string, val interface{}, ttl time.Duration) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c := p.pool.Get()
	defer c.Close()
	_, err = c.Do("SET", p.key(key), buf, "EX", int(ttl/time.Second))
	return err
}

// Get get cached value by key.
func (p *Cache) Get(key string, val interface{}) error {
	c := p.pool.Get()
	defer c.Close()
	buf, err := redis.Bytes(c.Do("GET", p.key(key)))
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, val)
}

// Status list all items
func (p *Cache) Status() (map[string]int, error) {
	c := p.pool.Get()
	defer c.Close()
	keys, err := redis.Strings(c.Do("KEYS", p.key("*")))
	if err != nil {
		return nil, err
	}
	items := make(map[string]int)
	for _, k := range keys {
		ttl, err := redis.Int(c.Do("TTL", k))
		if err != nil {
			return nil, err
		}
		items[k] = ttl
	}
	return items, nil
}

// Clear clear all cache.
func (p *Cache) Clear() error {
	c := p.pool.Get()
	defer c.Close()
	keys, err := redis.Strings(c.Do("KEYS", p.key("*")))
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	var args []interface{}
	for _, k := range keys {
		args = append(args, k)
	}
	_, err = c.Do("DEL", args...)
	return err
}

func (p *Cache) key(k string) string {
	return p.prefix + k
}
