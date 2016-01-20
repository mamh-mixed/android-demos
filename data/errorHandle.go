package data

import (
	"bufio"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
	"io"
	"os"
	"strings"
)

func ReadMerIdsByTxt(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	var merIds []string

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		merIds = append(merIds, line)
		if err != nil {
			if err == io.EOF {
				return merIds, nil
			}
			return nil, err
		}
	}
	return merIds, nil
}

func AddSettRoleToTrans(merIds []string) error {

	for _, mId := range merIds {
		trans, _, err := mongo.SpTransColl.Find(&model.QueryCondition{
			MerId:     mId,
			StartTime: "2015-10-15 00:00:00",
			EndTime:   "2015-10-15 23:59:59",
			Size:      10000,
			Page:      1,
			ChanCode:  "WXP",
			// TransStatus: []string{model.TransSuccess},
			// TransType:   model.PayTrans,
		})
		if err != nil {
			return err
		}

		if len(trans) > 0 {
			log.Infof("merId=%s,trans=%d", mId, len(trans))
		}

		// var noRole int
		// for _, t := range trans {
		// 	// err = mongo.SpTransColl.AddSettRole("CIL", t.MerId, t.OrderNum)
		// 	// if err != nil {
		// 	// 	log.Error(err)
		// 	// }
		// 	if t.SettRole != "CIL" {
		// 		noRole++
		// 	}
		// }
		// log.Infof("noRole trans %d", noRole)
		// return nil
	}

	return nil

}
