### xgpool

xdxduser@hotmail.com

---

[TOC]

#### Install

go get github.com/jollyburger/xgpool

#### Features

- go routine pool managment
- sync and async model of task

#### Examples

##### Sync task

```go
package main

import (
	"github.com/jollyburger/xgpool"
	"fmt"
)

func main() {
	stop := make(chan struct{})
	p := NewPool(10, stop)
	res, err := p.Sync(10, echo, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
```

##### Async task

```go
package main

import (
	"github.com/jollyburger/xgpool"
	"fmt"
)
func main() {
	stop := make(chan struct{})
	p := NewPool(10, stop)
	res, taskCh, err := p.Async(echo, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		taskCh <- struct{}{}
	}()
	select {
	case result := <-res:
		fmt.Println(result)
	}
}
```

#### Performance

```
Benchmark_SyncPool-8    	  200000	      9507 ns/op
Benchmark_AsyncPool-8   	 1000000	      4862 ns/op
```
