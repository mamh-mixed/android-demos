// 用于导入旧扫码系统交易记录
package data

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
	"math"
	"net/http"
	"strconv"
)

var (
	inprocess           = mongo.ScanPayRespCol.Get("INPROCESS")
	success             = mongo.ScanPayRespCol.Get("SUCCESS")
	CloseCode, CloseMsg = mongo.ScanPayRespCol.Get8583CodeAndMsg("ORDER_CLOSED")
)

const (
	crypto = "cilxl123$"
)

func init() {
	connect()
}

func ImportOldTrans(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if crypto != key {
		return
	}

	st := r.FormValue("st")
	et := r.FormValue("et")

	go func() {
		err := AddTransFromOldDB(st, et)
		log.Error(err)
	}()

	w.Write([]byte("已开始导入交易，请查看后台日志"))
}

// txn 交易表数据
type txn struct {
	Date    string `bson:"gw_date"`
	Time    string `bson:"gw_time"`
	Request struct {
		Busicd       string `bson:"busicd"`
		Txndir       string `bson:"txndir"`
		Inscd        string `bson:"inscd"`
		Chcd         string `bson:"chcd"`
		Mchntid      string `bson:"mchntid"`
		Terminalid   string `bson:"terminalid"`
		Txamt        string `bson:"txamt"`
		OrderNum     string `bson:"orderNum"`
		OrigOrderNum string `bson:"origOrderNum"`
		Currency     string `bson:"currency"`
		GoodsInfo    string `bson:"goodsInfo"`
	} `bson:"m_request"`
	Merchant struct {
		MerId string `bson:"clientid"`
		Inscd string `bson:"inscd"`
		Wxp   struct {
			MerFee   string `bson:"merFee"`
			MerId    string `bson:"mch_id"`
			SubMerId string `bson:"sub_mch_id"`
		} `bson:"WXP"`
		Alp struct {
			MerFee string `bson:"merFee"`
			MerId  string `bson:"partnerId"`
		} `bson:"ALP"`
	} `bson:"merchant"`
	ChanRespCode    string `bson:"resposeDetail"`
	RespCode        string `bson:"response"`
	ChannelOrderNum string `bson:"channelOrderNum"`
	Qrcode          string `bson:"qrcode"`
	ChcdDiscount    string `bson:"chcdDiscount"`
	MerDiscount     string `bson:"merDiscount"`
	ConsumerAccount string `bson:"consumerAccount"`
	ConsumerId      string `bson:"consumerId"`
	Status          string `bson:"status"`
}

