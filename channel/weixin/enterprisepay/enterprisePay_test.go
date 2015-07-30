package enterprisepay

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"testing"
)

func TestEnterprisePay(t *testing.T) {

	req := &model.ScanPayRequest{}
	req.ActTxamt = "1"
	req.AppID = "wx8854422b20240ed2"
	req.ChanMerId = "1230768802"
	req.CheckName = "NO_CHECK"
	req.Desc = "ipad mini 16G"
	req.OpenId = "omYJss7PyKb02j3Y5pnZLm2IL6F4"
	req.OrderNum = util.Millisecond()
	req.SignCert = "0F79610BEBC81F7BFD6212CD656FB467"

	DefaultClient.ProcessPay(req)

}
