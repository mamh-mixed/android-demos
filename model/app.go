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

	// register from
	SelfRegister       = 1
	PreRegister        = 2
	SalesToolsRegister = 3
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
	USER_ALREADY_IMPROVED   = NewAppResult(FAIL, "user_already_improved")
	MERID_NO_EXIST          = NewAppResult(FAIL, "merId_no_exist")
	USER_LOCK               = NewAppResult(FAIL, "user_lock")
	USER_THREE_TIMES        = NewAppResult(FAIL, "user_has_three_times")
	USER_TWO_TIMES          = NewAppResult(FAIL, "user_has_two_times")
	USER_ONE_TIMES          = NewAppResult(FAIL, "user_has_one_times")
	CODE_ERROR_CH           = NewAppResult(FAIL, "code码不存在")
	NO_PAY_MER              = NewAppResult(FAIL, "找不到支付商户")
	NO_TRANS                = NewAppResult(FAIL, "找不到交易")
	TIME_ERROR              = NewAppResult(FAIL, "日期格式错误")
	CODE_TIME_ERROR_CH      = NewAppResult(FAIL, "code码已过期")
	USERNAME_NO_EXIST_CH    = NewAppResult(FAIL, "用户名不存在")
	PARAMS_EMPTY_CH         = NewAppResult(FAIL, "参数为空")
	SYSTEM_ERROR_CH         = NewAppResult(FAIL, "系统错误")
	TOKEN_ERROR             = NewAppResult(FAIL, "accessToken_error")
	USER_DATA_ERROR         = NewAppResult(FAIL, "user_data_error")
	NOT_SUPPORT             = NewAppResult(FAIL, "请联系您的服务商为您修改清算信息。")
)

// AppUserContiditon app用户查找条件
type AppUserContiditon struct {
	SubAgentCode string
	RegisterFrom int
	Username     string
	StartTime    string
	EndTime      string
}

type AppResult struct {
	// 必填
	State string `json:"state"`           // 状态
	Error string `json:"error,omitempty"` // 错误消息

	// 可选
	User         *AppUser    `json:"user,omitempty"`
	Users        []*AppUser  `json:"users,omitempty"`
	TotalAmt     string      `json:"total,omitempty"`
	Count        int         `json:"count"`
	TotalRecord  int         `json:"totalRecord"`
	Size         int         `json:"size"`
	RefdCount    int         `json:"refdcount"`
	RefdTotalAmt string      `json:"refdtotal,omitempty"`
	SettInfo     *SettInfo   `json:"info,omitempty"`
	Txn          interface{} `json:"txn,omitempty"` // 交易，可存放数组或对象
	UploadToken  string      `json:"uploadToken,omitempty"`
	AccessToken  string      `json:"accessToken,omitempty"`
	DownloadUrl  string      `json:"downloadUrl,omitempty"`
}

type SettInfo struct {
	BankOpen   string `json:"bank_open,omitempty"`
	Payee      string `json:"payee,omitempty"`
	PayeeCard  string `json:"payee_card,omitempty"`
	PhoneNum   string `json:"phone_num,omitempty"`
	Province   string `json:"province,omitempty"`
	City       string `json:"city,omitempty"`
	BranchBank string `json:"branch_bank,omitempty"`
	BankNo     string `json:"bankNo,omitempty"`
}

// NewAppResult NewAppResult
func NewAppResult(state, err string) (ret AppResult) {
	return AppResult{
		State: state,
		Error: err,
	}
}

type AppTxn struct {
	Response        string `json:"response" bson:"response"`
	SystemDate      string `json:"system_date"`
	ConsumerAccount string `json:"consumerAccount,omitempty"`
	TransStatus     string `json:"transStatus,omitempty"`
	TicketNum       string `json:"receiptnum,omitempty"`
	RefundAmt       int64  `json:"refundAmt"`
	ReqData         struct {
		Busicd       string `json:"busicd,omitempty"`
		AgentCode    string `json:"inscd,omitempty"`
		Txndir       string `json:"txndir,omitempty"`
		Terminalid   string `json:"terminalid,omitempty"`
		OrigOrderNum string `json:"origOrderNum,omitempty"`
		OrderNum     string `json:"orderNum,omitempty"`
		MerId        string `json:"mchntid,omitempty"`
		TradeFrom    string `json:"tradeFrom,omitempty"`
		Txamt        string `json:"txamt,omitempty"`
		ChanCode     string `json:"chcd,omitempty"`
		Currency     string `json:"currency,omitempty"`
	} `json:"m_request"`
}

