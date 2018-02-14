package nut

import (
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
)

var (
	_cache     cache.Cache
	_cacheOnce sync.Once
)

// Cache get cache adapter.
func Cache() cache.Cache {
	_cacheOnce.Do(func() {
		var err error
		_cache, err = cache.NewCache(
			beego.AppConfig.String("cachedriver"),
			beego.AppConfig.String("cachesource"),
		)
		if err != nil {
			beego.Error(err)
		}
	})

	return _cache
}
