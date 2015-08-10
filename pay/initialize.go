package pay

import "github.com/CardInfoLink/quickpay/entrance"

// Initialize 执行系统初始化工作
func Initialize() {
	// 初始化卡 Bin 树
	// core.BuildTree()

	// 检查数据配置是否有变化
	// check.DoCheck()

	// 连接到 线下网关
	// cil.Connect()

	// tcp listen
	entrance.ListenScanPay()
}
