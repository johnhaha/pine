package pine

import "sync"

type CachedList[T any] struct {
	Data [][]T
	CacheLife
	Capacity int
	sync.RWMutex
}

func NewCachedList[T any](lifeTime int64, capacity int) *CachedList[T] {
	return &CachedList[T]{
		Data:      make([][]T, 0, capacity),
		Capacity:  capacity,
		CacheLife: CacheLife{Lifetime: lifeTime}}
}

func (list *CachedList[T]) Append(data []T, at int) {
	list.Lock()
	defer list.Unlock()
	if len(list.Data) != at {
		return
	}
	list.Data = append(list.Data, data)
	if len(list.Data) == 1 {
		list.Update()
	}
}

//clear cached list
func (list *CachedList[T]) Clear() {
	list.Data = make([][]T, 0, list.Capacity)
	list.UpdatedAt = 0
}

//get length of cached list
func (list *CachedList[T]) Length() int {
	return len(list.Data)
}

//return first element of cached list
func (list *CachedList[T]) First() []T {
	l := len(list.Data)
	if l == 0 {
		return nil
	}
	return list.Data[0]
}

//return last element of cached list
func (list *CachedList[T]) Last() []T {
	l := len(list.Data)
	if l == 0 {
		return nil
	}
	return list.Data[l-1]
}

//return last element of cached list
func (list *CachedList[T]) Count() (page int, num int) {
	var i int
	for _, n := range list.Data {
		i += len(n)
	}
	return len(list.Data), i
}

//get list from specific index
func (list *CachedList[T]) GetPage(page int, getData func(page int) ([]T, error)) (listData []T, err error) {
	if list.UpdatedAt > 0 && list.Expired() {
		list.Clear()
		return getData(page)
	}
	if page > list.Capacity {
		return getData(page)
	}
	listLen := len(list.Data)
	if listLen < page {
		d, err := getData(page)
		if err != nil {
			return nil, err
		}
		if listLen == page-1 {
			list.Append(d, page-1)
		}
		return d, nil
	}
	return list.Data[page-1], nil
}
