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

//  reportType
const (
	TransferReport       = 1 // 划款报表
	ReconciliationReport = 2 // 对账报表
	InsFlowReport        = 3 // 机构流水报表
	ChanMerReport        = 4 // 渠道商户报表
	// 分润报表
)

// DoScanpaySettReport 清算
func DoScanpaySettReport(settDate string) error {
	// 对账，入库
	tss, err := Reconciliation(settDate)
	if err != nil {
		return nil
	}

	// 出报表
	err = genReport(tss)
	if err != nil {
		return nil
	}

	// 发邮件?

	return nil
}

// Reconciliation 对账
func Reconciliation(settDate string) ([]model.TransSett, error) {

	var reconcilated []*model.Trans

	// TODO: 与渠道对账，然后将交易存进对账交易表里

	ts, _, err := mongo.SpTransColl.Find(&model.QueryCondition{
		StartTime:    settDate + " 00:00:00",
		EndTime:      settDate + " 23:59:59",
		TransStatus:  []string{model.TransSuccess},
		RefundStatus: model.TransRefunded,
		IsForReport:  true,
	})

	if err != nil {
		return nil, err
	}

	// 先已我们系统为准，默认都对上
	reconcilated = append(reconcilated, ts...)

	// 计算费率等，入清算表
	var tss []model.TransSett
	for _, t := range reconcilated {
		ts := model.TransSett{
			Trans: *t,
		}
		// TODO
		tss = append(tss, ts)
	}

	return tss, nil
}

func genReport([]model.TransSett) error {
	return nil
}

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
	for _, sg := range data {

		key := fmt.Sprintf(filename, sd, sg.SettRole, sd)

		// 查询该角色是否已出过报表
		rs, err := mongo.RoleSettCol.FindOne(sg.SettRole, settDate)
		if err != nil {
			rs = &model.RoleSett{SettRole: sg.SettRole, SettDate: settDate, ReportName: key}
		}

		rpData := settDataHandle(sg, rs)
		// 有数据才生成报表
		if len(rpData) != 0 {
			// 每一行就是一个报表
			excel := genSpTransferReportExcel(rpData, sd)

			var buf []byte
			bf := bytes.NewBuffer(buf)
			// 写到buf里
			excel.Write(bf)

			// 上传到七牛
			err = qiniu.Put(key, int64(bf.Len()), bf)
			if err != nil {
				log.Errorf("upload settReport excel err: %s", err)
				continue
			}
			err = mongo.RoleSettCol.Upsert(rs)
			if err != nil {
				log.Errorf("roleSett upsert error: %s", err)
			}
		}
	}

	return nil
}

type reportData struct {
	mg model.MerGroup
	m  model.Merchant
}

// settDataHandle 清算数据处理
func settDataHandle(sg model.SettRoleGroup, rs *model.RoleSett) []reportData {

	var rds []reportData
	if rs == nil {
		return rds
	}

	for _, mg := range sg.MerGroups {
		m, err := mongo.MerchantColl.Find(mg.MerId)
		if err != nil {
			// cmMap[mg.MerId] = 0 // 标识不成功
			// continue
			m = &model.Merchant{MerId: mg.MerId} // 兼容老系统数据，可能商户没同步到新系统
		}
		// 补充开户银行和支行
		if m.Detail.OpenBankName == "" {
			m.Detail.OpenBankName = m.Detail.BankName
		}
		if m.Detail.BankName == "" {
			m.Detail.BankName = m.Detail.OpenBankName
		}
		rds = append(rds, reportData{mg: mg, m: *m})
	}

	return rds
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
