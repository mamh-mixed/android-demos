// 用于从nodejs数据库导入商户数据
package data

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2"
	"os"
)

type merchant struct {
	Inscd         string  `bson:"inscd"`
	Clientid      string  `bson:"clientid"`
	CommodityName string  `bson:"commodityName"`
	ClientidName  string  `bson:"clientidName"`
	Alp           channel `bson:"ALP"`
	Wxp           channel `bson:"WXP"`
	City          string  `bson:"city"`
	AcctNum       string  `bson:"account"`
	AcctName      string  `bson:"accountName"`
	SignKey       string  `bson:"merchantMd5"`
}

type channel struct {
	Md5       string `bson:"md5"`
	MchId     string `bson:"mch_id"`
	AppId     string `bson:"appid"`
	AcqFee    string `bson:"acqfee"`
	MerFee    string `bson:"merfee"`
	Fee       string `bson:"fee"`
	PartnerId string `bson:"partnerId"`
	SubMchId  string `bson:"sub_mch_id"`
}

func AddMerchantFromOldDB() error {

	mers, err := readMerFromOldDB()
	if err != nil {
		return err
	}
	var count = 0
	for _, mer := range mers {

		if mer.Clientid == "" {
			continue
		}
		count++
		// 导入商户
		m := &model.Merchant{}
		m.Detail.AcctName = mer.AcctName
		m.Detail.AcctNum = mer.AcctNum
		m.Detail.City = mer.City
		m.Detail.CommodityName = mer.CommodityName
		m.Detail.MerName = mer.ClientidName
		m.MerId = mer.Clientid
		m.InsCode = mer.Inscd
		m.Permission = []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd}
		m.Remark = "老扫码系统商户"
		m.SignKey = mer.SignKey
		if m.SignKey != "" {
			m.IsNeedSign = true
		}
		err = mongo.MerchantColl.Insert(m)
		if err != nil {
			return err
		}

		if mer.Alp.PartnerId != "" {
			// 导入渠道商户
			alp := &model.ChanMer{}
			alp.ChanMerId = mer.Alp.PartnerId
			alp.SignCert = mer.Alp.Md5
			alp.ChanCode = "ALP"
			alp.AcqFee = mer.Alp.AcqFee
			alp.MerFee = mer.Alp.MerFee
			err = mongo.ChanMerColl.Add(alp)
			if err != nil {
				return err
			}
			// 路由策略
			ralp := &model.RouterPolicy{}
			ralp.MerId = m.MerId
			ralp.ChanCode = alp.ChanCode
			ralp.CardBrand = alp.ChanCode
			ralp.ChanMerId = alp.ChanMerId
			err = mongo.RouterPolicyColl.Insert(ralp)
			if err != nil {
				return err
			}
		}

		if mer.Wxp.MchId != "" {
			// 导入渠道商户
			wxp := &model.ChanMer{}
			wxp.SignCert = mer.Wxp.Md5
			wxp.ChanCode = "WXP"
			wxp.WxpAppId = mer.Wxp.AppId
			wxp.AcqFee = mer.Wxp.AcqFee
			wxp.MerFee = mer.Wxp.MerFee
			wxp.ChanMerId = mer.Wxp.MchId
			err = mongo.ChanMerColl.Add(wxp)
			if err != nil {
				return err
			}

			// 路由策略
			rwxp := &model.RouterPolicy{}
			rwxp.MerId = m.MerId
			rwxp.ChanCode = wxp.ChanCode
			rwxp.CardBrand = wxp.ChanCode
			rwxp.ChanMerId = wxp.ChanMerId

			// 代理商模式
			if mer.Wxp.SubMchId != "" {
				rwxp.IsAgent = true
				rwxp.SubMerId = mer.Wxp.SubMchId
			}
			err = mongo.RouterPolicyColl.Insert(rwxp)
			if err != nil {
				return err
			}
		}
	}

	log.Debugf("success add %d records", count)
	return nil
}

func readMerFromOldDB() ([]merchant, error) {
	connect()
	var mers []merchant
	err := saomaDB.C("merchant").Find(nil).All(&mers)
	return mers, err
}

var saomaDB *mgo.Database
var url = "mongodb://saoma:saoma@211.147.72.70:10001/online"
var dbname = "online"

// Connect 连接到nodejs扫码程序数据库
func connect() {

	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("unable connect to mongodb server %s\n", err)
		os.Exit(1)
	}

	session.SetMode(mgo.Eventual, true) //需要指定为Eventual
	session.SetSafe(&mgo.Safe{})

	saomaDB = session.DB(dbname)

	log.Infof("connected to mongodb host `%s` and database `%s`", url, dbname)
}
