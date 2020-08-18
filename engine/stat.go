package engine

type Stat struct {
	//查询次数
	Count int64
	//key数量
	KeySize int64
	//空间使用情况
	ValueSize int64
	//未命中统计
	MissCounter int64
}

func (s *Stat) add(key string, data []byte) {
	s.Count += 1
	s.KeySize += int64(len(key))
	s.ValueSize += int64(len(data))
}

func (s *Stat) del(key string, data []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(key))
	s.ValueSize -= int64(len(data))
}

/**
 *查询缓存未命中
 */
func (s *Stat) miss() {
	s.MissCounter += 1
}
