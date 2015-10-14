package downloader

import (
	middleware "MKWebCrawler/middleware"
	"errors"
	"fmt"
	"reflect"
)

// 生成网页下载器的函数类型
type GeneratePageDownloader func() MKPageDownloader

// 网页下载器池接口
type MKPageDownloaderPool interface {
	Take() (MKPageDownloader, error)          // 从池中取出一个网页下载器
	Return(downloader MKPageDownloader) error // 把一个网页下载器归还给池
	Total() uint32                            // 获得池的总容量
	Used() uint32                             // 获得正在被池使用的网页下载器的数量
}

// 创建网页下载器池
func NewPageDownloaderPool(total uint32, generator GeneratePageDownloader) (MKPageDownloaderPool, error) {
	etype := reflect.TypeOf(generator)
	generateEntity := func() middleware.MKEntity {
		return generator()
	}

	pool, err := middleware.NewPool(total, etype, generateEntity)
	if err != nil {
		return nil, err
	}

	downloaderPool := &mk_PageDownloaderPool{
		pool:  pool,
		etype: etype,
	}

	return downloaderPool, nil
}

// 网页下载器池实现类型
type mk_PageDownloaderPool struct {
	pool  middleware.MKPool // 实体池
	etype reflect.Type      // 池内实体的类型
}

func (pool *mk_PageDownloaderPool) Take() (MKPageDownloader, error) {
	entity, err := pool.pool.Take()
	if err != nil {
		return nil, err
	}

	downloader, ok := entity.(MKPageDownloader)
	if !ok {
		errMsg := fmt.Sprintf("实体的类型不是%s\n", pool.etype)
		panic(errors.New(errMsg))
	}

	return downloader, nil
}

func (pool *mk_PageDownloaderPool) Return(downloader MKPageDownloader) error {
	return pool.pool.Return()
}

func (pool *mk_PageDownloaderPool) Total() uint32 {
	return pool.pool.Total()
}

func (pool *mk_PageDownloaderPool) Used() uint32 {
	return pool.pool.Used()
}
