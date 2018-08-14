package xgpool

import (
	"testing"
	"time"
)

func echo(t int64) int64 {
	//fmt.Println(t)
	return t
}

func Test_SyncPool(t *testing.T) {
	stop := make(chan struct{})
	p := NewPool(10, stop)
	res, err := p.Sync(10, echo, time.Now().Unix())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func Test_AsyncPool(t *testing.T) {
	stop := make(chan struct{})
	p := NewPool(10, stop)
	res, taskCh, err := p.Async(echo, time.Now().Unix())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(p.QueryCount())
	go func() {
		taskCh <- struct{}{}
	}()
	select {
	case result := <-res:
		t.Log(result)
	}
}

func Benchmark_SyncPool(b *testing.B) {
	stop := make(chan struct{})
	p := NewPool(1024*1024*10, stop)
	for i := 0; i < b.N; i++ {
		p.Sync(10, echo, int64(100))
	}
}

func Benchmark_AsyncPool(b *testing.B) {
	stop := make(chan struct{})
	p := NewPool(1024*1024*10, stop)
	for i := 0; i < b.N; i++ {
		p.Async(echo, int64(100))
	}
}
