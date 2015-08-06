package cache

import (
	"sync"
	"time"

	"github.com/CardInfoLink/quickpay/goconf"
)

// 默认缓存失效时间
var DefaultExpiration = time.Duration(goconf.Config.App.DefaultCacheTime)

const (
	NoExpiration = 0
)

// 维护应用中的缓存对象
var Client cachePool

func init() {
	caches := make(map[string]*Cache)
	Client.caches = caches
}

type Item struct {
	Object     interface{}
	Expiration *time.Time
}

type Cache struct {
	mutex sync.RWMutex
	items map[string]*Item
	// ...
}

type cachePool struct {
	caches map[string]*Cache
}

func (c *Cache) Set(k string, o interface{}, d time.Duration) {

	var e *time.Time

	// 有设过期时间
	if d > NoExpiration {
		t := time.Now().Add(d)
		e = &t
	}
	c.mutex.Lock()
	c.items[k] = &Item{
		Object:     o,
		Expiration: e,
	}
	c.mutex.Unlock()
}

func (c *Cache) Get(k string) (interface{}, bool) {

	c.mutex.RLock()
	v, found := c.items[k]
	c.mutex.RUnlock()
	// TODO
	if !found || v.Expired() {
		// delete(c.items, k)
		return nil, false
	}

	return v.Object, true
}

// Items 读取缓存对象
func (c *Cache) Items() map[string]*Item {

	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.items
}

// Clear 清空
func (c *Cache) Clear() {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	for k, _ := range c.items {
		delete(c.items, k)
	}
}

// Expired 检查是否过期
func (i *Item) Expired() bool {

	if i.Expiration == nil {
		return false
	}

	return i.Expiration.Before(time.Now())
}

// New 创建一个新的cache
// name app中缓存唯一
func New(name string) *Cache {
	items := make(map[string]*Item)

	cache := &Cache{
		items: items,
	}

	// 注册到全局缓存池中
	Client.Add(name, cache)
	return cache
}

// Add 往缓存池里增加一个
func (c *cachePool) Add(name string, cache *Cache) {
	c.caches[name] = cache
}

// Get 取一个
func (c *cachePool) Get(name string) (*Cache, bool) {

	if cache, ok := c.caches[name]; ok {
		return cache, ok
	}
	return nil, false
}
