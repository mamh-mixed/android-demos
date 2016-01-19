package unionlive

import (
	"strconv"

	"github.com/CardInfoLink/quickpay/adaptor"
	"github.com/CardInfoLink/quickpay/channel/unionlive/coupon"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/log"
)

// unionliveScanPay 卡券接口
type unionliveScanPay struct{}

// DefaultUnionLiveScanPay 卡券默认实现
var DefaultClient unionliveScanPay

// ProcessPurchaseCoupons 电子券验证/刷卡活动券查询
func (u *unionliveScanPay) ProcessPurchaseCoupons(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	unionLiveReq := &coupon.PurchaseCouponsReq{
		Header: coupon.PurchaseCouponsReqHeader{
			Version:       Version,
			TransDirect:   TransDirectQ,
			TransType:     "W412",
			MerchantId:    req.ChanMerId,
			SubmitTime:    req.CreateTime,
			ClientTraceNo: req.OrderNum,
		},
		Body: coupon.PurchaseCouponsReqBody{
			CouponsNo: req.ScanCodeId,
			TermId:    req.Terminalid,
			// TermSn:    req.Terminalsn,
			ExtMercId:   req.Mchntid,
			ExtTermId:   req.Terminalsn,
			Amount:      req.IntVeriTime,
			TransAmount: req.IntTxamt,
		},
		SpReq: req,
	}
	unionLiveResp := &coupon.PurchaseCouponsResp{}
	err := Execute(unionLiveReq, unionLiveResp)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=PurchaseCoupons, channel=ULIVE", req.OrderNum)
		return nil, err
	}

	// 将渠道的错误应答码转为为系统应答码
	returncode, errDetail := transChanToSysCode(unionLiveResp.Header.Returncode, unionLiveResp.Header.Returnmessage)

	actualPayAmount := strconv.Itoa(unionLiveResp.Body.Price)
	if req.Txamt == "" {
		actualPayAmount = ""
	}

	// 处理结果返回
	scanPayResponse := &model.ScanPayResponse{
		Txndir:          unionLiveResp.Header.Transdirect,
		Busicd:          req.Busicd,
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
		Authcode:        unionLiveResp.Body.Authcode,
		ChannelTime:     unionLiveResp.Header.Hosttime,
		VoucherType:     strconv.Itoa(unionLiveResp.Body.VoucherType),
		SaleMinAmount:   strconv.Itoa(unionLiveResp.Body.SaleMinAmount),
		SaleDiscount:    strconv.Itoa(unionLiveResp.Body.SaleDiscount),
		ActualPayAmount: actualPayAmount,
		Txamt:           req.Txamt,
	}

	return scanPayResponse, nil
}

// ProcessPurchaseActCoupons 刷卡活动券验证
func (u *unionliveScanPay) ProcessPurchaseActCoupons(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {

	unionLiveReq := &coupon.PurchaseActCouponsReq{
		Header: coupon.PurchaseActCouponsReqHeader{
			Version:       Version,
			TransDirect:   TransDirectQ,
			TransType:     "W452",
			MerchantId:    req.ChanMerId,
			SubmitTime:    req.CreateTime,
			ClientTraceNo: req.OrderNum,
		},
		Body: coupon.PurchaseActCouponsReqBody{
			CouponsNo:      req.ScanCodeId,
			OldHostTraceNo: req.OrigChanOrderNum,
			TermId:         req.Terminalid,
			// TermSn:    req.Terminalsn,
			ExtMercId:   req.Mchntid,
			ExtTermId:   req.Terminalsn,
			Amount:      req.IntVeriTime,
			Cardbin:     req.Cardbin,
			TransAmount: req.IntTxamt,
			PayType:     req.IntPayType,
		},
		SpReq: req,
	}
	unionLiveResp := &coupon.PurchaseActCouponsResp{}
	err := Execute(unionLiveReq, unionLiveResp)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=PurchaseActCoupons, channel=ULIVE", req.OrderNum)
		return nil, err
	}

	// 将渠道的错误应答码转为为系统应答码
	returncode, errDetail := transChanToSysCode(unionLiveResp.Header.ReturnCode, unionLiveResp.Header.ReturnMessage)

	// 处理结果返回
	scanPayResponse := &model.ScanPayResponse{
		Txndir:          unionLiveResp.Header.TransDirect,
		Busicd:          req.Busicd,
		Respcd:          returncode,
		AgentCode:       req.AgentCode,
		Chcd:            req.Chcd,
		Mchntid:         req.Mchntid,
		ErrorDetail:     errDetail,
		OrderNum:        unionLiveResp.Header.ClientTraceNo,
		ScanCodeId:      unionLiveResp.Body.CouponsNo,
		VeriTime:        req.VeriTime,
		CardId:          unionLiveResp.Body.ProdName,
		CardInfo:        unionLiveResp.Body.ProdDesc,
		AvailCount:      strconv.Itoa(unionLiveResp.Body.AvailCount),
		ExpDate:         unionLiveResp.Body.ExpDate,
		ChanRespCode:    unionLiveResp.Header.ReturnCode,
		ChannelOrderNum: unionLiveResp.Header.HostTraceNo,
		// Terminalid:      req.Terminalsn,
		Authcode:        unionLiveResp.Body.AuthCode,
		ChannelTime:     unionLiveResp.Header.HostTime,
		VoucherType:     strconv.Itoa(unionLiveResp.Body.VoucherType),
		SaleMinAmount:   strconv.Itoa(unionLiveResp.Body.SaleMinAmount),
		SaleDiscount:    strconv.Itoa(unionLiveResp.Body.SaleDiscount),
		ActualPayAmount: strconv.Itoa(unionLiveResp.Body.ActualPayAmount),
		Txamt:           req.Txamt,
		PayType:         req.PayType,
		Cardbin:         req.Cardbin,
		OrigOrderNum:    req.OrigOrderNum,
	}

	return scanPayResponse, nil
}

