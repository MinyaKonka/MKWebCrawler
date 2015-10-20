package itempipeline

import (
	base "core/base"
	"errors"
	"fmt"
	"sync/atomic"
)

// 条目处理管理的接口类型
type MKItemPipeline interface {
	// 发送条目
	Send(item base.MKItem) []error

	// 返回一个布尔值，表示当前条目处理管道是否的快速失败的。
	// 这里的快速失败是指：只要对某个条目的处理流程在某个步骤上出错，
	// 那么条目处理管道就会忽略掉后续的所有处理步骤并报告错误
	FailFast() bool

	// 设置是否快速失败
	SetFailFast(failFast bool)

	// 获得已发送、已接受和已处理的条目的计数值
	// 更确切的说，作为结果值的切片总会有三个元素值，这三个值分别代表前述的三个计数
	Count() []uint64

	// 获取正在被处理的条目的数量
	ProcessingNumber() uint64

	// 获取摘要信息
	Summary() string
}

// 创建条目处理管道
func NewItemPipeline(itemProcessors []MKProcessItem) MKItemPipeline {
	if itemProcessors == nil {
		panic(errors.New(fmt.Sprintln("无效的条目处理器列表！")))
	}

	innerItemProcessors := make([]MKItemPipeline, 0)

	for i, itemProcessor := range itemProcessors {
		if itemProcessor == nil {
			panic(errors.New(fmt.Sprintf("无效的条目处理器[%d]\n"), i))
		}

		innerItemProcessors = append(innerItemProcessors, itemProcessor)
	}

	return &mk_itemPipeline{itemProcessors: itemProcessors}
}

type mk_itemPipeline struct {
	itemProcessors   []MKProcessItem // 条目处理器的列表
	failFast         bool            // 表示处理是否需要快速失败的标志位
	sent             uint64          // 已被发送的条目数量
	accepted         uint64          // 已被接受的条目数量
	processed        uint64          // 已被处理的条目数量
	processingNumber uint64          // 正在被处理的条目的数量
}

func (pipeline *mk_itemPipeline) Send(item base.MKItem) []error {
	atomic.AddUint64(&pipeline.processingNumber, 1)
	defer atomic.AddUint64(&pipeline.processingNumber, ^uint64(0))

	atomic.AddUint64(&pipeline.sent, 1)
	errs := make([]error, 0)
	if item == nil {
		errs = append(errs, errors.New("无效条目！"))
		return errs
	}

	atomic.AddUint64(&pipeline.accepted, 1)

	var currentItem base.MKItem = item
	for _, itemProcessor := range pipeline.itemProcessors {
		processedItem, err := itemProcessor(currentItem)
		if err != nil {
			errs = append(errs, err)
			if pipeline.failFast {
				break
			}
		}

		if processedItem != nil {
			currentItem = processedItem
		}
	}

	atomic.AddUint64(&pipeline.processed, 1)

	return errs
}

func (pipeline *mk_itemPipeline) FailFast() bool {
	return pipeline.failFast
}

func (pipeline *mk_itemPipeline) SetFailFast(failFast bool) {
	pipeline.failFast = failFast
}

func (pipeline *mk_itemPipeline) Count() []uint64 {
	counts := make([]uint64, 3)
	counts[0] = atomic.LoadUint64(&pipeline.sent)
	counts[1] = atomic.LoadUint64(&pipeline.accepted)
	counts[2] = atomic.LoadUint64(&pipeline.processed)
	return counts
}

func (pipeline *mk_itemPipeline) ProcessingNumber() uint64 {
	return atomic.LoadUint64(&pipeline.processingNumber)
}

var summaryTemplate = "failFast: %v, processorNumber: %d, " +
	" send: %d, accepted: %d, processed: %d, processingNumber: %d"

func (pipeline *mk_itemPipeline) Summary() string {
	counts := pipeline.Count()
	summary := fmt.Sprintf(summaryTemplate,
		pipeline.failFast,
		len(pipeline.itemProcessors),
		counts[0],
		counts[1],
		counts[2],
		pipeline.ProcessingNumber())

	return summary
}
