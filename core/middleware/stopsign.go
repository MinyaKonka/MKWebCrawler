package middleware

import (
	"fmt"
	"sync"
)

// 停止信号的接口
type MKStopSign interface {

	// 发出停止信号。如果先前已发出过停止信号，则返回false
	Stop() bool

	// 判断停止信号是否已被发出
	Signed() bool

	// 重置停止信号。相当于收回停止信号，并清除所有的停止信号处理记录
	Reset()

	// 处理停止信号。
	// @Param code: 停止信号的处理方。该代码会出现在停止信号的处理记录中
	Deal(code string)

	// 获取某一个停止信号处理方的代号。该代号会出现在停止信号的处理记录中。
	DealCount(code string) uint32

	// 获取停止信号被处理的总计数
	DealTotal() uint32

	// 获取摘要信息。其中应该包含所有的停止信号处理记录。
	Summary() string
}

// 创建停止信号
func NewStopSign() MKStopSign {
	sign := &mk_stopSign{
		dealCountMap: make(map[string]uint32),
	}

	return sign
}

type mk_stopSign struct {
	mutex        sync.RWMutex      // 读写锁
	signed       bool              // 表示信号是否已发出的标志位
	dealCountMap map[string]uint32 // 处理计数的字典
}

func (sign *mk_stopSign) Sign() bool {

	sign.mutex.Lock()
	defer sign.mutex.Unlock()

	if sign.signed {
		return false
	}

	sign.signed = true

	return true
}

func (sign *mk_stopSign) Signed() bool {

	return sign.signed
}

func (sign *mk_stopSign) Reset() {

	sign.mutex.Lock()
	defer sign.mutex.Unlock()

	sign.signed = false
	sign.dealCountMap = make(map[string]uint32)
}

func (sign *mk_stopSign) Deal(code string) {

	sign.mutex.Lock()
	defer sign.mutex.Unlock()

	if !sign.signed {
		return
	}

	if _, ok := sign.dealCountMap[code]; !ok {
		sign.dealCountMap[code] = 1
	} else {
		sign.dealCountMap[code] += 1
	}
}

func (sign *mk_stopSign) DealCount(code string) uint32 {

	sign.mutex.Lock()
	defer sign.mutex.Unlock()

	return sign.dealCountMap[code]
}

func (sign *mk_stopSign) DealTotal() uint32 {

	sign.mutex.Lock()
	defer sign.mutex.Unlock()

	var total uint32
	for _, v := range sign.dealCountMap {
		total += v
	}

	return total
}

func (sign *mk_stopSign) Summary() string {

	if sign.signed {
		return fmt.Sprintf("signed: true, dealCount: %v", sign.dealCountMap)
	} else {
		return "signed: false"
	}
}
