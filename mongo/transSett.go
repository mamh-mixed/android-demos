package mongo

import (
	"errors"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type transSettCollection struct {
	name string
}

// transSettLogCollection
// 记录某台机器某个时间点执行的方法
type transSettLogCollection struct {
	name string
}

var TransSettColl = transSettCollection{"transSett"}
var SpTransSettColl = transSettCollection{"transSett.sp"}

// settlement log
var TransSettLogColl = transSettLogCollection{"transSettLog"}

// AtomUpsert mongodb-findAndModify
func (col *transSettLogCollection) AtomUpsert(l *model.TransSettLog) (int, error) {
	q := bson.M{
		"method":    l.Method,
		"date":      l.Date,
		"$isolated": 1,
	}
	c := mgo.Change{}
	// 开始时status==0
	// 只更新method,date
	if l.Status == 0 {
		c.Update = bson.M{
			"$set": bson.M{
				"method": l.Method,
				"date":   l.Date,
			},
		}
	} else {
		// 结束时
		c.Update = l
	}
	c.Upsert = true

	result := new(model.TransSettLog)
	change, err := database.C(col.name).Find(q).Apply(c, result)

	return change.Updated, err

}

// Summary 清算信息汇总
func (col *transSettCollection) Summary(merId, transDate string) ([]model.SummarySettData, error) {

	//根据商户号、清算时间查找成功清算交易的汇总信息
	var s []model.SummarySettData
	//使用pipe统计
	err := database.C(col.name).Pipe([]bson.M{
		{"$match": bson.M{
			"merId": merId,
			"createTime": bson.M{"$gte": transDate,
				"$lt": util.NextDay(transDate),
			},
			"settFlag": 1,
		}},
		{"$group": bson.M{
			"_id":           "$transType",
			"totalTransAmt": bson.M{"$sum": "$transAmt"},
			"totalSettAmt":  bson.M{"$sum": "$merSettAmt"},
			"totalMerFee":   bson.M{"$sum": "$merFee"},
			"totalTransNum": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{"transType": "$_id",
			"totalTransAmt": 1,
			"totalSettAmt":  1,
			"totalMerFee":   1,
			"totalTransNum": 1}},
	}).All(&s)

	return s, err
}

// Add 增加一条清分记录
func (col *transSettCollection) Add(t *model.TransSett) error {
	if t == nil {
		return errors.New("transSett is nil")
	}
	t.SettDate = time.Now().Format("2006-01-02 15:04:05")
	return database.C(col.name).Insert(t)
}

// BatchAdd 批量添加
func (col *transSettCollection) BatchAdd(ts []model.TransSett) (err error) {
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
		// log.Infof("insert transSett [%d, %d)", s, e)
	}
	return nil
}

// BatchRemove 按照时间批量删除，多用于对账数据重新生成
func (col *transSettCollection) BatchRemove(date string) (err error) {

	q := bson.M{
		"settDate": date,
	}
	_, err = database.C(col.name).RemoveAll(q)
	return err
}

// Find 根据商户Id,清分时间查找交易明细
// 按照商户订单号降排序
func (col *transSettCollection) FindByDate(merId, transDate, nextOrderNum string) ([]model.TransSettInfo, error) {

	var transSettInfo []model.TransSettInfo

	p := []bson.M{
		//查找
		{"$match": bson.M{"merId": merId, "settFlag": 1,
			"createTime": bson.M{"$gte": transDate, "$lt": util.NextDay(transDate)}}},
		//排序
		{"$sort": bson.M{"orderNum": 1}},
	}
	//商户实际拉取为10
	limit := bson.M{"$limit": 11}
	if nextOrderNum != "" {
		p = append(p, bson.M{"$match": bson.M{"orderNum": bson.M{"$gte": nextOrderNum}}}, limit)
	} else {
		p = append(p, limit)
	}
	err := database.C(col.name).Pipe(p).All(&transSettInfo)
	return transSettInfo, err
}

// FindByOrderNum 根据渠道订单号查找
func (col *transSettCollection) FindByOrderNum(sysOrderNum string) (t *model.TransSett, err error) {
	// 订单是uuid 全局唯一
	t = new(model.TransSett)
	q := bson.M{
		"sysOrderNum": sysOrderNum,
	}
	err = database.C(col.name).Find(q).One(t)

	return
}

// FindOne 根据交易订单号、渠道订单号查找唯一记录
func (col *transSettCollection) FindOne(orderNum, relatedOrderNum string) (t *model.TransSett, err error) {
	// 订单是uuid 全局唯一
	t = new(model.TransSett)
	q := bson.M{
		"$or": []bson.M{
			bson.M{
				"trans.orderNum":     orderNum,
				"trans.chanOrderNum": relatedOrderNum,
			},
			bson.M{
				"trans.orderNum":     orderNum,
				"trans.origOrderNum": relatedOrderNum,
			},
		},
	}
	err = database.C(col.name).Find(q).One(t)

	return
}

// FindOrders 根据交易订单号、渠道订单号查询相应的订单和源订单号匹配的交易
func (col *transSettCollection) FindOrders(orderNum, chanOrderNum string) ([]model.TransSett, error) {
	// 订单是uuid 全局唯一
	var tss []model.TransSett
	q := bson.M{
		"$or": []bson.M{
			bson.M{"trans.orderNum": orderNum, "trans.chanOrderNum": chanOrderNum},
			bson.M{"trans.origOrderNum": orderNum, "trans.chanOrderNum": chanOrderNum},
		},
	}
	err := database.C(col.name).Find(q).All(&tss)

	return tss, err
}

// FindOne 根据交易订单号、渠道订单号删除记录
func (col *transSettCollection) RemoveOne(orderNum, chanOrderNum string) error {
	// 订单是uuid 全局唯一
	q := bson.M{
		"trans.orderNum":     orderNum,
		"trans.chanOrderNum": chanOrderNum,
	}
	return database.C(col.name).Remove(q)
}

// Update 更新
func (col *transSettCollection) Update(t *model.TransSett) error {
	if t == nil {
		return errors.New("transSett is nil")
	}
	// t.SettDate = time.Now().Format("2006-01-02 15:04:05")
	return database.C(col.name).Update(bson.M{"trans._id": t.Trans.Id}, t)
}

// Find
func (col *transSettCollection) Find(q *model.QueryCondition) ([]model.TransSett, error) {

	find := bson.M{}
	if q.MerId != "" {
		find["trans.merId"] = q.MerId
	}
	if q.AgentCode != "" {
		find["trans.agentCode"] = q.AgentCode
	}
	if q.SubAgentCode != "" {
		find["trans.subAgentCode"] = q.SubAgentCode
	}
	if q.GroupCode != "" {
		find["trans.groupCode"] = q.GroupCode
	}
	if q.SettRole != "" {
		find["settRole"] = q.SettRole
	}
	if q.StartTime != "" && q.EndTime != "" {
		find["trans.payTime"] = bson.M{"$gte": q.StartTime, "$lte": q.EndTime}
	}
	if q.Date != "" {
		find["settDate"] = q.Date
	}

	var ts []model.TransSett
	var err error

	if q.IsForReport {
		err = database.C(col.name).Find(find).Sort("-trans.payTime").All(&ts)
		return ts, err
	}

	err = database.C(col.name).Find(find).Sort("-trans.payTime").Skip((q.Page - 1) * q.Size).Limit(q.Size).All(&ts)
	return ts, err
}

// FindAndGroupBy 按照商户号分组统计
func (col *transSettCollection) FindAndGroupBy(q *model.QueryCondition) ([]model.TransGroup, []model.Channel, error) {

	var group []model.TransGroup

	find := bson.M{}
	if q.Date != "" {
		find["settDate"] = q.Date
	}
	if q.MerId != "" {
		find["trans.merId"] = q.MerId
	}
	if q.AgentCode != "" {
		find["trans.agentCode"] = q.AgentCode
	}
	if q.SubAgentCode != "" {
		find["trans.subAgentCode"] = q.SubAgentCode
	}
	if q.GroupCode != "" {
		find["trans.groupCode"] = q.GroupCode
	}
	find["blendType"] = 0 // match

	pipeline := []bson.M{
		{"$match": find},
		{"$group": bson.M{
			"_id": bson.M{"merId": "$trans.merId", "chanCode": "$trans.chanCode"},
			"transAmt": bson.M{"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$trans.transType", 1}}, "$trans.transAmt",
					bson.M{"$subtract": []interface{}{
						0, "$trans.transAmt",
					}},
				},
			}},
			"transNum": bson.M{"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$trans.transType", 1}}, 1, 0},
			}}, //只记录支付的笔数
			"merName":     bson.M{"$last": "$trans.merName"},
			"agentName":   bson.M{"$last": "$trans.agentName"},
			"groupCode":   bson.M{"$last": "$trans.groupCode"},
			"groupName":   bson.M{"$last": "$trans.groupName"},
			"companyName": bson.M{"$last": "$trans.subAgentName"},
			"fee": bson.M{"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$trans.transType", 1}}, "$merFee",
					bson.M{"$subtract": []interface{}{
						0, "$merFee",
					}},
				},
			}},
		}},
		{"$group": bson.M{
			"_id":         "$_id.merId",
			"transAmt":    bson.M{"$sum": "$transAmt"},
			"transNum":    bson.M{"$sum": "$transNum"},
			"merName":     bson.M{"$first": "$merName"},
			"agentName":   bson.M{"$first": "$agentName"},
			"groupCode":   bson.M{"$first": "$groupCode"},
			"companyName": bson.M{"$first": "$companyName"},
			"groupName":   bson.M{"$first": "$groupName"},
			"detail": bson.M{"$push": bson.M{"chanCode": "$_id.chanCode",
				"transNum": "$transNum",
				"transAmt": "$transAmt",
				"fee":      "$fee",
			}},
		}},
	}

	err := database.C(col.name).Pipe(pipeline).All(&group)
	if err != nil {
		return nil, nil, err
	}

	// 按渠道汇总所有符合条件数据
	var all []model.Channel
	err = database.C(col.name).Pipe([]bson.M{
		{"$match": find},
		{"$group": bson.M{"_id": "$trans.chanCode",
			"transAmt": bson.M{"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$trans.transType", 1}}, "$trans.transAmt",
					bson.M{"$subtract": []interface{}{
						0, "$trans.transAmt",
					}},
				},
			}},
			"transNum": bson.M{"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$trans.transType", 1}}, 1, 0},
			}}, //只记录支付的笔数
			"fee": bson.M{"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$trans.transType", 1}}, "$merFee",
					bson.M{"$subtract": []interface{}{
						0, "$merFee",
					}},
				},
			}},
		}},
		{"$project": bson.M{
			"chanCode": "$_id",
			"transAmt": 1,
			"transNum": 1,
			"fee":      1,
		}},
	}).All(&all)
	return group, all, err
}

