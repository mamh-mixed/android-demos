// 用于导入旧扫码系统交易记录
package data

import (
	"math"
	"net/http"
	"strconv"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2/bson"
)

var (
	inprocess              = mongo.ScanPayRespCol.Get("INPROCESS")
	success                = mongo.ScanPayRespCol.Get("SUCCESS")
	CloseCode, CloseMsg, _ = mongo.ScanPayRespCol.Get8583CodeAndMsg("ORDER_CLOSED")
)

const (
	crypto = "cilxl123$"
)

// func init() {
// 	url = "mongodb://saoma:saoma@211.147.72.70:10006/online"
// 	connect()
// }

func Import(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if crypto != key {
		return
	}

	st := r.FormValue("st")
	et := r.FormValue("et")
	t := r.FormValue("type")

	switch t {
	case "updTrans":
		go func() {
			err := UpdateTrans()
			if err != nil {
				log.Error(err)
			}
		}()
	case "trans":
		go func() {
			err := AddTransFromOldDB(st, et)
			if err != nil {
				log.Error(err)
			}
		}()
	case "merchant":
		go func() {
			err := DoSyncMerchant("/Users/migo/Desktop/product_pem/")
			if err != nil {
				log.Error(err)
			}
		}()
	case "app":
		go func() {
			err := AsyncAppUser()
			if err != nil {
				log.Error(err)
			}
		}()
	case "settInfo":
		go func() {
			err := AsyncSettInfo()
			if err != nil {
				log.Error(err)
			}
		}()
	}

	w.Write([]byte("已开始导入，请查看后台日志"))
}

// txn 交易表数据
type txn struct {
	Date    string `bson:"gw_date"`
	Time    string `bson:"gw_time"`
	Request struct {
		Attach       string `bson:"attach"`
		BackUrl      string `bson:"backUrl"`
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
		TradeFrom    string `bson:"tradeFrom"`
		Code         string `bson:"code"`
		VeriCode     string `bson:"veriCode"`
	} `bson:"m_request"`
	Merchant        merchant `bson:"merchant"`
	ChanRespCode    string   `bson:"resposeDetail"`
	RespCode        string   `bson:"response"`
	ChannelOrderNum string   `bson:"channelOrderNum"`
	Qrcode          string   `bson:"qrcode"`
	ChcdDiscount    string   `bson:"chcdDiscount"`
	MerDiscount     string   `bson:"merDiscount"`
	ConsumerAccount string   `bson:"consumerAccount"`
	ConsumerId      string   `bson:"consumerId"`
	Status          string   `bson:"status"`
	Payjson         struct {
		UserInfo struct {
			OpenID     string `bson:"openid"`
			NickName   string `bson:"nickname"`
			Sex        int    `bson:"sex"`
			Language   string `bson:"language"`
			City       string `bson:"city"`
			Province   string `bson:"province"`
			Country    string `bson:"country"`
			Headimgurl string `bson:"headimgurl"`
		} `bson:"userinfo"`
	} `bson:"payjson"`
}

// UpdateTrans 更新交易来源
func UpdateTrans() error {
	txns, err := findTransFromOldDB()
	if err != nil {
		return err
	}

	log.Infof("从旧系统找到有交易来源标识的交易 %d 条，正在更新新系统的交易...", len(txns))

	var updated int
	for _, ot := range txns {
		nt, err := mongo.SpTransColl.FindOne(ot.Merchant.Clientid, ot.Request.OrderNum)
		if err != nil {
			log.Errorf("没找到新系统交易 merId=%s, orderNum=%s", ot.Merchant.Clientid, ot.Request.OrderNum)
			continue
		}

		if ot.Request.Code != "" && nt.Busicd != model.Jszf {
			nt.Busicd = model.Jszf
		}
		nt.GatheringId = ot.Payjson.UserInfo.OpenID
		nt.NickName = ot.Payjson.UserInfo.NickName
		nt.HeadImgUrl = ot.Payjson.UserInfo.Headimgurl
		nt.VeriCode = ot.Request.VeriCode
		mongo.SpTransColl.Update(nt)
		updated++
	}

	log.Infof("已成功更新 %d 交易记录", updated)
	return nil
}

