package conf

import (
	"github.com/CardInfoLink/quickpay/channel/cil"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/mongo"
)

// Initialize 执行系统初始化工作
func Initialize() {
	// 连接到 MongoDB
	mongo.Connect()

	// 初始化卡 Bin 树
	core.BuildTree()

	// 连接到 线下网关
	cil.Connect()

	// 执行清分任务
	core.DoSettWork()

	// 检查数据配置是否有变化
	CheckConf()
}
