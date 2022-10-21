package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// 一致性哈希算法
type Map struct {
	// 传参数的时候好看一点？ - 定义hash算法
	hash Hash
	// 虚拟节点倍数 - 去创建几个与之对应的虚拟节点
	replicas int
	// 哈希环 - 存储所有虚拟节点的映射关系
	keys []int
	// [虚拟节点]: 真实节点
	hashMap map[int]string
}
// New 可以自定义虚拟节点倍数和Hash函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash: fn,
		hashMap: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}


func (m *Map) Add(keys ...string) {
	// 对于每一个key, 对应创建m.replicas虚拟节点
	for _, key := range keys {
		for i:=0; i < m.replicas; i ++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 记录所有虚拟节点名称
			m.keys = append(m.keys, hash)
			// 这些虚拟节点都指向一个Key
			m.hashMap[hash] = key
		}
	}
	// 哈希环 排序
	sort.Ints(m.keys)
}

// key应该由哪个节点处理
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	// 计算key的哈希值
	hash := int(m.hash([]byte(key)))
	// 找到第一个 i 使得m.keys[i] >= hash
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}