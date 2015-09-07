package cil

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

func send0(msg *model.CilMsg, timeout time.Duration) (back *model.CilMsg) {
	host := goconf.Config.CILOnline.Host
	port := goconf.Config.CILOnline.Port
	addr := host + ":" + strconv.Itoa(port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorf("can't connect to CIL-Online tcp://%s: %s", addr, err)
		return nil
	}
	// defer conn.Close()
	log.Infof("connected to CIL-Online %s", addr)

	sendOne(conn, msg)

	back, err = receiveOne(conn)
	if err != nil {
		log.Errorf("receive message from CIL-Online error: %s", err)
		return nil
	}

	return back
}

func sendOne(conn net.Conn, msg *model.CilMsg) (err error) {
	// log.Debugf("%#v", msg)

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		log.Error(err)
		return err
	}

	// xxxx | 60 0000 00 00 | 60 31 0 0 000000 | xxxxxxxx...
	mLen := uint16(len(jsonBytes) + 10 + 12)
	binary.Write(conn, binary.BigEndian, mLen)
	binary.Write(conn, binary.BigEndian, byte(6))
	binary.Write(conn, binary.BigEndian, byte(0))
	binary.Write(conn, binary.BigEndian, uint64(0))
	binary.Write(conn, binary.BigEndian, byte(6))
	binary.Write(conn, binary.BigEndian, byte(0))
	binary.Write(conn, binary.BigEndian, byte(3))
	binary.Write(conn, binary.BigEndian, byte(1))
	binary.Write(conn, binary.BigEndian, uint64(0))

	_, err = conn.Write(jsonBytes)
	if err != nil {
		log.Error("write len error", err)
		return err
	}

	log.Infof("write message: %04x | 60 0000 00 00 | 60 31 0 0 000000 | %s", mLen, jsonBytes)

	return nil
}

func receiveOne(conn net.Conn) (back *model.CilMsg, err error) {
	var mLen uint16
	binary.Read(conn, binary.BigEndian, mLen)

	tpduHeader := make([]byte, 22)
	_, err = io.ReadFull(conn, tpduHeader)
	log.Debugf("tpdu header: %s", tpduHeader)

	msg := make([]byte, mLen)
	_, err = io.ReadFull(conn, msg)
	if err != nil {
		return nil, err
	}

	log.Infof("recieve message: %04x | 60 0000 00 00 | 60 31 0 0 000000 | %s", mLen, msg)

	back = &model.CilMsg{}
	err = json.Unmarshal(msg, back)
	if err != nil {
		log.Warnf("msg(% x) can not unmarshal to object", msg, err)
	}

	return back, err
}
