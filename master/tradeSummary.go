package master

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/currency"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/query"
	"github.com/tealeg/xlsx"
	"net/http"
	"time"
)

const intFormat = "#,##0"

var headStyle = &xlsx.Style{
	Border: xlsx.Border{
		Left:        "thin",
		LeftColor:   "FF999999",
		Right:       "thin",
		RightColor:  "FF999999",
		Top:         "thin",
		TopColor:    "FF999999",
		Bottom:      "thin",
		BottomColor: "FF999999",
	},
	Fill: xlsx.Fill{
		PatternType: "solid",
		FgColor:     "FF00BCD4",
	},
	Font: xlsx.Font{
		Size:    10,
		Name:    "微软雅黑",
		Family:  2,
		Charset: 134,
	},
	Alignment: xlsx.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	},
}
var bodyStyle = &xlsx.Style{
	Font: xlsx.Font{
		Size:    10,
		Name:    "微软雅黑",
		Family:  2,
		Charset: 134,
	},
}

// tradeQueryStat 交易查询统计信息
func tradeQueryStats(q *model.QueryCondition) (result *model.ResultBody) {

	// 调用core方法统计
	s, total := query.TransStatistics(q)

	// 分页信息
	pagination := &model.Pagination{
		Page:  q.Page,
		Total: total,
		Size:  q.Size,
		Count: len(s.Data),
		Data:  s,
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}
	return result
}

// statTradeReport 交易统计报表
func statTradeReport(w http.ResponseWriter, q *model.QueryCondition) {
	// 语言模板
	reportLocale := GetLocale(q.Locale)

	// 调用core方法统计
	s, _ := query.TransStatistics(q)

	// 设置为东八区
	q.UtcOffset = 60 * 60 * 8

	// TODO 币种处理，导出
	genStatReport(s, q, reportLocale, reportLocale.Currency).Write(w)
}

