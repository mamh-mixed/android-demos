package settle

import (
	"bytes"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/CardInfoLink/log"
	"github.com/tealeg/xlsx"
	"time"
)

func upload(key string, excel *xlsx.File) error {
	// 空则忽略
	if excel == nil {
		return nil
	}

	bf := bytes.NewBuffer([]byte{})
	//写到buf里
	excel.Write(bf)
	// 上传到七牛
	err := qiniu.Put(key, int64(bf.Len()), bf)
	if err != nil {
		log.Errorf("upload report excel key=%s, err: %s, ", key, err)
	}
	return err
}

// genSpTransferSettleReportExcel 对账报表
func genReconciliatReportExcel(data reconciliationMap, date string) *xlsx.File {
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
		for _, chanMap := range elementMap {
			for _, d := range chanMap {
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
				cell.Value = fmt.Sprintf("%d", d.transNum)
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = fmt.Sprintf("%0.2f", float64(d.transAmt)/100)
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = fmt.Sprintf("%0.2f", float64(d.MerFee)/100)
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = fmt.Sprintf("%0.2f", float64(d.MerSettAmt)/100)
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = fmt.Sprintf("%0.2f", float64(d.AcqFee)/100)
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = fmt.Sprintf("%0.2f", float64(d.AcqSettAmt)/100)
			}
		}
	}

	return file
}

// genSpTransferReportExcel 划款报表
func genSpTransferReportExcel(data []reportData, date string) *xlsx.File {
	var file = xlsx.NewFile()
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	sheet, _ = file.AddSheet("商户清算划款表")

	// 第一行
	row = sheet.AddRow()
	row.SetHeightCM(0.91)
	cell = row.AddCell()
	cell.Merge(10, 0) // 11个单元格
	cell.SetValue("O2O商户划款报表(讯汇通)")
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
	cell.Value = "行号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "城市"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "银行名称"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "开户行名称"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "收款方姓名"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "收款方银行账号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "金额"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "备注"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "商户订单号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "渠道编号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "收支标识"

	// 接下来是数据填充
	for _, d := range data {
		row = sheet.AddRow()
		row.SetHeightCM(1.48)
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = d.m.Detail.BankId
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = d.m.Detail.City
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = d.m.Detail.BankName
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = d.m.Detail.OpenBankName
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = d.m.Detail.AcctName
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = d.m.Detail.AcctNum
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = fmt.Sprintf("%0.2f", float32(d.mg.TransAmt-d.mg.Fee)/100)
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = fmt.Sprintf("%s手续费%0.2f元", date, float32(d.mg.Fee)/100)
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = d.m.MerId
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = "05"
		cell = row.AddCell()
		cell.SetStyle(bodyStyle)
		cell.Value = "0"
	}

	return file
}

//生成C001
func genC001ReportExcel(data map[string]map[string][]model.TransSett, date string) *xlsx.File {
	var file = xlsx.NewFile()
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	sheet, _ = file.AddSheet("IC401")

	// 第一行
	row = sheet.AddRow()
	row.SetHeightCM(0.91)
	cell = row.AddCell()
	cell.Merge(9, 0) // 9个单元格
	cell.SetValue("云收银受理方可疑交易报表")
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
	cell.Value = "报表代码:"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "IC401"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "IC401"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	//cell.Merge(2, 0)
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "清算日期："
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = date
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "报表日期："
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	now := time.Now()
	year, mon, day := now.Date()
	cell.Value = fmt.Sprintf("%d-%d-%d", year, mon, day)

	//第三行
	row = sheet.AddRow()
	row.SetHeightCM(1.83)
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "渠道编号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "渠道名称"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "商户号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "虚拟商户号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "商户名称"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "订单号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "交易时间"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "交易类型"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "交易金额"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "对账标识"

	// 接下来是数据填充
	for _, elementMap := range data {
		for _, elementArray := range elementMap {
			for _, element := range elementArray {
				d := element.Trans
				row = sheet.AddRow()
				row.SetHeightCM(1.00)
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.ChanCode
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.ChanCode
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.MerId
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.ChanMerId
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.MerName
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.ChanOrderNum
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.PayTime
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				switch d.TransType {
				case 1:
					cell.Value = "支付"
				case 2:
					cell.Value = "退款"
				case 3:
					cell.Value = "预授权"
				case 4:
					cell.Value = "撤销"
				case 5:
					cell.Value = "关单"
				}
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = fmt.Sprintf("%0.2f", float64(d.TransAmt/100))
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = "渠道少清"
			}
		}
	}

	return file
}

//生成C002
func genC002ReportExcel(data map[string]map[string][]model.BlendElement, date string) *xlsx.File {
	var file = xlsx.NewFile()
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	sheet, _ = file.AddSheet("IC402")

	// 第一行
	row = sheet.AddRow()
	row.SetHeightCM(0.91)
	cell = row.AddCell()
	cell.Merge(8, 0) // 9个单元格
	cell.SetValue("云收银渠道方可疑交易报表")
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
	cell.Merge(1, 0)
	cell.Value = "报表代码:"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "IC402"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = ""
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = ""
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "清算日期："
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = date
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "报表日期："
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	now := time.Now()
	year, mon, day := now.Date()
	cell.Value = fmt.Sprintf("%d-%d-%d", year, mon, day)

	//第三行
	row = sheet.AddRow()
	row.SetHeightCM(1.83)
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "渠道编号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "渠道名称"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "渠道商户名称"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "渠道商户号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "订单号"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "交易时间"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "交易类型"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "交易金额"
	cell = row.AddCell()
	cell.SetStyle(bodyStyle)
	cell.Value = "对账标识"

	// 接下来是数据填充
	for _, elementMap := range data {
		for _, elementArray := range elementMap {
			for _, d := range elementArray {
				row = sheet.AddRow()
				row.SetHeightCM(1.00)
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.Chcd
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.ChcdName
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.MerName
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.ChanMerID
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.OrderID
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.OrderTime
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.OrderType
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = d.OrderAct
				cell = row.AddCell()
				cell.SetStyle(bodyStyle)
				cell.Value = "渠道多清"
				// _, ret := chanOrderMap[d.OrderID]
				// if !ret { //数据库没有,需要添加
				// 	var t model.TransSett
				// 	t.BlendType = model.SettChanRemain
				// 	t.Trans.ChanMerId = d.ChanMerID
				// 	t.Trans.ChanOrderNum = d.OrderID
				// 	switch d.OrderType {
				// 	case "在线支付", "交易成功":
				// 		t.Trans.TransType = model.PayTrans
				// 	case "交易退款", "退款":
				// 		t.Trans.TransType = model.RefundTrans
				// 	}
				// 	acct, _ := strconv.ParseFloat(d.OrderAct, 64)
				// 	t.Trans.TransAmt = int64(math.Floor(acct*100 + 0.5))
				// 	t.Trans.ChanCode = d.Chcd
				// 	t.Trans.MerId = d.MerID
				// 	t.Trans.MerName = d.MerName
				// 	t.Trans.PayTime = d.OrderTime
				// 	mongo.SpTransSettleColl.Add(&t)
				// }
			}
		}
	}

	return file
}
