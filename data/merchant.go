// 用于从nodejs数据库导入商户数据
package data

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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
	Group         struct {
		GroupCode string `bson:"merId"`
		GroupName string `bson:"commodityName"`
	} `bson:"headMerchant"`
}

type Agent struct {
	AgentCode string `bson:"inscd"`
	AgentName string `bson:"name"`
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

type merCert struct {
	MerId    string
	HttpCert string
	HttpKey  string
}

// AddHttpCertFromFile 导入文件
func AddHttpCertFromFile(root string) (map[string]merCert, error) {

	var filePaths []string
	// 读取该目录下所有文件名称
	filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			filename := f.Name()
			// log.Debug(filename)
			if matched, _ := regexp.MatchString(`^\d+$`, filename); matched {
				filePaths = append(filePaths, filename)
			}
		}
		return nil
	})
	// log.Debug(filePaths)
	certsMap := make(map[string]merCert)
	for _, merId := range filePaths {
		certPath := root + "/" + merId + "/apiclient_cert.pem"
		cert, err := os.Open(certPath)
		if err != nil {
			return certsMap, err
		}
		certBytes, err := ioutil.ReadAll(cert)
		if err != nil {
			return certsMap, err
		}
		cert.Close()
		keyPath := root + "/" + merId + "/apiclient_key.pem"
		key, err := os.Open(keyPath)
		if err != nil {
			return certsMap, err
		}
		keyBytes, err := ioutil.ReadAll(key)
		if err != nil {
			return certsMap, err
		}
		key.Close()
		c := merCert{MerId: merId, HttpCert: string(certBytes), HttpKey: string(keyBytes)}
		certsMap[merId] = c
	}

	return certsMap, nil
}

// AddMerchantFromOldDB 导入商户
func AddMerchantFromOldDB(path string) error {
	merCerts, err := AddHttpCertFromFile(path)
	if err != nil {
		return err
	}
	connect()
	agents, err := readAgentFromOldDB()
	if err != nil {
		return err
	}
	agentMap := make(map[string]string)
	for _, agent := range agents {
		agentMap[agent.AgentCode] = agent.AgentName
	}
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
		// 基本信息
		m := &model.Merchant{}
		m.Detail.AcctName = mer.AcctName
		m.Detail.AcctNum = mer.AcctNum
		m.Detail.City = mer.City
		m.Detail.CommodityName = mer.CommodityName
		m.Detail.MerName = mer.ClientidName
		m.MerId = mer.Clientid
		m.Permission = []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd, model.Jszf, model.Qyfk}
		m.Remark = "old_system_data"
		m.SignKey = mer.SignKey
		// 集团
		m.GroupCode = mer.Group.GroupCode
		m.GroupName = mer.Group.GroupName
		// 代理代码
		m.AgentCode = mer.AgentCode
		m.AgentName = agentMap[m.AgentCode]

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
			// 保存证书
			if merCert, ok := merCerts[wxpMerId]; ok {
				wxp.HttpCert = merCert.HttpCert
				wxp.HttpKey = merCert.HttpKey
			} else {
				log.Errorf("找不到商户：%s, 相应证书。", wxpMerId)
			}

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
	var mers []merchant
	err := saomaDB.C("merchant").Find(nil).All(&mers)
	return mers, err
}

func readAgentFromOldDB() ([]Agent, error) {
	var agents []Agent
	err := saomaDB.C("acq").Find(nil).All(&agents)
	return agents, err
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
