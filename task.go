package xgpool

import (
	"errors"
	"reflect"
)

type taskType int

const (
	TASK_SYNC taskType = iota + 1
	TASK_ASYNC
)

type Task struct {
	taskT  taskType
	fn     interface{}
	inputs interface{}

	result chan interface{}
	stop   chan interface{}
}

func NewTask(fn interface{}, inputs interface{}, result chan interface{}, stop chan interface{}, tt taskType) *Task {
	t := new(Task)
	t.taskT = tt
	t.fn = fn
	t.inputs = inputs
	t.result = result
	t.stop = stop
	return t
}

func (this *Task) Run() (interface{}, error) {
	for {
		select {
		case <-this.stop:
			return nil, errors.New("stopped")
		default:
			f := reflect.ValueOf(this.fn)
			ft := f.Type()
			if ft.Kind() != reflect.Func || ft.NumIn() == 0 || ft.NumOut() == 0 {
				return nil, errors.New("not find function")
			}
			res := f.Call([]reflect.Value{reflect.ValueOf(this.inputs)})
			// only one response
			return res[0].Interface(), nil
		}
	}
}
