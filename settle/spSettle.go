package settle

import (
	"bytes"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
	"strings"
)

const filePrefix = "sett/report/%s/" // 文件名：sett/report/20151012/IC202_99911888_20151012.xlsx

// doScanpaySettReport 扫码每天出清算报表
func doScanpaySettReport(settDate string) error {

	data, err := mongo.SpTransColl.GroupBySettRole(settDate)
	if err != nil {
		log.Errorf("fail to find trans group by settRole: %s", err)
		return err
	}

	// 报表日期显示格式
	sd := strings.Replace(settDate, "-", "", -1)
	filename := filePrefix + "IC202_%s_%s.xlsx"

	// 遍历数据
	for _, sr := range data {
		// 每一行就是一个报表
		sr.SettDate = sd
		excel := genSpSettReportExcel(sr)

		var buf []byte
		bf := bytes.NewBuffer(buf)
		// 写到buf里
		excel.Write(bf)

		// 上传到七牛
		err = qiniu.Upload(fmt.Sprintf(filename, sd, sr.SettRole, sd), int64(len(bf.Bytes())), bf)
		if err != nil {
			log.Errorf("upload settReport excel err: %s", err)
		}
	}

	return nil
}

func genSpSettReportExcel(sr model.SettRoleGroup) *xlsx.File {
	var file = xlsx.NewFile()
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	sheet, _ = file.AddSheet("商户清算划款表")

	// xlsx SDK 暂不支持Alignment
	// 第一行
	row = sheet.AddRow()
	row.SetHeightCM(0.91)
	cell = row.AddCell()
	cell.Merge(10, 0) // 11个单元格
	cell.SetValue("O2O商户划款报表(讯汇通)")
	style := xlsx.NewStyle()

	style.Alignment = xlsx.Alignment{Horizontal: "Center", Vertical: "Center"}
	style.Font = *xlsx.NewFont(20, "宋体")
	style.Border = *xlsx.NewBorder("Left", "Right", "Top", "Bottom")
	style.ApplyAlignment = true
	style.ApplyFont = true
	style.ApplyBorder = true
	cell.SetStyle(style)

	// 第二行
	twoStyle := xlsx.NewStyle()
	twoStyle.Font = *xlsx.NewFont(10, "Times New Roman")
	twoStyle.ApplyFont = true
	row = sheet.AddRow()
	row.SetHeightCM(1.83)
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "行号"
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "城市"
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "银行名称"
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "开户行名称"
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "收款方姓名"
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "收款方银行账号"
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "金额"
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "备注"
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "商户订单号"
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "渠道编号"
	cell = row.AddCell()
	cell.SetStyle(twoStyle)
	cell.Value = "收支标识"

	// 接下来是数据填充
	for _, mg := range sr.Mers {
		m, err := mongo.MerchantColl.Find(mg.MerId)
		if err != nil {
			log.Errorf("find merchant error: %s, merId=%s", err, mg.MerId)
			m = &model.Merchant{MerId: mg.MerId}
		}

		row = sheet.AddRow()
		row.SetHeightCM(1.48)
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = m.Detail.BankId
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = m.Detail.City
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = m.Detail.BankName
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = m.Detail.OpenBankName
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = m.Detail.AcctName
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = m.Detail.AcctNum
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = fmt.Sprintf("%0.2f", float32(mg.TransAmt-mg.RefundAmt)/100)
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = fmt.Sprintf("%s手续费%0.2f元", sr.SettDate, float32(mg.Fee)/100)
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = m.MerId
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = "05"
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = "0"
	}

	return file
}
