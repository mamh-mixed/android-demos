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
	sftpAddr     = "192.168.1.33"
	sftpUserName = "11233404"
	sftpPassword = "sgfdgdfg435345435fdg"
)

var tranFileName = "IA502-%s.csv"

type transFlow struct {
}

func (t *transFlow) GenerateTransFlow(date string, agentCode string) {

	var transSettsWXP []model.TransSett
	var transSettsALP []model.TransSett
	transSettsWXP, err := mongo.SpTransSettColl.Find(&model.QueryCondition{Date: date, AgentCode: agentCode, ChanCode: "WXP", BlendType: "0"})
	if err != nil {
		return
	}

	if len(transSettsWXP) == 0 {
		log.Infof("these is no trans flow in date:%s, agentcode:%s, chanCode:%s", date, agentCode, "WXP")
		return
	}

	dateStr := strings.Replace(date, "-", "", -1)
	tranFileName = fmt.Sprintf(tranFileName, dateStr)

	var strBuffer = ""
	strBuffer += "清算日期,交易类型,交易时间,支付时间,客户代码,商户编号,终端编号,交易金额,订单号,渠道订单号,收单币种,收单交易金额,收单手续费,商户币种,商户交易金额,商户手续费,商户清算金额,交易渠道\r\n"

	generateFile(transSettsWXP, &strBuffer, dateStr) //微信

	transSettsALP, err = mongo.SpTransSettColl.Find(&model.QueryCondition{Date: date, AgentCode: agentCode, ChanCode: "ALP"})
	if err != nil {
		log.Infof("these is no trans flow in date:%s, agentcode:%s, chanCode:%s", date, agentCode, "ALP")
		return
	}

	generateFile(transSettsALP, &strBuffer, dateStr) //支付宝

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

	sftpFile, err := client.Create(tranFileName)
	if err != nil {
		log.Errorf("create the sftp file fail, error detail :%s", err)
		return
	}
	defer sftpFile.Close()

	_, err = sftpFile.Write([]byte(strBuffer))
	if err != nil {
		log.Errorf("write the sftp fail, error detail :%s", err)
		return
	}
}

func generateFile(data []model.TransSett, sBuffer *string, dateStr string) {
	for _, v := range data {
		*sBuffer += dateStr
		*sBuffer += ","
		if v.Trans.Busicd == "PURC" { //下单支付
			*sBuffer += "下单支付"
		} else if v.Trans.Busicd == "PAUT" {
			*sBuffer += "预下单"
		} else if v.Trans.Busicd == "VOID" {
			*sBuffer += "撤销"
		} else if v.Trans.Busicd == "REFD" {
			*sBuffer += "退款"
		} else if v.Trans.Busicd == "CANC" {
			*sBuffer += "取消"
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
		*sBuffer += fmt.Sprintf(",%0.2f,", float64(v.Trans.TransAmt)/100)
		*sBuffer += v.Trans.OrderNum
		*sBuffer += ","
		*sBuffer += v.Trans.ChanOrderNum
		*sBuffer += fmt.Sprintf(",CNY,%0.2f,%0.2f,CNY,%0.2f,%0.2f,%0.2f,", float64(v.Trans.TransAmt)/100, float64(v.AcqFee+v.MerFee)/100,
			float64(v.Trans.TransAmt)/100, float64(v.MerFee)/100, float64(v.Trans.TransAmt-v.MerFee)/100)
		if v.Trans.ChanCode == "WXP" {
			*sBuffer += "微信"
		} else if v.Trans.ChanCode == "ALP" {
			*sBuffer += "支付宝"
		}
		*sBuffer += "\r\n"
	}
}
