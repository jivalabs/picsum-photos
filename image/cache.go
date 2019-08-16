package image

import (
	"context"

	"github.com/jivalabs/picsum-photos/cache"
	"github.com/jivalabs/picsum-photos/storage"
)

// Cache is an image cache
type Cache = cache.Auto

// NewCache instantiates a new cache
func NewCache(cacheProvider cache.Provider, storageProvider storage.Provider) *Cache {
	return &Cache{
		Provider: cacheProvider,
		Loader: func(ctx context.Context, key string) (data []byte, err error) {
			return storageProvider.Get(ctx, key)
		},
	}
}
