package pay

import (
	"github.com/CardInfoLink/quickpay/check"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/entrance"
	"github.com/CardInfoLink/quickpay/goconf"
)

// Initialize 执行系统初始化工作
func Initialize() {
	// 初始化卡 Bin 树
	core.BuildTree()

	// 检查数据配置是否有变化
	check.DoCheck()

	// 连接到 线下网关
	// cil.Connect()

	// 扫码 TCP 接口，UTF-8 编码传输，UTF-8 签名
	port := goconf.Config.App.TCPAddr
	entrance.ListenScanPay(port)

	// 扫码 TCP 接口，GBK 编码传输，UTF-8 签名
	port = goconf.Config.App.TCPGBKAddr
	entrance.ListenScanPay(port, true)
}
