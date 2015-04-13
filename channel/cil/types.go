package cil

/*
机构：00000050

商户 050310058120002
终端：00000001

IP：192.168.1.102
端口：端口7823（长连接端口）
*/
type CilMsg struct {
	Busicd          string `json:"busicd"`          //  6  业务代码
	Txndir          string `json:"txndir"`          //  6  交易方向
	Routchnl        string `json:"routchnl"`        //  8  支付路由渠道
	Posentrymode    string `json:"posentrymode"`    // 12  POS输入码
	Respcd          string `json:"respcd"`          //  6  应答码
	Inscd           string `json:"inscd"`           //  5  机构代码
	Chcd            string `json:"chcd"`            //  4  平台分配机构号
	Clisn           string `json:"clisn"`           //  5  终端流水号
	Clientstan      string `json:"clientstan"`      // 10  客户端流水号
	Stldt           string `json:"stldt"`           //  5  清算日期
	Mchntid         string `json:"mchntid"`         //  7  商户号
	Terminalid      string `json:"terminalid"`      // 10  终端号
	Trackdata2      string `json:"trackdata2"`      // 10  二磁信息
	Trackdata3      string `json:"trackdata3"`      // 10  三磁信息
	Spctrackdata    string `json:"spctrackdata"`    // 12  FY磁道信息
	Cardpin         string `json:"cardpin"`         //  7  PIN数据
	Spccardpin      string `json:"spccardpin"`      // 10  FY PIN数据
	Txamt           string `json:"txamt"`           //  5  交易金额
	Txcurrcd        string `json:"txcurrcd"`        //  8  交易币种
	Txdt            string `json:"txdt"`            //  4  交易时间
	Chname          string `json:"chname"`          //  6  持卡人姓名
	Chbank          string `json:"chbank"`          //  6  开户银行
	Cardcd          string `json:"cardcd"`          //  6  卡号
	Outgoingacct    string `json:"outgoingacct"`    // 12  借记卡卡号
	Incomingacct    string `json:"incomingacct"`    // 12  信用卡卡号
	Syssn           string `json:"syssn"`           //  5  检索参考号
	Setamt          string `json:"setamt"`          //  6  清算金额
	Setcurr         string `json:"setcurr"`         //  7  清算币种
	Billyymm        string `json:"billyymm"`        //  8  订单年月
	Custmrpin       string `json:"custmrpin"`       //  9  客户密码
	Custmrpinfrmt   string `json:"custmrpinfrmt"`   // 13  客户密码加密格式
	Custmrtp        string `json:"custmrtp"`        //  8  客户类型
	Custmracnt      string `json:"custmracnt"`      // 10  客户帐号
	Goodscd         string `json:"goodscd"`         //  7  商品编码
	Paymd           string `json:"paymd"`           //  5  支付方式
	Origclisn       string `json:"origclisn"`       //  9  原始交易流水号
	Origbusicd      string `json:"origbusicd"`      // 10  原始业务代码
	Origdt          string `json:"origdt"`          //  6  原始交易时间
	Regioncd        string `json:"regioncd"`        //  8  交易地区代码
	Authid          string `json:"authid"`          //  6  授权号
	Localdt         string `json:"localdt"`         //  7  本地时间
	Fxrate          string `json:"fxrate"`          //  6  汇率
	Nminfo          string `json:"nminfo"`          //  6  网络管理信息
	Expiredate      string `json:"expiredate"`      // 10  卡片有效期
	Mcc             string `json:"mcc"`             //  3  Mcc
	Mchntnm         string `json:"mchntnm"`         //  7  商户名称
	Issuerbank      string `json:"issuerbank"`      // 10  发卡银行代码
	Mac             string `json:"mac"`             //  3  MAC
	Balance         string `json:"balance"`         //  7  实际余额
	Newkey          string `json:"newkey"`          //  6  新密钥
	Receiveinsid    string `json:"receiveinsid"`    // 12  接收机构代码
	Retrievalnum    string `json:"retrievalnum"`    // 12  检索参考号(上行渠道返回)
	Phonenum        string `json:"phonenum"`        //  8
	Cvv2            string `json:"cvv2"`            //  4
	Psamcd          string `json:"psamcd"`          //  6
	Termserialcd    string `json:"termserialcd"`    // 12
	Paymethod       string `json:"paymethod"`       //  9
	Billinscd       string `json:"billinscd"`       //  9
	Barcd           string `json:"barcd"`           //  5
	Billamt         string `json:"billamt"`         //  7
	Billstat        string `json:"billstat"`        //  8
	Billingamt      string `json:"billingamt"`      // 10
	Billingrate     string `json:"billingrate"`     // 11
	Billingcurr     string `json:"billingcurr"`     // 11
	Retclisn        string `json:"retclisn"`        //  8
	Transfee        string `json:"transfee"`        //  8
	Convdate        string `json:"convdate"`        //  8
	Inchname        string `json:"inchname"`        //  8
	Cardseqnum      string `json:"cardseqnum"`      // 10
	Iccdata         string `json:"iccdata"`         //  7
	CardholderInfo  string `json:"cardholderinfo"`  // 14
	Dynamicauthcode string `json:"dynamicauthcode"` // 15
	Txnmode         string `json:"txnmode"`         //  7
	Termreadability string `json:"termreadability"` // 15
	Icccondcode     string `json:"icccondcode"`     // 11
	Usagetags       string `json:"usagetags"`       //  9
}
