package pine

import "time"

type ListElement interface {
	GetUID() string
}

type CacheLife struct {
	Lifetime  int64
	UpdatedAt int64
}

func (life *CacheLife) Update() {
	life.UpdatedAt = time.Now().Unix()
}

func (life *CacheLife) Expired() bool {
	now := time.Now().Unix()
	return now-life.UpdatedAt > life.Lifetime
}
