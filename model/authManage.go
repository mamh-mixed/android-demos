package model

// User 用户表
type User struct {
	UserName  string `json:"userName" bson:"userName"`
	NickName  string `json:"nickName" bson:"nickName"`
	Password  string `json:"password" bson:"password"`
	Mail      string `json:"mail" bson:"mail"`
	PhoneNum  string `json:"phoneNum" bson:"phoneNum"`
	UserType  string `json:"userType" bson:"userType"`
	AgentCode string `json:"agentCode" bson:"agentCode"`
	AgentName string `json:"agentName" bson:"agentName"`
	GroupCode string `json:"groupCode" bson:"groupCode"`
	GroupName string `json:"groupName" bson:"groupName"`
	MerId     string `json:"merId" bson:"merId"`
	MerName   string `json:"merName" bson:"merName"`
}

// Session Session表
type Session struct {
	SessionID string `json:"sessionId" bson:"sessionId"`
	User      *User  `json:"user" bson:"user"`
	Expires   string `json:"expires" bson:"expires"`
}
