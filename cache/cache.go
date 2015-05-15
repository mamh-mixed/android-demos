package cache

import (
	"sync"
	"time"
)

const (
	NoExpiration = -1
)

type Item struct {
	Object     interface{}
	Expiration *time.Time
}

type Cache struct {
	mutex sync.RWMutex
	items map[string]*Item
	// ...
}

func (c *Cache) Set(k string, o interface{}, d time.Duration) {

	var e *time.Time

	// 有设过期时间
	if d > NoExpiration {
		t := time.Now().Add(d)
		e = &t
	}

	c.items[k] = &Item{
		Object:     o,
		Expiration: e,
	}
}

func (c *Cache) Get(k string) (interface{}, bool) {

	v, found := c.items[k]
	// TODO
	if !found || v.Expired() {
		// delete(c.items, k)
		return nil, false
	}
	return v.Object, true
}

// Expired 检查是否过期
func (i *Item) Expired() bool {

	if i.Expiration == nil {
		return false
	}

	return i.Expiration.Before(time.Now())
}

// New 创建一个新的cache
func New() *Cache {
	items := make(map[string]*Item)
	return &Cache{
		items: items,
	}
}
