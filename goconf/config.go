package goconf

import (
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// Config 系统启动时先读取配置文件，绑定到这个 struct 上
var Config = &ConfigStruct{}

// ConfigStruct 对应于 config_<env>.js 文件
type ConfigStruct struct {
	App struct {
		LogLevel   log.Level
		EncryptKey string
		HTTPAddr   string
		TCPPort    string
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
		ClientCert string
		ClientKey  string
		URL        string
		NotifyURL  string
	}

	AlipayScanPay struct {
		AlipayPubKey string
		OpenAPIURL   string
		URL          string
		NotifyUrl    string
	}
}

// postProcess 后续处理
func (c *ConfigStruct) postProcess() {
	Config.CFCA.CCACert = util.WorkDir + "/" + Config.CFCA.CCACert
	Config.CFCA.Cert = util.WorkDir + "/" + Config.CFCA.Cert
	Config.CFCA.RootCert = util.WorkDir + "/" + Config.CFCA.RootCert

	Config.AlipayScanPay.AlipayPubKey = util.WorkDir + "/" + Config.AlipayScanPay.AlipayPubKey

	Config.WeixinScanPay.ClientCert = util.WorkDir + "/" + Config.WeixinScanPay.ClientCert
	Config.WeixinScanPay.ClientKey = util.WorkDir + "/" + Config.WeixinScanPay.ClientKey
}