// GroupBySettRole 根据清算角色出报表
func (col *transSettCollection) GroupBySettRole(settDate string) ([]model.SettRoleGroup, error) {

	var result []model.SettRoleGroup
	err := database.C(col.name).Pipe([]bson.M{
		{"$match": bson.M{
			"settDate": settDate,
		}},
		{"$group": bson.M{"_id": bson.M{"merId": "$trans.merId", "settRole": "$trans.settRole"},
			"transAmt": bson.M{"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$trans.transType", 1}}, "$trans.transAmt",
					bson.M{"$subtract": []interface{}{
						0, "$trans.transAmt", // 相当于将逆向交易的金额变为负数
					}},
				},
			}},
			"fee": bson.M{"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$trans.transType", 1}}, "$trans.fee",
					bson.M{"$subtract": []interface{}{
						0, "$trans.fee",
					}},
				},
			}},
		}},
		{"$group": bson.M{
			"_id": "$_id.settRole",
			"detail": bson.M{"$push": bson.M{"merId": "$_id.merId",
				"transAmt": "$transAmt",
				"fee":      "$fee",
			}},
		}},
		{"$project": bson.M{
			"settRole": "$_id",
			"mers":     "$detail",
		}},
	}).All(&result)
	return result, err
}

func (col *transSettCollection) FindAll() ([]model.TransSett, error) {
	var result []model.TransSett
	err := database.C(col.name).Find(nil).All(&result)
	return result, err
}
