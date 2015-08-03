package pay

import (
	"github.com/CardInfoLink/quickpay/check"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/entrance"
)

// Initialize 执行系统初始化工作
func Initialize() {
	// 初始化卡 Bin 树
	core.BuildTree()

	// 连接到 线下网关
	// cil.Connect()

	// 检查数据配置是否有变化
	check.DoCheck()

	// tcp listen
	entrance.ListenScanPay()
}
