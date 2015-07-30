package enterprisepay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"testing"
)

func TestEnterprisePay(t *testing.T) {

	req := &model.ScanPayRequest{}
	req.ActTxamt = "1"
	req.AppID = "wxaa785395d3d09403"
	req.ChanMerId = "1228767002"
	req.CheckName = "NO_CHECK"
	req.Desc = "ipad mini 16G"
	req.OpenId = "omYJssw14onUqv2tocdt0EID3dIc"
	req.OrderNum = util.Millisecond()
	req.SignCert = "dskskfasfsdsjdjqisi343sd99f9djfj"

	DefaultClient.ProcessPay(req)

}
