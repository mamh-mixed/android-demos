package mongo

import (
	"errors"
	"fmt"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2"
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
func (col *transCollection) BatchAdd(ts []*model.Trans) (err error) {

	l := len(ts)
	var temp []interface{}
	for _, t := range ts {
		temp = append(temp, t)
	}

	// 一次最多一W
	batch := 10000
	for s, e := 0, batch; s < l; s, e = e, e+batch {
		if e > l {
			e = l
		}
		err = database.C(col.name).Insert(temp[s:e]...)
		if err != nil {
			return err
		}
		log.Infof("insert coupon [%d, %d)", s, e)
	}

	return nil
}

// FindAndLock
// 该方法查找时交易将交易锁住
// 如果锁住成功，将返回最新的交易
func (col *transCollection) FindAndLock(merId, orderNum string) (*model.Trans, error) {

	query := bson.M{
		"merId":    merId,
		"orderNum": orderNum,
		"lockFlag": bson.M{"$ne": 1}, // 此处不直接写为 lockFlag=0是为了兼容以前数据
	}

	change := mgo.Change{}
	change.Update = bson.M{
		"$set": bson.M{"lockFlag": 1,
			"updateTime": time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	change.ReturnNew = true
	result := &model.Trans{}
	_, err := database.C(col.name).Find(query).Apply(change, result)
	if err != nil {
		log.Error(err)
	}
	return result, err
}

// UpdateAndUnlock 更新并解锁
func (col *transCollection) UpdateAndUnlock(t *model.Trans) error {
	t.LockFlag = 0
	return col.Update(t)
}

// Unlock 只做解锁操作
func (col *transCollection) Unlock(merId, orderNum string) {
	set := bson.M{"$set": bson.M{"lockFlag": 0}}
	database.C(col.name).Update(bson.M{"merId": merId, "orderNum": orderNum}, set)
}

// Update 通过Add时生成的Id来修改
func (col *transCollection) Update(t *model.Trans) error {

	t.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	// 查找条件
	update := bson.M{"_id": t.Id}
	err := database.C(col.name).Update(update, t)
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

// FindOneByOrigOrderNum 通过订单号、商户号查找一条交易记录
func (col *transCollection) FindOneByOrigOrderNum(q *model.QueryCondition) (t *model.Trans, err error) {
	match := bson.M{
		"busicd":       q.Busicd,
		"origOrderNum": q.OrigOrderNum,
		"transStatus":  "30",
	}
	t = new(model.Trans)
	err = database.C(col.name).Find(match).One(t)

	return t, err
}

// FindHandingTrans 找到三十分钟前的处理中的交易
func (col *transCollection) FindHandingTrans() ([]model.Trans, error) {
	q := bson.M{
		"updateTime":  bson.M{"$lte": time.Now().Add(-30 * time.Minute).Format("2006-01-02 15:04:05")},
		"lockFlag":    0,
		"transStatus": model.TransHandling,
		"transType":   model.PayTrans,
	}

	var ts []model.Trans
	err := database.C(col.name).Find(q).Limit(1000).All(&ts)

	return ts, err
}

// FindByAccount 通过订单号、商户号查找一条交易记录
func (col *transCollection) FindByAccount(account string) (t *model.Trans, err error) {

	q := bson.M{
		"consumerAccount": account,
		// "transType": transType,
		"busicd":      model.Jszf,
		"transStatus": model.TransSuccess,
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

// GetBySettOrder 由结算订单号随机获取一条交易信息
func (col *transCollection) GetBySettOrder(merId, settOrderNum string) (*model.Trans, error) {
	q := bson.M{
		"merId":        merId,
		"transType":    model.PayTrans,
		"transStatus":  model.TransSuccess,
		"settOrderNum": settOrderNum,
	}

	result := new(model.Trans)
	err := database.C(col.name).Find(q).One(result)
	return result, err
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
func (col *transCollection) Find(q *model.QueryCondition) ([]*model.Trans, int, error) {
	var trans []*model.Trans

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
	if q.AgentCode != "" {
		match["agentCode"] = q.AgentCode
	}
	if q.GroupCode != "" {
		match["groupCode"] = q.GroupCode
	}
	if q.Respcd != "" {
		match["respCode"] = q.Respcd
	}
	if q.TradeFrom != "" {
		match["tradeFrom"] = q.TradeFrom
	}
	if q.TransType != 0 {
		match["transType"] = q.TransType
	}
	if q.BindingId != "" {
		match["bindingId"] = q.BindingId
	}
	// or 退款的和成功的
	or := []bson.M{}
	if len(q.TransStatus) != 0 {
		or = append(or, bson.M{"transStatus": bson.M{"$in": q.TransStatus}})
	}
	if q.RefundStatus != 0 {
		or = append(or, bson.M{"refundStatus": q.RefundStatus})
	}
	if len(or) > 0 {
		match["$or"] = or
	}
	if q.StartTime != "" && q.EndTime != "" {
		match["createTime"] = bson.M{"$gte": q.StartTime, "$lte": q.EndTime}
	}

	// 将取消订单原交易不成功的过滤掉，如果原交易不成功则取消这笔订单的金额为0
	match["transAmt"] = bson.M{"$ne": 0}

	p := []bson.M{
		{"$match": match},
	}

	log.Debugf("find condition: %#v", match)
	// total
	total, err := database.C(col.name).Find(match).Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页
	skipRecord := 0
	if q.Skip != 0 {
		skipRecord = q.Skip
	} else {
		skipRecord = (q.Page - 1) * q.Size
	}
	skip := bson.M{"$skip": skipRecord}

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
	}
	if q.MerId != "" {
		find["merId"] = bson.RegEx{q.MerId, "."}
	}
	if q.AgentCode != "" {
		find["agentCode"] = q.AgentCode
	}
	if q.MerName != "" {
		find["merName"] = bson.RegEx{q.MerName, "."}
	}
	find["transType"] = q.TransType
	find["$or"] = []bson.M{bson.M{"transStatus": bson.M{"$in": q.TransStatus}}, bson.M{"refundStatus": q.RefundStatus}}

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
			"merName":   bson.M{"$last": "$merName"},
			"agentName": bson.M{"$last": "$agentName"},
			"netFee":    bson.M{"$sum": "$netFee"}, // !!!这里计算的是净手续费，不是fee字段。
		}},
		{"$group": bson.M{
			"_id":       "$_id.merId",
			"refundAmt": bson.M{"$sum": "$refundAmt"},
			"transAmt":  bson.M{"$sum": "$transAmt"},
			"transNum":  bson.M{"$sum": "$transNum"},
			"merName":   bson.M{"$first": "$merName"},
			"agentName": bson.M{"$first": "$agentName"},
			"detail": bson.M{"$push": bson.M{"chanCode": "$_id.chanCode",
				"transNum":  "$transNum",
				"transAmt":  "$transAmt",
				"refundAmt": "$refundAmt",
				"fee":       "$netFee",
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
			"fee":       bson.M{"$sum": "$netFee"}, // !!!这里计算的是净手续费，不是fee字段。
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

// MerBills 商户账单
func (col *transCollection) MerBills(q *model.QueryCondition) ([]model.TransTypeGroup, error) {
	find := bson.M{
		"createTime": bson.M{"$gte": q.StartTime, "$lt": q.EndTime},
	}
	find["merId"] = q.MerId

	var or []bson.M
	if len(q.TransStatus) != 0 {
		or = append(or, bson.M{"transStatus": bson.M{"$in": q.TransStatus}})
	}
	if q.RefundStatus != 0 {
		or = append(or, bson.M{"refundStatus": q.RefundStatus})
	}
	if len(or) > 0 {
		find["$or"] = or
	}

	// 过滤掉取消不成功的订单
	find["transAmt"] = bson.M{"$ne": 0}

	var results []model.TransTypeGroup
	err := database.C(col.name).Pipe([]bson.M{
		{"$match": find},
		{"$group": bson.M{"_id": "$transType",
			"transAmt": bson.M{"$sum": "$transAmt"},
			// "refundAmt": bson.M{"$sum": "$refundAmt"},
			"transNum": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"transType": "$_id",
			"transAmt":  1,
			"transNum":  1,
			// "refundAmt": 1,
		}},
	}).All(&results)
	return results, err
}
