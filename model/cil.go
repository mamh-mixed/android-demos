package model

/*
机构：00000050

商户 050310058120002
终端：00000001

IP：192.168.1.102
端口：端口7823（长连接端口）
*/
type CilMsg struct {
	UUID             string `json:"uuid,omitempty" bson:"uuid,omitempty"`               // 存储到数据库中的主键
	Busicd           string `json:"busicd" bson:"busicd,omitempty"`                     //  6  业务代码
	Txndir           string `json:"txndir" bson:"txndir,omitempty"`                     //  6  交易方向
	Routchnl         string `json:"routchnl" bson:"routchnl,omitempty"`                 //  8  支付路由渠道
	Posentrymode     string `json:"posentrymode" bson:"posentrymode,omitempty"`         // 12  POS输入码
	Respcd           string `json:"respcd" bson:"respcd,omitempty"`                     //  6  应答码
	Inscd            string `json:"inscd" bson:"inscd,omitempty"`                       //  5  机构代码
	Chcd             string `json:"chcd" bson:"chcd,omitempty"`                         //  4  平台分配机构号
	Clisn            string `json:"clisn" bson:"clisn,omitempty"`                       //  5  终端流水号
	Clientstan       string `json:"clientstan" bson:"clientstan,omitempty"`             // 10  客户端流水号
	Stldt            string `json:"stldt" bson:"stldt,omitempty"`                       //  5  清算日期
	Mchntid          string `json:"mchntid" bson:"mchntid,omitempty"`                   //  7  商户号
	Terminalid       string `json:"terminalid" bson:"terminalid,omitempty"`             // 10  终端号
	Trackdata2       string `json:"trackdata2" bson:"trackdata2,omitempty"`             // 10  二磁信息
	Trackdata3       string `json:"trackdata3" bson:"trackdata3,omitempty"`             // 10  三磁信息
	Spctrackdata     string `json:"spctrackdata" bson:"spctrackdata,omitempty"`         // 12  FY磁道信息
	Cardpin          string `json:"cardpin" bson:"cardpin,omitempty"`                   //  7  PIN数据
	Spccardpin       string `json:"spccardpin" bson:"spccardpin,omitempty"`             // 10  FY PIN数据
	Txamt            string `json:"txamt" bson:"txamt,omitempty"`                       //  5  交易金额
	Txcurrcd         string `json:"txcurrcd" bson:"txcurrcd,omitempty"`                 //  8  交易币种
	Txdt             string `json:"txdt" bson:"txdt,omitempty"`                         //  4  交易时间
	Chname           string `json:"chname" bson:"chname,omitempty"`                     //  6  持卡人姓名
	Chbank           string `json:"chbank" bson:"chbank,omitempty"`                     //  6  开户银行
	Cardcd           string `json:"cardcd" bson:"cardcd,omitempty"`                     //  6  卡号
	Outgoingacct     string `json:"outgoingacct" bson:"outgoingacct,omitempty"`         // 12  借记卡卡号
	Incomingacct     string `json:"incomingacct" bson:"incomingacct,omitempty"`         // 12  信用卡卡号
	Syssn            string `json:"syssn" bson:"syssn,omitempty"`                       //  5  检索参考号
	Setamt           string `json:"setamt" bson:"setamt,omitempty"`                     //  6  清算金额
	Setcurr          string `json:"setcurr" bson:"setcurr,omitempty"`                   //  7  清算币种
	Billyymm         string `json:"billyymm" bson:"billyymm,omitempty"`                 //  8  订单年月
	Custmrpin        string `json:"custmrpin" bson:"custmrpin,omitempty"`               //  9  客户密码
	Custmrpinfrmt    string `json:"custmrpinfrmt" bson:"custmrpinfrmt,omitempty"`       // 13  客户密码加密格式
	Custmrtp         string `json:"custmrtp" bson:"custmrtp,omitempty"`                 //  8  客户类型
	Custmracnt       string `json:"custmracnt" bson:"custmracnt,omitempty"`             // 10  客户帐号
	Goodscd          string `json:"goodscd" bson:"goodscd,omitempty"`                   //  7  商品编码
	Paymd            string `json:"paymd" bson:"paymd,omitempty"`                       //  5  支付方式
	Origclisn        string `json:"origclisn" bson:"origclisn,omitempty"`               //  9  原始交易流水号
	Origbusicd       string `json:"origbusicd" bson:"origbusicd,omitempty"`             // 10  原始业务代码
	Origdt           string `json:"origdt" bson:"origdt,omitempty"`                     //  6  原始交易时间
	Regioncd         string `json:"regioncd" bson:"regioncd,omitempty"`                 //  8  交易地区代码
	Authid           string `json:"authid" bson:"authid,omitempty"`                     //  6  授权号
	Localdt          string `json:"localdt" bson:"localdt,omitempty"`                   //  7  本地时间
	Fxrate           string `json:"fxrate" bson:"fxrate,omitempty"`                     //  6  汇率
	Nminfo           string `json:"nminfo" bson:"nminfo,omitempty"`                     //  6  网络管理信息
	Expiredate       string `json:"expiredate" bson:"expiredate,omitempty"`             // 10  卡片有效期
	Mcc              string `json:"mcc" bson:"mcc,omitempty"`                           //  3  Mcc
	Mchntnm          string `json:"mchntnm" bson:"mchntnm,omitempty"`                   //  7  商户名称
	Issuerbank       string `json:"issuerbank" bson:"issuerbank,omitempty"`             // 10  发卡银行代码
	Mac              string `json:"mac" bson:"mac,omitempty"`                           //  3  MAC
	Balance          string `json:"balance" bson:"balance,omitempty"`                   //  7  实际余额
	Newkey           string `json:"newkey" bson:"newkey,omitempty"`                     //  6  新密钥
	Receiveinsid     string `json:"receiveinsid" bson:"receiveinsid,omitempty"`         // 12  接收机构代码
	Retrievalnum     string `json:"retrievalnum" bson:"retrievalnum,omitempty"`         // 12  检索参考号(上行渠道返回)
	Phonenum         string `json:"phonenum" bson:"phonenum,omitempty"`                 //  8
	Cvv2             string `json:"cvv2" bson:"cvv2,omitempty"`                         //  4
	Psamcd           string `json:"psamcd" bson:"psamcd,omitempty"`                     //  6
	Termserialcd     string `json:"termserialcd" bson:"termserialcd,omitempty"`         // 12
	Paymethod        string `json:"paymethod" bson:"paymethod,omitempty"`               //  9
	Billinscd        string `json:"billinscd" bson:"billinscd,omitempty"`               //  9
	Barcd            string `json:"barcd" bson:"barcd,omitempty"`                       //  5
	Billamt          string `json:"billamt" bson:"billamt,omitempty"`                   //  7
	Billstat         string `json:"billstat" bson:"billstat,omitempty"`                 //  8
	Billingamt       string `json:"billingamt" bson:"billingamt,omitempty"`             // 10
	Billingrate      string `json:"billingrate" bson:"billingrate,omitempty"`           // 11
	Billingcurr      string `json:"billingcurr" bson:"billingcurr,omitempty"`           // 11
	Retclisn         string `json:"retclisn" bson:"retclisn,omitempty"`                 //  8
	Transfee         string `json:"transfee" bson:"transfee,omitempty"`                 //  8
	Convdate         string `json:"convdate" bson:"convdate,omitempty"`                 //  8
	Inchname         string `json:"inchname" bson:"inchname,omitempty"`                 //  8
	Cardseqnum       string `json:"cardseqnum" bson:"cardseqnum,omitempty"`             // 10
	Iccdata          string `json:"iccdata" bson:"iccdata,omitempty"`                   //  7
	CardholderInfo   string `json:"cardholderinfo" bson:"cardholderinfo,omitempty"`     // 14
	Dynamicauthcode  string `json:"dynamicauthcode" bson:"dynamicauthcode,omitempty"`   // 15
	Txnmode          string `json:"txnmode" bson:"txnmode,omitempty"`                   //  7
	Termreadability  string `json:"termreadability" bson:"termreadability,omitempty"`   // 15
	Icccondcode      string `json:"icccondcode" bson:"icccondcode,omitempty"`           // 11
	Usagetags        string `json:"usagetags" bson:"usagetags,omitempty"`               //  9
	EciIndicator     string `json:"eciindicator" bson:"eciindicator,omitempty"`         //  2  线上3D交易发卡行验证结果，applePay必填
	Transactionid    string `json:"transactionid" bson:"transactionid,omitempty"`       //  M20  交易订单号，applePay必填
	Onlinesecuredata string `json:"onlinesecuredata" bson:"onlinesecuredata,omitempty"` //  M50  3DSecure数据，applePay必填
}
