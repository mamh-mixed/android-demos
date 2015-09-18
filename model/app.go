package model

const (
	KEY     = "eu1dr0c8znpa43blzy1wirzmk8jqdaon"
	SUCCESS = "success"
	FAIL    = "fail"
	// USERNAME_EXIST          = "username_exist"
	// SYSTEM_ERROR            = "system_error"
	// SIGN_FAIL               = "sign_fail"
	// USER_NO_ACTIVATE        = "user_no_activate"
	// USERNAME_PASSWORD_ERROR = "username_password_error"
	// USERNAME_NO_EXIST       = "username_no_exist"
	// OLD_PASSWORD_ERROR      = "old_password_error"
	// PARAMS_EMPTY            = "params_empty"
	JSON_ERROR = `{"state":"fail","error","system_error"}`
)

var (
	SUCCESS1                = NewAppResult("success", "")
	USERNAME_EXIST          = NewAppResult(FAIL, "username_exist")
	SYSTEM_ERROR            = NewAppResult(FAIL, "system_error")
	SIGN_FAIL               = NewAppResult(FAIL, "sign_fail")
	USER_NO_ACTIVATE        = NewAppResult(FAIL, "user_no_activate")
	USERNAME_PASSWORD_ERROR = NewAppResult(FAIL, "username_password_error")
	USERNAME_NO_EXIST       = NewAppResult(FAIL, "username_no_exist")
	OLD_PASSWORD_ERROR      = NewAppResult(FAIL, "old_password_error")
	PARAMS_EMPTY            = NewAppResult(FAIL, "params_empty")
	CODE_ERROR              = NewAppResult(FAIL, "code_error")
	NO_PAY_MER              = NewAppResult(FAIL, "找不到支付商户")
	NO_TRANS                = NewAppResult(FAIL, "找不到交易")
	TIME_ERROR              = NewAppResult(FAIL, "日期格式错误")
)

type AppResult struct {
	// 必填
	State string `json:"state"` // 状态
	Error string `json:"error"` // 错误消息

	// 可选
	User         *AppUser    `json:"user,omitempty"`
	TotalAmt     string      `json:"total,omitempty"`
	Count        int         `json:"count,omitempty"`
	Size         int         `json:"size,omitempty"`
	RefdCount    int         `json:"refdcount,omitempty"`
	RefdTotalAmt string      `json:"refdtotal,omitempty"`
	BankOpen     string      `json:"bank_open,omitempty"`
	Payee        string      `json:"payee,omitempty"`
	PayeeCard    string      `json:"payee_card,omitempty"`
	PhoneNum     string      `json:"phone_num,omitempty"`
	Txn          interface{} `json:"txn,omitempty"` // 交易，可存放数组或对象
}

// NewAppResult NewAppResult
func NewAppResult(state, err string) (ret *AppResult) {
	return &AppResult{
		State: state,
		Error: err,
	}
}

type AppTxn struct {
	Response        string `json:"response" bson:"response"`
	SystemDate      string `json:"system_date"`
	ConsumerAccount string `json:"consumerAccount"`
	ReqData         struct {
		Busicd       string `json:"busicd"`
		AgentCode    string `json:"inscd"`
		Txndir       string `json:"txndir"`
		Terminalid   string `json:"terminalid"`
		OrigOrderNum string `json:"origOrderNum"`
		OrderNum     string `json:"orderNum"`
		MerId        string `json:"mchntid"`
		TradeFrom    string `json:"tradeFrom"`
		Txamt        string `json:"txamt"`
		ChanCode     string `json:"chcd"`
		Currency     string `json:"currency"`
	} `json:"m_request"`
}

// AppUser 云收银用户
type AppUser struct {
	UserName  string `json:"userName,omitempty" bson:"userName,omitempty"`
	Password  string `json:"-" bson:"password,omitempty"`
	Activate  string `json:"activate,omitempty" bson:"activate,omitempty"`
	MerId     string `json:"clientid,omitempty" bson:"merId,omitempty"`
	Limit     string `json:"limit,omitempty" bson:"limit,omitempty"`
	SignKey   string `json:"signKey,omitempty" bson:"signKey,omitempty"`
	AgentCode string `json:"inscd,omitempty" bson:"inscd,omitempty"`
	UniqueId  string `json:"objectId,omitempty" bson:"-"` // 不存
	BankOpen  string `json:"bank_open,omitempty" bson:"bank_open,omitempty"`
	Payee     string `json:"payee,omitempty" bson:"payee,omitempty"`
	PayeeCard string `json:"payee_card,omitempty" bson:"payee_card,omitempty`
	PhoneNum  string `json:"phone_num,omitempty" bson:"phone_num,omitempty"`
}

// Email 发送email记录
type Email struct {
	UserName  string `json:"username,omitempty" bson:"userName,omitempty"`
	Code      string `json:"code,omitempty" bson:"code,omitempty"`
	Success   bool   `json:"success,omitempty" bson:"success,omitempty"`
	Timestamp string `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}
