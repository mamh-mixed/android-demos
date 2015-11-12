package mongo

import "github.com/CardInfoLink/quickpay/model"

type masterLogCollection struct {
	name string
}

// MasterLogColl 平台日志 Collection
var MasterLogColl = masterLogCollection{"masterlog"}

func (col *masterLogCollection) Insert(log *model.MasterLog) error {
	err := database.C(col.name).Insert(log)
	return err
}
