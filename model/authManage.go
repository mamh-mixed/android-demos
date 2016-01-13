package model

import "time"

// 密码加密用到的随机字符串
const RAND_PWD = "cnIjnlJbXsN2WAdpjV6AZJKKSorRt23"

// 云收银平台默认密码
const DEFAULT_PWD = "Yun#1016"

// 用户类型
const (
	UserTypeCIL      = "admin"
	UserTypeGenAdmin = "genAdmin"
	UserTypeAgent    = "agent"
	UserTypeCompany  = "subAgent"
	UserTypeMerchant = "group"
	UserTypeShop     = "merchant"
)

// User 用户表
type User struct {
	UserName     string `json:"userName" bson:"userName"`
	NickName     string `json:"nickName" bson:"nickName"`
	Password     string `json:"password" bson:"password"`
	Mail         string `json:"mail" bson:"mail,omitempty"`
	PhoneNum     string `json:"phoneNum" bson:"phoneNum,omitempty"`
	UserType     string `json:"userType" bson:"userType"` // 用户类型： admin,genAdmin,agent, subAgent, group, merchant
	AgentCode    string `json:"agentCode" bson:"agentCode,omitempty"`
	AgentName    string `json:"agentName" bson:"agentName,omitempty"`
	SubAgentCode string `json:"subAgentCode" bson:"subAgentCode,omitempty"`
	SubAgentName string `json:"subAgentName" bson:"subAgentName,omitempty"`
	GroupCode    string `json:"groupCode" bson:"groupCode,omitempty"`
	GroupName    string `json:"groupName" bson:"groupName,omitempty"`
	MerId        string `json:"merId" bson:"merId,omitempty"`
	MerName      string `json:"merName" bson:"merName,omitempty"`
	RelatedEmail string `json:"relatedEmail" bson:"relatedEmail,omitempty"`       // 关联email
	AreaCode     string `json:"areaCode" bson:"areaCode,omitempty"`               // 地区代码
	CreateTime   string `bson:"createTime,omitempty" json:"createTime,omitempty"` // 创建时间
	UpdateTime   string `bson:"updateTime,omitempty" json:"updateTime,omitempty"` // 更新时间
	LoginTime    string `bson:"loginTime,omitempty" json:"loginTime,omitempty"`   // 记录登陆时间
	LockTime     string `bson:"lockTime,omitempty" json:"lockTime,omitempty"`     // 记录锁定时间
}

// Session Session表
type Session struct {
	SessionID  string    `json:"sessionId" bson:"sessionId"`
	User       *User     `json:"-" bson:"user"`
	UserType   string    `json:"userType" bson:"userType"`
	NickName   string    `json:"nickName" bson:"nickName"`
	CreateTime time.Time `json:"-" bson:"createTime"`
	UpdateTime time.Time `json:"-" bson:"updateTime"`
	Expires    time.Time `json:"-" bson:"expires"`
	Locale     string    `json:"-" bson:"locale"` // default is zh_CN
}

type UserPwd struct {
	UserName string `json:'userName'`
	Password string `json:'password'`
	NewPwd   string `json:'newPwd'`
}

type ResetUser struct {
	UserName  string `json:'username'`
	PassWord  string `json:'password'`
	CheckCode string `json:'checkCode'`
}
