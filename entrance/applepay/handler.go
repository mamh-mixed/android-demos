package applepay

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"

	"github.com/omigo/log"
)

// ApplePayHandle Apple Pay
func ApplePayHandle(data []byte, merId string) (ret *model.BindingReturn) {
	ap := new(model.ApplePay)

	err := json.Unmarshal(data, ap)
	if err != nil {
		log.Errorf("接卸报文错误: %s", err)
		return mongo.RespCodeColl.Get("200020")
	}

	ap.MerId = merId

	if ret = validateApplePay(ap); ret != nil {
		return ret
	}

	ret = core.ProcessApplePay(ap)
	return ret
}
