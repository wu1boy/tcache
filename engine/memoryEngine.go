package engine

import (
	"fmt"
	"hash/adler32"
	"math"
	"sync"
	"tcache/conf"
	"tcache/eliminate"
	"tcache/timewheel"
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
	data      map[uint32]map[string]Value
	mutex     []sync.RWMutex
	timewheel *timewheel.TimeWheel
	config    *conf.Config
	Stat
}

/**
 * 动态调整最大内存可能会导致数据被大量淘汰,这里放到配置文件内实现该功能.
 */
func (m *MemoryData) SetMaxMemory(size string) bool {
	//1KB，100KB，1MB，2MB，1GB
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

	if expire.Seconds() > 0 {
		//有设置过期时间,添加时间轮过期任务
		m.timewheel.AddTimer(expire, key, key)
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
		//查询未命中.更新统计
		m.Stat.miss()
		return nil, false
	}

	//秒级别ttl
	if data.ttl < time.Now().Unix() {
		//数据超时已失效,删除数据,定时器可以不删除
		delete(m.data[index], key)
		m.Stat.miss()
		return nil, false
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

	m.mutex[index].Lock()
	defer m.mutex[index].Unlock()

	data, exist := m.data[index][key]

	if !exist {
		return false
	}

	//秒级别ttl
	if data.ttl < time.Now().Unix() {
		//数据超时已失效,删除数据
		delete(m.data[index], key)
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
 * 定时任务
 */
func (m *MemoryData) Hz() {
	interval := int(math.Ceil(1000 / 10))

	tick := time.Tick(2 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("hello, tick")
		}
	}
}

/**
 * 内存达到指定最大内存时开始淘汰.
 * 这里可以定时触发，也可以手动触发.
 */
func (m *MemoryData) Eliminate() {
	var e eliminate.Eliminate

	switch m.config.EliminatePolic {
	case 1:
		//随机删除
		e = new(eliminate.Random)
	case 2:
		//lru算法
		e = new(eliminate.LRUCache)
	default:
		return
	}

	e.Remove()
}

/**
 * 创建内存缓存对象
 */
func MemoryEngineNew(config *conf.Config) *MemoryData {
	mem := &MemoryData{
		data:   make(map[uint32]map[string]Value),
		mutex:  []sync.RWMutex{},
		config: config,
		Stat:   Stat{},
	}

	mem.timewheel = timewheel.New(1*time.Second, 60, func(data interface{}) {
		//过期操作
		mem.Del(data.(string))
	})

	//时间轮.用于key过期和内存超限后淘汰.
	go mem.timewheel.Start()

	return mem
}
