package goPool

import (
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

// worker实例池，目的实现worker的复用
var workerPool sync.Pool

type worker struct {
	//worker关联的协程池
	pool *pool
}

// 创建空worker
func newWorker() interface{} {
	return &worker{}
}

// worker关联创建协程
func (w *worker) run() {
	go func() {
		for {
			var t *task

			//取出任务
			w.pool.taskLock.Lock()
			if w.pool.taskHead != nil {
				t = w.pool.taskHead
				w.pool.taskHead = w.pool.taskHead.next
				atomic.AddInt32(&w.pool.taskCount, -1)
			}
			if t == nil {
				// if there's no task to do, exit
				w.close()
				w.pool.taskLock.Unlock()
				w.Recycle()
				return
			}
			w.pool.taskLock.Unlock()

			func() {
				//错误处理
				defer func() {
					if r := recover(); r != nil {
						if w.pool.panicHandler != nil {
							w.pool.panicHandler(t.ctx, r)
						} else {
							msg := fmt.Sprintf("GOPOOL: panic in pool: %s: %v: %s", w.pool.name, r, debug.Stack())
							logger.CtxErrorf(t.ctx, msg)
						}
					}
				}()
				//执行任务
				t.f()
			}()
			//清空已经执行的任务，并放入任务池中
			t.Recycle()
		}
	}()
}

// 关闭协程或者worker
func (w *worker) close() {
	w.pool.decWorkerCount()
}

func (w *worker) zero() {
	w.pool = nil
}

func (w *worker) Recycle() {
	w.zero()
	workerPool.Put(w)
}
