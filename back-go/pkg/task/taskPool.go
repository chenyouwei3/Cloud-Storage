package task

import (
	"errors"
	"sync"
)

// 默认的任务队列的大小
const defaultTaskSize = 1024

// 任务池struct
type TaskPool struct {
	goroutines    int           //当前正在运行的工作协程数
	goroutineSize int           //工作协程的最大数量
	Lock          sync.Mutex    //mutex
	taskChan      chan Task     //任务队列
	die           chan struct{} //关闭信号chan
	dieOnce       sync.Once     //确保任务池只能停止一次(放在结构体是因为让他和结构体的生命周期一致)
}

func (this *TaskPool) NumWorker() int {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	return this.goroutineSize
}

// 返回
func (this *TaskPool) NumTask() int {
	return len(this.taskChan)
}

// NewTaskPool   workerSize(goroutine数量)   taskSize(任务队列大小)
func NewTaskPool(goroutineSize, taskSize int) *TaskPool {
	//防止任务数量过大
	if taskSize < defaultTaskSize {
		taskSize = defaultTaskSize
	}
	// goroutineSize > 0 , 限制goroutine的数量; goroutineSize = 0 , 不限制
	if goroutineSize < 0 {
		goroutineSize = 0
	}
	//初始化变量
	pool := new(TaskPool)
	pool.taskChan = make(chan Task, taskSize)
	pool.die = make(chan struct{})
	pool.goroutineSize = goroutineSize

	return pool
}

// Stop 停止任务池，关闭 die 通道，通知所有工作 goroutine 停止
func (this *TaskPool) Stop() {
	this.dieOnce.Do(func() {
		close(this.die) // 只关闭一次 die 通道，避免多次关闭
	})
}

// 控制并发工作协程的数量,保证不会创建过多的goroutine
func (this *TaskPool) submit(task Task, fullReturn bool) error {
	select {
	case <-this.die:
		return errors.New("taskPool:Submit pool is stopped")
	default:
	}

	var taskChan chan Task
	if this.goroutineSize == 0 {
		taskChan = make(chan Task, 1)
	} else {
		taskChan = this.taskChan
	}
	if fullReturn {
		select {
		case taskChan <- task:
		default:
			return errors.New("taskPool:Submit task channel is full")
		}
	} else {
		taskChan <- task
	}
	this.Lock.Lock()
	defer this.Lock.Unlock()

	if this.goroutineSize == 0 || this.goroutines < this.goroutineSize {
		this.goroutines++
		this.goWorker(taskChan)
	}
	return nil
}

// 执行函数任务
func (this *TaskPool) goWorker(taskC chan Task) {
	go func() {
		defer func() {
			this.Lock.Lock()
			this.goroutines--
			defer this.Lock.Unlock()
		}()
		for {
			select {
			case task := <-taskC:
				task.Do()
			default:
				return
			}
		}
	}()
}

// 向协程池提交
func (this *TaskPool) Submit(fn interface{}, args ...interface{}) error {
	return this.submit(&funcTask{fn: fn, args: args}, false)
}

func (this *TaskPool) SubmitTask(task Task, fullRet ...bool) error {
	fullReturn := false
	if len(fullRet) > 0 && fullRet[0] {
		fullReturn = true
	}
	return this.submit(task, fullReturn)
}
