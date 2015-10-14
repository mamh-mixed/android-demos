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
			excel := genSpSettReportExcel(rpData, sd)

			var buf []byte
			bf := bytes.NewBuffer(buf)
			// 写到buf里
			excel.Write(bf)

			// 上传到七牛
			err = qiniu.Put(key, int64(len(bf.Bytes())), bf)
			if err != nil {
				log.Errorf("upload settReport excel err: %s", err)
				continue
			}
		}

		err = mongo.RoleSettCol.Upsert(rs)
		if err != nil {
			log.Errorf("roleSett upsert error: %s", err)
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

	var cmMap = make(map[string]int)
	for _, cm := range rs.ContainMers {
		cmMap[cm.MerId] = cm.Status
	}

	for _, mg := range sg.MerGroups {
		if status, ok := cmMap[mg.MerId]; ok {
			if status == 1 {
				continue
			}
			// 存在，但状态不成功
			delete(cmMap, mg.MerId)
		}

		m, err := mongo.MerchantColl.Find(mg.MerId)
		if err != nil {
			cmMap[mg.MerId] = 0 // 标识不成功
			continue
		}

		if m.Detail.BankId == "" || m.Detail.AcctNum == "" || m.Detail.AcctName == "" ||
			m.Detail.OpenBankName == "" || m.Detail.BankName == "" || m.Detail.City == "" {
			log.Warnf("settinfo not found , gen report skip , merId=%s", mg.MerId)
			// 清算信息缺一不可
			cmMap[mg.MerId] = 0
			continue
		}

		cmMap[mg.MerId] = 1 // 清算成功
		rds = append(rds, reportData{mg: mg, m: *m})
	}

	var cms []model.MerSettStatus
	for k, v := range cmMap {
		cms = append(cms, model.MerSettStatus{MerId: k, Status: v})
	}
	rs.ContainMers = cms

	return rds
}

func genSpSettReportExcel(data []reportData, date string) *xlsx.File {
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
	for _, d := range data {
		row = sheet.AddRow()
		row.SetHeightCM(1.48)
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = d.m.Detail.BankId
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = d.m.Detail.City
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = d.m.Detail.BankName
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = d.m.Detail.OpenBankName
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = d.m.Detail.AcctName
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = d.m.Detail.AcctNum
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = fmt.Sprintf("%0.2f", float32(d.mg.TransAmt-d.mg.RefundAmt-d.mg.Fee)/100)
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = fmt.Sprintf("%s手续费%0.2f元", date, float32(d.mg.Fee)/100)
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = d.m.MerId
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = "05"
		cell = row.AddCell()
		cell.SetStyle(twoStyle)
		cell.Value = "0"
	}

	return file
}
