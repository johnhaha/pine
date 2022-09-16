package pine

import (
	"container/heap"
	"context"
	"log"
	"sync"
	"time"
)

var (
	lifeHeap = &DataLifeHeap{}
	vault    = make(map[string]any)
)

var mtx sync.RWMutex

type DataLife struct {
	ID       string
	LifeTime time.Time
}

type CacheData[T any] struct {
	Data      T
	FromCache bool
}

type DataLifeHeap []DataLife

func (h DataLifeHeap) Len() int           { return len(h) }
func (h DataLifeHeap) Less(i, j int) bool { return h[i].LifeTime.Before(h[j].LifeTime) }
func (h DataLifeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *DataLifeHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(DataLife))
}

func (h *DataLifeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Get[T any](id string, lifeTime time.Duration, execution func() (T, error)) (*CacheData[T], error) {

	//read data from cache
	mtx.RLock()
	if v, ok := vault[id]; ok {
		mtx.RUnlock()
		return &CacheData[T]{
			Data:      v.(T),
			FromCache: true,
		}, nil
	}
	mtx.RUnlock()

	//execute if no cache exists
	res, err := execution()
	if err != nil {
		return nil, err
	}

	//write new data to cache
	mtx.Lock()
	vault[id] = res
	registerLife(DataLife{
		ID:       id,
		LifeTime: time.Now().Add(lifeTime),
	})
	mtx.Unlock()

	return &CacheData[T]{
		Data:      res,
		FromCache: false,
	}, nil
}

func registerLife(life DataLife) {
	heap.Push(lifeHeap, life)
}

func StartCleaner(ctx context.Context, step time.Duration) {
	log.Println("pine cleaner is running")
	go func() {
		for {
			if len(*lifeHeap) > 0 && (*lifeHeap)[0].LifeTime.Before(time.Now()) {
				mtx.Lock()
				tm := heap.Pop(lifeHeap)
				delete(vault, tm.(DataLife).ID)
				mtx.Unlock()
				log.Println("data", tm.(DataLife).ID, "is removed from cache")
			}
			time.Sleep(step)
		}
	}()
	<-ctx.Done()
}
