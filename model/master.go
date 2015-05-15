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
