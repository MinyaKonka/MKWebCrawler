package analyzer

import (
	middleware "core/middleware"
	"errors"
	"fmt"
	"reflect"
)

// 生成分析器的函数类型
type GenerateAnalyzer func() MKAnalyzer

// 分析器池的接口
type MKAnalyzerPool interface {
	Take() (MKAnalyzer, error)        // 从池中取出一个分析器
	Return(analyzer MKAnalyzer) error // 把一个分析器归还给池
	Total() uint32                    // 获得池的总容量
	Used() uint32                     // 获得正在被使用的分析器的数量
}

func NewAnalyzerPool(total uint32, generator GenerateAnalyzer) MKAnalyzerPool {

	etype := reflect.TypeOf(generator())
	generatorEntity := func() middleware.MKEntity {
		return generator()
	}

	pool, err := middleware.NewPool(total, etype, generatorEntity)
	if err != nil {
		return nil, err
	}

	downloaderPool := &mk_analyzerPool{pool: pool, etype: etype}

	return downloaderPool, nil
}

type mk_analyzerPool struct {
	pool  middleware.MKPool // 实体池
	etype reflect.Type      // 池内实体的类型
}

func (pool *mk_analyzerPool) Take() (MKAnalyzer, error) {

	entity, err := pool.pool.Take()
	if err != nil {
		return nil, err
	}

	analyzer, ok := entity.(MKAnalyzer)
	if !ok {
		errMsg := fmt.Sprintf("实体的类型不是%s\n", pool.etype)
		panic(errors.New(errMsg))
	}

	return analyzer, nil
}

func (pool *mk_analyzerPool) Return(analyzer MKAnalyzer) error {
	return pool.pool.Return(analyzer)
}

func (pool *mk_analyzerPool) Total() uint32 {
	return pool.pool.Total()
}

func (pool *mk_analyzerPool) Used() uint32 {
	return pool.pool.Used()
}
