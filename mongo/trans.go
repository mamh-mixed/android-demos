package mongo

import (
	"errors"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/tools"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
)

type transCollection struct {
	name string
}

// TransColl 交易 Collection
var TransColl = transCollection{"trans"}

// Add 添加一笔交易
func (col *transCollection) Add(t *model.Trans) error {
	// default
	t.Id = bson.NewObjectId()
	t.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	t.TransStatus = "00"
	err := database.C(col.name).Insert(t)
	if err != nil {
		log.Error("add trans(%+v) fail: %s", t, err)
	}
	return err
}

// Update 通过Add时生成的Id来修改
func (col *transCollection) Update(t *model.Trans) error {
	t.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	err := database.C(col.name).Update(bson.M{"_id": t.Id}, t)
	if err != nil {
		log.Error("update trans(%+v) fail: %s", t, err)
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
func (col *transCollection) Find(merId, orderNum string) (t *model.Trans, err error) {

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
			"$lte": tools.NextDay(time),
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
		"transType":      model.RefundTrans,
		"merId":          merId,
		"refundOrderNum": origOrderNum,
		"transStatus":    model.TransSuccess,
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
