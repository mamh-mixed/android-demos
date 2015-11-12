package mongo

import (
	"time"

	"github.com/CardInfoLink/quickpay/model"
)

type masterLogCollection struct {
	name string
}

// MasterLogColl 平台日志 Collection
var MasterLogColl = masterLogCollection{"masterlog"}

func (col *masterLogCollection) Insert(log *model.MasterLog) error {
	log.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	err := database.C(col.name).Insert(log)
	return err
}
