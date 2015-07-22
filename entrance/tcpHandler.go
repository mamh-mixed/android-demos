// Package entrance 主要为了兼容 支付宝/微信 扫码 TCP 接口
package entrance

import (
	"io"
	"net"
	"strconv"

	"github.com/CardInfoLink/quickpay/entrance/scanpay"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/omigo/log"
)

// ListenScanPay 启动扫码支付端口监听
func ListenScanPay() {
	port := goconf.Config.App.TCPPort
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Errorf("fail to listen %s port: %s ", port, err)
		return
	}
	log.Infof("ScanPay is listening on %s", port)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Errorf("listener fail to accept: %s ", err)
				return
			}
			go handleConnection(conn)
		}
	}()

}

func handleConnection(conn net.Conn) {
	log.Debugf("%s connected", conn.RemoteAddr())

	for {
		reqBytes, err := read(conn)
		if err != nil {
			if err == io.EOF {
				return
			}
		}
		// process scanpay
		respBytes := scanpay.TcpScanPayHandle(reqBytes)

		// return
		_, err = conn.Write(respBytes)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func read(conn net.Conn) ([]byte, error) {
	mLenByte := make([]byte, 4)

	_, err := conn.Read(mLenByte)
	if err != nil {
		log.Errorf("read length error: %s", err)
		return nil, err
	}

	mlen, err := strconv.Atoi(string(mLenByte))
	if err != nil {
		log.Errorf("can not convert string %s to int: %s", mLenByte, err)
		return nil, err
	}

	switch {
	case mlen > 9999 || mlen < 0:
		log.Errorf("read error message length %d", mlen)
	case mlen == 0:
		log.Debugf("read keepalive length %d", mlen)
		return nil, nil
	}

	log.Debugf("message length %d", mlen)

	msg := make([]byte, mlen)
	var size int
	for size < mlen {
		rlen, err := conn.Read(msg[size:])
		if err != nil {
			if err == io.EOF {
				// read end
				break
			}
			log.Error(err)
			break
		}
		size += rlen
	}

	log.Debugf("recieve message: %d %s", size, msg)

	return msg, err
}
