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

	submitTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.CreateTime, time.Local)
	if err != nil {
		log.Errorf("format req.CreateTime err,%s", err)
		return nil, err
	}

	amount, err := strconv.Atoi(req.VeriTime)
	if err != nil {
		log.Errorf("format req.VeriTime to int err,%s", err)
		return nil, err
	}

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
			// TermId:    req.Terminalid,
			// TermSn:    req.Terminalsn,
			ExtMercId: req.Mchntid,
			ExtTermId: req.Terminalid,
			Amount:    amount,
		},
		SpReq: req,
	}
	unionLiveResp := &coupon.PurchaseCouponsResp{}
	err = Execute(unionLiveReq, unionLiveResp)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=PurchaseCoupons, channel=UNIONLIVE", req.OrderNum)
		return nil, err
	}

	// 将渠道的错误应答码转为为系统应答码
	returncode, ok := ChanSysRespCode[unionLiveResp.Header.Returncode]
	if !ok {
		log.Warnf("chan Returncode(%s) is not in ChanSysRespCode,", unionLiveResp.Header.Returncode)
		// 未知应答
		returncode = "58"
	}
	errDetail, ok := SysRespCode[returncode]
	if !ok {
		log.Warnf("ChanSysRespCode(key=%s) is not in SysRespCode", returncode)
		errDetail = unionLiveResp.Header.Returnmessage
	}

	// 处理结果返回
	scanPayResponse := &model.ScanPayResponse{
		Txndir:          unionLiveResp.Header.Transdirect,
		Busicd:          model.Veri,
		Respcd:          returncode,
		AgentCode:       req.AgentCode,
		Chcd:            req.Chcd,
		Mchntid:         req.Mchntid,
		ErrorDetail:     errDetail,
		OrderNum:        unionLiveResp.Header.Clienttraceno,
		ScanCodeId:      unionLiveResp.Body.Couponsno,
		VeriTime:        req.VeriTime,
		CardId:          unionLiveResp.Body.Prodname,
		CardInfo:        unionLiveResp.Body.Proddesc,
		AvailCount:      strconv.Itoa(unionLiveResp.Body.AvailCount),
		ExpDate:         unionLiveResp.Body.ExpDate,
		ChanRespCode:    unionLiveResp.Header.Returncode,
		ChannelOrderNum: unionLiveResp.Header.Hosttraceno,
		// Terminalid:      req.Terminalsn,
		Authcode: unionLiveResp.Body.Authcode,
	}

	return scanPayResponse, nil
}
