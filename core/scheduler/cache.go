package scheduler

import (
	base "core/base"
	"fmt"
	"sync"
)

// 状态字典
var statusMap = map[byte]string{
	0: "running",
	1: "closed",
}

// 请求缓存接口
type requestCache interface {

	// 将请求放入请求缓存
	put(request *base.MKRequest) bool

	// 从请求缓存获取最早被放入且仍在其中的请求
	get() *base.MKRequest

	// 获得请求缓存的容量
	capacity() int

	// 获得请求缓存的实时长度，即：其中的请求的即时数量
	length() int

	// 关闭请求缓存
	close()

	// 获取请求缓存的摘要信息
	summary() string
}

// 创建请求缓存
func newRequestCache() requestCache {
	cache := &mk_requestCache{
		cache: make([]*base.MKRequest, 0),
	}

	return cache
}

// 请求缓存实现类型
type mk_requestCache struct {
	cache  []*base.MKRequest // 请求存储切片
	mutex  sync.Mutex        // 互斥锁
	status byte              // 缓存状态。0表示正在运行，1表示已关闭
}

func (cache *mk_requestCache) put(request *base.MKRequest) bool {

	if reqeust == nil {
		return false
	}

	if request.status == 1 {
		return false
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.cache = append(cache.cache, request)

	return true
}

func (cache *mk_requestCache) get() *base.MKRequest {

	if cache.length() == 0 {
		return nil
	}

	if cache.status == 1 {
		return nil
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	request := cache.cache[0]
	cache.cache = cache.cache[1:]

	return request
}

func (cache *mk_requestCache) capacity() int {
	return cap(cache.cache)
}

func (cache *mk_requestCache) length() int {
	return len(cache.cache)
}

func (cache *mk_requestCache) close() {
	if cache.status == 1 {
		return
	}

	cache.status = 1
}

// 摘要信息模板
var summaryTemplate = "status: %s, " + "length: %d, " + "capacity: %d"

func (cache *mk_requestCache) summary() string {

	summary := fmt.Sprintf(summaryTemplate,
		statusMap[cache.status],
		cache.length(),
		cache.capacity())

	return summary
}
