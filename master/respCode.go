package master

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

type respCode struct{}

var RespCode respCode

// FindOne 根据指定的应答查询应答码信息
func (r *respCode) FindOne(code string) (result *model.ResultBody) {
	log.Debugf("respCode=%s", code)

	respCode, err := mongo.RespCodeColl.FindOne(code)
	if err != nil {
		log.Errorf("find response code(%s) error: %s", code, err)
		return model.NewResultBody(1, "查询失败")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    respCode,
	}

	return result
}
