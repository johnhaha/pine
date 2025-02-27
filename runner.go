package pine

import (
	"context"
	"encoding/json"
	"time"

	"github.com/allegro/bigcache/v3"
)

var Cache *bigcache.BigCache

type CacheData[T any] struct {
	Data      T
	FromCache bool
}

func (data *CacheData[T]) Encode() ([]byte, error) {
	res, err := json.Marshal(data)
	return res, err
}

func (data *CacheData[T]) Decode(cache []byte) error {
	err := json.Unmarshal(cache, data)
	return err
}

// get from cache if exists, store cache with execution if not exists
func Get[T any](id string, withCache bool, execution func() (T, error)) (*CacheData[T], error) {

	if withCache {
		if res, err := Cache.Get(id); err == nil {
			data := new(CacheData[T])
			err = data.Decode(res)
			if err == nil {
				data.FromCache = true
				return data, nil
			}
		}
	}

	//execute if no cache exists
	res, err := execution()
	if err != nil {
		return nil, err
	}

	cache := &CacheData[T]{
		Data: res,
	}

	data, err := cache.Encode()
	if err != nil {
		return nil, err
	}
	Cache.Set(id, data)

	return cache, nil
}

// remove cache with id
func Remove(id string) {
	Cache.Delete(id)
}

// normally 10*time.Minute for eviction
func InitCache(ctx context.Context, eviction time.Duration) error {
	var err error
	Cache, err = bigcache.New(ctx, bigcache.DefaultConfig(eviction))
	if err != nil {
		return err
	}
	return nil

}
