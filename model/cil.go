package model

type CilMsg struct {
	Busicd           string `json:"busicd" bson:"busicd,omitempty"`                               //  6  业务代码
	Txndir           string `json:"txndir" bson:"txndir,omitempty"`                               //  6  交易方向
	Routchnl         string `json:"routchnl,omitempty" bson:"routchnl,omitempty"`                 //  8  支付路由渠道
	Posentrymode     string `json:"posentrymode,omitempty" bson:"posentrymode,omitempty"`         // 12  POS输入码
	Respcd           string `json:"respcd,omitempty" bson:"respcd,omitempty"`                     //  6  应答码
	Inscd            string `json:"inscd,omitempty" bson:"inscd,omitempty"`                       //  5  机构代码
	Chcd             string `json:"chcd,omitempty" bson:"chcd,omitempty"`                         //  4  平台分配机构号
	Clisn            string `json:"clisn,omitempty" bson:"clisn,omitempty"`                       //  5  终端流水号
	Clientstan       string `json:"clientstan,omitempty" bson:"clientstan,omitempty"`             // 10  客户端流水号
	Stldt            string `json:"stldt,omitempty" bson:"stldt,omitempty"`                       //  5  清算日期
	Mchntid          string `json:"mchntid,omitempty" bson:"mchntid,omitempty"`                   //  7  商户号
	Terminalid       string `json:"terminalid,omitempty" bson:"terminalid,omitempty"`             // 10  终端号
	Trackdata2       string `json:"trackdata2,omitempty" bson:"trackdata2,omitempty"`             // 10  二磁信息
	Trackdata3       string `json:"trackdata3,omitempty" bson:"trackdata3,omitempty"`             // 10  三磁信息
	Spctrackdata     string `json:"spctrackdata,omitempty" bson:"spctrackdata,omitempty"`         // 12  FY磁道信息
	Cardpin          string `json:"cardpin,omitempty" bson:"cardpin,omitempty"`                   //  7  PIN数据
	Spccardpin       string `json:"spccardpin,omitempty" bson:"spccardpin,omitempty"`             // 10  FY PIN数据
	Txamt            string `json:"txamt,omitempty" bson:"txamt,omitempty"`                       //  5  交易金额
	Txcurrcd         string `json:"txcurrcd,omitempty" bson:"txcurrcd,omitempty"`                 //  8  交易币种
	Txdt             string `json:"txdt,omitempty" bson:"txdt,omitempty"`                         //  4  交易时间
	Chname           string `json:"chname,omitempty" bson:"chname,omitempty"`                     //  6  持卡人姓名
	Chbank           string `json:"chbank,omitempty" bson:"chbank,omitempty"`                     //  6  开户银行
	Cardcd           string `json:"cardcd,omitempty" bson:"cardcd,omitempty"`                     //  6  卡号
	Outgoingacct     string `json:"outgoingacct,omitempty" bson:"outgoingacct,omitempty"`         // 12  借记卡卡号
	Incomingacct     string `json:"incomingacct,omitempty" bson:"incomingacct,omitempty"`         // 12  信用卡卡号
	Syssn            string `json:"syssn,omitempty" bson:"syssn,omitempty"`                       //  5  检索参考号
	Setamt           string `json:"setamt,omitempty" bson:"setamt,omitempty"`                     //  6  清算金额
	Setcurr          string `json:"setcurr,omitempty" bson:"setcurr,omitempty"`                   //  7  清算币种
	Billyymm         string `json:"billyymm,omitempty" bson:"billyymm,omitempty"`                 //  8  订单年月
	Custmrpin        string `json:"custmrpin,omitempty" bson:"custmrpin,omitempty"`               //  9  客户密码
	Custmrpinfrmt    string `json:"custmrpinfrmt,omitempty" bson:"custmrpinfrmt,omitempty"`       // 13  客户密码加密格式
	Custmrtp         string `json:"custmrtp,omitempty" bson:"custmrtp,omitempty"`                 //  8  客户类型
	Custmracnt       string `json:"custmracnt,omitempty" bson:"custmracnt,omitempty"`             // 10  客户帐号
	Goodscd          string `json:"goodscd,omitempty" bson:"goodscd,omitempty"`                   //  7  商品编码
	Paymd            string `json:"paymd,omitempty" bson:"paymd,omitempty"`                       //  5  支付方式
	Origclisn        string `json:"origclisn,omitempty" bson:"origclisn,omitempty"`               //  9  原始交易流水号
	Origbusicd       string `json:"origbusicd,omitempty" bson:"origbusicd,omitempty"`             // 10  原始业务代码
	Origdt           string `json:"origdt,omitempty" bson:"origdt,omitempty"`                     //  6  原始交易时间
	Regioncd         string `json:"regioncd,omitempty" bson:"regioncd,omitempty"`                 //  8  交易地区代码
	Authid           string `json:"authid,omitempty" bson:"authid,omitempty"`                     //  6  授权号
	Localdt          string `json:"localdt,omitempty" bson:"localdt,omitempty"`                   //  7  本地时间
	Fxrate           string `json:"fxrate,omitempty" bson:"fxrate,omitempty"`                     //  6  汇率
	Nminfo           string `json:"nminfo,omitempty" bson:"nminfo,omitempty"`                     //  6  网络管理信息
	Expiredate       string `json:"expiredate,omitempty" bson:"expiredate,omitempty"`             // 10  卡片有效期
	Mcc              string `json:"mcc,omitempty" bson:"mcc,omitempty"`                           //  3  Mcc
	Mchntnm          string `json:"mchntnm,omitempty" bson:"mchntnm,omitempty"`                   //  7  商户名称
	Issuerbank       string `json:"issuerbank,omitempty" bson:"issuerbank,omitempty"`             // 10  发卡银行代码
	Mac              string `json:"mac,omitempty" bson:"mac,omitempty"`                           //  3  MAC
	Balance          string `json:"balance,omitempty" bson:"balance,omitempty"`                   //  7  实际余额
	Newkey           string `json:"newkey,omitempty" bson:"newkey,omitempty"`                     //  6  新密钥
	Receiveinsid     string `json:"receiveinsid,omitempty" bson:"receiveinsid,omitempty"`         // 12  接收机构代码
	Retrievalnum     string `json:"retrievalnum,omitempty" bson:"retrievalnum,omitempty"`         // 12  检索参考号(上行渠道返回)
	Phonenum         string `json:"phonenum,omitempty" bson:"phonenum,omitempty"`                 //  8
	Cvv2             string `json:"cvv2,omitempty" bson:"cvv2,omitempty"`                         //  4
	Psamcd           string `json:"psamcd,omitempty" bson:"psamcd,omitempty"`                     //  6
	Termserialcd     string `json:"termserialcd,omitempty" bson:"termserialcd,omitempty"`         // 12
	Paymethod        string `json:"paymethod,omitempty" bson:"paymethod,omitempty"`               //  9
	Billinscd        string `json:"billinscd,omitempty" bson:"billinscd,omitempty"`               //  9
	Barcd            string `json:"barcd,omitempty" bson:"barcd,omitempty"`                       //  5
	Billamt          string `json:"billamt,omitempty" bson:"billamt,omitempty"`                   //  7
	Billstat         string `json:"billstat,omitempty" bson:"billstat,omitempty"`                 //  8
	Billingamt       string `json:"billingamt,omitempty" bson:"billingamt,omitempty"`             // 10
	Billingrate      string `json:"billingrate,omitempty" bson:"billingrate,omitempty"`           // 11
	Billingcurr      string `json:"billingcurr,omitempty" bson:"billingcurr,omitempty"`           // 11
	Retclisn         string `json:"retclisn,omitempty" bson:"retclisn,omitempty"`                 //  8
	Transfee         string `json:"transfee,omitempty" bson:"transfee,omitempty"`                 //  8
	Convdate         string `json:"convdate,omitempty" bson:"convdate,omitempty"`                 //  8
	Inchname         string `json:"inchname,omitempty" bson:"inchname,omitempty"`                 //  8
	Cardseqnum       string `json:"cardseqnum,omitempty" bson:"cardseqnum,omitempty"`             // 10
	Iccdata          string `json:"iccdata,omitempty" bson:"iccdata,omitempty"`                   //  7
	CardholderInfo   string `json:"cardholderinfo,omitempty" bson:"cardholderinfo,omitempty"`     // 14
	Dynamicauthcode  string `json:"dynamicauthcode,omitempty" bson:"dynamicauthcode,omitempty"`   // 15
	Txnmode          string `json:"txnmode,omitempty" bson:"txnmode,omitempty"`                   //  7
	Termreadability  string `json:"termreadability,omitempty" bson:"termreadability,omitempty"`   // 15
	Icccondcode      string `json:"icccondcode,omitempty" bson:"icccondcode,omitempty"`           // 11
	Usagetags        string `json:"usagetags,omitempty" bson:"usagetags,omitempty"`               //  9
	EciIndicator     string `json:"eciindicator,omitempty" bson:"eciindicator,omitempty"`         //  2  线上3D交易发卡行验证结果，applePay必填
	Transactionid    string `json:"transactionid,omitempty" bson:"transactionid,omitempty"`       //  M20  交易订单号，applePay必填
	Onlinesecuredata string `json:"onlinesecuredata,omitempty" bson:"onlinesecuredata,omitempty"` //  M50  3DSecure数据，applePay必填
}