// AppUser 云收银用户
type AppUser struct {
	UserName       string `json:"username,omitempty" bson:"username,omitempty"`
	Password       string `json:"password,omitempty" bson:"password,omitempty"`
	Activate       string `json:"activate,omitempty" bson:"activate,omitempty"`
	MerId          string `json:"clientid,omitempty" bson:"merId,omitempty"`
	Limit          string `json:"limit,omitempty" bson:"limit,omitempty"`
	CreateTime     string `json:"createTime,omitempty" bson:"createTime,omitempty"`
	RegisterFrom   int    `json:"-" bson:"registerFrom"` // 0-自行注册 1-预先分配 2-公司地推
	Remark         string `json:"-" bson:"remark,omitempty"`
	UpdateTime     string `json:"-" bson:"updateTime,omitempty"`
	SubAgentCode   string `json:"-" bson:"subAgentCode,omitempty"` // 隶属某个公司发展的
	MerName        string `json:"merName,omitempty" bson:"merName,omitempty"`
	BelongsTo      string `json:"-" bson:"belongsTo,omitempty"` // 属于哪个公司人员发展的
	InvitationCode string `json:"-" bson:"invitationCode,omitempty"`
	LoginTime      string `json:"loginTime,omitempty" bson:"loginTime,omitempty"`       //记录第一次登陆时间
	LockTime       string `json:"lockTime,omitempty" bson:"lockTime,omitempty"`         //记录锁定时间
	Device_type    string `json:"device_type,omitempty" bson:"device_type,omitempty"`   //设备类型
	Device_token   string `json:"device_token,omitempty" bson:"device_token,omitempty"` //app唯一标识

	// 清算相关信息不存
	BankOpen  string `json:"bank_open,omitempty" bson:"-"`
	Payee     string `json:"payee,omitempty" bson:"-"`
	PayeeCard string `json:"payee_card,omitempty" bson:"-"`
	PhoneNum  string `json:"phone_num,omitempty" bson:"-"`
	// 商户里的不存
	SignKey   string   `json:"signKey,omitempty" bson:"-"`
	AgentCode string   `json:"inscd,omitempty" bson:"-"`
	UniqueId  string   `json:"objectId,omitempty" bson:"-"`
	PayUrl    string   `json:"payUrl,omitempty" bson:"-"`
	Images    []string `json:"images,omitempty" bson:"-"`
}

// Email 发送email记录
type Email struct {
	UserName  string `json:"username,omitempty" bson:"username,omitempty"`
	Code      string `json:"code,omitempty" bson:"code,omitempty"`
	Success   bool   `json:"success,omitempty" bson:"success,omitempty"`
	Timestamp string `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

//推送请求
type PushMessageReq struct {
	MerID        string
	UserName     string //app的终端号为用户名
	Title        string
	Message      string
	Device_token string
	MsgType      string
	To           string
}

//推送应答
type PushMessageRsp struct {
	UserName string `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Title    string `json:"title,omitempty" bson:"title,omitempty"`
	Message  string `json:"message,omitempty" bson:"message,omitempty"`
	PushTime string `json:"pushtime,omitempty" bson:"pushtime,omitempty"`
	Index    string `json:"index,omitempty" bson:"index,omitempty"`
	Size     string `json:"size,omitempty" bson:"size,omitempty"`
	LastTime string `json:"lastTime,omitempty" bson:"lastTime,omitempty"`
}

//推送表结构
type PushInfo struct {
	UserName string `json:"username,omitempty" bson:"username,omitempty"`
	Title    string `json:"title,omitempty" bson:"title,omitempty"`
	Message  string `json:"message,omitempty" bson:"message,omitempty"`
	PushTime string `json:"pushtime,omitempty" bson:"pushtime,omitempty"`
}
