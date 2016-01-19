package cil

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	// "github.com/CardInfoLink/log"
)

// transformResp 转换应答内容
// TODO 完善应答码转换内容
func transformResp(respcd string) (ret *model.BindingReturn) {
	// default
	ret = mongo.RespCodeColl.Get("000001")
	if respcd == "" {
		return
	}

	// switch respcd {
	// case "00":
	// 	ret = mongo.RespCodeColl.Get("000000")
	// default:
	// 	ret = &model.BindingReturn{
	// 		RespCode: respcd,
	// 		RespMsg:  "渠道返回" + respcd,
	// 	}
	// }
	ret = mongo.RespCodeColl.GetByCIL(respcd)

	ret.ChanRespCode = respcd

	return ret
}
