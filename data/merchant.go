// 用于从nodejs数据库导入商户数据
package data

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type merchant struct {
	UniqueId      bson.ObjectId `bson:"_id"`
	AgentCode     string        `bson:"inscd"`
	Clientid      string        `bson:"clientid"`
	CommodityName string        `bson:"commodityName"`
	ClientidName  string        `bson:"clientidName"`
	Alp           channel       `bson:"ALP"`
	Wxp           channel       `bson:"WXP"`
	City          string        `bson:"city"`
	AcctNum       string        `bson:"account"`
	AcctName      string        `bson:"accountName"`
	BankName      string        `bson:"bankName"`
	BankId        string        `bson:"bankNum"`
	OpenBank      string        `bson:"openBank"`
	SignKey       string        `bson:"merchantMd5"`
	SignRule      string        `bson:"signRule"`
	Group         struct {
		GroupCode string `bson:"merId"`
		GroupName string `bson:"commodityName"`
		TitleOne  string `bson:"title_one"`
		TitleTwo  string `bson:"title_two"`
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
	// log.Debug(certsMap)
	return certsMap, nil
}

// AddMerchantFromOldDB 导入商户
func AddMerchantFromOldDB(path string) error {
	// 从某个目录获取证书信息
	merCerts, err := AddHttpCertFromFile(path)
	if err != nil {
		return err
	}

	// 建立连接
	connect()

	// 获取老系统代理信息
	agents, err := readAgentFromOldDB()
	if err != nil {
		return err
	}

	// 更新新系统代理信息
	aRec := 0
	agentMap := make(map[string]string)
	for _, agent := range agents {
		agentMap[agent.AgentCode] = agent.AgentName
		a := &model.Agent{AgentCode: agent.AgentCode, AgentName: agent.AgentName}
		err = mongo.AgentColl.Upsert(a)
		if err == nil {
			aRec++
		}
	}
	log.Infof("在旧系统查找到 %d 条代理信息，成功插入新系统 %d 条。", len(agents), aRec)

	// 集团信息
	groupMap := make(map[string]*model.Group)

	// 读取商户
	mers, err := readMerFromOldDB()
	if err != nil {
		return err
	}
	var count = 0
	for _, mer := range mers {
		mer.Clientid = strings.TrimSpace(mer.Clientid)
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
		m.Detail.BankId = mer.BankId
		m.Detail.BankName = mer.BankName
		m.Detail.OpenBankName = mer.OpenBank
		m.Detail.TitleOne = mer.Group.TitleOne
		m.Detail.TitleTwo = mer.Group.TitleTwo
		m.MerId = mer.Clientid
		m.Permission = []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd, model.Jszf, model.Qyzf}
		m.Remark = "old_system_data"
		m.SignKey = mer.SignKey
		m.UniqueId = mer.UniqueId.Hex()
		// 代理代码
		m.AgentCode = mer.AgentCode
		m.AgentName = agentMap[m.AgentCode]

		// 集团
		m.GroupCode = mer.Group.GroupCode
		m.GroupName = mer.Group.GroupName
		if m.GroupCode != "" {
			if _, ok := groupMap[m.GroupCode]; !ok {
				// 没有则存放
				groupMap[m.GroupCode] = &model.Group{
					GroupCode: m.GroupCode,
					GroupName: m.GroupName,
					AgentCode: m.AgentCode,
					AgentName: m.AgentName,
				}
			}
		}
		if m.SignKey != "" {
			m.IsNeedSign = true
		}
		if mer.SignRule != "" && mer.SignRule == "0" {
			m.IsNeedSign = false
		}
		err = mongo.MerchantColl.Insert(m)
		if err != nil {
			return err
		}

		if mer.Alp.PartnerId != "" {
			// 导入渠道商户
			alp := &model.ChanMer{}
			alp.ChanMerId = mer.Alp.PartnerId
			alp.SignKey = mer.Alp.Md5
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
			// 只保存子渠道商户
			wxp.SignKey = mer.Wxp.Md5
			wxp.ChanCode = "WXP"
			wxp.WxpAppId = mer.Wxp.AppId
			acqFee, _ := strconv.ParseFloat(mer.Wxp.AcqFee, 32)
			merFee, _ := strconv.ParseFloat(mer.Wxp.MerFee, 32)
			wxp.AcqFee = float32(acqFee)
			wxp.MerFee = float32(merFee)
			// 非受理商模式
			wxpMerId := ""
			if mer.Wxp.SubMchId != "" {
				wxpMerId = mer.Wxp.SubMchId
				a, err := mongo.ChanMerColl.Find("WXP", mer.Wxp.MchId)
				if err != nil {
					log.Errorf("受理商模式下，没找到受理商商户，商户ID为：%s", mer.Wxp.MchId)
				}
				wxp.AgentMer = a
				wxp.SignKey = "" // 清空证书以及appid，这时的数据是大商户的。
				wxp.WxpAppId = ""
				wxp.IsAgentMode = true
			} else {
				wxpMerId = mer.Wxp.MchId
				// 保存证书
				if merCert, ok := merCerts[mer.Clientid]; ok {
					wxp.HttpCert = merCert.HttpCert
					wxp.HttpKey = merCert.HttpKey
				} else {
					log.Errorf("找不到商户：%s, 相应证书。", mer.Clientid)
				}
			}
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

	gRec := 0
	for _, g := range groupMap {
		err = mongo.GroupColl.Upsert(g)
		if err == nil {
			gRec++
		}
	}
	log.Infof("在旧系统查找到 %d 条集团信息，成功插入新系统 %d 条。", len(groupMap), gRec)
	log.Infof("在旧系统查找到 %d 条商户信息，成功插入新系统 %d 条。", len(mers), count)

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

	session.SetMode(mgo.Strong, true) //需要指定为Eventual
	session.SetSafe(nil)
	session.SetSocketTimeout(time.Hour * 1)

	saomaDB = session.DB(dbname)

	log.Infof("connected to mongodb host `%s` and database `%s`", url, dbname)
}
