package goroutinePool

import (
	"context"
	"sync"
)

func init() {
	workerPool.New = newWorker
}

/*
 *任务相关，使用实例池来达到task的复用
 */
type task struct {
	ctx  context.Context
	f    func()
	next *task
}

// 任务池，复用task对象
var taskPool sync.Pool

func init() {
	taskPool.New = newTask
}

// 清空任务内容
func (t *task) zero() {
	t.ctx = nil
	t.f = nil
	t.next = nil
}

// 实现任务内容清空并放入对象池
func (t *task) Recycle() {
	t.zero()
	taskPool.Put(t)
}

// 创建一个空的task
func newTask() interface{} {
	return &task{}
}

// 任务列表
type taskList struct {
	sync.Mutex
	taskHead *task
	taskTail *task
}