func AddTransFromOldDB(st, et string) error {
	// TODO 判断新系统是否包含该天的数据，如果包含，那么报错，避免数据紊乱

	agents, err := readAgentFromOldDB()
	if err != nil {
		return err
	}
	agentsMap := make(map[string]string)
	for _, a := range agents {
		agentsMap[a.AgentCode] = a.AgentName
	}

	txns, err := readTransFromOldDB(st, et)
	if err != nil {
		return err
	}

	if len(txns) == 0 {
		log.Warn("没有找到符合条件数据。。")
		return nil
	}
	log.Infof("从老系统取出 %d 条符合条件交易数据，正在逻辑处理。。。", len(txns))

	// 存放退款、撤销、取消的交易
	var reversalTrans []*model.Trans
	// 存放正常交易
	payTransMap := make(map[string]*model.Trans)

	// 先整理数据并且放在内存里
	for _, t := range txns {
		tran := &model.Trans{}
		tran.Id = bson.NewObjectId()
		tran.Attach = t.Request.Attach
		tran.NotifyUrl = t.Request.BackUrl
		tran.TradeFrom = t.Request.TradeFrom
		tran.MerId = t.Merchant.Clientid
		tran.Terminalid = t.Request.Terminalid
		tran.AgentCode = t.Merchant.AgentCode
		tran.OrderNum = t.Request.OrderNum
		tran.OrigOrderNum = t.Request.OrigOrderNum
		tran.CreateTime = t.Date + " " + t.Time
		tran.UpdateTime = tran.CreateTime
		tran.MerDiscount = t.MerDiscount
		tran.ChanDiscount = t.ChcdDiscount
		tran.ChanCode = t.Request.Chcd
		tran.Remark = "old_system_trans"
		tran.GoodsInfo = t.Request.GoodsInfo
		tran.TransCurr = t.Request.Currency
		tran.ChanOrderNum = t.ChannelOrderNum
		tran.RespCode = t.RespCode
		tran.MerName = t.Merchant.ClientidName
		tran.AgentName = agentsMap[tran.AgentCode]
		tran.GroupCode = t.Merchant.Group.GroupCode
		tran.GroupName = t.Merchant.Group.GroupName

		// JSZF一些字段
		tran.GatheringId = t.Payjson.UserInfo.OpenID
		tran.NickName = t.Payjson.UserInfo.NickName
		tran.HeadImgUrl = t.Payjson.UserInfo.Headimgurl
		tran.VeriCode = t.Request.VeriCode

		// 旧系统没有JSZF类型
		tran.Busicd = t.Request.Busicd
		if tran.Busicd == model.Purc {
			// 如果是PURC类型，并且code不为空，那么就是JSZF
			if t.Request.Code != "" {
				tran.Busicd = model.Jszf
			}
		}

		if tran.ChanCode == "ALP" {
			tran.ChanMerId = t.Merchant.Alp.PartnerId
			// 讯联清算
			if t.Merchant.Alp.Type == "1" {
				tran.SettRole = "CIL"
			} else {
				tran.SettRole = "ALP"
			}
		} else {
			if t.Merchant.Wxp.SubMchId != "" {
				tran.ChanMerId = t.Merchant.Wxp.SubMchId
			} else {
				tran.ChanMerId = t.Merchant.Wxp.MchId
			}
			if t.Merchant.Wxp.Type == "1" {
				tran.SettRole = "CIL"
			} else {
				tran.SettRole = "WXP"
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
		case model.Qyzf:
			tran.TransType = model.EnterpriseTrans
			payTransMap[tran.MerId+tran.OrderNum] = tran
		case model.Purc, model.Paut, model.Jszf:
			tran.TransType = model.PayTrans
			// 计算费率
			var merFee float64
			if tran.ChanCode == "ALP" {
				merFee, err = strconv.ParseFloat(t.Merchant.Alp.MerFee, 64)
				if err != nil {
					log.Errorf("商户号：%s，支付宝手续费：%s，转换错误：%s", tran.MerId, t.Merchant.Alp.MerFee, err)
				}
			} else {
				merFee, err = strconv.ParseFloat(t.Merchant.Wxp.MerFee, 64)
				if err != nil {
					log.Errorf("商户号：%s，微信手续费：%s，转换错误：%s", tran.MerId, t.Merchant.Wxp.MerFee, err)
				}
			}
			tran.MerFee = merFee
			tran.Fee = int64(math.Floor(float64(tran.TransAmt)*merFee + 0.5))
			tran.NetFee = tran.Fee
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
	var effectTrans []*model.Trans
	for _, t := range reversalTrans {
		var orig *model.Trans
		var isGetFromDB bool
		// 拿到原交易
		if tran, ok := payTransMap[t.MerId+t.OrigOrderNum]; ok {
			orig = tran
		} else {
			// 从数据库获取
			// log.Infof("从数据库获取原订单，商户号：%s，订单号：%s", t.MerId, t.OrigOrderNum)
			tran, err := mongo.SpTransColl.FindOne(t.MerId, t.OrigOrderNum)
			if err != nil {
				// log.Errorf("从内存和数据库里获取不到原交易，商户号：%s，订单号：%s", t.MerId, t.OrigOrderNum)
				continue
			}
			isGetFromDB = true
			orig = tran
		}

		var handled bool
		if orig.TransStatus == model.TransClosed || orig.RefundStatus != 0 {
			handled = true
		}

		// 计算手续费
		if orig.TransStatus == model.TransSuccess {
			t.Fee = int64(math.Floor(float64(t.TransAmt)*orig.MerFee + 0.5))
			orig.NetFee = orig.NetFee - t.Fee
		}

		// 具体处理
		switch t.Busicd {
		case model.Refd:
			// 累计退款
			refundedAmt := t.TransAmt + orig.RefundAmt
			// 对原交易逻辑处理
			if refundedAmt >= orig.TransAmt {
				// 全额退款
				// orig.RespCode = CloseCode
				// orig.ErrorDetail = CloseMsg
				orig.RefundAmt = orig.TransAmt
				orig.RefundStatus = model.TransRefunded
				orig.TransStatus = model.TransClosed
				// orig.Fee = 0
			} else if refundedAmt < orig.TransAmt {
				// 部分退款
				orig.RefundAmt = refundedAmt
				orig.RefundStatus = model.TransPartRefunded
			}
		case model.Canc:
			// 判断原交易是否成功
			if orig.TransStatus == model.TransSuccess {
				t.TransAmt = orig.TransAmt
				orig.RefundStatus = model.TransRefunded
				orig.RefundAmt = orig.TransAmt
				orig.TransStatus = model.TransClosed
				break
			}
			// orig.RespCode = CloseCode
			// orig.ErrorDetail = CloseMsg
			orig.TransStatus = model.TransClosed
			t.TransAmt = 0 // 如果原交易不成功，则交易金额为0
			// orig.Fee = 0

		case model.Void:
			// 相当于全额退款
			t.TransAmt = orig.TransAmt
			// orig.RespCode = CloseCode
			// orig.ErrorDetail = CloseMsg
			orig.RefundAmt = orig.TransAmt
			orig.RefundStatus = model.TransRefunded
			orig.TransStatus = model.TransClosed
			// orig.Fee = 0
		}

		// 如果是从数据库拿的，那么更新数据库里的数据
		if isGetFromDB {
			if !handled {
				log.Infof("更新数据库原订单，商户号：%s，订单号：%s", t.MerId, t.OrigOrderNum)
				mongo.SpTransColl.Update(orig)
			} else {
				log.Infof("数据已被更新，商户号：%s，订单号：%s", t.MerId, t.OrigOrderNum)
			}
		}

		effectTrans = append(effectTrans, t)
	}
	// 将map里的数据放到数组里
	for _, v := range payTransMap {
		effectTrans = append(effectTrans, v)
	}

	log.Infof("从老系统拿出%d条数据，成功处理%d条数据。正在导入数据库。。。", len(txns), len(effectTrans))
	// log.Infof("%+v", reversalTrans)
	// 批量入库
	return mongo.SpTransColl.BatchAdd(effectTrans)
	// return nil
}

func findTransFromOldDB() ([]txn, error) {
	q := bson.M{"gw_date": bson.M{"$lte": "2015-10-20"}, "payjson.userinfo.headimgurl": bson.M{"$exists": true}, "response": "00"}
	var txns []txn
	err := saomaDB.C("txn").Find(q).All(&txns)
	return txns, err
}

func readTransFromOldDB(startTime, endTime string) ([]txn, error) {

	q := bson.M{"gw_date": bson.M{"$gte": startTime, "$lte": endTime}, "response": "00"}

	var txns []txn
	err := saomaDB.C("txn").Find(q).All(&txns)
	return txns, err
}
