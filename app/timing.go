package app

import (
	"github.com/CardInfoLink/quickpay/goconf"
	"net/http"
	"time"
)

var riskEmail = goconf.Config.App.RiskMail

func AppTimingTaskProcess() {

	date := time.Now().Format("2006-01-02")
	NotifySalesman(date)      // 销售工具通知
	InvitationSummary(date)   // 邀请码汇总
	PromoteLimitSummary(date) // 限额发风控
	ClearQiniuResource()
}

// ClearQiniuResource TODO 删除七牛7天之前供外部下载的商户zip
func ClearQiniuResource() {}

func TestSendMail(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("key") == "cil123$" {
		AppTimingTaskProcess()
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("wrong key"))
	}
}
