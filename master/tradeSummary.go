package master

import (
	// "encoding/json"
	"fmt"
	"net/http"
	// "time"

	// "github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/query"
	// "github.com/CardInfoLink/quickpay/mongo"
	// "github.com/omigo/log"
	"github.com/tealeg/xlsx"
)

// tradeQueryStat 交易查询统计信息
func tradeQueryStats(q *model.QueryCondition) (result *model.ResultBody) {

	// 调用core方法统计
	qr := query.TransStatistics(q)

	// 分页信息
	pagination := &model.Pagination{
		Page:  qr.Page,
		Total: qr.Total,
		Size:  qr.Size,
		Count: qr.Size,
		Data:  qr.Rec,
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}
	return result
}

// tradeQueryStatReport 交易汇总报表
func tradeQueryStatsReport(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	filename := params.Get("filename")

	var file = xlsx.NewFile()

	q := &model.QueryCondition{
		MerId:     params.Get("merId"),
		AgentCode: params.Get("agentCode"),
		MerName:   params.Get("merName"),
		StartTime: params.Get("startTime"),
		EndTime:   params.Get("endTime"),
		Page:      1,
		Size:      maxReportRec,
	}

	qr := query.TransStatistics(q)

	if summarys, ok := qr.Rec.(model.Summary); ok {
		genQueryStatReport(file, summarys, q)
	}

	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`, fmt.Sprintf(`attachment; filename="%s"`, filename))
	file.Write(w)
}

// TODO: 优化
func genQueryStatReport(file *xlsx.File, result model.Summary, cond *model.QueryCondition) {

	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	sheet = file.AddSheet("商户交易报表汇总")

	// 表头样式
	genHead(sheet, row, cell, cond)

	// 填充数据
	// 详细数据
	for _, d := range result.Data {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = d.MerId
		cell.Merge(1, 0)
		row.AddCell()
		cell = row.AddCell()
		cell.Value = d.MerName
		cell.Merge(1, 0)
		row.AddCell()
		cell = row.AddCell()
		cell.SetInt(d.TotalTransNum)
		cell = row.AddCell()
		cell.SetFloat(float64(d.TotalTransAmt))
		cell = row.AddCell()
		cell.SetFloat(float64(d.TotalFee))
		cell = row.AddCell()
		cell.SetInt(d.Alp.TransNum)
		cell = row.AddCell()
		cell.SetFloat(float64(d.Alp.TransAmt))
		cell = row.AddCell()
		cell.SetFloat(float64(d.Alp.Fee))
		cell = row.AddCell()
		cell.SetInt(d.Wxp.TransNum)
		cell = row.AddCell()
		cell.SetFloat(float64(d.Wxp.TransAmt))
		cell = row.AddCell()
		cell.SetFloat(float64(d.Wxp.Fee))
		cell = row.AddCell()
		cell.Value = d.AgentName
		cell.Merge(1, 0)
	}

	// 最后填写汇总
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "总计："
	cell.Merge(3, 0)
	for i := 0; i < 3; i++ {
		row.AddCell()
	}
	cell = row.AddCell()
	cell.SetInt(result.TotalTransNum)
	cell = row.AddCell()
	cell.SetFloat(float64(result.TotalTransAmt))
	cell = row.AddCell()
	cell.SetFloat(float64(result.TotalFee))
	cell = row.AddCell()
	cell.SetInt(result.Alp.TransNum)
	cell = row.AddCell()
	cell.SetFloat(float64(result.Alp.TransAmt))
	cell = row.AddCell()
	cell.SetFloat(float64(result.Alp.Fee))
	cell = row.AddCell()
	cell.SetInt(result.Wxp.TransNum)
	cell = row.AddCell()
	cell.SetFloat(float64(result.Wxp.TransAmt))
	cell = row.AddCell()
	cell.SetFloat(float64(result.Wxp.Fee))
	row.AddCell().Merge(1, 0)
}

func genHead(sheet *xlsx.Sheet, row *xlsx.Row, cell *xlsx.Cell, cond *model.QueryCondition) {
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "开始日期："
	cell = row.AddCell()
	cell.Value = cond.StartTime
	cell = row.AddCell()
	cell.Value = "结束日期："
	cell = row.AddCell()
	cell.Value = cond.EndTime
	cell = row.AddCell()
	cell.Value = "注：手续费为每笔单笔计算后四舍五入精确到分，跟总额计算手续费略有误差。因本表仅统计了讯联数据系统的数据，数据仅供参考"
	cell.Merge(10, 0)
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "商户号"
	cell.Merge(1, 1)
	row.AddCell()
	cell = row.AddCell()
	cell.Value = "商户名称"
	cell.Merge(1, 1)
	row.AddCell()
	cell = row.AddCell()
	cell.Value = "汇总"
	cell.Merge(2, 0)
	row.AddCell()
	row.AddCell()
	cell = row.AddCell()
	cell.Value = "支付宝"
	cell.Merge(2, 0)
	row.AddCell()
	row.AddCell()
	cell = row.AddCell()
	cell.Value = "微信"
	cell.Merge(2, 0)
	row.AddCell()
	row.AddCell()
	cell = row.AddCell()
	cell.Value = "代理名称"
	cell.Merge(1, 1)
	row.AddCell()
	row = sheet.AddRow()
	for i := 0; i < 4; i++ {
		row.AddCell()
	}
	cell = row.AddCell()
	cell.Value = "总笔数"
	cell = row.AddCell()
	cell.Value = "总金额"
	cell = row.AddCell()
	cell.Value = "手续费"
	cell = row.AddCell()
	cell.Value = "笔数"
	cell = row.AddCell()
	cell.Value = "金额"
	cell = row.AddCell()
	cell.Value = "手续费"
	cell = row.AddCell()
	cell.Value = "笔数"
	cell = row.AddCell()
	cell.Value = "金额"
	cell = row.AddCell()
	cell.Value = "手续费"
}
