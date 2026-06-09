package cache

import (
	"time"

	go_cache "github.com/patrickmn/go-cache"
)

func NewMemoryCache(defaultExpiration, cleanupInterval time.Duration) Cache {
	return go_cache.New(defaultExpiration, cleanupInterval)
}
