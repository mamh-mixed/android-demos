package cil

/*
机构：00000050

商户 050310058120002
终端：00000001

IP：192.168.1.102
端口：端口7823（长连接端口）
*/
type CilMsg struct {
	Busicd          string //  6  业务代码
	Txndir          string //  6  交易方向
	Routchnl        string //  8  支付路由渠道
	Posentrymode    string // 12  POS输入码
	Respcd          string //  6  应答码
	Inscd           string //  5  机构代码
	Chcd            string //  4  平台分配机构号
	Clisn           string //  5  终端流水号
	Clientstan      string // 10  客户端流水号
	Stldt           string //  5  清算日期
	Mchntid         string //  7  商户号
	Terminalid      string // 10  终端号
	Trackdata2      string // 10  二磁信息
	Trackdata3      string // 10  三磁信息
	Spctrackdata    string // 12  FY磁道信息
	Cardpin         string //  7  PIN数据
	Spccardpin      string // 10  FY PIN数据
	Txamt           string //  5  交易金额
	Txcurrcd        string //  8  交易币种
	Txdt            string //  4  交易时间
	Chname          string //  6  持卡人姓名
	Chbank          string //  6  开户银行
	Cardcd          string //  6  卡号
	Outgoingacct    string // 12  借记卡卡号
	Incomingacct    string // 12  信用卡卡号
	Syssn           string //  5  检索参考号
	Setamt          string //  6  清算金额
	Setcurr         string //  7  清算币种
	Billyymm        string //  8  订单年月
	Custmrpin       string //  9  客户密码
	Custmrpinfrmt   string // 13  客户密码加密格式
	Custmrtp        string //  8  客户类型
	Custmracnt      string // 10  客户帐号
	Goodscd         string //  7  商品编码
	Paymd           string //  5  支付方式
	Origclisn       string //  9  原始交易流水号
	Origbusicd      string // 10  原始业务代码
	Origdt          string //  6  原始交易时间
	Regioncd        string //  8  交易地区代码
	Authid          string //  6  授权号
	Localdt         string //  7  本地时间
	Fxrate          string //  6  汇率
	Nminfo          string //  6  网络管理信息
	Expiredate      string // 10  卡片有效期
	Mcc             string //  3  Mcc
	Mchntnm         string //  7  商户名称
	Issuerbank      string // 10  发卡银行代码
	Mac             string //  3  MAC
	Balance         string //  7  实际余额
	Newkey          string //  6  新密钥
	Receiveinsid    string // 12  接收机构代码
	Retrievalnum    string // 12  检索参考号(上行渠道返回)
	Phonenum        string // 8
	Cvv2            string // 4
	Psamcd          string // 6
	Termserialcd    string // 12
	Paymethod       string // 9
	Billinscd       string // 9
	Barcd           string // 5
	Billamt         string // 7
	Billstat        string // 8
	Billingamt      string // 10
	Billingrate     string // 11
	Billingcurr     string // 11
	Retclisn        string // 8
	Transfee        string // 8
	Convdate        string // 8
	Inchname        string // 8
	Cardseqnum      string // 10
	Iccdata         string // 7
	CardholderInfo  string // 14
	Dynamicauthcode string // 15
	Txnmode         string // 7
	Termreadability string // 15
	Icccondcode     string // 11
	Usagetags       string // 9
}

func NewConsumeCilMsg() (m *CilMsg) {
	m = &CilMsg{
		Busicd:       "000000",
		Txndir:       "Q",
		Posentrymode: "",
		Chcd:         "00000050",
		Clisn:        "",
		Mchntid:      "",
		Terminalid:   "",
		Txamt:        "",
		Txcurrcd:     "",
		Cardcd:       "",
		Syssn:        "",
		Localdt:      "",
	}
	return m
}
