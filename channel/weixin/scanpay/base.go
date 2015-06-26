package scanpay

import (
	"github.com/omigo/log"
	"github.com/omigo/validator"
)

// BaseReq 只是为了注入签名方便
type BaseReq interface {
	GenSign()
}

// BaseResp 只是为了传参方便
type BaseResp interface{}

func base(d BaseReq, r BaseResp) (err error) {
	if err := validator.Validate(d); err != nil {
		log.Errorf("validate error, %s", err)
		return err
	}

	err = sendRequest(d, r)
	if err != nil {
		log.Errorf("weixin device scan phone request error: %s", err)
		return err
	}
	return nil
}
