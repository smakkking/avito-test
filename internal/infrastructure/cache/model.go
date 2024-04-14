package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	inmemorycache "github.com/patrickmn/go-cache"
	"github.com/smakkking/avito_test/internal/models"
)

const (
	cleanupInterval = 24 * time.Hour
)

type InMemoryCache struct {
	cacheImpl *inmemorycache.Cache
}

func NewCache(expirationTime time.Duration) *InMemoryCache {
	return &InMemoryCache{
		cacheImpl: inmemorycache.New(expirationTime, cleanupInterval),
	}
}

func (c *InMemoryCache) GetUserBanner(ctx context.Context, tagID, featureID int) (models.BannerContent, error) {
	if value, found := c.cacheImpl.Get(fmt.Sprint(tagID) + ":" + fmt.Sprint(featureID)); found {
		return value.(models.BannerContent), nil
	}

	return models.BannerContent{}, errors.New("cache miss")
}

func (c *InMemoryCache) SaveUserBanner(ctx context.Context, tagID, featureID int, value models.BannerContent) {
	c.cacheImpl.SetDefault(fmt.Sprint(tagID)+":"+fmt.Sprint(featureID), value)
}
