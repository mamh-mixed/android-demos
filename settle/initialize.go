package settle

// Initialize 初始化管理平台
func Initialize() {

	// 连接mongo
	// mongo.Connect()

	// 清分任务
	doSettWork()
}
