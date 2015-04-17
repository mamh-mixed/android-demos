package alp

// type Request
import (
	"encoding/xml"
)

// AlpResponse 支付宝接口返回信息
type AlpResponse struct {
	XMLName   xml.Name `xml:"alipay"`
	IsSuccess string   `xml:"is_success,omitempty"`
	Sign      string   `xml:"sign,omitempty"`
	SignType  string   `xml:"sign_type,omitempty"`
	Error     string   `xml:"error,omitempty"`
	// Alipay
	Response AlpBody `xml:"response,omitempty"`
}

type AlpBody struct {
	Alipay AlpDetail `xml:"alipay,omitempty"`
}

// AlpDetail response节点字段
type AlpDetail struct {
	BuyerLogonId    string          `xml:"buyer_logon_id,omitempty"`    //买家支付宝账号
	BuyerUserId     string          `xml:"buyer_user_id,omitempty"`     //买家支付宝用户号 以 2088 开头的纯 16 位数字
	OutTradeNo      string          `xml:"out_trade_no,omitempty"`      //商户网站唯一订单号
	ResultCode      string          `xml:"result_code,omitempty"`       //查询处理结果响应码。SUCCESS:查询成功 FAIL:查询失败 PROCESS_EXCEPTION:处理异常
	TradeNo         string          `xml:"trade_no,omitempty"`          //支付宝交易号 最短16位,最长64位
	DetailErrorCode string          `xml:"detail_error_code,omitempty"` //详细错误码
	DetailErrorDes  string          `xml:"detail_error_des,omitempty"`
	ExtendInfo      string          `xml:"extend_info,omitempty"`
	TradeStatus     string          `xml:"trade_status,omitempty"`                 //交易状态
	Partner         string          `xml:"partner,omitempty"`                      //合作者身份ID
	FundBillList    []TradeFundBill `xml:"fund_bill_list>TradeFundBill,omitempty"` //本次交易支付单据信息集合
	TotalFee        string          `xml:"total_fee,omitempty"`                    //订单金额
	SendPayDate     string          `xml:"send_pay_date,omitempty"`                //本次交易打款到卖家账户的时间,格式为 yyyy-MM-dd HH:mm:ss
}

// TradeFundBill 支付单据信息
type TradeFundBill struct {
	XMLName     xml.Name `xml:"TradeFundBill"`
	Amount      string   `xml:"amount,omitempty"`       //支付金额
	FundChannel string   `xml:"fund_channel,omitempty"` //支付渠道
}
