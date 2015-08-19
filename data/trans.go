// 用于导入旧扫码系统交易记录
package data

// import (
// 	"github.com/CardInfoLink/quickpay/model"
// 	"github.com/CardInfoLink/quickpay/mongo"
// )

// {
//     "_id" : ObjectId("55d2feedc3509bdc0408d3b6"),
//     "gw_date" : "2015-08-18",
//     "gw_time" : "17:46:21",
//     "system_date" : "20150818174621",
//     "current_time" : 1.439891181095E12,
//     "m_request" : {
//         "busicd" : "PAUT",
//         "txndir" : "Q",
//         "inscd" : "90711888",
//         "chcd" : "WXP",
//         "mchntid" : "907118885840003",
//         "terminalid" : "90710004",
//         "txamt" : "000000000001",
//         "orderNum" : "230174115301"
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
//     "front_response_to_g" : "<xml><return_code><![CDATA[SUCCESS]]></return_code>\n<return_msg><![CDATA[OK]]></return_msg>\n<appid><![CDATA[wx8c2e6e8f0f46d469]]></appid>\n<mch_id><![CDATA[10013970]]></mch_id>\n<nonce_str><![CDATA[7wDQ8lxR2a1Zvdeg]]></nonce_str>\n<sign><![CDATA[08DCA991907BA2416B9A192F938CB9D1]]></sign>\n<result_code><![CDATA[SUCCESS]]></result_code>\n<prepay_id><![CDATA[wx201508181746360f0ef1ab640954478079]]></prepay_id>\n<trade_type><![CDATA[NATIVE]]></trade_type>\n<code_url><![CDATA[weixin://wxpay/bizpayurl?pr=WXgOrqo]]></code_url>\n</xml>",
//     "resposeDetail" : "ORDER_SUCCESS_PAY_INPROCESS",
//     "response" : "09",
//     "channelOrderNum" : "",
//     "qrcode" : "weixin://wxpay/bizpayurl?pr=WXgOrqo"
// }

// txn 交易表数据
type txn struct {
	Date    string
	Time    string
	Request struct {
		Busicd     string `bson:"busicd"`
		Txndir     string `bson:"txndir"`
		Inscd      string `bson:"inscd"`
		Chcd       string `bson:"chcd"`
		Mchntid    string `bson:"mchntid"`
		Terminalid string `bson:"terminalid"`
		Txamt      string `bson:"txamt"`
		OrderNum   string `bson:"orderNum"`
	}
}

func readTransFromOldDB() {

}