// ProcessQueryPurchaseCouponsResult 电子券验证结果查询
func (u *unionliveScanPay) ProcessQueryPurchaseCouponsResult(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	unionLiveReq := &coupon.QueryPurchaseCouponsResultReq{
		Header: coupon.QueryPurchaseCouponsResultReqHeader{
			Version:       Version,
			TransDirect:   TransDirectQ,
			TransType:     "W394",
			MerchantId:    req.ChanMerId,
			SubmitTime:    req.CreateTime,
			ClientTraceNo: req.OrderNum,
		},
		Body: coupon.QueryPurchaseCouponsResultReqBody{
			CouponsNo: req.ScanCodeId,
			TermId:    req.Terminalid,
			// TermSn:    req.Terminalsn,
			ExtMercId:        req.Mchntid,
			ExtTermId:        req.Terminalsn,
			Amount:           req.IntVeriTime,
			TransAmount:      req.IntTxamt,
			PayType:          req.IntPayType,
			OldClientTraceNo: req.OrigOrderNum,
			OldSubmitTime:    req.OrigSubmitTime,
		},
		SpReq: req,
	}
	unionLiveResp := &coupon.QueryPurchaseCouponsResultResp{}
	err := Execute(unionLiveReq, unionLiveResp)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=QueryPurchaseCouponsResult, channel=ULIVE", req.OrderNum)
		return nil, err
	}

	// 将渠道的错误应答码转为为系统应答码
	returncode, errDetail := transChanToSysCode(unionLiveResp.Header.ReturnCode, unionLiveResp.Header.ReturnMessage)

	// 将原渠道的错误应答码转为为系统应答码
	origReturnCode, origErrDetail := transChanToSysCode(unionLiveResp.Body.OldReturnCode, unionLiveResp.Body.OldReturnMessage)

	// 处理结果返回
	scanPayResponse := &model.ScanPayResponse{
		Txndir:          unionLiveResp.Header.TransDirect,
		Busicd:          req.Busicd,
		Respcd:          returncode,
		AgentCode:       req.AgentCode,
		Chcd:            req.Chcd,
		Mchntid:         req.Mchntid,
		ErrorDetail:     errDetail,
		OrderNum:        unionLiveResp.Header.ClientTraceNo,
		ScanCodeId:      unionLiveResp.Body.CouponsNo,
		VeriTime:        req.VeriTime,
		CardId:          unionLiveResp.Body.ProdName,
		CardInfo:        unionLiveResp.Body.ProdDesc,
		AvailCount:      strconv.Itoa(unionLiveResp.Body.AvailCount),
		ExpDate:         unionLiveResp.Body.ExpDate,
		ChanRespCode:    unionLiveResp.Header.ReturnCode,
		ChannelOrderNum: unionLiveResp.Header.HostTraceNo,
		// Terminalid:      req.Terminalsn,
		Authcode:        unionLiveResp.Body.AuthCode,
		ChannelTime:     unionLiveResp.Header.HostTime,
		VoucherType:     strconv.Itoa(unionLiveResp.Body.VoucherType),
		SaleMinAmount:   strconv.Itoa(unionLiveResp.Body.SaleMinAmount),
		SaleDiscount:    strconv.Itoa(unionLiveResp.Body.SaleDiscount),
		ActualPayAmount: strconv.Itoa(unionLiveResp.Body.ActualPayAmount),
		OrigOrderNum:    req.OrigOrderNum,
		OrigRespcd:      origReturnCode,
		OrigErrorDetail: origErrDetail,
	}

	return scanPayResponse, nil
}

