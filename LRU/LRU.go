package LRU

import "container/list"

type Cache struct {
	// 具体存储数据
	ll *list.List
	// 最大容量
	maxBytes int64
	// 已用容量
	nbytes int64
	// 字符串: 双向链表的指针
	cache map[string]*list.Element
	// ？ 某个节点被删除的时候，执行的方法
	OnEvicted func(key string, value Value)
}

// 双向链表中的数据类型
// key的目的是淘汰首节点的时候，从字典删除对应的映射？
// 为啥不直接用map
type entry struct {
	key string
	value Value
}
// 允许值是实现了Val的任意类型
type Value interface {
	Len() int
}

// 实例化Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache{
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}
func (c *Cache) Get(key string) (value Value, ok bool) {
	// chace和list具体的作用？？
	// cache 中存在数据的话
	if ele, ok := c.cache[key]; ok {
		// 移动到链表头部
		c.ll.MoveToFront(ele)
		// 断言--
		kv := ele.Value.(*entry)

		return kv.value,  true
	}
	return
}
func  (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		// 从链表当中删除
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 从cache map 中删除
		delete(c.cache, kv.key)
		// 重新计算大小
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 执行回调函数? 回调-有时候可以去执行
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	// cache当中存在？更新键对应的值
	if ele, ok := c.cache[key]; ok {

		c.ll.MoveToFront(ele)

		kv := ele.Value.(*entry)
		// 更新新旧value的差值
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())

		kv.value = value
	} else {
		// 不存在，list增加新节点
		ele := c.ll.PushFront(&entry{key, value})
		// 增加映射关系
		c.cache[key] = ele
		// 更新长度
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}