// tradeSettReport 交易清算汇总报表
func tradeSettReport(w http.ResponseWriter, q *model.QueryCondition) {
	// 语言模板
	rl := GetLocale(q.Locale)

	var filename string
	reportName := rl.ReportName.SettleSummary
	switch q.UserType {
	case model.UserTypeCIL, model.UserTypeGenAdmin:
		filename = reportName
	case model.UserTypeAgent:
		filename = rl.Role.Agent + reportName
	case model.UserTypeMerchant:
		filename = rl.Role.Group + reportName
	case model.UserTypeCompany:
		filename = rl.Role.Company + reportName
	case model.UserTypeShop:
		filename = rl.Role.Mer + reportName
	}
	filename += ".xlsx"

	// 调用core方法统计
	s := query.TransSettStatistics(q)

	// 设置为东八区
	q.UtcOffset = 60 * 60 * 8

	// 设置content-type
	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`, fmt.Sprintf(`attachment; filename="%s"`, filename))

	// TODO 暂时设置为日币
	// 导出
	genStatReport(s, q, rl, "JPY").Write(w)
}

// TODO: 优化
func genStatReport(result model.Summary, cond *model.QueryCondition, locale *LocaleTemplate, curr string) (file *xlsx.File) {

	// 语言模板
	reportLocale := GetLocale(cond.Locale).StatReport

	// 币种转换
	cur := currency.Get(curr)

	file = xlsx.NewFile()
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	sheet, _ = file.AddSheet(reportLocale.Title)

	// 表头样式
	genHead(sheet, row, cell, cond)

	// 金额显示格式
	var floatFormat = "#,##0"
	for i := 0; i < cur.Precision; i++ {
		if i == 0 {
			floatFormat += "."
		}
		floatFormat += "0"
	}

	// 填充数据
	for _, d := range result.Data {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = d.MerId
		cell.SetStyle(bodyStyle)
		cell.Merge(1, 0)
		row.AddCell()
		cell = row.AddCell()
		cell.Value = d.MerName
		cell.SetStyle(bodyStyle)
		cell.Merge(1, 0)
		row.AddCell()
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(d.TotalTransNum), intFormat)
		cell.SetStyle(bodyStyle)
		cell = row.AddCell()
		cell.SetFloatWithFormat(cur.F64(d.TotalTransAmt-d.TotalFee), floatFormat)
		cell.SetStyle(bodyStyle)
		cell = row.AddCell()
		cell.SetFloatWithFormat(cur.F64(d.TotalFee), floatFormat)
		cell.SetStyle(bodyStyle)
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(d.Alp.TransNum), intFormat)
		cell.SetStyle(bodyStyle)
		cell = row.AddCell()
		cell.SetFloatWithFormat(cur.F64(d.Alp.TransAmt-d.Alp.Fee), floatFormat)
		cell.SetStyle(bodyStyle)
		cell = row.AddCell()
		cell.SetFloatWithFormat(cur.F64(d.Alp.Fee), floatFormat)
		cell.SetStyle(bodyStyle)
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(d.Wxp.TransNum), intFormat)
		cell.SetStyle(bodyStyle)
		cell = row.AddCell()
		cell.SetFloatWithFormat(cur.F64(d.Wxp.TransAmt-d.Wxp.Fee), floatFormat)
		cell.SetStyle(bodyStyle)
		cell = row.AddCell()
		cell.SetFloatWithFormat(cur.F64(d.Wxp.Fee), floatFormat)
		cell.SetStyle(bodyStyle)
		cell = row.AddCell()
		cell.Value = d.AgentName
		cell.SetStyle(bodyStyle)
		cell.Merge(1, 0)
	}

	// 最后填写汇总
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = reportLocale.Total
	cell.SetStyle(bodyStyle)
	cell.Merge(3, 0)
	for i := 0; i < 3; i++ {
		row.AddCell()
	}
	cell = row.AddCell()
	cell.SetInt(result.TotalTransNum)
	cell.SetStyle(bodyStyle)
	cell = row.AddCell()
	cell.SetFloatWithFormat(cur.F64(result.TotalTransAmt-result.TotalFee), floatFormat)
	cell.SetStyle(bodyStyle)
	cell = row.AddCell()
	cell.SetFloatWithFormat(cur.F64(result.TotalFee), floatFormat)
	cell.SetStyle(bodyStyle)
	cell = row.AddCell()
	cell.SetFloatWithFormat(float64(result.Alp.TransNum), intFormat)
	cell.SetStyle(bodyStyle)
	cell = row.AddCell()
	cell.SetFloatWithFormat(cur.F64(result.Alp.TransAmt-result.Alp.Fee), floatFormat)
	cell.SetStyle(bodyStyle)
	cell = row.AddCell()
	cell.SetFloatWithFormat(cur.F64(result.Alp.Fee), floatFormat)
	cell.SetStyle(bodyStyle)
	cell = row.AddCell()
	cell.SetFloatWithFormat(float64(result.Wxp.TransNum), intFormat)
	cell.SetStyle(bodyStyle)
	cell = row.AddCell()
	cell.SetFloatWithFormat(cur.F64(result.Wxp.TransAmt-result.Wxp.Fee), floatFormat)
	cell.SetStyle(bodyStyle)
	cell = row.AddCell()
	cell.SetFloatWithFormat(cur.F64(result.Wxp.Fee), floatFormat)
	cell.SetStyle(bodyStyle)
	row.AddCell().Merge(1, 0)
	return file
}

func genHead(sheet *xlsx.Sheet, row *xlsx.Row, cell *xlsx.Cell, cond *model.QueryCondition) {

	// 语言模板
	reportLocale := GetLocale(cond.Locale).StatReport

	// 时区
	z := &Zone{cond.UtcOffset, time.Local}

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = reportLocale.StartDate
	cell = row.AddCell()
	cell.Value = z.GetTime(cond.StartTime)
	cell.SetStyle(bodyStyle)
	cell.Merge(1, 0)
	row.AddCell()

	cell = row.AddCell()
	cell.Value = reportLocale.EndDate
	cell = row.AddCell()
	cell.Value = z.GetTime(cond.EndTime)
	cell.SetStyle(bodyStyle)
	cell.Merge(1, 0)
	row.AddCell()

	cell = row.AddCell()
	cell.Value = reportLocale.Remark
	cell.SetStyle(bodyStyle)
	cell.Merge(8, 0)

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = reportLocale.MerId
	cell.SetStyle(headStyle)
	cell.Merge(1, 1)
	row.AddCell()
	cell = row.AddCell()
	cell.Value = reportLocale.MerName
	cell.SetStyle(headStyle)
	cell.Merge(1, 1)
	row.AddCell()
	cell = row.AddCell()
	cell.Value = reportLocale.Summary
	cell.SetStyle(headStyle)
	cell.Merge(2, 0)
	row.AddCell()
	row.AddCell()
	cell = row.AddCell()
	cell.Value = reportLocale.ALP
	cell.SetStyle(headStyle)
	cell.Merge(2, 0)
	row.AddCell()
	row.AddCell()
	cell = row.AddCell()
	cell.Value = reportLocale.WXP
	cell.SetStyle(headStyle)
	cell.Merge(2, 0)
	row.AddCell()
	row.AddCell()
	cell = row.AddCell()
	cell.Value = reportLocale.AgentName
	cell.SetStyle(headStyle)
	cell.Merge(1, 1)
	row.AddCell()

	row = sheet.AddRow()
	for i := 0; i < 4; i++ {
		row.AddCell()
	}
	cell = row.AddCell()
	cell.Value = reportLocale.TotalCount
	cell.SetStyle(headStyle)
	cell = row.AddCell()
	cell.Value = reportLocale.TotalAmt
	cell.SetStyle(headStyle)
	cell = row.AddCell()
	cell.Value = reportLocale.Fee
	cell.SetStyle(headStyle)
	cell = row.AddCell()
	cell.Value = reportLocale.Count
	cell.SetStyle(headStyle)
	cell = row.AddCell()
	cell.Value = reportLocale.Amt
	cell.SetStyle(headStyle)
	cell = row.AddCell()
	cell.Value = reportLocale.Fee
	cell.SetStyle(headStyle)
	cell = row.AddCell()
	cell.Value = reportLocale.Count
	cell.SetStyle(headStyle)
	cell = row.AddCell()
	cell.Value = reportLocale.Amt
	cell.SetStyle(headStyle)
	cell = row.AddCell()
	cell.Value = reportLocale.Fee
	cell.SetStyle(headStyle)
}