// ProcessUndoPurchaseActCoupons 刷卡活动券验证撤销
func (u *unionliveScanPay) ProcessUndoPurchaseActCoupons(req *model.ScanPayRequest) (*model.ScanPayResponse, error) {
	unionLiveReq := &coupon.UndoPurchaseActCouponsReq{
		Header: coupon.UndoPurchaseActCouponsReqHeader{
			Version:       Version,
			TransDirect:   TransDirectQ,
			TransType:     "W492",
			MerchantId:    req.ChanMerId,
			SubmitTime:    req.CreateTime,
			ClientTraceNo: req.OrderNum,
		},
		Body: coupon.UndoPurchaseActCouponsReqBody{
			CouponsNo: req.ScanCodeId,
			TermId:    req.Terminalid,
			// TermSn:    req.Terminalsn,
			ExtMercId:        req.Mchntid,
			ExtTermId:        req.Terminalsn,
			OldTransAmount:   req.OrigVeriTime,
			OldSubmitTime:    req.OrigSubmitTime,
			OldClientTraceNo: req.OrigOrderNum,
			OldHostTraceNo:   req.OrigChanOrderNum,
		},
		SpReq: req,
	}
	unionLiveResp := &coupon.UndoPurchaseActCouponsResp{}
	err := Execute(unionLiveReq, unionLiveResp)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=UndoPurchaseActCoupons, channel=ULIVE", req.OrderNum)
		return nil, err
	}

	// 将渠道的错误应答码转为为系统应答码
	returncode, errDetail := transChanToSysCode(unionLiveResp.Header.ReturnCode, unionLiveResp.Header.ReturnMessage)

	// 处理结果返回
	scanPayResponse := &model.ScanPayResponse{
		Txndir:          unionLiveResp.Header.TransDirect,
		Busicd:          req.Busicd,
		Respcd:          returncode,
		AgentCode:       req.AgentCode,
		Chcd:            req.Chcd,
		Mchntid:         req.Mchntid,
		ErrorDetail:     errDetail,
		OrderNum:        unionLiveResp.Header.ClientTraceNo,
		ScanCodeId:      unionLiveResp.Body.CouponsNo,
		ChanRespCode:    unionLiveResp.Header.ReturnCode,
		ChannelOrderNum: unionLiveResp.Header.HostTraceNo,
		// Terminalid:      req.Terminalsn,
		Authcode:     unionLiveResp.Body.AuthCode,
		ChannelTime:  unionLiveResp.Header.HostTime,
		OrigOrderNum: req.OrigOrderNum,
	}

	return scanPayResponse, nil
}

// transChanToSysCode 将渠道的错误应答码转为为系统应答码
func transChanToSysCode(chanReturnCode, chanErrMessage string) (returnCode, errDetail string) {
	if chanReturnCode == "" {
		return "", ""
	}
	returnCode, ok := ChanSysRespCode[chanReturnCode]
	if !ok {
		log.Warnf("chan Returncode(%s) is not in ChanSysRespCode,", chanReturnCode)
		// 未知应答
		returnCode = "58"
	}
	errDetail, ok = SysRespCode[returnCode]
	if !ok {
		log.Warnf("ChanSysRespCode(key=%s) is not in SysRespCode", returnCode)
		errDetail = chanErrMessage
	}
	return returnCode, errDetail
}

