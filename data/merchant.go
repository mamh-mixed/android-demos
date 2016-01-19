// 用于从nodejs数据库导入商户数据
package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type merchant struct {
	UniqueId      bson.ObjectId `bson:"_id"`
	AgentCode     string        `bson:"inscd"`
	Clientid      string        `bson:"clientid"`
	CommodityName string        `bson:"commodityName"`
	ClientidName  string        `bson:"clientidName"`
	Alp           channel       `bson:"ALP"`
	Wxp           channel       `bson:"WXP"`
	Area          string        `bson:"area"`
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
	Type      string `bson:"type"`
	GoodsTag  string `bson:"goods_tag"`
}

type merCert struct {
	MerId    string
	HttpCert string
	HttpKey  string
}

func UpdateMerchant() error {
	ms, err := mongo.MerchantColl.FindNoUniqueId()
	if err != nil {
		return err
	}

	log.Infof("find merchant %d ", len(ms))

	for _, m := range ms {
		if m.UniqueId == "" {
			m.UniqueId = util.Confuse(m.MerId)
			mongo.MerchantColl.Update(m)
		}
	}

	return nil
}

// DoSyncMerchant 同步旧系统和新系统的商户
func DoSyncMerchant(path string) error {
	// connect()
	mers, err := readMerFromOldDB()
	if err != nil {
		return err
	}

	var addMers []merchant
	var updateCount int
	for _, om := range mers {
		_, err = mongo.MerchantColl.Find(strings.TrimSpace(om.Clientid))
		if err != nil {
			// add
			addMers = append(addMers, om)
		} else {
			// update
			// err = updateMerchantFromOldDB(om, nm)
			// if err != nil {
			// 	return err
			// }
			updateCount++
		}
	}

	log.Infof("修改：成功更新 %d 调数据", updateCount)

	// add
	if len(addMers) > 0 {
		err = addMerchantFromOldDB(addMers, path)
		if err != nil {
			return err
		}
	}

	return nil
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

func updateMerchantFromOldDB(om merchant, nm *model.Merchant) error {
	// 商户更新字段
	if nm.Detail.AcctName == "" {
		nm.Detail.AcctName = strings.TrimSpace(om.AcctName)
	}
	if nm.Detail.AcctNum == "" {
		nm.Detail.AcctNum = strings.TrimSpace(om.AcctNum)
	}
	if nm.Detail.BankId == "" {
		nm.Detail.BankId = strings.TrimSpace(om.BankId)
	}
	if nm.Detail.BankName == "" {
		nm.Detail.BankName = strings.TrimSpace(om.BankName)
	}
	if nm.Detail.OpenBankName == "" {
		nm.Detail.OpenBankName = strings.TrimSpace(om.OpenBank)
	}
	if nm.Detail.City == "" {
		nm.Detail.City = strings.TrimSpace(om.City)
	}
	if nm.Detail.Area == "" {
		nm.Detail.Area = strings.TrimSpace(om.Area)
	}
	if nm.Detail.TitleOne == "" {
		nm.Detail.TitleOne = strings.TrimSpace(om.Group.TitleOne)
	}
	if nm.Detail.TitleTwo == "" {
		nm.Detail.TitleTwo = strings.TrimSpace(om.Group.TitleTwo)
	}
	if nm.UniqueId == "" {
		nm.UniqueId = om.UniqueId.Hex()
	}

	if om.Alp.PartnerId != "" {
		updALPRouterPolicy(om)
	}

	if om.Wxp.MchId != "" {
		updWXPRouterPolicy(om)
	}

	mongo.MerchantColl.Update(nm)
	return nil
}

func updALPRouterPolicy(om merchant) {
	// 更新路由策略，费率
	r := mongo.RouterPolicyColl.Find(strings.TrimSpace(om.Clientid), "ALP")
	if r == nil {
		log.Errorf("找不到商户(%s)支付宝路由策略，请检查", om.Clientid)
		return
	}

	if r.ChanMerId != om.Alp.PartnerId {
		log.Errorf("商户(%s)对应的支付宝渠道商户号不一致，新=%s, 旧=%s", om.Clientid, r.ChanMerId, om.Alp.PartnerId)
		return
	}

	if r.MerFee == 0 {
		r.MerFee, _ = strconv.ParseFloat(om.Alp.MerFee, 64)
	}
	if r.AcqFee == 0 {
		r.AcqFee, _ = strconv.ParseFloat(om.Alp.AcqFee, 64)
	}

	// if strings.TrimSpace(om.Alp.Type) == "1" {
	// 	r.SettFlag = "CIL"
	// 	r.SettRole = "CIL"
	// } else {
	// 	r.SettFlag = "CHANNEL"
	// 	r.SettRole = "ALP"
	// }
}

func updWXPRouterPolicy(om merchant) {
	// 更新路由策略，费率
	r := mongo.RouterPolicyColl.Find(strings.TrimSpace(om.Clientid), "WXP")
	if r == nil {
		log.Errorf("找不到商户(%s)微信路由策略，请检查", om.Clientid)
		return
	}

	var chanMerId string
	if om.Wxp.SubMchId != "" {
		chanMerId = om.Wxp.SubMchId

	} else {
		chanMerId = om.Wxp.MchId
	}
	if r.ChanMerId != chanMerId {
		log.Errorf("商户(%s)对应的微信渠道商户号不一致，新=%s, 旧=%s", om.Clientid, r.ChanMerId, chanMerId)
		return
	}

	if r.MerFee == 0 {
		r.MerFee, _ = strconv.ParseFloat(om.Wxp.MerFee, 64)
	}
	if r.AcqFee == 0 {
		r.AcqFee, _ = strconv.ParseFloat(om.Wxp.AcqFee, 64)
	}

	// if strings.TrimSpace(om.Wxp.Type) == "1" {
	// 	r.SettFlag = "CIL"
	// 	r.SettRole = "CIL"
	// } else {
	// 	r.SettFlag = "CHANNEL"
	// 	r.SettRole = "WXP"
	// }

	// mongo.RouterPolicyColl.Insert(r)
}

// addMerchantFromOldDB 导入商户
func addMerchantFromOldDB(mers []merchant, path string) error {
	// 从某个目录获取证书信息
	merCerts, err := AddHttpCertFromFile(path)
	if err != nil {
		return err
	}

	// 更新新系统代理信息
	// aRec := 0
	agentMap := make(map[string]string)
	// for _, agent := range agents {
	// 	agentMap[agent.AgentCode] = agent.AgentName
	// 	a := &model.Agent{AgentCode: agent.AgentCode, AgentName: agent.AgentName}
	// 	err = mongo.AgentColl.Upsert(a)
	// 	if err == nil {
	// 		aRec++
	// 	}
	// }
	// log.Infof("在旧系统查找到 %d 条代理信息，成功插入新系统 %d 条。", len(agents), aRec)

	// 集团信息
	groupMap := make(map[string]*model.Group)

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
		m.Detail.Area = mer.Area
		m.Detail.TitleOne = mer.Group.TitleOne
		m.Detail.TitleTwo = mer.Group.TitleTwo
		m.Detail.GoodsTag = mer.Wxp.GoodsTag
		m.MerId = mer.Clientid
		m.Permission = []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd, model.Jszf, model.Qyzf}
		m.Remark = "old_system_data_" + time.Now().Format("20060102")
		m.SignKey = mer.SignKey
		m.UniqueId = mer.UniqueId.Hex()
		// 代理代码
		m.AgentCode = mer.AgentCode
		if agentName, ok := agentMap[mer.AgentCode]; ok {
			m.AgentName = agentName
		} else {
			agent, err := mongo.AgentColl.Find(mer.AgentCode)
			if err == nil {
				m.AgentName = agent.AgentName
				agentMap[mer.AgentCode] = agent.AgentName
			} else {
				log.Errorf("没有找到代理代码: %s, 商户号: %s", mer.AgentCode, m.MerId)
			}
		}

		// 集团
		m.GroupCode = mer.Group.GroupCode
		if m.GroupCode != "" {
			if g, ok := groupMap[m.GroupCode]; ok {
				m.GroupName = g.GroupName
			} else {
				g, err := mongo.GroupColl.Find(m.GroupCode)
				if err == nil {
					m.GroupName = g.GroupName
					groupMap[m.GroupCode] = g
				} else {
					log.Errorf("没有找到集团代码: %s, 商户号: %s", m.GroupCode, m.MerId)
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
			alp, err := mongo.ChanMerColl.Find("ALP", mer.Alp.PartnerId)
			if err != nil {
				// 导入渠道商户
				alp = &model.ChanMer{}
				alp.ChanMerId = mer.Alp.PartnerId
				alp.SignKey = mer.Alp.Md5
				alp.ChanCode = "ALP"
				err = mongo.ChanMerColl.Add(alp)
				if err != nil {
					return err
				}
			}

			// 路由策略
			ralp := &model.RouterPolicy{}
			ralp.MerId = m.MerId
			ralp.ChanCode = alp.ChanCode
			ralp.CardBrand = alp.ChanCode
			ralp.ChanMerId = alp.ChanMerId
			ralp.MerFee, _ = strconv.ParseFloat(mer.Alp.MerFee, 64)
			ralp.AcqFee, _ = strconv.ParseFloat(mer.Alp.AcqFee, 64)
			if strings.TrimSpace(mer.Alp.Type) == "1" {
				ralp.SettFlag = "CIL"
				ralp.SettRole = "CIL"
			} else {
				ralp.SettFlag = "CHANNEL"
				ralp.SettRole = "ALP"
			}
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
			// acqFee, _ := strconv.ParseFloat(mer.Wxp.AcqFee, 32)
			// merFee, _ := strconv.ParseFloat(mer.Wxp.MerFee, 32)
			// wxp.AcqFee = float32(acqFee)
			// wxp.MerFee = float32(merFee)
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
			_, err = mongo.ChanMerColl.Find("WXP", wxpMerId)
			if err != nil {
				err = mongo.ChanMerColl.Add(wxp)
				if err != nil {
					return err
				}
			}

			// 路由策略
			rwxp := &model.RouterPolicy{}
			rwxp.MerId = m.MerId
			rwxp.ChanCode = wxp.ChanCode
			rwxp.CardBrand = wxp.ChanCode
			rwxp.ChanMerId = wxp.ChanMerId
			rwxp.MerFee, _ = strconv.ParseFloat(mer.Wxp.MerFee, 64)
			rwxp.AcqFee, _ = strconv.ParseFloat(mer.Wxp.AcqFee, 64)
			if strings.TrimSpace(mer.Wxp.Type) == "1" {
				rwxp.SettFlag = "CIL"
				rwxp.SettRole = "CIL"
			} else {
				rwxp.SettFlag = "CHANNEL"
				rwxp.SettRole = "WXP"
			}
			err = mongo.RouterPolicyColl.Insert(rwxp)
			if err != nil {
				return err
			}
		}
	}

	// gRec := 0
	// for _, g := range groupMap {
	// 	err = mongo.GroupColl.Upsert(g)
	// 	if err == nil {
	// 		gRec++
	// 	}
	// }
	// log.Infof("新增：在旧系统查找到 %d 条集团信息，成功插入新系统 %d 条。", len(groupMap), gRec)
	log.Infof("新增：在旧系统查找到 %d 条商户信息，成功插入新系统 %d 条。", len(mers), count)

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
		fmt.Printf("unable connect to mongo %s: %s\n", url, err)
		os.Exit(1)
	}

	session.SetMode(mgo.Strong, true) //需要指定为Eventual
	session.SetSafe(nil)
	session.SetSocketTimeout(time.Hour * 1)

	saomaDB = session.DB(dbname)

	log.Infof("connected to mongodb host `%s` and database `%s`", url, dbname)
}
