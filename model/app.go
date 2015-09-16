package model

const (
	KEY = "eu1dr0c8znpa43blzy1wirzmk8jqdaon"
	// register
	SUCCESS        = "success"
	FAIL           = "fail"
	USERNAME_EXIST = "username_exist"
	SYSTEM_ERROR   = "system_error"
	SIGN_FAIL      = "sign_fail"
	// login
	USER_NO_ACTIVATE        = "user_no_activate"
	USERNAME_PASSWORD_ERROR = "username_password_error"
	USERNAME_NO_EXIST       = "username_no_exist"
)

type AppResult struct {
	State string `json:"state"` // 状态
	Error string `json:"error"` // 错误消息
}

// NewAppResult NewAppResult
func NewAppResult(state, err string) (ret *AppResult) {
	return &AppResult{
		State: state,
		Error: err,
	}
}

type AppUser struct {
	UserName  string `json:"userName" bson:"userName,omitempty"`
	Password  string `json:"password" bson:"password,omitempty"`
	Transtime string `json:"transtime" bson:"transtime,omitempty"`
}
