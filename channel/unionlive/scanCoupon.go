package unionlive

import (
	"strconv"
	"time"

	"github.com/CardInfoLink/quickpay/channel/unionlive/coupon"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
)

// unionliveScanPay 卡券接口
type unionliveScanPay struct{}

// DefaultUnionLiveScanPay 卡券默认实现
var DefaultClient unionliveScanPay

// ProcessPurchaseCoupons 卡券核销
func (u *unionliveScanPay) ProcessPurchaseCoupons(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	submitTime, _ := time.ParseInLocation("2006-01-02 15:04:05", req.CreateTime, time.Local)
	amount, _ := strconv.Atoi(req.VeriTime)

	unionLiveReq := &coupon.PurchaseCouponsReq{
		Header: coupon.PurchaseCouponsReqHeader{
			Version:       Version,
			TransDirect:   TransDirectQ,
			TransType:     "W412",
			MerchantId:    req.ChanMerId,
			SubmitTime:    submitTime.Format("20060102150405"),
			ClientTraceNo: req.OrderNum,
		},
		Body: coupon.PurchaseCouponsReqBody{
			CouponsNo: req.ScanCodeId,
			TermId:    req.Terminalid,
			TermSn:    req.Terminalsn,
			Amount:    amount,
		},
	}
	unionLiveResp := &coupon.PurchaseCouponsResp{}
	err := Execute(unionLiveReq, unionLiveResp)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=PurchaseCoupons, channel=UnionLive", req.OrderNum)
		return nil, err
	}

	// 处理结果返回
	scanPayResponse := &model.ScanPayResponse{
		Txndir:          unionLiveResp.Header.Transdirect,
		Busicd:          model.Veri,
		Respcd:          unionLiveResp.Header.Returncode, // 这个暂时填渠道响应码，之后会改为对应我们系统的码
		AgentCode:       req.AgentCode,
		Chcd:            "UNIONLIVE",
		Mchntid:         req.Mchntid,
		ErrorDetail:     unionLiveResp.Header.Returnmessage,
		OrderNum:        unionLiveResp.Header.Clienttraceno,
		ScanCodeId:      unionLiveResp.Body.Couponsno,
		VeriTime:        strconv.Itoa(unionLiveResp.Body.Amount),
		Sign:            req.Sign,
		CardId:          unionLiveResp.Body.Prodname,
		CardInfo:        unionLiveResp.Body.Proddesc,
		AvailCount:      strconv.Itoa(unionLiveResp.Body.AvailCount),
		ExpDate:         unionLiveResp.Body.ExpDate,
		ChanRespCode:    unionLiveResp.Header.Returncode,
		ChannelOrderNum: unionLiveResp.Header.Hosttraceno,
		Terminalid:      req.Terminalsn,
	}

	return scanPayResponse, nil
}
