package goconf

import (
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
)

// Config 系统启动时先读取配置文件，绑定到这个 struct 上
var Config = &ConfigStruct{}

// ConfigStruct 对应于 config.ini 文件
type ConfigStruct struct {
	App struct {
		LogLevel   log.Level
		EncryptKey string
		TcpPort    string
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

	WeixinScanPay struct {
		URL string
	}

	AlipayScanPay struct {
		URL       string
		NotifyUrl string
	}
}

// postProcess 后续处理
func (c *ConfigStruct) postProcess() {
	Config.CFCA.CCACert = tools.WorkDir + "/" + Config.CFCA.CCACert
	Config.CFCA.Cert = tools.WorkDir + "/" + Config.CFCA.Cert
	Config.CFCA.RootCert = tools.WorkDir + "/" + Config.CFCA.RootCert
}
