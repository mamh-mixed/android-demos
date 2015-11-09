// dmf1.0外海接口
package oversea

import (
	"github.com/CardInfoLink/quickpay/channel/alipay/scanpay1"
	"github.com/CardInfoLink/quickpay/model"
)

var DefaultClient alp

type alp struct{}

// ProcessBarcodePay 下单
func (a *alp) ProcessBarcodePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	b := &PayReq{}
	p := &PayResp{}
	if err := scanpay1.Execute(b, p); err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *alp) ProcessQrCodeOfflinePay(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, nil
}

// ProcessRefund 退款
func (a *alp) ProcessRefund(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, nil
}

// ProcessEnquiry 查询
func (a *alp) ProcessEnquiry(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, nil
}

// ProcessCancel 撤销
func (a *alp) ProcessCancel(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, nil
}

// ProcessClose 关闭
func (a *alp) ProcessClose(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, nil
}

// ProcessRefundQuery 退款查询
func (a *alp) ProcessRefundQuery(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	return nil, nil
}
