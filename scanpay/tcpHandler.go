package scanpay

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
)

// 兼容 支付宝/微信 扫码 TCP 接口

// ListenScanPay 启动扫码支付端口监听
func ListenScanPay(addr string, useGBK ...bool) {
	gbk := false
	if len(useGBK) > 0 {
		gbk = useGBK[0]
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("fail to listen %s port: %s", addr, err)
		return
	}
	log.Infof("ScanPay TCP is listening, addr=%s, GBK? %t", addr, gbk)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Errorf("listener fail to accept: %s ", err)
				return
			}
			go handleConnection(conn, gbk)
		}
	}()
}

var errPing = errors.New("ping message")

// 处理这个连接，无论遇到任何错误，都立即断开连接
// TODO 为便于调试跟踪问题，断开前，可以返回 JSON，告知通信异常
func handleConnection(conn net.Conn, gbk bool) {
	log.Infof("%s connected, GBK? %t", conn.RemoteAddr(), gbk)
	defer conn.Close()

	for {
		reqBytes, err := read(conn)
		if err != nil {
			if err == errPing {
				continue
			}

			// read end
			if err == io.EOF {
				log.Infof("read EOF from %s, close connection", conn.RemoteAddr())
				conn.Close()
				return
			}

			log.Error(err)
			return
		}

		var msg string
		var ok bool
		if !gbk {
			// UTF-8 编码
			msg = string(ScanPayHandle([]byte(reqBytes), false)) // 测试中文编码
		} else {
			// 数据是以 GBK 编码传输的，需要解码，把 GBK 转成 UTF-8
			msg, ok = util.GBKTranscoder.Decode(string(reqBytes))
			if !ok {
				log.Error("decode failed")
				return
			}

			// process scanpay
			respBytes := ScanPayHandle([]byte(msg), true)

			// 数据是以 GBK 编码传输的，发送时需要编码，把 UTF-8 转成 GBK
			msg, ok = util.GBKTranscoder.Encode(string(respBytes))
			if !ok {
				log.Error("encode failed")
				return
			}
		}
		err = write(conn, msg)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func read(conn net.Conn) ([]byte, error) {
	mLenByte := make([]byte, 4)

	_, err := io.ReadFull(conn, mLenByte)
	if err != nil {
		if err != io.EOF {
			log.Errorf("read length error: %s", err)
		}
		return nil, err
	}

	mlen, err := strconv.Atoi(string(mLenByte))
	if err != nil {
		log.Errorf("can not convert string %s to int: %s", mLenByte, err)
		return nil, err
	}
	log.Debugf("message length %d", mlen)

	if mlen < 0 {
		log.Errorf("read error message length %d", mlen)
		return nil, fmt.Errorf("error message length %d", mlen)
	}

	// 长度为 0 ，表明读取到 ping 消息
	if mlen == 0 {
		log.Debugf("read keepalive length %d", mlen)
		return nil, errPing
	}

	msg := make([]byte, mlen)
	size, err := io.ReadFull(conn, msg)
	if err != nil {
		return nil, err
	}
	log.Debugf("recieve message: %d%s", size, msg)

	return msg, err
}

func write(conn net.Conn, msg string) error {
	mlen := fmt.Sprintf("%04d", len(msg))

	_, err := conn.Write([]byte(mlen + msg))
	return err
}
