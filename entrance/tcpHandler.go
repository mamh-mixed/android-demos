// Package entrance 主要为了兼容 支付宝/微信 扫码 TCP 接口
package entrance

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/CardInfoLink/quickpay/entrance/scanpay"
	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/util"
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
	log.Infof("ScanPay TCP is listening at %s", port)

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

var errPing = errors.New("ping message")

// 处理这个连接，无论遇到任何错误，都立即断开连接
// TODO 为便于调试跟踪问题，断开前，可以返回 JSON，告知通信异常
func handleConnection(conn net.Conn) {
	log.Debugf("%s connected", conn.RemoteAddr())
	defer conn.Close()

	for {
		reqBytes, err := read(conn)
		if err != nil {
			if err == errPing {
				continue
			}

			if err == io.EOF {
				// read end
				return
			}
			log.Error(err)
			return
		}

		// gbk := string(scanpay.ScanPayHandle([]byte(reqBytes))) // 测试中文编码

		// 数据是以 GBK 编码传输的，需要解码，把 GBK 转成 UTF-8
		utf8, ok := util.GBKTranscoder.Decode(string(reqBytes))
		if !ok {
			log.Error("decode failed")
			return
		}

		// process scanpay
		respBytes := scanpay.ScanPayHandle([]byte(utf8))

		// 数据是以 GBK 编码传输的，发送时需要编码，把 UTF-8 转成 GBK
		gbk, ok := util.GBKTranscoder.Encode(string(respBytes))
		if !ok {
			log.Error("encode failed")
			return
		}

		err = write(conn, gbk)
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
		log.Errorf("read length error: %s", err)
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
