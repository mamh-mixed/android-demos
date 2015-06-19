package goconf

import "github.com/omigo/log"

//
var Config = &ConfigStruct{}

// ConfigStruct 对应于 config.ini 文件
type ConfigStruct struct {
	App struct {
		LogLevel   log.Level
		EncryptKey string
	}

	Mongo struct {
		URL string
		DB  string
	}

	CILOnline struct {
		Host             string
		Port             int
		QueueSize        int
		InitWindowSize   int // 全双工 pipeline 方式，接收数据后，找到对应的请求者
		KeepaliveTime    Duration
		ReconnectTimeout Duration
	}

	CFCA struct {
		Cert     string
		CCACert  string
		RootCert string
		URL      string
	}
}
