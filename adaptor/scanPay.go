package adaptor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
)

// 扫码支付所有与判断渠道相关的逻辑统一写在这里，为上层屏蔽各个渠道的差异

// ProcessBarcodePay 扫条码下单
func ProcessBarcodePay(req *model.ScanPay) (*model.ScanPayResponse, error) {
	// 渠道选择
	// 根据扫码Id判断走哪个渠道
	if req.Chcd == "" {
		if strings.HasPrefix(req.ScanCodeId, "1") {
			req.Chcd = "WXP"
		} else if strings.HasPrefix(req.ScanCodeId, "2") {
			req.Chcd = "ALP"
		} else {
			// 不送，返回 TODO check error code
			return nil, errors.New("SYSTEM_ERROR")
		}
	}

	// TODO 这里直接字符串转换可能方便很多，判断前 n-3 位，去除前缀 0，然后再倒数第二位前加一个小数点
	// 金额单位转换 txamt:000000000010 分
	f, err := strconv.ParseInt(req.Txamt, 10, 64)
	if err != nil {
		return nil, errors.New("SYSTEM_ERROR")
	}
	// // 参数处理
	switch req.ChanCode {
	case "ALP":
		req.ActTxamt = fmt.Sprintf("%d.%d", f/100, f%100)
	case "WXP":
		req.ActTxamt = fmt.Sprintf("%d", f)
	default:
		req.ActTxamt = req.Txamt
	}

	// 获得渠道实例，请求
	sp := channel.GetScanPayChan(req.Chcd)
	if sp == nil {
		return nil, errors.New("SYSTEM_ERROR")
	}

	return sp.ProcessBarcodePay(req)
}
