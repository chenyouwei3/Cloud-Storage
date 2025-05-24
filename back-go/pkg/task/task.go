package task

import (
	"runtime"
	"sync"
)

type Task interface {
	Do()
}

type funcTask struct {
	fn   interface{}   //要执行的函数
	args []interface{} //函数的参数
}

func (this *funcTask) Do() {
	_, _ = CallFunc(this.fn, this.args...)
}

var (
	defaultTaskPool *TaskPool //默认协程池
	createOnce      sync.Once
)

func Default() *TaskPool {
	createOnce.Do(func() {
		defaultTaskPool = NewTaskPool(runtime.NumCPU()*2, defaultTaskSize)
	})
	return defaultTaskPool
}
