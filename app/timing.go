package app

import "time"

var riskEmail = "Risk_management@cardinfolink.com"

func AppTimingTaskProcess() {

	date := time.Now().Format("2006-01-02")
	NotifySalesman(date)      // 销售工具通知
	InvitationSummary(date)   // 邀请码汇总
	PromoteLimitSummary(date) // 限额发风控
	ClearQiniuResource()
}

// ClearQiniuResource TODO 删除七牛7天之前供外部下载的商户zip
func ClearQiniuResource() {}
