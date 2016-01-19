package query

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
	"time"
)

// BpTransQuery 绑定支付交易查询
func BpTransQuery(q *model.QueryCondition) (ret *model.ResultBody) {

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
	trans, total, err := mongo.TransColl.Find(q)
	if err != nil {
		log.Errorf("find trans error: %s", err)
	}

	// 没有数组的话返回空数据
	count := len(trans)
	if count == 0 {
		trans = make([]*model.Trans, 0, 0)
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  q.Page,
		Total: total,
		Size:  q.Size,
		Count: count,
		Data:  trans,
	}

	ret = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return ret
}
