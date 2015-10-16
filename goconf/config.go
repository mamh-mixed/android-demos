package goconf

import (
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

// Config 系统启动时先读取配置文件，绑定到这个 struct 上
var Config = &configStruct{}

// configStruct 对应于 config_<env>.js 文件
type configStruct struct {
	App struct {
		LogLevel           log.Level
		EncryptKey         string
		HTTPAddr           string
		TCPAddr            string
		TCPGBKAddr         string
		DefaultCacheTime   Duration
		NotifyURL          string
		OrderCloseTime     Duration
		SessionExpiredTime Duration
	}

	Qiniu struct {
		Bucket string
		Domain string
	}

	Mongo struct {
		Encrypt    bool
		URL        string
		EncryptURL string
		DB         string
	}

	CILOnline struct {
		Host       string
		Port       int
		ServerCert string
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
func (c *configStruct) postProcess() {
	Config.CILOnline.ServerCert = util.WorkDir + "/" + Config.CILOnline.ServerCert

	Config.CFCA.CheckSignPublicKey = util.WorkDir + "/" + Config.CFCA.CheckSignPublicKey

	Config.AlipayScanPay.AlipayPubKey = util.WorkDir + "/" + Config.AlipayScanPay.AlipayPubKey
}
