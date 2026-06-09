package cache

import (
	"fmt"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type CacheConfig struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

type CACHE_TYPE string

const (
	NO_CACHE_TYPE     CACHE_TYPE = "NO_CACHE"
	MEMORY_CACHE_TYPE CACHE_TYPE = "MEMORY"
)

func NewCache(cacheType CACHE_TYPE, config *CacheConfig) (Cache, error) {
	switch cacheType {
	case NO_CACHE_TYPE:
		return NewNoCache(), nil
	case MEMORY_CACHE_TYPE:
		return NewMemoryCache(config.DefaultExpiration, config.CleanupInterval), nil
	default:
		return nil, fmt.Errorf("cannot initialize %s, cache not implemented", cacheType)
	}
}
