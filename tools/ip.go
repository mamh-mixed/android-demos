package tools

import (
	"fmt"
	"net"
	"strings"
)

var LocalIP string

func init() {

	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	LocalIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	fmt.Printf("local ip: %s\n", LocalIP)
}
