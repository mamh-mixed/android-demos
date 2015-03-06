package channel

import (
	"encoding/xml"
	"github.com/omigo/g"
	"io/ioutil"
	"net/http"
	"net/url"
	"quickpay/domain"
	"quickpay/model"
	"quickpay/tools"
)

var requestUrl = "https://test.china-clearing.com/Gateway/InterfaceII"

// 中金渠道
type Chinapay struct {
}

// 绑定关系请求
func (c *Chinapay) CreateBinding(data domain.BindingCreateRequest) *model.ChannelRes {

	// 将参数转化为Request
	request := Request{}

	//包装
	reqHead := requestHead{}
	reqHead.InstitutionID = "001405" //测试ID
	reqHead.TxCode = "2501"
	request.Version = "2.0"
	reqBody := requestBody{}
	reqBody.AccountName = data.AcctName
	reqBody.AccountNumber = data.AcctNum

	// 提交
	response := ChinaPayRequestHandler(request)

	//对返回值处理...
	// handledResponse(response)
	//TODO
	channelRes := model.ChannelRes{"000000", response}
	return &channelRes
}

// 查询绑定关系
func (c *Chinapay) QueryBinding() *model.ChannelRes {
	// 将参数转化为Request
	request := Request{}

	// 提交
	response := ChinaPayRequestHandler(request)

	//对返回值处理...
	// handledResponse(response)
	//TODO
	channelRes := model.ChannelRes{"000000", response}
	return &channelRes
}

// 快捷支付
func (c *Chinapay) QuickPay() *model.ChannelRes {
	// 将参数转化为Request
	request := Request{}

	// 提交
	response := ChinaPayRequestHandler(request)

	//对返回值处理...
	// handledResponse(response)
	//TODO
	channelRes := model.ChannelRes{"000000", response}
	return &channelRes
}

// 快捷支付查询
func (c *Chinapay) QuickPayQuery() *model.ChannelRes {
	// 将参数转化为Request
	request := Request{}

	// 提交
	response := ChinaPayRequestHandler(request)

	//对返回值处理...
	// handledResponse(response)
	//TODO
	channelRes := model.ChannelRes{"000000", response}
	return &channelRes
}

// 快捷支付退款
func (c *Chinapay) QuickPayRefund() *model.ChannelRes {
	// 将参数转化为Request
	request := Request{}

	// 提交
	response := ChinaPayRequestHandler(request)

	//对返回值处理...
	// handledResponse(response)
	//TODO
	channelRes := model.ChannelRes{"000000", response}
	return &channelRes
}

// 快捷支付退款查询
func (c *Chinapay) QuickPayRefundQuery() *model.ChannelRes {
	// 将参数转化为Request
	request := Request{}

	// 提交
	response := ChinaPayRequestHandler(request)

	//对返回值处理...
	// handledResponse(response)
	//TODO
	channelRes := model.ChannelRes{"000000", response}
	return &channelRes
}

// 交易对账单
func (c *Chinapay) TradePayments() *model.ChannelRes {
	// 将参数转化为Request
	request := Request{}

	// 提交
	response := ChinaPayRequestHandler(request)

	//对返回值处理...
	// handledResponse(response)
	//TODO
	channelRes := model.ChannelRes{"000000", response}
	return &channelRes
}

// ChinaPaySignature 中金支付渠道签名
// message  采用Base64编码
// signature 采用Sha1WithRsa签名后用Hex编码
func ChinaPaySignature(data Request) (message, signature string) {
	// to xml
	xmlBytes := tools.ToXML(data)

	g.Debug("transfer data into xml : (%s)", xmlBytes)

	return tools.EncodeBase64(xmlBytes), tools.EncodeHex(tools.SignatureUseSha1WithRsa(xmlBytes))
}

