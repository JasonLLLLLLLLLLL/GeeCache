package LRU

import "container/list"

type Cache struct {
	// 具体存储数据
	l1 *list.List
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
	len() int
}

// 实例化Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache{
	return &Cache{
		maxBytes: maxBytes,
		l1: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted
	}
}
func (c *Cache) Get(key string) (value Value, ok bool) {
	// chace和list具体的作用？？
	// cache 中存在数据的话
	if ele, ok := c.cache[key]; ok {
		// 移动到链表头部
		c.l1.MoveToFront(ele)
		// 断言--
		kv := ele.Value.(*entry)

		return kv.value,  true
	}
	return
}
