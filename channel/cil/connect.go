package cil

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/omigo/log"
)

var sendQueue = make(chan *CilMsg, 100)
var recvMap = make(map[string]chan *CilMsg, 400)
var mapMutex sync.RWMutex

// send 方法会同步返回线下处理结果，它最大的好处是把一个异步 TCP 请求响应变成同步的，无需回调。
// 这对调用者来说是透明的，调用者无需关心与上游网关的通信方式和通信过程，按照正常的顺序流程编写代码，
// 注意：如果上游请求延迟较大，这个方法会阻塞。
func send(msg *CilMsg) (back *CilMsg) {
	// 串行写入，以免写入错乱
	sendQueue <- msg

	// 交易唯一流水号
	sn := fmt.Sprintf("%s_%s_%s_%s", msg.Chcd, msg.Mchntid, msg.Terminalid, msg.Clisn)
	log.Debugf("send %s", sn)

	// 结果会异步写入到这个管道中
	c := make(chan *CilMsg)

	mapMutex.Lock()
	recvMap[sn] = c
	mapMutex.Unlock()

	// 等待结果返回
	back = <-c
	log.Debugf("recv %s", sn)

	// 返回后，清除这个 key
	mapMutex.Lock()
	delete(recvMap, sn)
	mapMutex.Unlock()
	log.Debugf("delete %s", sn)

	return back
}

// var addr = "localhost:8080"
var addr = "192.168.1.102:7823"
var conn net.Conn

func init() {
	var err error
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	// defer conn.Close()
	log.Infof("connect to cil channel %s", addr)

	conn.SetDeadline(time.Now().Add(24 * time.Hour))

	go func() {
		for {
			msg, err := read()
			if err != nil {
				log.Error(err)
			}
			//  如果 msg == nil && err == nil, 表示 keepalive
			if msg == nil {
				log.Info("recv keepalive")
				continue
			}

			sn := fmt.Sprintf("%s_%s_%s_%s", msg.Chcd, msg.Mchntid, msg.Terminalid, msg.Clisn)
			log.Debugf("read: %s", sn)

			// 根据交易流水号取到存放结果的管道
			mapMutex.RLock()
			c := recvMap[sn]
			mapMutex.RUnlock()

			c <- msg
		}
	}()
	go func() {

		for {
			select {
			case msg := <-sendQueue:
				write(msg)
				sn := fmt.Sprintf("%s_%s_%s_%s", msg.Chcd, msg.Mchntid, msg.Terminalid, msg.Clisn)
				log.Debugf("write: %s", sn)
			case <-time.After(60 * time.Second):
				log.Info("send keepalive")
				keepalive()
			}
		}
	}()
}

func read() (back *CilMsg, err error) {
	mLenByte := make([]byte, 4)

	_, err = conn.Read(mLenByte)
	log.Warn("read length error: ", err)

	mlen, err := strconv.Atoi(string(mLenByte))
	if err != nil {
		log.Errorf("can not convert string %s to int: %s", mLenByte, err)
	}

	switch {
	case mlen > 9999 || mlen < 0:
		log.Errorf("read error message length %d", mlen)
	case mlen == 0:
		log.Errorf("read keepalive length %d", mlen)
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

	back = &CilMsg{}
	err = json.Unmarshal(msg, back)
	log.Warnf("msg(% x) can not unmarshal to object", msg, err)

	return back, err
}

func write(msg *CilMsg) (err error) {
	// log.Debugf("%#v", msg)

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		log.Error(err)
		return err
	}

	mLen := len(jsonBytes)
	mLenStr := fmt.Sprintf("%04d", mLen)

	_, err = io.WriteString(conn, mLenStr)
	if err != nil {
		log.Error("write len error", err)
		return err
	}

	_, err = conn.Write(jsonBytes)
	if err != nil {
		log.Error("write len error", err)
		return err
	}

	log.Debugf("write message: %s %s", mLenStr, jsonBytes)

	return nil
}

func keepalive() {
	_, err := io.WriteString(conn, "0000")
	if err != nil {
		log.Error("write len error", err)
	}
}
