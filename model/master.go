package model

type ResultBody struct {
	Status  int         `json:"status"`  // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// NewBindingReturn ResultBody构造函数
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
