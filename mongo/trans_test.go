package mongo

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/CardInfoLink/log"
	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2/bson"
)

func TestFindLastRecord(t *testing.T) {
	lt, err := SpTransColl.FindLastRecord(&model.QueryCondition{MerId: "999118880000003", StartTime: "2016-01-01 00:00:00"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%+v", lt)
}

func TestFindToSett(t *testing.T) {
	SpTransColl.FindToSett("2015-12-07")
}

func TestBeginWith(t *testing.T) {
	var target = "GC-12u972398789"
	var prefix = "GC-"

	var result = strings.HasPrefix(target, prefix)

	t.Logf("hasPrefx result is %s", result)

	t.Logf("substring is %s", target[3:len(target)])
}

func TestTransFindAndGroupBy(t *testing.T) {

	q := &model.QueryCondition{
		StartTime:   "2016-01-01 00:00:00",
		EndTime:     "2016-01-05 23:59:59",
		TransStatus: []string{model.TransSuccess},
		// TransType:   model.PayTrans,
		// IsAggregateByGroup: true,
		// RefundStatus:       model.TransRefunded,
		// MerIds:       []string{"999118880000312"},
		Page: 1,
		Size: 20,
	}
	t.Logf("%+v", q)
	ss, all, total, err := SpTransColl.FindAndGroupBy(q)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%+v	%+v", ss, all)
	t.Log(len(ss), total)
}

func TestFindTransQuery(t *testing.T) {

	q := &model.QueryCondition{
		StartTime: "2015-12-01 00:00:00",
		EndTime:   "2015-12-31 23:59:59",
		MerId:     "999118880000018",
		Skip:      0,
		Page:      1,
		Size:      100,
		// TransStatus: []string{model.TransSuccess},
	}

	transInfo, total, err := SpTransColl.Find(q)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("total : %d", total)
	t.Logf("act total: %d", len(transInfo))
	for _, v := range transInfo {
		t.Logf("%s,%s", v.OrderNum, v.CreateTime)
	}

	// t.Logf("%d", len(transInfo))
}

func TestTransAdd(t *testing.T) {
	if debug {
		trans := &model.Trans{
			TransStatus: transStatus,
			MerId:       merId,
			OrderNum:    orderNum,
			TransType:   int8(transType),
		}
		log.Debugf("%+v", TransColl)

		err := TransColl.Add(trans)
		if err != nil {
			t.Errorf("add trans unsunccessful: %s", err)
			t.FailNow()
		}
		log.Debugf("add trans success %s", trans)
	}
}

func TestTransUpdate(t *testing.T) {
	hi := "56248e894fde83cf36000001"
	objectId := bson.ObjectIdHex(hi)
	on := time.Now().Unix()
	t.Logf("OrderNo is %d", on)
	trans := &model.Trans{
		// CreateTime:  time.Now().Unix(),
		Id:          objectId,
		MerId:       merId,
		OrderNum:    "100000000182979",
		TransType:   int8(transType),
		TransStatus: transStatus,
		DiscountAmt: 1,
	}
	err := TransColl.Update(trans)
	if err != nil {
		t.Errorf("modify trans unsunccessful: %s", err)
		t.FailNow()
	}
	log.Debugf("modify trans success %+v", trans)

}

func TestFindByAccount(t *testing.T) {
	trans, err := SpTransColl.FindByAccount("obFqRjp1bVfXrDYWT1JLf-k2vAek")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%s %s", trans.MerId, trans.OrderNum)
	n, err := NotifyRecColl.FindOne(trans.MerId, trans.OrderNum)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%+v", trans)
	t.Logf("%+v", n.FromChanMsg)
}

func TestCountTrans(t *testing.T) {

	c, err := TransColl.Count(merId, orderNum)
	if err != nil {
		t.Errorf("count trans unsunccessful: %s", err)
		t.FailNow()
	}
	log.Debugf("count trans success %d", c)
}

func TestFindTrans(t *testing.T) {
	trans, err := TransColl.FindOne(merId, orderNum)
	if err != nil {
		t.Errorf("find trans unsunccessful: %s", err)
		t.FailNow()
	}
	log.Debugf("find trans success %s", trans)
}

func TestFindByTime(t *testing.T) {
	trans, err := TransColl.FindByTime(createTime)
	if err != nil {
		t.Errorf("find trans unsunccessful: %s", err)
		t.FailNow()
	}
	log.Debugf("find trans success %s", trans)
}

func TestFindTransRefundAmt(t *testing.T) {
	total, err := TransColl.FindTransRefundAmt("1000000000002", "DqfTuPvvvTWDfD0Ke9DGOqbT")
	if err != nil {
		t.FailNow()
	}
	log.Debug(total)
}

func TestAgentProfit(t *testing.T) {
	data, err := SpTransColl.ExportAgentProfit("2015-11-31 23:59:59", "2015-11-01 00:00:00")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("find %d data", len(data))
	// t.Logf("%+v", data)

	excel := xlsx.NewFile()
	var row *xlsx.Row
	var cell *xlsx.Cell

	sheet, _ := excel.AddSheet("扫码数据统计")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "交易日期"
	cell = row.AddCell()
	cell.Value = "机构代码"
	cell = row.AddCell()
	cell.Value = "商户编号"
	cell = row.AddCell()
	cell.Value = "渠道商户编号"
	cell = row.AddCell()
	cell.Value = "受理商"
	cell = row.AddCell()
	cell.Value = "交易笔数"
	cell = row.AddCell()
	cell.Value = "交易金额"
	cell = row.AddCell()
	cell.Value = "交易渠道"

	var wxpMers = make(map[string]*model.ChanMer)

	for _, d := range data {

		var agentMerId string
		if cm, ok := wxpMers[d.ID.ChanMerId]; ok {
			if cm.IsAgentMode && cm.AgentMer != nil {
				agentMerId = cm.AgentMer.ChanMerId
			} else {
				agentMerId = cm.ChanMerId
			}
		} else {
			cm, err := ChanMerColl.Find(d.ID.ChanCode, d.ID.ChanMerId)
			if err == nil {
				wxpMers[d.ID.ChanMerId] = cm
			}
		}

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = d.ID.Date
		cell = row.AddCell()
		cell.Value = d.AgentCode
		cell = row.AddCell()
		cell.Value = d.ID.MerId
		cell = row.AddCell()
		cell.Value = d.ID.ChanMerId
		cell = row.AddCell()
		cell.Value = agentMerId
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%d", d.TransNum)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%0.2f", float64(d.TransAmt-d.RefundAmt)/100)
		cell = row.AddCell()
		cell.Value = d.ID.ChanCode
	}

	var buf []byte
	bf := bytes.NewBuffer(buf)
	// 写到buf里
	excel.Write(bf)

	// 上传到七牛
	err = qiniu.Put("201511_trans.xlsx", int64(bf.Len()), bf)
	if err != nil {
		t.Error(err)
	}

}
