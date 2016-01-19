package query

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
	"time"
)

// CouponTransQuery 卡券查询
func CouponTransQuery(q *model.QueryCondition) (ret *model.ResultBody) {

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

	log.Debugf("condition is %#v", q)

	results, total, err := mongo.CouTransColl.Find(q)
	if err != nil {
		log.Errorf("find trans error: %s", err)
		return model.NewResultBody(101, "查询卡券核销列表失败")
	}

	count := len(results)
	log.Debugf("total is %d", count)
	if count == 0 {
		results = make([]*model.Trans, 0, 0)
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  q.Page,
		Total: total,
		Size:  q.Size,
		Count: count,
		Data:  results,
	}

	ret = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return ret
}
