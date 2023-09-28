package pine

import "golang.org/x/sync/singleflight"

var g singleflight.Group

func SingleRun[T any](key string, execution func() (T, error)) (*CacheData[T], error) {
	res, err, shared := g.Do(key, func() (interface{}, error) {
		return execution()
	})
	if err != nil {
		return nil, err
	}
	return &CacheData[T]{
		Data:      res.(T),
		FromCache: shared,
	}, nil
}
