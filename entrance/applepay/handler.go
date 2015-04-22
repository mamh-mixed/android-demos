package applepay

import (
	"encoding/json"

	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"

	"github.com/omigo/log"
)

// ApplePayHandle Apple Pay 支付入口
func ApplePayHandle(data []byte, merId string) (ret *model.BindingReturn) {
	ap := new(model.ApplePay)

	if err := json.Unmarshal(data, ap); err != nil {
		log.Errorf("can't unmarshal `%s` to json: %s", string(data), err)
		return mongo.RespCodeColl.Get("200020")
	}

	ap.MerId = merId

	if ret = validateApplePay(ap); ret != nil {
		return ret
	}

	ret = core.ProcessApplePay(ap)
	return ret
}
