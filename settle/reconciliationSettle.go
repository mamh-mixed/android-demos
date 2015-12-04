package settle

import (
	"bytes"
	"fmt"
	"github.com/CardInfoLink/quickpay/mongo"
	"time"
	//"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
)

type reconciliatReportData struct {
	insCode     string  //客户代码
	insName     string  //客户名称
	chcd        string  //渠道名称
	role        string  //清算角色
	orderNum    int     //交易笔数
	orderAct    float64 //交易金额
	orderNetFee float64 //商户手续费
	orderNetAct float64 //商户应收金额
	orderAcqFee float64 //讯联手续费
	orderAcqAct float64 //讯联应收金额
}

type errorReportData struct {
	chanCode    string  //渠道编号
	chanName    string  //渠道名称
	merCode     string  //商户编号
	chanMerCode string  //虚拟商户号
	merName     string  //商户名称
	orderId     string  //订单号
	transTime   string  //交易时间
	transType   string  //交易类型
	transAct    float64 //交易金额
	compareType string  //对账标识
}

func Reconciliat(startTime string, endTime string) error {

	coll, err := mongo.SpTransSettleColl.FindBySettleTime(startTime, endTime)
	if err != nil {
		log.Errorf("fail to find trans: %s", err)
		return err
	}

	blendMMap := make(map[string]map[string][]model.TransSett) //外部key为商户号，内部key为订单号

	//BlendRecord(coll)

	merMap := make(map[string]string)

	//报表名称
	filename := "账务对账报表IC002.xlsx"

	reconciliatMMap := make(map[string]map[string]reconciliatReportData) //外key为机构号，内map的key为渠道编号

	//处理数据
	for _, v := range coll {
		element := v.Trans
		if element.ChanCode == "ALP" {
			//fmt.Println(element)
			merMap[element.ChanMerId] = element.MerId
		}
		if element.RespCode == "00" { //成功则对账
			reportDataMap, ret := reconciliatMMap[element.AgentCode]
			act := float64(element.TransAmt / 100)
			//计算每条记录的讯联手续费
			chanMer, err := mongo.ChanMerColl.Find(element.ChanCode, element.ChanMerId)
			if err != nil {
				log.Errorf("can't find acqFee, chanMerID:%s, ChanCode:%s, err:%s", element.ChanMerId, element.ChanCode, err)
				continue
			}
			acqFee := float64(chanMer.AcqFee) * act
			blendMap, ret1 := blendMMap[element.ChanMerId]
			if !ret1 {
				blendMap = make(map[string][]model.TransSett)
			}
			blendArray, ret1 := blendMap[element.ChanOrderNum]
			if !ret1 {
				blendArray = make([]model.TransSett, 0)
			}

			if ret { //如果map存在该机构
				//清算角色
				var reportData reconciliatReportData
				ret1 := false
				reportData, ret1 = reportDataMap[element.SettRole]
				if !ret1 {
					reportData.insCode = element.AgentCode
					reportData.insName = element.AgentName
					reportData.chcd = element.ChanCode
					reportData.role = element.SettRole
				}
				//下单和预下单相加，退款、取消和撤销相减
				if (element.Busicd == "PURC") || (element.Busicd == "PAUT") {
					reportData.orderAct += act
					reportData.orderNetFee += float64(element.NetFee) / 100
					reportData.orderNetAct += (act - float64(element.NetFee)/100)
					reportData.orderAcqFee += acqFee
					reportData.orderAcqAct += (act - acqFee)
					reportData.orderNum++
					v.SettFlag = model.SettOK
					v.SettDate = time.Now().Format("2015-12-03")
					blendArray = append(blendArray, v)
				} else if (element.Busicd == "VOID") || (element.Busicd == "REFD") || (element.Busicd == "CANC") {
					if (element.Busicd == "CANC") && (element.TransAmt == 0) {
						continue //取消，如果金额为0，是没有支付
					}
					reportData.orderAct -= act
					reportData.orderNetFee -= float64(element.NetFee) / 100
					reportData.orderNetAct -= (act - float64(element.NetFee)/100)
					reportData.orderAcqFee -= acqFee
					reportData.orderAcqAct -= (act - acqFee)
					reportData.orderNum++
					v.SettFlag = model.SettOK
					v.SettDate = time.Now().Format("2015-12-03")
					blendArray = append(blendArray, v)
				}

				if !ret1 {
					reportDataMap[element.SettRole] = reportData
				}

			} else { //map中没有查到该机构
				reportDataMap = make(map[string]reconciliatReportData)
				var reportData reconciliatReportData
				reportData.insCode = element.AgentCode
				reportData.insName = element.AgentName
				reportData.chcd = element.ChanCode
				reportData.role = element.SettRole
				//下单和预下单相加，退款、取消和撤销相减
				if (element.Busicd == "PURC") || (element.Busicd == "PAUT") {
					reportData.orderAct = act
					reportData.orderNetFee = float64(element.NetFee) / 100
					reportData.orderNetAct = (act - float64(element.NetFee)/100)
					reportData.orderAcqFee = acqFee
					reportData.orderAcqAct = (act - acqFee)
					reportData.orderNum = 1
					v.SettFlag = model.SettOK
					v.SettDate = time.Now().Format("2015-12-03")
					blendArray = append(blendArray, v)
				} else if (element.Busicd == "VOID") || (element.Busicd == "REFD") || (element.Busicd == "CANC") {
					if (element.Busicd == "CANC") && (element.TransAmt == 0) {
						continue //取消，如果金额为0，是没有支付
					}
					reportData.orderAct = act * (-1.0)
					reportData.orderNetFee = float64(element.NetFee) / 100 * (-1.0)
					reportData.orderNetAct = (act - float64(element.NetFee)/100) * (-1.0)
					reportData.orderAcqFee = acqFee * (-1.0)
					reportData.orderAcqAct = (act - acqFee) * (-1.0)
					reportData.orderNum = 1
					v.SettFlag = model.SettOK
					v.SettDate = time.Now().Format("2015-12-03")
					blendArray = append(blendArray, v)
				}

				reportDataMap[element.SettRole] = reportData
				reconciliatMMap[element.AgentCode] = reportDataMap
			}

			if v.SettFlag == model.SettOK {
				blendMap[element.ChanOrderNum] = blendArray
				blendMMap[element.ChanMerId] = blendMap
			}
		}
	}
	//勾兑
	BlendRecord(blendMMap, merMap, startTime, endTime)

	date := []byte(startTime)[:10]

	if reconciliatMMap != nil {
		excel := genReconciliatReportExcel(reconciliatMMap, string(date))

		var buf []byte
		bf := bytes.NewBuffer(buf)
		//写到buf里
		excel.Write(bf)
		spath := "E:\\"
		spath += filename
		excel.Save(spath)
		/*
			// 上传到七牛
			err = qiniu.Put(filename, int64(bf.Len()), bf)
			if err != nil {
				log.Errorf("upload settReport excel err: %s", err)
			}*/
	}

	return nil
}

