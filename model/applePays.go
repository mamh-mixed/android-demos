package model

// Appyle pay
type ApplePay struct {
	MerId         string       `json:"merId"`                //商户ID
	TransType     string       `json:"transType"`            //子交易类型
	SubMerId      string       `json:"subMerId,omitempty"`   //子商户号
	TerminalId    string       `json:"terminalId,omitempty"` //终端号
	MerOrderNum   string       `json:"merOrderNum"`          //商户订单号
	TransactionId string       `json:"transactionId"`        //ApplePay标识
	ApplePayData  ApplePayData `json:"applePayData"`         //ApplePay数据
	SignKey       string       `json:"SignKey,omitempty"`    //签名密钥
	Chcd          string       `json:"chcd,omitempty"`       //下游商户配置的渠道机构号
	Mchntid       string       `json:"mchntid,omitempty"`    //下游商户配置的渠道商户号
	CliSN         string       `json:"cliSN,omitempty"`      //商户的终端在当天对应的一个序列号
	SysSN         string       `json:"sysSN,omitempty"`      //系统序列号
}

// ApplePayData applePay数据
type ApplePayData struct {
	ApplicationPrimaryAccountNumber string      `json:"applicationPrimaryAccountNumber"` // 主账号
	ApplicationExpirationDate       string      `json:"applicationExpirationDate"`       // 有效期截止日
	CurrencyCode                    string      `json:"currencyCode"`                    // 货币代码
	TransactionAmount               int64       `json:"transactionAmount"`               // 交易金额
	DeviceManufacturerIdentifier    string      `json:"deviceManufacturerIdentifier"`    // 设备制造商标识符
	PaymentDataType                 string      `json:"paymentDataType"`                 // 支付数据类型(EMV或3D Secure)
	PaymentData                     PaymentData `json:"paymentData"`                     // 支付数据内容
}

// PaymentData 支付数据
type PaymentData struct {
	OnlinePaymentCryptogram string `json:"onlinePaymentCryptogram,omitempty"` // 3D Secure类型的在线支付密码
	EciIndicator            string `json:"eciIndicator,omitempty"`            // 线上 3D Secure 交易发卡行验证结果
	EmvData                 string `json:"emvData,omitempty"`                 // EMV类型的支付数据，到线下网关的时候存到iccdata里面
}
