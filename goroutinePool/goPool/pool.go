package goPool

import (
	"context"
	"sync"
	"sync/atomic"
)

type pool struct {
	//名称
	name string
	//容量，也就是goroutine的最大数量
	capacity int32
	//配置
	config *Config

	//整个协程池的任务队列的头指针与尾指针
	taskHead *task
	taskTail *task
	//任务锁。保证任务相关操作的原子性
	taskLock sync.Mutex
	//任务数量
	taskCount int32

	//worker数量或者叫协程数量
	workerCount int32

	//当执行任务出现painc的处理
	panicHandler func(context.Context, interface{})
}

// Pool 接口定义
type Pool interface {
	//获取协程名称
	Name() string
	//设置容量大小
	SetCap(cap int32)

	//执行任务
	Go(f func())

	//执行任务并传入上线文
	CtxGo(ctx context.Context, f func())

	//设置panic处理逻辑
	SetPanicHandler(f func(context.Context, interface{}))

	//获取work也就是goroutine数量
	WorkerCount() int32
}

func NewPool(name string, capacity int32, config *Config) Pool {
	p := &pool{
		name:     name,
		capacity: capacity,
		config:   config,
	}
	return p
}

// 获取协程池名称
func (p *pool) Name() string {
	return p.name
}

// 设置容量
func (p *pool) SetCap(capacity int32) {
	atomic.StoreInt32(&p.capacity, capacity)
}

// 传入任务，执行任务
func (p *pool) Go(f func()) {
	p.CtxGo(context.Background(), f)
}

// 传入上线文执行任务
func (p *pool) CtxGo(ctx context.Context, f func()) {
	//从实体池中获取未初始化的任务并赋值初始化
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.f = f
	//将任务加入任务队列中，并将任务数量+1
	p.taskLock.Lock()
	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}
	p.taskLock.Unlock()
	atomic.AddInt32(&p.taskCount, 1)
	//如果任务的数量大于配置的任务数量并且work数量大于配置数量，或者work数量为0，创建work(goroutine)
	if (atomic.LoadInt32(&p.taskCount) >= p.config.ScaleThreshold && p.WorkerCount() < atomic.LoadInt32(&p.capacity)) || p.WorkerCount() == 0 {
		p.incWorkerCount()
		w := workerPool.Get().(*worker)
		w.pool = p
		w.run()
	}
}

// 设置panic后的处理逻辑
func (p *pool) SetPanicHandler(f func(context.Context, interface{})) {
	p.panicHandler = f
}

// 查询work数量
func (p *pool) WorkerCount() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

// work数量+1
func (p *pool) incWorkerCount() {
	atomic.AddInt32(&p.workerCount, 1)
}

// work数量-1
func (p *pool) decWorkerCount() {
	atomic.AddInt32(&p.workerCount, -1)
}
