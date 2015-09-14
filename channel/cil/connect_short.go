package cil

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/CardInfoLink/quickpay/goconf"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

var addr = goconf.Config.CILOnline.Host + ":" + strconv.Itoa(goconf.Config.CILOnline.Port)

func send(msg *model.CilMsg, timeout time.Duration) (back *model.CilMsg) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorf("can't connect to CIL-Online tcp://%s: %s", addr, err)
		return nil
	}
	log.Infof("connected to CIL-Online %s", addr)
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(30 * time.Second))

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

	w := bufio.NewWriter(conn)

	// xxxx | 60 0000 0000 | 60 31 0 0 000000 | xxxxxxxx...
	mLen := uint16(len(jsonBytes) + 11)
	binary.Write(w, binary.BigEndian, mLen)

	tpduHeader := []byte{0x60, 0x00, 0x00, 0x00, 0x00, 0x60, 0x31, 0x00, 0x00, 0x00, 0x00}
	_, err = w.Write(tpduHeader)
	if err != nil {
		log.Error("write len error", err)
		return err
	}

	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Error("write len error", err)
		return err
	}

	w.Flush()

	log.Infof("write message: %04x | %+x | %s", mLen, tpduHeader, jsonBytes)
	return nil
}

func receiveOne(conn net.Conn) (back *model.CilMsg, err error) {
	var mLen uint16
	err = binary.Read(conn, binary.BigEndian, &mLen)
	if err != nil {
		return nil, err
	}
	log.Debugf("length: %d", mLen)
	if mLen <= 0 {
		return nil, errors.New("read nothing from CIL-Online")
	}

	tpduHeader := make([]byte, 11)
	_, err = io.ReadFull(conn, tpduHeader)
	// log.Debugf("tpdu and header: %s", tpduHeader)

	msg := make([]byte, mLen-11)
	_, err = io.ReadFull(conn, msg)
	if err != nil {
		return nil, err
	}

	log.Infof("recieve message: %04x | %+x | %s", mLen, tpduHeader, msg)

	back = &model.CilMsg{}
	err = json.Unmarshal(msg, back)
	if err != nil {
		log.Warnf("msg(% x) can not unmarshal to object", msg, err)
	}

	return back, err
}
