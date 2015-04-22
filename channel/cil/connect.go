package cil

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

// send 方法会同步返回线下处理结果，它最大的好处是把一个异步 TCP 请求响应变成同步的，无需回调。
// 这对调用者来说是透明的，调用者无需关心与上游网关的通信方式和通信过程，按照正常的顺序流程编写代码，
// 注意：如果上游请求延迟较大，这个方法会阻塞。
func send(msg *model.CilMsg) (back *model.CilMsg) {
	// 串行写入，以免写入错乱
	sendQueue <- msg

	// 交易唯一流水号
	sn := fmt.Sprintf("%s_%s_%s_%s", msg.Chcd, msg.Mchntid, msg.Terminalid, msg.Clisn)
	log.Tracef("send: %s", sn)

	// 结果会异步写入到这个管道中
	c := make(chan *model.CilMsg)

	mapMutex.Lock()
	recvMap[sn] = c
	mapMutex.Unlock()

	// 等待结果返回
	select {
	case back = <-c:
		log.Debug("received request normally")
	case <-time.After(reversalTime * time.Second):
		// 超时处理
		log.Warn("request timeout")
		back = &model.CilMsg{
			Respcd: reversalFlag,
		}
	}
	log.Debugf("rcvd: %s", sn)

	// 清除这个 key 和 管道
	mapMutex.Lock()
	delete(recvMap, sn)
	mapMutex.Unlock()
	log.Debugf("delete: %s", sn)

	return back
}

var sendQueue = make(chan *model.CilMsg, 100)
var recvMap = make(map[string]chan *model.CilMsg, 400)
var mapMutex sync.RWMutex

var addr = "192.168.1.102:7823"
var conn net.Conn

func Connect() {
	var err error
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	// defer conn.Close()
	log.Infof("connect to cil channel %s", addr)

	// 循环接收
	go recv0()

	// 循环发送
	go send0()
}

func restart() {
	log.Warn("connection error, connecting...")
	conn.Close()
	Connect()
}

func recv0() {
	defer restart()
	for {
		msg, err := read()
		if err != nil {
			log.Error(err)
			return
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

		//  如果读取错误或链接断开 EOF 直接返回
		if err != nil {
			return
		}
	}
}

func send0() {
	defer restart()
	for {
		select {
		case msg := <-sendQueue:
			err := write(msg)
			//  如果写入错误或链接断开 EOF 直接返回
			if err != nil {
				return
			}
			sn := fmt.Sprintf("%s_%s_%s_%s", msg.Chcd, msg.Mchntid, msg.Terminalid, msg.Clisn)
			log.Debugf("write: %s", sn)
		case <-time.After(60 * time.Second):
			log.Info("send keepalive")
			err := keepalive()
			if err != nil {
				return
			}
		}
	}
}

func read() (back *model.CilMsg, err error) {
	mLenByte := make([]byte, 4)

	_, err = conn.Read(mLenByte)
	if err != nil {
		log.Debug("read length error: ", err)
		return nil, err
	}

	mlen, err := strconv.Atoi(string(mLenByte))
	if err != nil {
		log.Errorf("can not convert string %s to int: %s", mLenByte, err)
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

	back = &model.CilMsg{}
	err = json.Unmarshal(msg, back)
	if err != nil {
		log.Warnf("msg(% x) can not unmarshal to object", msg, err)
	}

	return back, err
}

func write(msg *model.CilMsg) (err error) {
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

func keepalive() (err error) {
	_, err = io.WriteString(conn, "0000")
	if err != nil {
		log.Error("write len error", err)
	}
	return err
}
