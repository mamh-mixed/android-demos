package settle

import (
	"github.com/CardInfoLink/quickpay/mongo"
)

// Initialize 初始化管理平台
func Initialize() {

	// 连接mongo
	mongo.Connect()

	// 清分任务
	doSettWork()
}
