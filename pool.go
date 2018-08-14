package xgpool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const DEFAULT_POOL_SIZE = 100

var (
	LenErr = errors.New("pool is full")
	RunErr = errors.New("task timed out")
)

type Pool struct {
	size int
	stop <-chan struct{}

	tasks chan *Task
	wg    sync.WaitGroup

	syncTaskCount  int32
	asyncTaskCount int32
}

func NewPool(size int, s <-chan struct{}) *Pool {
	p := new(Pool)
	if size == 0 {
		p.size = DEFAULT_POOL_SIZE
	} else {
		p.size = size
	}
	p.stop = s
	p.tasks = make(chan *Task, p.size)

	go p.run()

	return p
}

func (this *Pool) run() {
	for {
		select {
		case t := <-this.tasks:
			go func() {
				this.wg.Add(1)
				res, err := t.Run()
				if err != nil {
					//TODO err handle
					t.result <- nil
				}
				t.result <- res
				if t.taskT == TASK_SYNC {
					atomic.AddInt32(&this.syncTaskCount, -1)
				} else if t.taskT == TASK_ASYNC {
					atomic.AddInt32(&this.asyncTaskCount, -1)
				}
				this.wg.Done()
			}()
		case <-this.stop:
			this.Close()
		}
	}
}

func (this *Pool) checkLen() bool {
	return len(this.tasks) == this.size
}

func (this *Pool) Sync(timeout int, fn interface{}, input interface{}) (res interface{}, err error) {
	if this.checkLen() {
		err = LenErr
		return
	}

	result := make(chan interface{})
	stop := make(chan interface{})

	go func() {
		this.tasks <- NewTask(fn, input, result, stop, TASK_SYNC)
	}()

	atomic.AddInt32(&this.syncTaskCount, 1)

	ticker := time.NewTicker(time.Duration(timeout) * time.Second)
	defer ticker.Stop()

	select {
	case res = <-result:
		return
	case <-ticker.C:
		err = RunErr
		return
	}
}

func (this *Pool) Async(fn interface{}, input interface{}) (res chan interface{}, stop chan interface{}, err error) {
	if this.checkLen() {
		err = LenErr
		return
	}

	res = make(chan interface{})
	stop = make(chan interface{})

	go func() {
		this.tasks <- NewTask(fn, input, res, stop, TASK_ASYNC)
	}()

	atomic.AddInt32(&this.asyncTaskCount, 1)

	return
}

// query for amount of task (sync/async)
func (this *Pool) QueryCount() (int32, int32) {
	return this.syncTaskCount, this.asyncTaskCount
}

func (this *Pool) Close() {
	close(this.tasks)
	this.wg.Wait()
}
