package query

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"time"
)

var noMerCode, noMerMsg = mongo.ScanPayRespCol.Get8583CodeAndMsg("NO_MERCHANT")
var sysErrCode, sysErrMsg = mongo.ScanPayRespCol.Get8583CodeAndMsg("SYSTEM_ERROR")
var sucCode, sucMsg = mongo.ScanPayRespCol.Get8583CodeAndMsg("SUCCESS")

func GetBills(q *model.QueryCondition) (result *model.QueryResult) {

	result = &model.QueryResult{RespCode: sucCode, RespMsg: sucMsg}
	// 限制最大1000条
	var maxRec = 1000

	// 拉取1000+1条用于返回最后记录订单号
	q.Size = maxRec + 1

	trans, err := mongo.SpTransColl.FindByNextRecord(q)
	if err != nil {
		result.RespCode, result.RespMsg = sysErrCode, sysErrMsg
		return result
	}

	type rec struct {
		OrderNum  string `json:"orderNum"`
		TransType int8   `json:"transType"`
		TransTime string `json:"transTime"`
		TransAmt  int64  `json:"transAmt"`
	}

	tSize := len(trans)
	var recs []rec
	if tSize > 0 {
		result.Count = tSize
		// 取满一页
		if tSize == q.Size {
			result.Count = maxRec
			result.NextOrderNum = trans[maxRec].OrderNum
			trans = trans[:maxRec]
		}
		for _, t := range trans {

			// 交易类型
			var transType int8
			switch t.Busicd {
			case model.Purc:
				transType = 1
			case model.Paut:
				transType = 2
			case model.Jszf:
				transType = 3
			case model.Refd:
				transType = 5
			case model.Void:
				transType = 6
			case model.Qyzf:
				transType = 7
			case model.Canc:
				transType = 8
			}

			r := rec{t.OrderNum, transType, t.CreateTime, t.TransAmt}
			recs = append(recs, r)
		}
		result.Rec = recs
	}

	return result
}

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
		TradeFrom:    "wap",
		TransStatus:  []string{model.TransSuccess},
		Busicd:       model.Jszf,
		MerId:        m.MerId,
		RefundStatus: model.TransRefunded,
		Size:         150,
		Page:         1,
	})

	if err != nil {
		response.Response = sysErrCode
		response.ErrorDetail = sysErrMsg
		return response
	}

	var data []scanFixedData
	for _, t := range trans {
		fd := scanFixedData{}
		fd.Amount = fmt.Sprintf("%0.2f", float64(t.TransAmt)/100)
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

func GetSpTransLogs(q *model.QueryCondition, msgType int) ([]model.SpTransLogs, int, error) {

	var spLogs []model.SpTransLogs
	var err error
	var total int

	switch msgType {
	case 1:
		// total
		total, err = mongo.SpMerLogsCol.Count(q)
		if err != nil {
			return nil, 0, err
		}

		// 先查来的报文
		q.Direction = "in"
		inSpLogs, err := mongo.SpMerLogsCol.Find(q)
		if err != nil {
			return nil, 0, err
		}
		spLogs = append(spLogs, inSpLogs...)

		// 再查返回的报文
		if len(inSpLogs) > 0 {
			q.Direction = "out"
			var reqIds []string
			for _, l := range inSpLogs {
				reqIds = append(reqIds, l.ReqId)
			}
			q.ReqIds = reqIds
			outSpLogs, err := mongo.SpMerLogsCol.Find(q)
			if err != nil {
				return nil, 0, err
			}
			spLogs = append(spLogs, outSpLogs...)
		}

	case 2:
		spLogs, err = mongo.SpChanLogsCol.Find(q)
		total = len(spLogs)
	}
	return spLogs, total, err
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
	response.AgentCode = m.AgentCode
	return response
}

// SpTransQuery 交易查询
func SpTransQuery(q *model.QueryCondition) (ret *model.QueryResult) {

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

// SpTransFindOne 交易查询
func SpTransFindOne(q *model.QueryCondition) (ret *model.ResultBody) {

	// mongo统计
	trans, err := mongo.SpTransColl.FindOneByOrigOrderNum(q)
	if err != nil {
		log.Errorf("find trans error: %s", err)
		ret = model.NewResultBody(1, "查询数据库失败")
		return ret
	}

	ret = &model.ResultBody{
		Status:  0,
		Message: "",
		Data:    trans,
	}
	return ret
}

// TransStatistics 交易统计
func TransStatistics(q *model.QueryCondition) (ret *model.QueryResult) {

	errResult := &model.QueryResult{RespCode: "000001", RespMsg: "系统错误，请重试。"}

	// 设置条件过滤
	q.TransStatus = []string{model.TransSuccess}
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
	AgentCode   string          `json:"inscd,omitempty"`
	TitleOne    string          `json:"title_one"`
	TitleTwo    string          `json:"title_two"`
	ErrorDetail string          `json:"errorDetail,omitempty"`
	Data        []scanFixedData `json:"data,omitempty"`
	Count       int             `json:"count,omitempty"`
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
