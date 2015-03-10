package mongo

import (
	"gopkg.in/mgo.v2/bson"
)

// GetRespCode 根据传入的code类型得到Resp对象
// 如：code : "270001",codeType : "cfca"
// 表示将中金应答码转为系统应答码
// codeType : sys,cfca,....
func GetRespCode(code string, codeType string) (resp *Resp) {

	resp = &Resp{}
	switch codeType {
	case "sys":
		db.respCode.Find(bson.M{"respcode": code}).One(resp)
	case "cfca":
		db.respCode.Find(bson.M{"cfca.code": code}).One(resp)
	default:
		resp = nil
	}

	return
}
