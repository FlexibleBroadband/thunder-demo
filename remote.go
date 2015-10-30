package main

import (
	"net"
	"sync/atomic"

	"util"
)

func main() {
	var (
		outTraffic int64
		inTraffic  int64
		conns      int64
	)
	go util.RecordTraffic(&inTraffic, &outTraffic, &conns)
	l, _ := net.Listen("tcp", "127.0.0.1:9091")
	for {
		rc, _ := l.Accept()
		go func(conn net.Conn) {
			// read addr.
			b := make([]byte, 4096)
			n, err := conn.Read(b)
			if err != nil {
				println(err.Error())
				return
			}
			b = b[:n]
			addrAndOther := b[1:]
			addr := string(addrAndOther[:int(b[0])])

			c, err := net.Dial("tcp", addr)
			if err != nil {
				println("connect err:", err.Error())
				return
			}
			atomic.AddInt64(&conns, 1)
			if len(addrAndOther[int(b[0]):]) > 0 {
				println("I need write.", len(addrAndOther[int(b[0]):]))
				c.Write(addrAndOther[int(b[0]):])
			}
			go util.Copy(c, conn, &outTraffic)
			util.Copy(conn, c, &outTraffic)
			atomic.AddInt64(&conns, -1)
		}(rc)
	}
}
