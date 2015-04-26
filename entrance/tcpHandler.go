// Package entrance 主要为了兼容 支付宝 扫码 TCP 接口
package entrance

import (
	"github.com/CardInfoLink/quickpay/entrance/scanpay"
	"github.com/omigo/log"
	"io"
	"net"
	"strconv"
)

func Listen() {
	port := ":3000"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Errorf("fail to listen %s port: %s ", port, err)
		return
	}

	go func(l net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Errorf("listener fail to accept: %s ", err)
				return
			}
			go handleConnection(conn)
		}
	}(ln)

}

func handleConnection(conn net.Conn) {

	log.Debugf("%s connected", conn.RemoteAddr())

	for {
		reqBytes, err := read(conn)
		if err != nil {
			return
		}
		// process scanpay
		respBytes := scanpay.Router(reqBytes)

		// return
		_, err = conn.Write(respBytes)
		if err != nil {
			log.Error(err)
			return
		}
	}

}

func read(conn net.Conn) (back []byte, err error) {
	mLenByte := make([]byte, 4)

	_, err = conn.Read(mLenByte)
	if err != nil {
		log.Debug("read length error: ", err)
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
		log.Debug(size)
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

	return back, err
}
