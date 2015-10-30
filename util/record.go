package util

import (
	"fmt"
	"sync/atomic"
	"time"
)

func RecordTraffic(in, out, conns *int64) {
	c := time.Tick(time.Second)
	for {
		select {
		case <-c:
			tin, tout := atomic.AddInt64(in, 0), atomic.AddInt64(out, 0)
			atomic.AddInt64(in, -tin)
			atomic.AddInt64(out, -tout)
			fmt.Printf("Traffic in:%v KB/s,out:%v KB/s,coons:%v \n", tin/1024, tout/1024, *conns)
		}
	}
}
