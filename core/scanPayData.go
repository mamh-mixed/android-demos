package core

import (
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"time"
)

// TransQuery 交易查询
func TransQuery(q *model.QueryCondition) (ret *model.QueryResult) {

	now := time.Now().Format("2006-01-02")
	// 默认当天开始
	if q.StartTime == "" {
		q.StartTime = now + " 00:00:00"
	}
	// 默认当天结束
	if q.EndTime == "" {
		q.EndTime = now + " 23:59:59"
	}

	// mongo统计
	trans, total, err := mongo.SpTransColl.Find(q)
	if err != nil {
		log.Errorf("find trans error: %s", err)
	}

	size := len(trans)
	ret = &model.QueryResult{
		Page:     q.Page,
		Size:     size,
		Total:    total,
		RespCode: "000000",
		RespMsg:  "成功",
		Rec:      trans,
		Count:    size,
	}

	return ret
}

// TransStatistics 交易统计
func TransStatistics(q *model.QueryCondition) (ret *model.QueryResult) {

	errResult := &model.QueryResult{RespCode: "000001", RespMsg: "系统错误，请重试。"}

	// 先找商户所有商户号
	mers, err := mongo.MerchantColl.FuzzyFind(q)
	if err != nil {
		log.Errorf("find merchant error: %s", err)
		return errResult
	}
	var merIds []string
	m := make(map[string]*model.Merchant)
	for _, mer := range mers {
		merIds = append(merIds, mer.MerId)
		m[mer.MerId] = mer
	}

	q.TransStatus = model.TransSuccess
	q.TransType = model.PayTrans
	q.MerIds = merIds
	// 查询交易
	group, all, total, err := mongo.SpTransColl.FindAndGroupBy(q)
	if err != nil {
		log.Errorf("find trans error: %s", err)
		return errResult
	}
	var data []model.Summary
	// 将数据合并
	for _, d := range group {
		if mer, ok := m[d.MerId]; ok {
			s := model.Summary{
				MerId:     d.MerId,
				AgentName: mer.AgentName,
				MerName:   mer.Detail.MerName,
			}
			// 遍历渠道，合并数据
			combine(&s, d.Detail)
			data = append(data, s)
		}
	}

	// 汇总数据
	summary := model.Summary{Data: data}
	combine(&summary, all)

	size := len(data)
	ret = &model.QueryResult{
		Page:  q.Page,
		Size:  size,
		Total: total,
		Rec:   summary,
		Count: size,
	}

	return ret

}

func combine(s *model.Summary, detail []model.Channel) {
	for _, d := range detail {
		switch d.ChanCode {
		case channel.ChanCodeAlipay:
			s.Alp.TransAmt = float64(d.TransAmt-d.RefundAmt) / 100
			s.Alp.TransNum = d.TransNum
			s.TotalTransAmt += s.Alp.TransAmt
			s.TotalTransNum += s.Alp.TransNum
		case channel.ChanCodeWeixin:
			s.Wxp.TransAmt = float64(d.TransAmt-d.RefundAmt) / 100
			s.Wxp.TransNum = d.TransNum
			s.TotalTransAmt += s.Wxp.TransAmt
			s.TotalTransNum += s.Wxp.TransNum
		}
	}
}
