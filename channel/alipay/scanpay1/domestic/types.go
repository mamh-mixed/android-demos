package domestic

// type Request
import (
	"encoding/xml"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"strconv"
)

// AlpRequest 请求参数
type alpRequest struct {
	Service          string                `bson:"service,omitempty"`
	Partner          string                `bson:"partner,omitempty"`
	Charset          string                `bson:"charset,omitempty"`
	NotifyUrl        string                `bson:"notifyUrl,omitempty"`
	OutTradeNo       string                `bson:"outTradeNo,omitempty"`
	RefundAmount     string                `bson:"refundAmount,omitempty"`
	OutRequestNo     string                `bson:"outRequestNo,omitempty"`
	PassbackParams   string                `bson:"passBackParams,omitempty"`
	Subject          string                `bson:"subject,omitempty"`
	GoodsDetail      string                `bson:"goodsDetail,omitempty"`
	ProductCode      string                `bson:"productCode,omitempty"`
	TotalFee         string                `bson:"totalFee,omitempty"`
	SellerId         string                `bson:"sellerId,omitempty"`
	Currency         string                `bson:"currency,omitempty"`
	ExtendParams     string                `bson:"extendParams,omitempty"`
	ItBPay           string                `bson:"itBPay,omitempty"`
	DynamicIdType    string                `bson:"dynamicIdType,omitempty"`
	DynamicId        string                `bson:"dynamicId,omitempty"`
	Key              string                `bson:"-"`
	SpReq            *model.ScanPayRequest `bson:"-"`
	Gmt_create_start string                `bson:"gmt_create_start,omitempty"`
	Gmt_create_end   string                `bson:"gmt_create_end,omitempty"`
}

// AlpResponse 支付宝接口返回信息
type alpResponse struct {
	XMLName   xml.Name `xml:"alipay" bson:"-"`
	IsSuccess string   `xml:"is_success,omitempty"`
	Sign      string   `xml:"sign,omitempty"`
	SignType  string   `xml:"sign_type,omitempty"`
	Error     string   `xml:"error,omitempty"`
	Request   []Param  `xml:"request>param"`
	// Alipay
	Response alpBody `xml:"response,omitempty" bson:"response,omitempty"`
}

type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",innerxml"`
}

type alpBody struct {
	Alipay alpDetail `xml:"alipay,omitempty"`
}

// AlpDetail response节点字段
type alpDetail struct {
	BuyerLogonId    string          `xml:"buyer_logon_id,omitempty" bson:"buyer_logon_id,omitempty"` //买家支付宝账号
	BuyerUserId     string          `xml:"buyer_user_id,omitempty" bson:"buyer_user_id,omitempty"`   //买家支付宝用户号 以 2088 开头的纯 16 位数字
	OutTradeNo      string          `xml:"out_trade_no,omitempty" bson:"out_trade_no,omitempty"`     //商户网站唯一订单号
	ResultCode      string          `xml:"result_code,omitempty" bson:"result_code,omitempty"`       //查询处理结果响应码。SUCCESS:查询成功 FAIL:查询失败 PROCESS_EXCEPTION:处理异常
	TradeNo         string          `xml:"trade_no,omitempty" bson:"trade_no,omitempty"`             //支付宝交易号 最短16位,最长64位
	QrCode          string          `xml:"qr_code,omitempty" bson:"qr_code,omitempty"`               //二维码
	PicUrl          string          `xml:"pic_url,omitempty" bson:"pic_url,omitempty"`               //二维码图片地址
	SmallPicUrl     string          `xml:"small_pic_url,omitempty" bson:"small_pic_url,omitempty"`
	VoucherType     string          `xml:"voucher_type,omitempty" bson:"voucher_type,omitempty"`
	DetailErrorCode string          `xml:"detail_error_code,omitempty" bson:"detail_error_code,omitempty"` //详细错误码
	DetailErrorDes  string          `xml:"detail_error_des,omitempty" bson:"detail_error_des,omitempty"`
	ExtendInfo      string          `xml:"extend_info,omitempty" bson:"extend_info,omitempty"`
	TradeStatus     string          `xml:"trade_status,omitempty" bson:"trade_status,omitempty"`                   //交易状态
	Partner         string          `xml:"partner,omitempty" bson:"partner,omitempty"`                             //合作者身份ID
	FundBillList    []tradeFundBill `xml:"fund_bill_list>TradeFundBill,omitempty" bson:"fund_bill_list,omitempty"` //本次交易支付单据信息集合
	TotalFee        string          `xml:"total_fee,omitempty" bson:"total_fee,omitempty"`                         //订单金额
	SendPayDate     string          `xml:"send_pay_date,omitempty" bson:"send_pay_date,omitempty"`                 //本次交易打款到卖家账户的时间,格式为 yyyy-MM-dd HH:mm:ss
}

// TradeFundBill 支付单据信息
type tradeFundBill struct {
	XMLName     xml.Name `xml:"TradeFundBill" bson:"-"`
	Amount      string   `xml:"amount,omitempty"`       //支付金额
	FundChannel string   `xml:"fund_channel,omitempty"` //支付渠道
}

// DisCount 计算商户、渠道折扣
func (alp *alpDetail) DisCount() (string, string) {

	merf, chcdf := 0.00, 0.00
	for _, v := range alp.FundBillList {
		f, _ := strconv.ParseFloat(v.Amount, 64)
		switch v.FundChannel {
		// 渠道
		case "00", "30", "40":
			merf += f
		// 商户
		case "101", "102":
			chcdf += f
		}
	}
	return fmt.Sprintf("%0.2f", merf), fmt.Sprintf("%0.2f", chcdf)
}

//查询账单接口
type alpSettleResponse struct {
	XMLName   xml.Name `xml:"alipay" bson:"-"`
	IsSuccess string   `xml:"is_success,omitempty"`
	Sign      string   `xml:"sign,omitempty"`
	SignType  string   `xml:"sign_type,omitempty"`
	Error     string   `xml:"error,omitempty"`
	Request   []Param  `xml:"request>param"`
	// Alipay
	Response csvResult `xml:"response,omitempty" bson:"response,omitempty"`
}

type csvResult struct {
	CsvResult csvDetail `xml:"csv_result,omitempty"`
}

type csvDetail struct {
	Count  string `xml:"count,omitempty" bson:"count,omitempty"`
	CsvStr string `xml:"csv_data,omitempty" bson:"csv_data,omitempty"`
}
