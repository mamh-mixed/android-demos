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
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
)

const keepaliveTime = 60 * time.Second
const reconnectTimeout = 5 * time.Second

var (
	host = "$CIL_HOST | 192.168.1.102"
	port = "$CIL_PORT | 7823"
)

var defualtClient *CilOnlinePay

// Connect 连接到线下
func Connect() {
	fevHost := tools.FirstExistValue(host)
	fevPort := tools.FirstExistValue(port)
	addr := fevHost + ":" + fevPort

	defualtClient = NewCilOnlinePay(addr, 10000, 100)

	defualtClient.Connect()
}

// send 方法会同步返回线下处理结果，它最大的好处是把一个异步 TCP 请求响应变成同步的，无需回调。
// 这对调用者来说是透明的，调用者无需关心与上游网关的通信方式和通信过程，按照正常的顺序流程编写代码，
// 注意：如果上游请求延迟较大，这个方法会阻塞。
func send(msg *model.CilMsg, timeout time.Duration) (back *model.CilMsg) {
	return defualtClient.Send(msg, timeout)
}

// CilOnlinePay 线下联机系统
type CilOnlinePay struct {
	Addr string
	conn net.Conn

	closed bool

	sendQueue chan *model.CilMsg

	recvMap  map[string]chan *model.CilMsg
	mapMutex sync.RWMutex
}

// NewCilOnlinePay 创建一个线下联机系统客户端
func NewCilOnlinePay(addr string, queueSize, initRecvMapSize int) (c *CilOnlinePay) {
	c = &CilOnlinePay{Addr: addr}
	c.sendQueue = make(chan *model.CilMsg, queueSize)
	c.recvMap = make(map[string]chan *model.CilMsg, initRecvMapSize)

	return c
}

func (c *CilOnlinePay) reconnect() {
	log.Info("CIL-Online connect error, reconnect...")
	c.closed = true
	c.conn.Close()
	c.Connect()
}

// Connect 建立连接
func (c *CilOnlinePay) Connect() {
	// 异步建立连接，以免线下连接问题造成快捷支付无法启动
	go func() {
		var err error
		c.conn, err = net.Dial("tcp", c.Addr)
		if err != nil {
			log.Errorf("can't connect to CIL-Online tcp://%s: %s", c.Addr, err)
			log.Infof("sleep %s to reconnect...", reconnectTimeout)
			time.Sleep(reconnectTimeout)
			Connect()
			return
		}
		// defer conn.Close()
		log.Infof("connect to CIL-Online %s", c.Addr)

		go c.WaitAndReceive()
		go c.LoopToSend()
	}()
}

// WaitAndReceive 等待接收消息
func (c *CilOnlinePay) WaitAndReceive() {
	for {
		msg, err := c.ReceiveOne()
		if err != nil {
			log.Error(err)
			c.reconnect()
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
		c.mapMutex.RLock()
		retChan := c.recvMap[sn]
		c.mapMutex.RUnlock()

		retChan <- msg

		// 如果读取错误或链接断开 EOF 直接返回
		if err != nil {
			c.reconnect()
			return
		}
	}
}

// LoopToSend 循环从队列中取报文，发送
func (c *CilOnlinePay) LoopToSend() {
	for {
		select {
		case msg := <-c.sendQueue:
			if c.closed {
				// 连接关闭后，取出来的消息无法写入，只能再放回去
				c.sendQueue <- msg
				return
			}

			err := c.SendOne(msg)
			//  如果写入错误或链接断开 EOF 直接返回
			if err != nil {
				c.reconnect()
				return
			}
			sn := fmt.Sprintf("%s_%s_%s_%s", msg.Chcd, msg.Mchntid, msg.Terminalid, msg.Clisn)
			log.Debugf("write: %s", sn)

		case <-time.After(keepaliveTime):
			if c.closed {
				return
			}

			log.Info("send keepalive")
			err := c.Keepalive()
			if err != nil {
				c.reconnect()
				return
			}

		}
	}
}

// Close 关闭连接
func (c *CilOnlinePay) Close() {
	c.conn.Close()
}

// Keepalive 避免长时间连接空闲自动断开，每隔一段时间需调用一次这个方法
func (c *CilOnlinePay) Keepalive() (err error) {
	_, err = io.WriteString(c.conn, "0000")
	if err != nil {
		log.Error("write len error", err)
	}
	return err
}

// SendOne 发生报文
func (c *CilOnlinePay) SendOne(msg *model.CilMsg) (err error) {
	// log.Debugf("%#v", msg)

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		log.Error(err)
		return err
	}

	mLen := len(jsonBytes)
	mLenStr := fmt.Sprintf("%04d", mLen)

	_, err = io.WriteString(c.conn, mLenStr)
	if err != nil {
		log.Error("write len error", err)
		return err
	}

	_, err = c.conn.Write(jsonBytes)
	if err != nil {
		log.Error("write len error", err)
		return err
	}

	log.Debugf("write message: %s %s", mLenStr, jsonBytes)

	return nil
}

// ReceiveOne 接收异步返回的报文
func (c *CilOnlinePay) ReceiveOne() (back *model.CilMsg, err error) {
	mLenByte := make([]byte, 4)

	_, err = c.conn.Read(mLenByte)
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
		// log.Debug(size)
		rlen, err := c.conn.Read(msg[size:])
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

// Send 方法会同步返回线下处理结果，它最大的好处是把一个异步 TCP 请求响应变成同步的，无需回调。
// 这对调用者来说是透明的，调用者无需关心与上游网关的通信方式和通信过程，按照正常的顺序流程编写代码，
// 注意：如果上游请求延迟较大，这个方法会阻塞。
func (c *CilOnlinePay) Send(msg *model.CilMsg, timeout time.Duration) (back *model.CilMsg) {
	// 串行写入，以免写入错乱
	c.sendQueue <- msg

	// 交易唯一流水号
	sn := fmt.Sprintf("%s_%s_%s_%s", msg.Chcd, msg.Mchntid, msg.Terminalid, msg.Clisn)
	log.Tracef("send: %s", sn)

	// 结果会异步写入到这个管道中
	retChan := make(chan *model.CilMsg)

	c.mapMutex.Lock()
	c.recvMap[sn] = retChan
	c.mapMutex.Unlock()

	// 等待结果返回
	select {
	case back = <-retChan:
		log.Debug("received request normally")
	case <-time.After(timeout):
		// 超时处理
		log.Warn("request timeout")
		back = &model.CilMsg{
			Respcd: reversalFlag,
		}
	}
	log.Debugf("rcvd: %s", sn)

	// 清除这个 key 和 管道
	c.mapMutex.Lock()
	delete(c.recvMap, sn)
	c.mapMutex.Unlock()
	log.Debugf("delete: %s", sn)

	return back
}
