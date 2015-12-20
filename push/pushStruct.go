package push

type UmengBody struct {
	Ticker     string `json:"ticker,omitempty" bson:"ticker,omitempty"`
	Title      string `json:"title,omitempty" bson:"title,omitempty"`
	Text       string `json:"text,omitempty" bson:"text,omitempty"`
	After_open string `json:"after_open,omitempty" bson:"after_open,omitempty"`
}

type UmengPayload struct {
	Body         UmengBody `json:"body,omitempty" bson:"body,omitempty"`
	Display_type string    `json:"display_type,omitempty" bson:"display_type,omitempty"`
}

type UmengMessage struct {
	Appkey        string       `json:"appkey,omitempty" bson:"appkey,omitempty"`
	Timestamp     int64        `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	Device_tokens string       `json:"device_tokens,omitempty" bson:"device_tokens,omitempty"`
	Type          string       `json:"type,omitempty" bson:"type,omitempty"`
	Payload       UmengPayload `json:"payload,omitempty" bson:"payload,omitempty"`
}

type UmengRsp struct {
	Ret  string       `json:"ret,omitempty" bson:"ret,omitempty"`
	Data UmengRspData `json:"data,omitempty" bson:"data,omitempty"`
}

type UmengRspData struct {
	Msg_id        string `json:"msg_id,omitempty" bson:"msg_id,omitempty"`
	Task_id       string `json:"task_id,omitempty" bson:"task_id,omitempty"`
	Error_code    string `json:"error_code,omitempty" bson:"error_code,omitempty"`
	Thirdparty_id string `json:"thirdparty_id,omitempty" bson:"thirdparty_id,omitempty"`
}

type PushInfoRsp struct {
	Count   int         `json:"count"`
	Error   string      `json:"error"`
	Message interface{} `json:"message,omitempty"`
}
