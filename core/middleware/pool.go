package middleware

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// 实体接口
type MKEntity interface {
	ID() uint32 // ID的获取方法
}

// 实体池接口类型
type MKPool interface {
	Take() (MKEntity, error)      // 取出实体
	Return(entity MKEntity) error // 归还实体
	Total() uint32                // 实体池容量
	Used() uint32                 // 实体池中已被使用的实体数量
}

// 创建实体池
func NewPool(total uint32, entityType reflect.Type, generateEntity func() MKEntity) (MKPool, error) {

	if total == 0 {
		errMsg := fmt.Sprintf("无法初始化池【total = %d】\n", total)
		return nil, errors.New(errMsg)
	}

	size := int(total)
	container := make(chan MKEntity, size)
	idContainer := make(map[uint32]bool)

	for i := 0; i < size; i++ {
		newEntity := generateEntity()
		if entityType != reflect.TypeOf(newEntity) {
			errMsg := fmt.Sprintf("函数generateEntity的结果类型不是%s\n", entityType)
			return nil, errors.New(errMsg)
		}

		container <- newEntity
		idContainer[newEntity.ID()] = true
	}

	pool := &mk_pool{
		total:          total,
		etype:          entityType,
		generateEntity: generateEntity,
		container:      container,
		idContainer:    idContainer,
	}

	return pool, nil
}

// 实体池实现类型
type mk_pool struct {
	total          uint32          // 池的总容量
	etype          reflect.Type    // 池中实体的类型
	generateEntity func() MKEntity // 池中实体的生成函数
	container      chan MKEntity   // 实体容器
	idContainer    map[uint32]bool // 实体ID的容器
	mutex          sync.Mutex      // 针对实体ID容器操作的互斥锁
}

func (pool *mk_pool) Take() (MKEntity, error) {
	entity, ok := <-pool.container
	if !ok {
		return nil, errors.New("内部容器不可用")
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	pool.idContainer[entity.ID()] = false

	return entity, nil
}

func (pool *mk_pool) Return(entity MKEntity) error {

	if entity == nil {
		return errors.New("返回无效实体")
	}

	if pool.etype != reflect.TypeOf(entity) {
		errMsg := fmt.Sprintf("返回实体的类型不是 %s\n", pool.etype)
		return errors.New(errMsg)
	}

	entityID := entity.ID()
	compareResult := pool.compareAndSetForIDContainer(entityID, false, true)
	if compareResult == 1 {
		pool.container <- entity
		return nil
	} else if compareResult == 0 {
		errMsg := fmt.Sprintf("池中存在实体【id = %d】", entityID)
		return errors.New(errMsg)
	} else {
		errMsg := fmt.Sprintf("实体【id = %d】无效", entityID)
		return errors.New(errMsg)
	}
}

func (pool *mk_pool) Total() uint32 {
	return pool.total
}

func (pool *mk_pool) Used() uint32 {
	return pool.total - uint32(len(pool.container))
}

// 比较并设置实体ID容器中与给定实体ID对应的键值对的元素值
// 结果值：
// -1: 表示键值对不存在
//  0: 表示操作失败
//  1: 表示操作成功
func (pool *mk_pool) compareAndSetForIDContainer(entityID uint32, oldValue bool, newValue bool) int8 {

	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	v, ok := pool.idContainer[entityID]
	if !ok {
		return -1
	}

	if v != oldValue {
		return 0
	}

	pool.idContainer[entityID] = newValue

	return 1
}
