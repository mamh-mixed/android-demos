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

	// 先找商户，按商户分页
	mers, total, err := mongo.MerchantColl.FuzzyFind(q)
	if err != nil {
		log.Errorf("find merchant error: %s", err)
		return errResult
	}
	var merIds []string
	var summarys []*model.Summary
	m := make(map[string]*model.Summary)
	for _, mer := range mers {
		merIds = append(merIds, mer.MerId)
		summary := &model.Summary{MerId: mer.MerId}
		m[mer.MerId] = summary
		summarys = append(summarys, summary)
	}

	q.TransStatus = model.TransSuccess
	q.TransType = model.PayTrans
	q.MerIds = merIds
	// log.Debugf("%+v", q)
	// 查询交易
	data, err := mongo.SpTransColl.FindAndGroupBy(q)
	if err != nil {
		log.Errorf("find trans error: %s", err)
		return errResult
	}
	log.Debug(data)
	// 将数据合并
	for _, d := range data {
		if s, ok := m[d.Key.MerId]; ok {
			switch d.Key.ChanCode {
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
	size := len(merIds)
	ret = &model.QueryResult{
		Page:     q.Page,
		Size:     size,
		Total:    total,
		RespCode: "000000",
		RespMsg:  "成功",
		Rec:      summarys,
		Count:    size,
	}

	return ret

}
