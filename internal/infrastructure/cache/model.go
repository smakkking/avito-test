package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	inmemorycache "github.com/patrickmn/go-cache"
)

type InMemoryCache struct {
	cacheImpl *inmemorycache.Cache
}

func NewCache(expirationTime time.Duration) *InMemoryCache {
	return &InMemoryCache{
		cacheImpl: inmemorycache.New(expirationTime, 24*time.Hour),
	}
}

func (c *InMemoryCache) GetUserBanner(ctx context.Context, tagID int, featureID int) (interface{}, error) {
	if value, found := c.cacheImpl.Get(fmt.Sprint(tagID) + ":" + fmt.Sprint(featureID)); found {
		return value, nil
	}

	return struct{}{}, errors.New("cache miss")
}

func (c *InMemoryCache) SaveUserBanner(ctx context.Context, tagID int, featureID int, value interface{}) {
	c.cacheImpl.SetDefault(fmt.Sprint(tagID)+":"+fmt.Sprint(featureID), value)
}
