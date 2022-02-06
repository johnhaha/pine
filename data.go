package pine

type CachedData[T any] struct {
	Data *T
	CacheLife
}

func NewCachedData[T any](lifeTime int64) *CachedData[T] {
	return &CachedData[T]{CacheLife: CacheLife{Lifetime: lifeTime}}

}

func (data *CachedData[T]) Set(d *T) {
	data.Data = d
	data.Update()
}

func (data *CachedData[T]) Get(onExpire func() (*T, error)) (*T, error) {
	if data.Expired() {
		d, err := onExpire()
		if err != nil {
			return nil, err
		}
		data.Set(d)
		return d, nil
	}
	return data.Data, nil
}

func (data *CachedData[T]) Clear() {
	data.Data = nil
}
