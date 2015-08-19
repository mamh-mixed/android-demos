package mongo

import (
	"errors"
	"fmt"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type transCollection struct {
	name string
}

// TransColl 交易 Collection
var TransColl = transCollection{"trans"}

// SpTransColl 扫码交易 Collection
var SpTransColl = transCollection{"trans.sp"}

// Add 添加一笔交易
func (col *transCollection) Add(t *model.Trans) error {
	// default
	t.Id = bson.NewObjectId()
	t.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	if t.TransStatus == "" {
		t.TransStatus = model.TransFail
	}
	err := database.C(col.name).Insert(t)
	if err != nil {
		log.Errorf("add trans(%+v) fail: %s", t, err)
	}
	return err
}

// BatchAdd 批量添加
func (col *transCollection) BatchAdd(ts []*model.Trans) error {
	var temp []interface{}
	for _, t := range ts {
		temp = append(temp, t)
	}
	return database.C(col.name).Insert(temp...)
}

// Update 通过Add时生成的Id来修改
func (col *transCollection) Update(t *model.Trans) error {
	t.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	err := database.C(col.name).Update(bson.M{"_id": t.Id}, t)
	if err != nil {
		log.Errorf("update trans(%+v) fail: %s", t, err)
	}
	return err
}

// Count 通过订单号、商户号查找交易数量
func (col *transCollection) Count(merId, orderNum string) (count int, err error) {

	if merId == "" {
		return 0, errors.New("商户Id为空。")
	}

	if orderNum == "" {
		return 0, errors.New("订单号为空。")
	}
	// 不需要交易类型，对同一个商户下的订单号不能重复
	// if transType != 1 && transType != 2 {
	// 	return 0, errors.New("交易类型错误。")
	// }

	q := bson.M{
		"orderNum": orderNum,
		"merId":    merId,
		// "transType": transType,
	}
	count, err = database.C(col.name).Find(q).Count()
	if err != nil {
		log.Errorf("find trans(%s,%s) fail: (%s)", merId, orderNum, err)
	}
	return
}

// FindPay 通过订单号、商户号查找一条交易记录
func (col *transCollection) FindOne(merId, orderNum string) (t *model.Trans, err error) {

	q := bson.M{
		"orderNum": orderNum,
		"merId":    merId,
		// "transType": transType,
	}
	t = new(model.Trans)
	err = database.C(col.name).Find(q).One(t)
	return
}

// FindByTime 查找某天的交易记录
func (col *transCollection) FindByTime(time string) ([]*model.Trans, error) {

	var ts []*model.Trans
	q := bson.M{
		"createTime": bson.M{
			"$gt":  time,
			"$lte": util.NextDay(time),
		},
	}
	err := database.C(col.name).Find(q).All(&ts)
	return ts, err
}

// FindRefundTrans 查找某个订单成功的退款
func (col *transCollection) FindTransRefundAmt(merId, origOrderNum string) (int64, error) {

	var s = &struct {
		Amt int64 `bson:"refundedAmt"`
	}{}
	q := bson.M{
		"transType":    model.RefundTrans,
		"merId":        merId,
		"origOrderNum": origOrderNum,
		"transStatus":  model.TransSuccess,
	}

	err := database.C(col.name).Pipe([]bson.M{
		{"$match": q},
		{"$group": bson.M{
			"_id":         "$merId",
			"refundedAmt": bson.M{"$sum": "$transAmt"},
		}},
		{"$project": bson.M{"refundedAmt": 1}},
	}).One(s)

	if err != nil {
		if err.Error() == "not found" {
			err = nil
		} else {
			log.Errorf("find refund trans error : %s", err)
		}
	}
	return s.Amt, err
}

// FindByOrderNum 根据渠道订单号查找
func (col *transCollection) FindByOrderNum(sysOrderNum string) (t *model.Trans, err error) {
	// 订单是uuid 全局唯一
	t = new(model.Trans)
	q := bson.M{
		"sysOrderNum": sysOrderNum,
	}
	err = database.C(col.name).Find(q).One(t)

	return
}

// UpdateFields 更新指定字段
func (col *transCollection) UpdateFields(t *model.Trans) error {

	if t.Id == "" {
		return fmt.Errorf("%s", "id is null!")
	}
	// update fields
	fields := bson.M{
		"updateTime": time.Now().Format("2006-01-02 15:04:05"),
	}
	if t.MerDiscount != "" {
		fields["merDiscount"] = t.MerDiscount
	}
	if t.ChanDiscount != "" {
		fields["chanDiscount"] = t.ChanDiscount
	}
	// more fields

	return database.C(col.name).Update(bson.M{"_id": t.Id}, bson.M{"$set": fields})
}

// Find 根据商户Id,清分时间查找交易明细
// 按照商户订单号降排序
func (col *transCollection) Find(q *model.QueryCondition) ([]model.Trans, int, error) {

	log.Debugf("condition is %+v", q)

	var trans []model.Trans

	// 根据条件查找
	match := bson.M{}
	if q.OrderNum != "" {
		match["orderNum"] = q.OrderNum
	}
	if q.MerId != "" {
		match["merId"] = q.MerId
	}
	if q.Busicd != "" {
		match["busicd"] = q.Busicd
	}
	if q.OrigOrderNum != "" {
		match["origOrderNum"] = q.OrigOrderNum
	}
	// or 退款的和成功的
	or := []bson.M{}
	if q.TransStatus != "" {
		or = append(or, bson.M{"transStatus": q.TransStatus})
	}
	if q.RefundStatus != 0 {
		or = append(or, bson.M{"refundStatus": q.RefundStatus})
	}
	if len(or) > 0 {
		match["$or"] = or
	}
	match["createTime"] = bson.M{"$gte": q.StartTime, "$lt": q.EndTime}

	// 将取消订单原交易不成功的过滤掉，如果原交易不成功则取消这笔订单的金额为0
	match["transAmt"] = bson.M{"$ne": 0}

	p := []bson.M{
		{"$match": match},
	}

	// total
	total, err := database.C(col.name).Find(match).Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页
	skip := bson.M{"$skip": (q.Page - 1) * q.Size}

	// 不同类型排序
	sort := bson.M{"$sort": bson.M{"createTime": -1}}

	// 商户实际拉取为Size+1
	limit := bson.M{"$limit": q.Size}

	// 如果是导出报表
	if q.IsForReport {
		sortByChan := bson.M{"$sort": bson.M{"chanCode": 1}}
		sort = bson.M{"$sort": bson.M{"busicd": 1}}
		p = append(p, sort, skip, limit, sortByChan)
	} else {
		p = append(p, sort, skip, limit)
	}

	err = database.C(col.name).Pipe(p).All(&trans)
	return trans, total, err
}

// FindAndGroupBy 统计
func (col *transCollection) FindAndGroupBy(q *model.QueryCondition) ([]model.TransGroup, []model.Channel, int, error) {

	var group []model.TransGroup

	find := bson.M{
		"createTime": bson.M{"$gte": q.StartTime, "$lt": q.EndTime},
		"merId":      bson.M{"$in": q.MerIds},
		"transType":  q.TransType,
		"$or":        []bson.M{bson.M{"transStatus": q.TransStatus}, bson.M{"refundStatus": q.RefundStatus}},
	}

	// 计算total
	var total = struct {
		Value int `bson:"total"`
	}{}
	database.C(col.name).Pipe([]bson.M{
		{"$match": find},
		{"$group": bson.M{"_id": "$merId"}},
		{"$group": bson.M{"_id": "null", "total": bson.M{"$sum": 1}}},
	}).One(&total)

	//使用pipe统计
	err := database.C(col.name).Pipe([]bson.M{
		{"$match": find},
		{"$group": bson.M{
			"_id":       bson.M{"merId": "$merId", "chanCode": "$chanCode"},
			"transAmt":  bson.M{"$sum": "$transAmt"},
			"refundAmt": bson.M{"$sum": "$refundAmt"},
			"transNum":  bson.M{"$sum": 1},
			"fee":       bson.M{"$sum": "$fee"},
		}},
		{"$group": bson.M{
			"_id":       "$_id.merId",
			"refundAmt": bson.M{"$sum": "$refundAmt"},
			"transAmt":  bson.M{"$sum": "$transAmt"},
			"transNum":  bson.M{"$sum": "$transNum"},
			"detail": bson.M{"$push": bson.M{"chanCode": "$_id.chanCode",
				"transNum":  "$transNum",
				"transAmt":  "$transAmt",
				"refundAmt": "$refundAmt",
				"fee":       "$fee",
			}},
		}},
		{"$sort": bson.M{"transNum": -1}},
		{"$skip": (q.Page - 1) * q.Size},
		{"$limit": q.Size},
	}).All(&group)

	// 按渠道汇总所有符合条件数据
	var all []model.Channel
	err = database.C(col.name).Pipe([]bson.M{
		{"$match": find},
		{"$group": bson.M{"_id": "$chanCode",
			"transAmt":  bson.M{"$sum": "$transAmt"},
			"refundAmt": bson.M{"$sum": "$refundAmt"},
			"transNum":  bson.M{"$sum": 1},
			"fee":       bson.M{"$sum": "$fee"},
		}},
		{"$project": bson.M{
			"chanCode":  "$_id",
			"transAmt":  1,
			"transNum":  1,
			"refundAmt": 1,
			"fee":       1,
		}},
	}).All(&all)

	return group, all, total.Value, err
}
