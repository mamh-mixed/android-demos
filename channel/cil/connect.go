package cil

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/omigo/log"
)

var addr = "192.168.1.102:7823"

// var addr = "localhost:8080"

func init() {
	connect()

}

func connect() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Error(err.Error())
	}
	defer conn.Close()
	log.Infof("connect to cil channel %s", addr)

	conn.SetDeadline(time.Now().Add(24 * time.Hour))

	go read(conn)
	go write(conn)

	select {}
}

func read(conn net.Conn) {

	for {
		// 初始为空
		mLenByte := make([]byte, 4)
		msg := make([]byte, 9999)

		_, err := conn.Read(mLenByte)
		log.Check("read length error: ", err)

		mlen, err := strconv.Atoi(string(mLenByte))
		if err != nil {
			log.Errorf("can not convert string %s to int: %s", mLenByte, err)
		}

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

		// read err or read end
		if err != nil {
			break
		}
	}
}

func write(conn net.Conn) {
	msg := NewConsumeCilMsg()

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		log.Error(err)
		return
	}
	jsonBytes = []byte(`{"busicd": "500000", "txndir": "Q", "posentrymode": "022", "txamt": "000000001000", "txdt": "0926115934", "localdt": "0926115934", "cardcd": "9559970030000000215", "trackdata2": "9559970030000000215=00002101815546", "trackdata3": "", "cardpin": "", "syssn": "101213113013", "clisn": "115934", "inscd": "30512900", "chcd": "04012900", "mchntid": "0002220F0002804", "terminalid": "60000005", "mcc": "4816", "txcurrcd": "156", "billingcurr": "156", "regioncd": "0156", "mchntnm": "shanghai test                           ", "nminfo": "PKE", "cardseqnum": "001", "iccdata": "", "termreadability": "5", "icccondcode": "0", "outgoingacct": "9559970030000000215", "incomingacct": "4682030210337444", "custmrtp": "01", "custmracnt": "130412", "paymd": "01", "goodscd": "19100059", "billyymm": "201201", "chname": "", "inchname": "涓婃捣娴嬭瘯", "phonenum": "13611111111", "cvv2": "111", "paymethod": "3", "billinscd": "888880000502900", "barcd": "539100060832536001034816", "psamcd": "1234567890123456", "txnmode": "1", "termserialcd": "1234567890123", "expiredate": "1605", "usagetags": "12"}`)

	mLen := len(jsonBytes)
	mLenStr := fmt.Sprintf("%04d", mLen)

	_, err = io.WriteString(conn, mLenStr)
	if err != nil {
		log.Error("write len error", err)
		return
	}

	_, err = conn.Write(jsonBytes)
	if err != nil {
		log.Error("write len error", err)
		return
	}

	log.Debugf("write message: %s %s", mLenStr, jsonBytes)
}
