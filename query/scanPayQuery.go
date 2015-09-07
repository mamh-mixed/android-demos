package query

import (
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"time"
)

var noMerCode, noMerMsg = mongo.ScanPayRespCol.Get8583CodeAndMsg("NO_MERCHANT")
var sysErrCode, sysErrMsg = mongo.ScanPayRespCol.Get8583CodeAndMsg("SYSTEM_ERROR")

// GetOrderInfo 扫固定码订单信息
func GetOrderInfo(uniqueId string) scanFixedResponse {

	var response = scanFixedResponse{Response: "00"}
	m, err := mongo.MerchantColl.FindByUniqueId(uniqueId)
	if err != nil {
		response.Response, response.ErrorDetail = noMerCode, noMerMsg
		return response
	}

	response.MerID = m.MerId
	response.TitleOne = m.Detail.TitleOne
	response.TitleTwo = m.Detail.TitleTwo

	// find
	trans, count, err := mongo.SpTransColl.Find(&model.QueryCondition{
		TradeFrom:   "wap",
		TransStatus: model.TransSuccess,
		Busicd:      model.Jszf,
		MerId:       m.MerId,
		Size:        150,
		Page:        1,
	})

	if err != nil {
		response.Response = sysErrCode
		response.ErrorDetail = sysErrMsg
		return response
	}

	data := make([]scanFixedData, len(trans))
	for _, t := range trans {
		fd := scanFixedData{}
		fd.Amount = ""
		fd.Chcd = t.ChanCode
		fd.Headimgurl = t.HeadImgUrl
		fd.Nickname = t.NickName
		fd.Transtime = t.CreateTime
		fd.VeriCode = t.VeriCode
		fd.OrderNum = t.OrderNum
		data = append(data, fd)
	}

	response.Data = data
	response.Count = count

	return response
}

// GetMerInfo 扫固定码获取用户信息
func GetMerInfo(merId string) scanFixedResponse {

	var response = scanFixedResponse{Response: "00", MerID: merId}

	m, err := mongo.MerchantColl.Find(merId)
	if err != nil {
		response.Response = noMerCode
		response.ErrorDetail = noMerMsg
		response.MerID = ""
		return response
	}
	response.TitleOne = m.Detail.TitleOne
	response.TitleTwo = m.Detail.TitleTwo
	return response
}

// TransQuery 交易查询
func TransQuery(q *model.QueryCondition) (ret *model.QueryResult) {

	now := time.Now().Format("2006-01-02")
	// 默认当天开始
	if q.StartTime == "" {
		q.StartTime = now + " 00:00:00"
	} else {
		q.StartTime += " 00:00:00"
	}
	// 默认当天结束
	if q.EndTime == "" {
		q.EndTime = now + " 23:59:59"
	} else {
		q.EndTime += " 23:59:59"
	}

	// mongo统计
	trans, total, err := mongo.SpTransColl.Find(q)
	if err != nil {
		log.Errorf("find trans error: %s", err)
	}

	// 没有数组的话返回空数据
	if len(trans) == 0 {
		trans = make([]*model.Trans, 0, 0)
	}

	count := len(trans)
	ret = &model.QueryResult{
		Page:     q.Page,
		Size:     q.Size,
		Total:    total,
		RespCode: "000000",
		RespMsg:  "成功",
		Rec:      trans,
		Count:    count,
	}

	return ret
}

// TransStatistics 交易统计
func TransStatistics(q *model.QueryCondition) (ret *model.QueryResult) {

	errResult := &model.QueryResult{RespCode: "000001", RespMsg: "系统错误，请重试。"}

	// 设置条件过滤
	q.TransStatus = model.TransSuccess
	q.TransType = model.PayTrans
	q.RefundStatus = model.TransRefunded
	// q.MerIds = merIds
	q.StartTime += " 00:00:00"
	q.EndTime += " 23:59:59"

	// 查询交易
	now := time.Now()
	group, all, total, err := mongo.SpTransColl.FindAndGroupBy(q)
	if err != nil {
		log.Errorf("find trans error: %s", err)
		return errResult
	}
	after := time.Now()
	log.Debugf("spent %s", after.Sub(now))
	var data = make([]model.Summary, 0)

	// 将数据合并
	for _, d := range group {
		s := model.Summary{
			MerId:     d.MerId,
			AgentName: d.AgentName,
			MerName:   d.MerName,
		}
		// 遍历渠道，合并数据
		combine(&s, d.Detail)
		data = append(data, s)
	}

	// 汇总数据
	summary := model.Summary{Data: data}
	combine(&summary, all)

	// 组装返回报文
	count := len(data)
	ret = &model.QueryResult{
		Page:  q.Page,
		Size:  q.Size,
		Total: total,
		Rec:   summary,
		Count: count,
	}

	return ret

}

func combine(s *model.Summary, detail []model.Channel) {
	for _, d := range detail {
		switch d.ChanCode {
		case channel.ChanCodeAlipay:
			s.Alp.TransAmt = float32(d.TransAmt-d.RefundAmt) / 100
			s.Alp.TransNum = d.TransNum
			s.Alp.Fee = float32(d.Fee) / 100
			s.TotalTransAmt += s.Alp.TransAmt
			s.TotalTransNum += s.Alp.TransNum
			s.TotalFee += s.Alp.Fee
		case channel.ChanCodeWeixin:
			s.Wxp.TransAmt = float32(d.TransAmt-d.RefundAmt) / 100
			s.Wxp.TransNum = d.TransNum
			s.Wxp.Fee = float32(d.Fee) / 100
			s.TotalTransAmt += s.Wxp.TransAmt
			s.TotalTransNum += s.Wxp.TransNum
			s.TotalFee += s.Wxp.Fee
		}
	}
}

type scanFixedResponse struct {
	Response    string          `json:"response"`
	MerID       string          `json:"merID"`
	TitleOne    string          `json:"title_one"`
	TitleTwo    string          `json:"title_two"`
	ErrorDetail string          `json:"errorDetail,omitempty"`
	Data        []scanFixedData `json:"data"`
	Count       int             `json:"count"`
}

type scanFixedData struct {
	Transtime  string `json:"transtime,omitempty"`
	VeriCode   string `json:"veriCode,omitempty"`
	Nickname   string `json:"nickname,omitempty"`
	Headimgurl string `json:"headimgurl,omitempty"`
	Amount     string `json:"amount,omitempty"`
	OrderNum   string `json:"orderNum,omitempty"`
	Chcd       string `json:"chcd,omitempty"` //ALP,WXP
}
