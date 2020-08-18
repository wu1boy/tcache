package engine

import (
	"hash/adler32"
	"sync"
	"tcache/conf"
	"time"
)

type Cache interface {
	//size 是⼀一个字符串串。⽀支持以下参数: 1KB，100KB，1MB，2MB，1GB 等
	SetMaxMemory(size string) bool
	// 设置⼀一个缓存项，并且在expire时间之后过期
	Set(key string, val interface{}, expire time.Duration)
	// 获取⼀一个值
	Get(key string) (interface{}, bool)
	// 删除一个值
	Del(key string) bool
	// 检测⼀一个值 是否存在
	Exists(key string) bool
	// 情况所有值
	Flush() bool
	// 返回所有的key 多少
	Keys() int64
}

type Value struct {
	v   []byte
	ttl int64
}

type MemoryData struct {
	data   map[uint32]map[string]Value
	mutex  []sync.RWMutex
	config *conf.Config
	Stat
}

func (m *MemoryData) SetMaxMemory(size string) bool {

	return true
}

/**
 * 设置key
 */
func (m *MemoryData) Set(key string, val interface{}, expire time.Duration) {
	//判断是否超过了maxmemory
	//todo

	//根据key mod 看落入哪个index.
	index := adler32.Checksum([]byte(key)) % m.config.Shark

	m.mutex[index].Lock()
	defer m.mutex[index].Unlock()

	tmp, exist := m.data[index][key]

	if exist {
		m.del(key, tmp.v)
	}

	value := Value{
		v:   val.([]byte),
		ttl: time.Now().Unix() + int64(expire.Seconds()),
	}
	m.data[index][key] = value
	m.del(key, value.v)
}

/**
 * 获取key
 */
func (m *MemoryData) Get(key string) (interface{}, bool) {

	index := adler32.Checksum([]byte(key)) % m.config.Shark

	m.mutex[index].Lock()
	defer m.mutex[index].Unlock()

	data, exist := m.data[index][key]

	if !exist {
		return nil,false
	}

	//秒级别ttl
	if data.ttl < time.Now().Unix(){
		//数据超时已失效,删除数据
		delete(m.data[index], key)
		return nil,false
	}

	return data.v, true
}

/**
 * 删除key
 */
func (m *MemoryData) Del(key string) bool {
	index := adler32.Checksum([]byte(key)) % m.config.Shark

	m.mutex[index].Lock()
	defer m.mutex[index].Unlock()

	tmp, exist := m.data[index][key]

	if !exist {
		return false
	}

	delete(m.data[index], key)
	m.del(key, tmp.v)

	return true
}

/**
 * 查看key是否存在
 */
func (m *MemoryData) Exists(key string) bool {
	index := adler32.Checksum([]byte(key)) % m.config.Shark

	m.mutex[index].RLock()
	defer m.mutex[index].RUnlock()

	data, exist := m.data[index][key]

	if !exist {
		return false
	}

	//秒级别ttl
	if data.ttl < time.Now().Unix(){
		//数据超时已失效
		return false
	}

	return true
}

/**
 * 清空缓存
 */
func (m *MemoryData) Flush() bool {
	mutex := sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	//清空map. golang 1.11版本后可以使用这个方法. go会把for循环删除优化为一个runtime.mapclear调用,高效删除
	//1.11前的版本使用make 重新分配内存
	//m.data = make(map[uint32]map[string]Value)

	for index, _ := range m.data {
		delete(m.data, index)
	}

	return true
}

/**
 * 返回key的数量
 */
func (m *MemoryData) Keys() int64 {
	mutex := sync.RWMutex{}
	mutex.RLock()
	defer mutex.RUnlock()

	return m.Stat.Count
}

/**
 * 创建内存缓存对象
 */
func MemoryEngineNew(config *conf.Config) *MemoryData {
	return &MemoryData{
		data:   make(map[uint32]map[string]Value),
		mutex:  []sync.RWMutex{},
		config: config,
		Stat:   Stat{},
	}
}
