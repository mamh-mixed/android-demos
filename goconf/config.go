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
		LogLevel         log.Level
		EncryptKey       string
		HTTPAddr         string
		TCPAddr          string
		TCPGBKAddr       string
		DefaultCacheTime Duration
		NotifyURL        string
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
		URL                string
		CheckSignPublicKey string
	}

	WeixinScanPay struct {
		URL                 string
		NotifyURL           string
		DNSCacheRefreshTime Duration
	}

	AlipayScanPay struct {
		AlipayPubKey string
		OpenAPIURL   string
		URL          string
		NotifyUrl    string
		AgentId      string
	}
}

// postProcess 后续处理
func (c *ConfigStruct) postProcess() {
	Config.CFCA.CheckSignPublicKey = util.WorkDir + "/" + Config.CFCA.CheckSignPublicKey

	Config.AlipayScanPay.AlipayPubKey = util.WorkDir + "/" + Config.AlipayScanPay.AlipayPubKey
}
