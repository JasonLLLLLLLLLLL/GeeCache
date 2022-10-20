package GeeCache

import (
	lru "GeeCache/LRU"
	"sync"
)
// 带个锁，支持并发，再封装一下
// 所有对lru的访问都加上锁呗



type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	// ?
	cacheBytes int64
}
//
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// lazy load - 第一次使用的时候再初始化
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}
//
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}