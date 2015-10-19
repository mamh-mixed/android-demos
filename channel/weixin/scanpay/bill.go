package scanpay

import (
	"github.com/CardInfoLink/quickpay/channel/weixin"
)

// DownloadBillReq 下载对账单
type DownloadBillReq struct {
	weixin.CommonParams

	DeviceInfo string `xml:"device_info,omitempty" url:"device_info,omitempty"` // 设备号
	BillDate   string `xml:"bill_date" url:"bill_date"`
	BillType   string `xml:"bill_type,omitempty" url:"bill_type,omitempty"`
}

// GetURI 取接口地址
func (r *DownloadBillReq) GetURI() string {
	return "/pay/downloadbill"
}

// DownloadBillResp
type DownloadBillResp struct {
	weixin.CommonBody
}
