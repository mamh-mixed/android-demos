package settle

import (
	"fmt"
	"strings"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	sftpAddr     = "120.26.78.97:22"
	sftpUserName = "webapp"
	sftpPassword = "Cilxl123$"
)

var filePath = "/report/%s/in/%s"
var fileName = "IA502-%s.csv"

type transFlow struct {
}

func (t *transFlow) GenerateTransFlow(date string, agentCode string) {

	agent, err := mongo.AgentColl.Find(agentCode)
	if err != nil {
		log.Errorf("cann't find the agent agentCode is %s", agentCode)
		return
	}
	var transSettsWXP []model.TransSett
	var transSettsALP []model.TransSett
	transSettsWXP, err = mongo.SpTransSettColl.Find(&model.QueryCondition{Date: date, AgentCode: agentCode, ChanCode: "WXP", BlendType: "0"})
	if err != nil {
		return
	}

	if len(transSettsWXP) == 0 {
		log.Infof("these is no trans flow in date:%s, agentcode:%s, chanCode:%s", date, agentCode, "WXP")
		return
	}

	dateStr := strings.Replace(date, "-", "", -1)
	filePath = fmt.Sprintf(filePath, agentCode, dateStr)
	fileName = fmt.Sprintf(fileName, dateStr)

	var strBuffer = ""
	strBuffer += "清算日期,交易类型,交易时间,支付时间,客户代码,商户编号,终端编号,交易金额,订单号,渠道订单号,收单币种,收单交易金额,收单手续费,商户币种,商户交易金额,商户手续费,商户清算金额,交易渠道\r\n"

	generateFile(transSettsWXP, dateStr, agent.WxpCost, &strBuffer) //微信

	transSettsALP, err = mongo.SpTransSettColl.Find(&model.QueryCondition{Date: date, AgentCode: agentCode, ChanCode: "ALP"})
	if err != nil {
		log.Infof("these is no trans flow in date:%s, agentcode:%s, chanCode:%s", date, agentCode, "ALP")
		return
	}

	generateFile(transSettsALP, dateStr, agent.AlpCost, &strBuffer) //支付宝

	fmt.Println(strBuffer)

	var authMethods []ssh.AuthMethod
	// add password
	authMethods = append(authMethods, ssh.Password(sftpPassword))
	config := &ssh.ClientConfig{
		User: sftpUserName,
		Auth: authMethods,
	}

	// 建立连接
	conn, err := ssh.Dial("tcp", sftpAddr, config)
	if err != nil {
		log.Errorf("fail to connect sftp service, error: %s", err)
		return
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Errorf("fail to get sftp client, error: %s", err)
		return
	}
	defer client.Close()

	err = client.Mkdir(filePath)
	if err != nil {
		log.Errorf("create dir fail , dir is %s, error is %s", filePath, err)
		return
	}

	filePath += "/" + fileName
	sftpFile, err := client.Create(filePath)
	if err != nil {
		log.Errorf("create the sftp file fail, error detail :%s, filename:%s", err, filePath)
		return
	}
	defer sftpFile.Close()

	_, err = sftpFile.Write([]byte(strBuffer))
	if err != nil {
		log.Errorf("write the sftp fail, error detail :%s", err)
		return
	}
}

func generateFile(data []model.TransSett, dateStr string, agentFee float64, sBuffer *string) {
	for _, v := range data {
		amt := float64(v.Trans.TransAmt) / 100
		*sBuffer += dateStr
		*sBuffer += ","
		isReverse := false
		if v.Trans.Busicd == "PURC" { //下单支付
			*sBuffer += "下单支付"
		} else if v.Trans.Busicd == "PAUT" {
			*sBuffer += "预下单"
		} else if v.Trans.Busicd == "VOID" {
			*sBuffer += "撤销"
			isReverse = true
		} else if v.Trans.Busicd == "REFD" {
			*sBuffer += "退款"
			isReverse = true
		} else if v.Trans.Busicd == "CANC" {
			*sBuffer += "取消"
			isReverse = true
		}
		*sBuffer += ","
		*sBuffer += v.Trans.CreateTime
		*sBuffer += ","
		*sBuffer += v.Trans.UpdateTime
		*sBuffer += ","
		*sBuffer += v.Trans.AgentCode
		*sBuffer += ","
		*sBuffer += v.Trans.MerId
		*sBuffer += ","
		*sBuffer += v.Trans.Terminalid
		*sBuffer += fmt.Sprintf(",%0.2f,", amt)
		*sBuffer += v.Trans.OrderNum
		*sBuffer += ","
		*sBuffer += v.Trans.ChanOrderNum
		if isReverse {
			*sBuffer += fmt.Sprintf(",CNY,-%0.2f,-%0.2f,CNY,-%0.2f,-%0.2f,-%0.2f,", amt, amt*agentFee+float64(v.AcqFee/100), //收单币种,收单交易金额,收单手续费,商户币种,
				amt, float64(v.MerFee)/100, float64(v.Trans.TransAmt-v.MerFee)/100) //商户交易金额,商户手续费,商户清算金额
		} else {
			*sBuffer += fmt.Sprintf(",CNY,%0.2f,%0.2f,CNY,%0.2f,%0.2f,%0.2f,", amt, amt*agentFee+float64(v.AcqFee/100), //收单币种,收单交易金额,收单手续费,商户币种,
				amt, float64(v.MerFee)/100, float64(v.Trans.TransAmt-v.MerFee)/100) //商户交易金额,商户手续费,商户清算金额
		}

		if v.Trans.ChanCode == "WXP" {
			*sBuffer += "微信"
		} else if v.Trans.ChanCode == "ALP" {
			*sBuffer += "支付宝"
		}
		*sBuffer += "\r\n"
	}
}
