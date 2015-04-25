// Package entrance 主要为了兼容 支付宝 扫码 TCP 接口
package entrance

import (
	"github.com/omigo/log"
	"net"
)

func Listen() {
	port := ":3000"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Errorf("fail to listen %s port: %s ", port, err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}
