package cache

import "time"

type NoCache struct{}

func NewNoCache() Cache {
	return NoCache{}
}

func (n NoCache) Get(_ string) (interface{}, bool) {
	return nil, false
}

func (n NoCache) Set(_ string, _ interface{}, _ time.Duration) {}

func (n NoCache) Delete(_ string) {}
