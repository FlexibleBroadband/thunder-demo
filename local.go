package main

import (
	"log"
	"net"
	"os"
	"sync/atomic"

	"util"

	"github.com/FlexibleBroadband/proxy"
)

func main() {
	// record traffic.
	var (
		outTraffic int64
		inTraffic  int64

		conns int64

		// remote addr(host and port).
		remote = "127.0.0.1:9091"

		rand = -1
	)
	go util.RecordTraffic(&inTraffic, &outTraffic, &conns)

	logger := log.New(os.Stdout, "", log.Ldate|log.Lshortfile)
	listen, err := net.Listen("tcp", "127.0.0.1:9090")

	if err != nil {
		panic(err)
	}
	socks5 := proxy.Socks5Listen{
		HandleConnect: func(addr string) (*net.TCPConn, error) {
			// logger.Println("connet addr:=", addr)
			rand *= -1
			// rand = -1
			if rand == 1 {
				conn, err := net.Dial("tcp", addr)
				if err != nil {
					logger.Println("connect error:=", err)
					return nil, err
				}
				return conn.(*net.TCPConn), nil
			} else {
				conn, err := net.Dial("tcp", remote)
				if err != nil {
					logger.Println("connect remote error:=", err)
					return nil, err
				}
				data := []byte{byte(len(addr))}
				_, err = conn.Write(append(data, []byte(addr)...))
				if err != nil {
					logger.Println("conn write error::", err)
				}
				return conn.(*net.TCPConn), nil
			}

		},
		Transport: func(target net.Conn, client net.Conn) error {
			atomic.AddInt64(&conns, 1)
			defer atomic.AddInt64(&conns, -1)
			go util.Copy(client, target, &inTraffic)
			_, err := util.Copy(target, client, &outTraffic)
			return err
		},
		Auth: func(id, pwd []byte) bool {
			logger.Println(len(id), len(pwd))
			logger.Printf("user(%s) pwd(%s)", id, pwd)
			return true
		},
		HandleAssociate: proxy.DefaultHandleAssociate,
		TransportUdp:    proxy.DefaultTransportUdp,

		AddrForClient: "127.0.0.1",

		RawListen: listen,
	}
	socks5.Listen()
}
