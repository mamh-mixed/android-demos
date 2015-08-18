// 用于从nodejs数据库导入商户数据
package data

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2"
	"os"
	"strconv"
)

type merchant struct {
	AgentCode     string  `bson:"inscd"`
	Clientid      string  `bson:"clientid"`
	CommodityName string  `bson:"commodityName"`
	ClientidName  string  `bson:"clientidName"`
	Alp           channel `bson:"ALP"`
	Wxp           channel `bson:"WXP"`
	City          string  `bson:"city"`
	AcctNum       string  `bson:"account"`
	AcctName      string  `bson:"accountName"`
	SignKey       string  `bson:"merchantMd5"`
	Agent         struct {
		AgentCode string `bson:"merId"`
		AgentName string `bson:"commodityName"`
	} `bson:"headMerchant"`
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
		m.AgentCode = mer.AgentCode
		m.Permission = []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd, model.Jszf, model.Qyfk}
		m.Remark = "old_system_data"
		m.SignKey = mer.SignKey
		// 代理代码
		if mer.Agent.AgentCode != "" {
			m.AgentCode = mer.Agent.AgentCode
			m.AgentName = mer.Agent.AgentName
		} else {
			m.AgentCode = "99911888"
			m.AgentName = "讯联O2O机构"
		}
		// TODO:集团
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
			acqFee, _ := strconv.ParseFloat(mer.Alp.AcqFee, 32)
			merFee, _ := strconv.ParseFloat(mer.Alp.MerFee, 32)
			alp.AcqFee = float32(acqFee)
			alp.MerFee = float32(merFee)
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
			// 非受理商模式
			wxpMerId := ""
			if mer.Wxp.SubMchId != "" {
				wxpMerId = mer.Wxp.SubMchId
				a, err := mongo.ChanMerColl.Find("WXP", mer.Wxp.MchId)
				if err != nil {
					log.Errorf("受理商模式下，没找到受理商商户，商户ID为：%s", mer.Wxp.MchId)
				}
				wxp.AgentMer = a
				wxp.IsAgentMode = true
			} else {
				wxpMerId = mer.Wxp.MchId
			}
			// 只保存子渠道商户
			wxp.SignCert = mer.Wxp.Md5
			wxp.ChanCode = "WXP"
			wxp.WxpAppId = mer.Wxp.AppId
			acqFee, _ := strconv.ParseFloat(mer.Wxp.AcqFee, 32)
			merFee, _ := strconv.ParseFloat(mer.Wxp.MerFee, 32)
			wxp.AcqFee = float32(acqFee)
			wxp.MerFee = float32(merFee)
			wxp.ChanMerId = wxpMerId
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
