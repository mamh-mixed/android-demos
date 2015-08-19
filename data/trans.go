// 用于导入旧扫码系统交易记录
package data

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
	"math"
	"strconv"
)

var (
	inprocess           = mongo.ScanPayRespCol.Get("INPROCESS")
	success             = mongo.ScanPayRespCol.Get("SUCCESS")
	CloseCode, CloseMsg = mongo.ScanPayRespCol.Get8583CodeAndMsg("ORDER_CLOSED")
)

// {
//     "_id" : ObjectId("55d3d89fcfa871957f82f4dd"),
//     "gw_date" : "2015-08-19",
//     "gw_time" : "09:15:11",
//     "system_date" : "20150819091511",
//     "current_time" : 1.439946911748E12,
//     "m_request" : {
//         "busicd" : "PAUT",
//         "txndir" : "Q",
//         "inscd" : "90711888",
//         "chcd" : "ALP",
//         "mchntid" : "907118885840003",
//         "terminalid" : "90710004",
//         "txamt" : "000000000001",
//         "orderNum" : "231091115456"
//     },
//     "merchant" : {
//         "_id" : ObjectId("55bb08a82309f012f8b00a6e"),
//         "clientid" : "907118885840003",
//         "commodityName" : "上海讯联测试商",
//         "WXP" : {
//             "md5" : "42Ugz8i5OAI44MDJqjfXyX7juIGzP4Es",
//             "mch_id" : "10013970",
//             "appid" : "wx8c2e6e8f0f46d469",
//             "acqfee" : "0.02",
//             "merfee" : "0.03",
//             "fee" : "0.01"
//         },
//         "ALP" : {
//             "partnerId" : "2088811767473826",
//             "md5" : "eu1dr0c8znpa43blzy1wirzmk8jqdaon",
//             "acqfee" : "0.02",
//             "merfee" : "0.03",
//             "fee" : "0.03"
//         },
//         "clientidName" : "test",
//         "inscd" : "10134001",
//         "signRule" : "0",
//         "merchantMd5" : "",
//         "insMd5" : "",
//         "signType" : "sha1"
//     },
//     "front_response_to_g" : "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<alipay><is_success>T</is_success><request><param name=\"subject\">上海讯联测试商</param><param name=\"sign_type\">MD5</param><param name=\"notify_url\">http://211.147.72.70:10003/back?schema=ALP-907118885840003231091115456</param><param name=\"out_trade_no\">231091115456</param><param name=\"currency\">156</param><param name=\"sign\">e128a62b2856c36f65a6f145dba0f607</param><param name=\"_input_charset\">utf-8</param><param name=\"it_b_pay\">20m</param><param name=\"product_code\">QR_CODE_OFFLINE</param><param name=\"total_fee\">0.01</param><param name=\"service\">alipay.acquire.precreate</param><param name=\"seller_id\">2088811767473826</param><param name=\"partner\">2088811767473826</param></request><response><alipay><big_pic_url>https://mobilecodec.alipay.com/show.htm?code=bal9or45m1vwrug092&amp;d&amp;picSize=L</big_pic_url><out_trade_no>231091115456</out_trade_no><pic_url>https://mobilecodec.alipay.com/show.htm?code=bal9or45m1vwrug092&amp;d&amp;picSize=M</pic_url><qr_code>https://qr.alipay.com/bal9or45m1vwrug092</qr_code><result_code>SUCCESS</result_code><small_pic_url>https://mobilecodec.alipay.com/show.htm?code=bal9or45m1vwrug092&amp;d&amp;picSize=S</small_pic_url><voucher_type>qrcode</voucher_type></alipay></response><sign>2ec3bb7b9524eed7fb0732912666529f</sign><sign_type>MD5</sign_type></alipay>",
//     "resposeDetail" : "ORDER_SUCCESS_PAY_INPROCESS",
//     "response" : "00",
//     "channelOrderNum" : "2015081921001004970078059703",
//     "qrcode" : "https://qr.alipay.com/bal9or45m1vwrug092",
//     "chcdDiscount" : "0.00",
//     "merDiscount" : "0.00",
//     "back_response_to_g" : {
//         "schema" : "ALP-907118885840003231091115456",
//         "subject" : "上海讯联测试商",
//         "trade_no" : "2015081921001004970078059703",
//         "paytools_pay_amount" : "[{\"ALIPAYACCOUNT\":\"0.01\"}]",
//         "buyer_email" : "15071440565@163.com",
//         "gmt_create" : "2015-08-19 09:15:59",
//         "notify_type" : "trade_status_sync",
//         "quantity" : "1",
//         "out_trade_no" : "231091115456",
//         "notify_time" : "2015-08-19 09:16:05",
//         "seller_id" : "2088811767473826",
//         "trade_status" : "TRADE_SUCCESS",
//         "total_fee" : "0.01",
//         "gmt_payment" : "2015-08-19 09:16:05",
//         "seller_email" : "andy.li@cardinfolink.com",
//         "price" : "0.01",
//         "buyer_id" : "2088702949897971",
//         "notify_id" : "d294a517c778130dc662adc5b693dc7nhg",
//         "sign_type" : "MD5",
//         "sign" : "eb9175c85596598838f828cfafa81310"
//     },
//     "notify_action_type" : null,
//     "consumerAccount" : "15071440565@163.com",
//     "consumerId" : "2088702949897971",
//     "status" : "1"
// }

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

func AddTransFromOldDB() error {
	connect()

	// TODO 判断新系统是否包含该天的数据，如果包含，那么报错，避免数据紊乱

	txns, err := readTransFromOldDB("2015-08-18", "2015-08-19")
	if err != nil {
		return err
	}
	log.Debugf("%+v", txns[0])

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
