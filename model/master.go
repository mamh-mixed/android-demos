package model

type ResultBody struct {
	Status  int         `json:"status"`  // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// NewResultBody ResultBody 构造函数
func NewResultBody(status int, msg string) (ret *ResultBody) {
	return &ResultBody{
		Status:  status,
		Message: msg,
	}
}

// Pagination 分页对象
type Pagination struct {
	Page  int         `json:"page"`  // 当前页
	Total int         `json:"total"` // 总记录数
	Size  int         `json:"size"`  // 每页记录数
	Count int         `json:"count"` // 当前页记录数
	Data  interface{} `json:"data"`  // 数据
}

// MasterLog 平台操作日志
type MasterLog struct {
	Time       string `bson:"time" json:"time"`
	UserName   string `bson:"userName" json:"userName"`
	IP         string `bson:"ip" json:"ip"`
	Path       string `bson:"path" json:"path"`
	Method     string `bson:"method" json:"method"`
	Query      string `bson:"query" json:"query"`
	Body       string `bson:"body" json:"body"`
	CreateTime string `bson:"createTime,omitempty" json:"createTime,omitempty"` // 创建时间
}
