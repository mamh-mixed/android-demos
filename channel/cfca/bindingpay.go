package cfca

import "quickpay/model"

var requestURL = "https://test.china-clearing.com/Gateway/InterfaceII"

// CfcaRequest 中金渠道请求报文
type CfcaRequest struct {
	Version string `xml:"version,attr,omitempty"`
	Head    requestHead
	Body    requestBody
}

// CfcaResponse 中金渠道返回报文
type CfcaResponse struct {
	Head respHead
	Body respBody
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

// common resp head
type respHead struct {
	Code    string
	Message string
}

// common resp body
type respBody struct {
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

//Tx 1810 交易对账单
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

// ProcessBindingEnquiry 查询绑定关系
func ProcessBindingEnquiry(be *model.BindingEnquiry) (ret *model.BindingReturn) {
	// 将参数转化为CfcaRequest
	req := &CfcaRequest{
		Version: "2.0",
		Head: requestHead{
			InstitutionID: "001405", //测试ID
			TxCode:        "2502",
		},
		Body: requestBody{
			TxSNBinding: be.BindingId,
		},
	}

	// 向中金发起请求
	resp := sendRequest(req)

	// 应答码转换。。。

	ret = &model.BindingReturn{
		RespCode: resp.Head.Code,
		RespMsg:  resp.Head.Message,
	}
	return ret
}