// genSpTransferSettleReportExcel 对账报表
func genReconciliatReportExcel(data map[string]map[string]reconciliatReportData, date string) *xlsx.File {
	var file = xlsx.NewFile()
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	sheet, _ = file.AddSheet("财务对账报表")

	// 第一行
	row = sheet.AddRow()
	row.SetHeightCM(0.91)
	cell = row.AddCell()
	cell.Merge(8, 0) // 9个单元格
	cell.SetValue("云收银资金划拨财务报表")
	style := xlsx.NewStyle()

	style.Alignment = xlsx.Alignment{Horizontal: "center", Vertical: "center"}
	style.Font = *xlsx.NewFont(20, "宋体")
	style.ApplyAlignment = true
	style.ApplyFont = true
	cell.SetStyle(style)

	// 第二行
	bodyStyle := xlsx.NewStyle()
	bodyStyle.Font = xlsx.Font{
		Size: 10,
		Name: "Times New Roman",
		Bold: true,
	}
	bodyStyle.Alignment = xlsx.Alignment{Horizontal: "center", Vertical: "center"}
	bodyStyle.Border = xlsx.Border{
		Left:        "thin",
		LeftColor:   "FF999999",
		Right:       "thin",
		RightColor:  "FF999999",
		Top:         "thin",
		TopColor:    "FF999999",
		Bottom:      "thin",
		BottomColor: "FF999999",
	}
	bodyStyle.ApplyFont = true
	bodyStyle.ApplyAlignment = true
	bodyStyle.ApplyBorder = true

	row = sheet.AddRow()
	row.SetHeightCM(1.83)
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "清算日期"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = date
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "报表代码"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Merge(1, 0)
	cell.Value = "IC002"

	//第三行
	row = sheet.AddRow()
	row.SetHeightCM(1.83)
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "客户代码"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "客户名称"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "渠道名称"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "清算角色"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "交易笔数"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "交易金额"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "商户手续费"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "商户应收额"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "讯联手续费"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "讯联应收金额"

	// 接下来是数据填充
	for _, elementMap := range data {
		for _, d := range elementMap {
			row = sheet.AddRow()
			row.SetHeightCM(1.00)
			cell = row.AddCell()
			cell.SetStyle(bodyStyle)
			cell.Value = d.insCode
			cell = row.AddCell()
			cell.SetStyle(bodyStyle)
			cell.Value = d.insName
			cell = row.AddCell()
			cell.SetStyle(bodyStyle)
			cell.Value = d.chcd
			cell = row.AddCell()
			cell.SetStyle(bodyStyle)
			cell.Value = d.role
			cell = row.AddCell()
			cell.SetStyle(bodyStyle)
			cell.Value = fmt.Sprintf("%d", d.orderNum)
			cell = row.AddCell()
			cell.SetStyle(bodyStyle)
			cell.Value = fmt.Sprintf("%0.2f", d.orderAct)
			cell = row.AddCell()
			cell.SetStyle(bodyStyle)
			cell.Value = fmt.Sprintf("%0.2f", d.orderNetFee)
			cell = row.AddCell()
			cell.SetStyle(bodyStyle)
			cell.Value = fmt.Sprintf("%0.2f", d.orderNetAct)
			cell = row.AddCell()
			cell.SetStyle(bodyStyle)
			cell.Value = fmt.Sprintf("%0.2f", d.orderAcqFee)
			cell = row.AddCell()
			cell.SetStyle(bodyStyle)
			cell.Value = fmt.Sprintf("%0.2f", d.orderAcqAct)
		}
	}

	return file
}