// CheckChinaPaySignature 中金支付渠道验签
func CheckChinaPaySignature(data string, signature string) bool {
	// encode base64
	message := tools.DecodeBase64(data)
	// ecode hex
	sign := tools.DecodeHex(signature)
	// verify
	err := tools.CheckSignatureUseSha1WithRsa(message, sign)

	return err == nil
}

// 对中金接口访问的统一处理
func ChinaPayRequestHandler(request Request) *Response {

	// 数据处理、加密、签名
	message, signature := ChinaPaySignature(request)

	// 准备参数、提交
	param := url.Values{}
	param.Add("message", message)
	param.Add("signature", signature)

	resp, err := http.PostForm(requestUrl, param)
	if err != nil {
		g.Error("unable to connect ChinaPay gratway  (%s)", err)
	}

	// handler result
	// read from response
	bodys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		g.Error("unable to read from response (%s)", err)
	}
	g.Debug("response data : (%s)", bodys)
	// result map
	response := Response{}
	err = xml.Unmarshal(bodys, &response)
	if err != nil {
		g.Error("unable to unmarshal xml (%s)", err)
	}

	if err != nil {
		return nil
	}

	return &response
}

// 中金渠道请求报文
type Request struct {
	Version string `xml:"version,attr,omitempty"`
	Head    requestHead
	Body    requestBody
}

// 中金渠道返回报文
type Response struct {
	Head responseHead
	Body responseBody
}

// common request head
type requestHead struct {
	TxCode        string
	InstitutionID string `xml:",omitempty"`
}

// common request body
type requestBody struct {
	TxSNBinding          string `xml:",omitempty"` //绑定流水号
	BankID               string `xml:",omitempty"` //银行 ID
	AccountName          string `xml:",omitempty"` //账户名称
	AccountNumber        string `xml:",omitempty"` //账户号码
	IdentificationType   string `xml:",omitempty"` //开户证件类型
	IdentificationNumber string `xml:",omitempty"` //证件号码
	PhoneNumber          string `xml:",omitempty"` //手机号
	CardType             string `xml:",omitempty"` //卡类型
	ValidDate            string `xml:",omitempty"` //信用卡有效期
	CVN2                 string `xml:",omitempty"` //信用卡背面的末 3 位数字
	TxSNUnBinding        string `xml:",omitempty"` //解绑流水号
	PaymentNo            string `xml:",omitempty"` //支付交易流水号
	Amount               int64  `xml:",omitempty"` //支付金额,单位:分
	SettlementFlag       string `xml:",omitempty"` //结算标识
	Remark               string `xml:",omitempty"` //备注
	TxSN                 string `xml:",omitempty"` //退款交易流水号
	InstitutionID        string `xml:",omitempty"` //机构编号
	Date                 string `xml:",omitempty"` //对账日期,格式:YYYY-MM-DD
}

// common response head
type responseHead struct {
	Code    string
	Message string
}

// common response body
type responseBody struct {
	InstitutionID   string //机构编号
	TxSNBinding     string //绑定流水号
	Status          int8   //交易状态
	ResponseCode    string //响应代码
	ResponseMessage string //响应消息
	IssInsCode      string //发卡机构代码
	PayCardType     string //支付卡类型
	BankTxTime      int64  //银行处理时间
	TxSNUnBinding   string //解绑流水号
	PaymentNo       string //支付交易流水号
	TxSN            string //退款交易流水号
	Amount          int64  //退款金额,单位:分
}

//1810 交易对账单
type Tx struct {
	TxType               string //交易类型
	TxSN                 string //退款交易流水号
	TxAmount             int64  //交易金额,单位:分
	InstitutionAmount    int64  //机构应收的金额,单位:分
	PaymentAmount        int64  //支付平台应收的金额,单位:分
	PayerFee             int64  //付款人手续费,单位:分
	InstitutionFee       int64  //机构手续费,单位:分
	Fee                  int64  //手续费,单位:分
	Remark               string //备注
	BankNotificationTime string //支付平台收到银行通知时间,格式: YYYYMMDDhhmmss
	SettlementFlag       string //结算标识
}