// ProcessPurchaseCouponsSingle 电子券验证
func (u *unionliveScanPay) ProcessPurchaseCouponsSingle(req *model.ScanPayRequest) *model.ScanPayResponse {
	unionLiveReq := &coupon.PurchaseCouponsSingleReq{
		Header: coupon.PurchaseCouponsSingleReqHeader{
			Version:       Version,
			TransDirect:   TransDirectQ,
			TransType:     "W462",
			MerchantId:    req.ChanMerId,
			SubmitTime:    req.CreateTime,
			ClientTraceNo: req.OrderNum,
		},
		Body: coupon.PurchaseCouponsSingleReqBody{
			CouponsNo: req.ScanCodeId,
			TermId:    req.Terminalid,
			// TermSn:    req.Terminalsn,
			ExtMercId:   req.Mchntid,
			ExtTermId:   req.Terminalsn,
			Amount:      req.IntVeriTime,
			Cardbin:     req.Cardbin,
			TransAmount: req.IntTxamt,
			PayType:     req.IntPayType,
		},
		SpReq: req,
	}
	unionLiveResp := &coupon.PurchaseCouponsSingleResp{}
	err := Execute(unionLiveReq, unionLiveResp)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=PurchaseCouponsSingle, channel=ULIVE", req.OrderNum)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 将渠道的错误应答码转为为系统应答码
	returncode, errDetail := transChanToSysCode(unionLiveResp.Header.Returncode, unionLiveResp.Header.Returnmessage)

	actualPayAmount := strconv.Itoa(unionLiveResp.Body.Price)
	if req.Txamt == "" {
		actualPayAmount = ""
	}

	// 处理结果返回
	scanPayResponse := &model.ScanPayResponse{
		Txndir:          unionLiveResp.Header.Transdirect,
		Busicd:          req.Busicd,
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
		Authcode:        unionLiveResp.Body.Authcode,
		ChannelTime:     unionLiveResp.Header.Hosttime,
		VoucherType:     strconv.Itoa(unionLiveResp.Body.VoucherType),
		SaleMinAmount:   strconv.Itoa(unionLiveResp.Body.SaleMinAmount),
		SaleDiscount:    strconv.Itoa(unionLiveResp.Body.SaleDiscount),
		ActualPayAmount: actualPayAmount,
		Txamt:           req.Txamt,
		PayType:         req.PayType,
		Cardbin:         req.Cardbin,
	}

	return scanPayResponse
}

// ProcessRecoverCoupons 电子券验证冲正
func (u *unionliveScanPay) ProcessRecoverCoupons(req *model.ScanPayRequest) *model.ScanPayResponse {
	unionLiveReq := &coupon.RecoverCouponsReq{
		Header: coupon.RecoverCouponsReqHeader{
			Version:       Version,
			TransDirect:   TransDirectQ,
			TransType:     "W493",
			MerchantId:    req.ChanMerId,
			SubmitTime:    req.CreateTime,
			ClientTraceNo: req.OrderNum,
		},
		Body: coupon.RecoverCouponsReqBody{
			CouponsNo: req.OrigScanCodeId,
			TermId:    req.Terminalid,
			// TermSn:    req.Terminalsn,
			ExtMercId:        req.Mchntid,
			ExtTermId:        req.Terminalsn,
			Amount:           req.OrigVeriTime,
			Cardbin:          req.OrigCardbin,
			TransAmount:      req.IntTxamt,
			PayType:          req.IntPayType,
			OldSubmitTime:    req.OrigSubmitTime,
			OldClientTraceNo: req.OrigOrderNum,
		},
		SpReq: req,
	}
	unionLiveResp := &coupon.RecoverCouponsResp{}
	err := Execute(unionLiveReq, unionLiveResp)
	if err != nil {
		log.Errorf("sendRequest fail, orderNum=%s, service=RecoverCoupons, channel=ULIVE", req.OrderNum)
		return adaptor.ReturnWithErrorCode("SYSTEM_ERROR")
	}

	// 将渠道的错误应答码转为为系统应答码
	returncode, errDetail := transChanToSysCode(unionLiveResp.Header.ReturnCode, unionLiveResp.Header.ReturnMessage)

	// 处理结果返回
	scanPayResponse := &model.ScanPayResponse{
		Txndir:          unionLiveResp.Header.TransDirect,
		Busicd:          req.Busicd,
		Respcd:          returncode,
		AgentCode:       req.AgentCode,
		Chcd:            req.Chcd,
		Mchntid:         req.Mchntid,
		ErrorDetail:     errDetail,
		OrderNum:        unionLiveResp.Header.ClientTraceNo,
		ChanRespCode:    unionLiveResp.Header.ReturnCode,
		ChannelOrderNum: unionLiveResp.Header.HostTraceNo,
		// Terminalid:      req.Terminalsn,
		Authcode:     unionLiveResp.Body.AuthCode,
		ChannelTime:  unionLiveResp.Header.HostTime,
		OrigOrderNum: req.OrigOrderNum,
	}

	return scanPayResponse
}
