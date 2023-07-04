package pine

import (
	"context"
	"encoding/json"
	"time"

	"github.com/allegro/bigcache/v3"
)

// var (
// 	lifeHeap = &DataLifeHeap{}
// 	vault    = make(map[string]any)
// )

// var mtx sync.RWMutex

// type DataLife struct {
// 	ID       string
// 	LifeTime time.Time
// }

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

// type DataLifeHeap []DataLife

// func (h DataLifeHeap) Len() int           { return len(h) }
// func (h DataLifeHeap) Less(i, j int) bool { return h[i].LifeTime.Before(h[j].LifeTime) }
// func (h DataLifeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// func (h *DataLifeHeap) Push(x any) {
// 	// Push and Pop use pointer receivers because they modify the slice's length,
// 	// not just its contents.
// 	*h = append(*h, x.(DataLife))
// }

// func (h *DataLifeHeap) Pop() any {
// 	old := *h
// 	n := len(old)
// 	x := old[n-1]
// 	*h = old[0 : n-1]
// 	return x
// }

// get from cache if exists, store cache with execution if not exists
func Get[T any](id string, withCache bool, execution func() (T, error)) (*CacheData[T], error) {

	//read data from cache
	// mtx.RLock()
	if withCache {
		if res, err := Cache.Get(id); err == nil {
			data := new(CacheData[T])
			err := data.Decode(res)
			if err != nil {
				Cache.Delete(id)
				return nil, err
			}
			data.FromCache = true
			return data, nil
		}
	}
	// if v, ok := vault[id]; ok {
	// 	mtx.RUnlock()
	// 	return &CacheData[T]{
	// 		Data:      v.(T),
	// 		FromCache: true,
	// 	}, nil
	// }
	// mtx.RUnlock()

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

	//write new data to cache
	// mtx.Lock()
	// vault[id] = res
	// registerLife(DataLife{
	// 	ID:       id,
	// 	LifeTime: time.Now().Add(lifeTime),
	// })
	// mtx.Unlock()

	return cache, nil
}

// remove cache with id
func Remove(id string) {
	// mtx.Lock()
	// delete(vault, id)
	// mtx.Unlock()
	Cache.Delete(id)
}

// func registerLife(life DataLife) {
// 	heap.Push(lifeHeap, life)
// }

// normally 10*time.Minute for eviction
func InitCache(ctx context.Context, eviction time.Duration) error {
	var err error
	Cache, err = bigcache.New(ctx, bigcache.DefaultConfig(eviction))
	if err != nil {
		return err
	}
	return nil
	// go func() {
	// 	for {
	// 		if len(*lifeHeap) > 0 && (*lifeHeap)[0].LifeTime.Before(time.Now()) {
	// 			mtx.Lock()
	// 			tm := heap.Pop(lifeHeap)
	// 			delete(vault, tm.(DataLife).ID)
	// 			mtx.Unlock()
	// 			log.Println("data", tm.(DataLife).ID, "is removed from cache")
	// 		}
	// 		time.Sleep(step)
	// 	}
	// }()
	// <-ctx.Done()
}