func AddTransFromOldDB(st, et string) error {
	// TODO 判断新系统是否包含该天的数据，如果包含，那么报错，避免数据紊乱

	txns, err := readTransFromOldDB(st, et)
	if err != nil {
		return err
	}

	if len(txns) == 0 {
		log.Warn("没有找到符合条件数据。。")
		return nil
	}
	log.Debugf("从老系统取出 %d 条符合条件交易数据，正在逻辑处理。。。", len(txns))

	// 存放退款、撤销、取消的交易
	var reversalTrans []*model.Trans
	// 存放正常交易
	payTransMap := make(map[string]*model.Trans)

	// 先整理数据并且放在内存里
	for _, t := range txns {
		tran := &model.Trans{}
		tran.Id = bson.NewObjectId()
		tran.MerId = t.Merchant.MerId
		tran.AgentCode = t.Merchant.Inscd
		tran.OrderNum = t.Request.OrderNum
		tran.OrigOrderNum = t.Request.OrigOrderNum
		tran.CreateTime = t.Date + " " + t.Time
		tran.UpdateTime = tran.CreateTime
		tran.MerDiscount = t.MerDiscount
		tran.ChanDiscount = t.ChcdDiscount
		tran.Busicd = t.Request.Busicd
		tran.ChanCode = t.Request.Chcd
		tran.Remark = "old_system_trans"
		tran.GoodsInfo = t.Request.GoodsInfo
		tran.TransCurr = t.Request.Currency
		tran.ChanOrderNum = t.ChannelOrderNum
		if tran.ChanCode == "ALP" {
			tran.ChanMerId = t.Merchant.Alp.MerId
		} else {
			if t.Merchant.Wxp.SubMerId != "" {
				tran.ChanMerId = t.Merchant.Wxp.SubMerId
			} else {
				tran.ChanMerId = t.Merchant.Wxp.MerId
			}
		}

		// 金额转换
		toInt, _ := strconv.ParseInt(t.Request.Txamt, 10, 64)
		tran.TransAmt = toInt

		switch t.RespCode {
		case success.ISO8583Code:
			tran.ErrorDetail = success.ISO8583Msg
			tran.TransStatus = model.TransSuccess
		case inprocess.ISO8583Code:
			tran.ErrorDetail = inprocess.ISO8583Msg
			tran.TransStatus = model.TransHandling
		default:
			// 用原来的应答码，失败的应答码没什么意义
			tran.TransStatus = model.TransFail
		}

		// 交易类型
		switch tran.Busicd {
		case model.Qyfk:
			tran.TransType = model.EnterpriseTrans
			payTransMap[tran.MerId+tran.OrderNum] = tran
		case model.Purc, model.Paut, model.Jszf:
			tran.TransType = model.PayTrans
			// 计算费率
			var merFee float64
			if tran.ChanCode == "ALP" {
				merFee, _ = strconv.ParseFloat(t.Merchant.Alp.MerFee, 64)
			} else {
				merFee, _ = strconv.ParseFloat(t.Merchant.Wxp.MerFee, 64)
			}
			tran.MerFee = merFee
			tran.Fee = int64(math.Floor(float64(tran.TransAmt)*merFee + 0.5))
			payTransMap[tran.MerId+tran.OrderNum] = tran
		case model.Refd:
			tran.TransType = model.RefundTrans
			reversalTrans = append(reversalTrans, tran)
		case model.Canc:
			tran.TransType = model.CloseTrans
			reversalTrans = append(reversalTrans, tran)
		case model.Void:
			tran.TransType = model.CancelTrans
			reversalTrans = append(reversalTrans, tran)
		default:
			log.Errorf("未知的交易类型，交易数据：%+v", tran)
		}

	}

	// 对原交易处理，因为原交易没有退款等标识，所以得用逆向交易去实现逻辑
	for _, t := range reversalTrans {
		var orig *model.Trans
		var isGetFromDB bool
		// 拿到原交易
		if tran, ok := payTransMap[t.MerId+t.OrigOrderNum]; ok {
			orig = tran
		} else {
			// 从数据库获取
			log.Infof("从数据库获取原订单，商户号：%s，订单号：%s", t.MerId, t.OrigOrderNum)
			tran, err := mongo.SpTransColl.FindOne(t.MerId, t.OrigOrderNum)
			if err != nil {
				log.Errorf("从内存和数据库里获取不到原交易，商户号：%s，订单号：%s", t.MerId, t.OrigOrderNum)
				break
			}
			isGetFromDB = true
			orig = tran
		}
		// 具体处理
		switch t.Busicd {
		case model.Refd:

			// 累计退款
			refundedAmt := t.TransAmt + orig.RefundAmt
			// 对原交易逻辑处理
			if refundedAmt >= orig.TransAmt {
				// 全额退款
				orig.RespCode = CloseCode
				orig.ErrorDetail = CloseMsg
				orig.RefundAmt = orig.TransAmt
				orig.RefundStatus = model.TransRefunded
				orig.TransStatus = model.TransClosed
				orig.Fee = 0
			} else if refundedAmt < orig.TransAmt {
				// 部分退款
				orig.RefundAmt = refundedAmt
				orig.RefundStatus = model.TransPartRefunded
				// 重新计算手续费
				orig.Fee = int64(math.Floor(float64(orig.TransAmt-orig.RefundAmt))*orig.MerFee + 0.5)
			}

		case model.Canc:
			// 判断原交易是否成功
			if orig.TransStatus == model.TransSuccess {
				t.TransAmt = orig.TransAmt
				orig.RefundStatus = model.TransRefunded
				orig.RefundAmt = orig.TransAmt
			}
			orig.RespCode = CloseCode
			orig.ErrorDetail = CloseMsg
			orig.TransStatus = model.TransClosed
			orig.Fee = 0

		case model.Void:
			// 相当于全额退款
			t.TransAmt = orig.TransAmt
			orig.RespCode = CloseCode
			orig.ErrorDetail = CloseMsg
			orig.RefundAmt = orig.TransAmt
			orig.RefundStatus = model.TransRefunded
			orig.TransStatus = model.TransClosed
			orig.Fee = 0
		}

		// 如果是从数据库拿的，那么更新数据库里的数据
		if isGetFromDB {
			mongo.SpTransColl.Update(orig)
		}
	}
	// 将map里的数据放到数组里
	for _, v := range payTransMap {
		reversalTrans = append(reversalTrans, v)
	}

	log.Infof("从老系统拿出%d条数据，成功处理%d条数据。正在导入数据库。。。", len(txns), len(reversalTrans))

	// 批量入库
	return mongo.SpTransColl.BatchAdd(reversalTrans)
}

func readTransFromOldDB(startTime, endTime string) ([]txn, error) {

	q := bson.M{"gw_date": bson.M{"$gte": startTime, "$lte": endTime}}
	var txns []txn
	err := saomaDB.C("txn").Find(q).All(&txns)
	return txns, err
}
