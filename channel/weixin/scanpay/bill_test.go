package scanpay

import (
	"github.com/CardInfoLink/quickpay/channel/weixin"
	"github.com/CardInfoLink/quickpay/model"
	"testing"
)

func TestDownloadBill(t *testing.T) {

	spReq := &model.ScanPayRequest{}

	p := &DownloadBillReq{
		CommonParams: testCommonParams,
		BillDate:     "20151219",
		BillType:     "ALL",
	}

	p.CommonParams.Req = spReq

	d := &DownloadBillResp{}
	err := weixin.Execute(p, d)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}
