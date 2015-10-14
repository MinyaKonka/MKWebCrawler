package middleware

import (
	"math"
	"sync"
)

// 32位ID生成器的接口类型
type MKIDGenerator32 interface {
	GetUint32() uint32 // 获取一个uint32类型的ID
}

// 创建32位ID生成器
func NewIDGenerator() MKIDGenerator32 {
	return &mk_IDGenerator32{}
}

// 32位ID生成器的实现类型
type mk_IDGenerator32 struct {
	sn    uint32     // 当前的ID
	ended bool       // 前一个ID是否已经为其类型所能表示的最大值
	mutex sync.Mutex // 互斥锁
}

func (generator *mk_IDGenerator32) GetUint32() uint32 {
	generator.mutex.Lock()
	defer generator.mutex.Unlock()

	if generator.ended {
		defer func() {
			generator.ended = false
		}()

		generator.sn = 0
		return generator.sn
	}

	id := generator.sn
	if id < math.MaxUint32 {
		generator.sn++
	} else {
		generator.ended = true
	}

	return id
}

// 64位ID生成器接口
type MKIDGenerator64 interface {
	GetUint64() uint64 // 获得一个uint64类型的ID
}

// 创建64位ID生成器
func NewIDGenerator64() MKIDGenerator64 {
	return &mk_IDGenerator64{}
}

// 64位ID生成器的实现类型
type mk_IDGenerator64 struct {
	base       mk_IDGenerator32 // 基本的ID生成器
	cycleCount uint64           // 基于uint32类型的取值范围的周期计算
}

func (generator *mk_IDGenerator64) GetUint64() int64 {
	var id uint64

	if generator.cycleCount%2 == 1 {
		id += math.MaxUint32
	}

	id32 := generator.base.GetUint32()
	if id32 == math.MaxUint32 {
		generator.cycleCount++
	}

	id += uint64(id32)

	return id
}
