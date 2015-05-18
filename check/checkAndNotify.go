// Package check 定时检查数据库 checkAndNotify 文档所有记录，如果发生变化，
// 通知响应业务模块，重新加载缓存，或者重新构建业务逻辑
package check

// CheckAndNotify 检查并通知
type CheckAndNotify struct {
	BizType string
	CurTag  string
	PrevTag string
	App1Tag string
	App2Tag string
}

func init() {
	// do something

	go checking()
}

func checking() {
	// tick ...
	for {
		// now  <- tick
		// 取 checkAndNotify 文档所有记录

		// 遍历
		for {
			// switch  bizType
			// case xxx
			// notify yyy
			// update tag
		}
	}
}
