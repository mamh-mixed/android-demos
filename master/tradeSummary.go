package master

import (
	// "encoding/json"
	"fmt"
	"net/http"
	// "time"

	// "github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	// "github.com/CardInfoLink/quickpay/mongo"
	// "github.com/omigo/log"
	"github.com/tealeg/xlsx"
)

// tradeQueryStat 交易查询统计信息
func tradeQueryStats(q *model.QueryCondition) (result *model.ResultBody) {

	// 调用core方法统计
	qr := core.TransStatistics(q)

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

	qr := core.TransStatistics(q)

	if summarys, ok := qr.Rec.(model.Summary); ok {
		genQueryStatReport(file, summarys)
	}

	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`,
		fmt.Sprintf(`attachment; filename="%s";  filename*=utf-8''%s`, filename, filename))
	file.Write(w)
}

// TODO: 优化
func genQueryStatReport(file *xlsx.File, result model.Summary) {

	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	sheet = file.AddSheet("商户交易报表汇总")

	// 表头样式
	genHead(sheet, row, cell)

	// 填充数据
	// 先填写汇总
	row = sheet.AddRow()
	for i := 0; i < 4; i++ {
		row.AddCell()
	}
	cell = row.AddCell()
	cell.SetInt(result.TotalTransNum)
	cell = row.AddCell()
	cell.SetFloat(float64(result.TotalTransAmt))
	cell = row.AddCell()
	cell.SetInt(result.Alp.TransNum)
	cell = row.AddCell()
	cell.SetFloat(float64(result.Alp.TransAmt))
	cell = row.AddCell()
	cell.SetInt(result.Wxp.TransNum)
	cell = row.AddCell()
	cell.SetFloat(float64(result.Wxp.TransAmt))

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
		cell.SetInt(d.Alp.TransNum)
		cell = row.AddCell()
		cell.SetFloat(float64(d.Alp.TransAmt))
		cell = row.AddCell()
		cell.SetInt(d.Wxp.TransNum)
		cell = row.AddCell()
		cell.SetFloat(float64(d.Wxp.TransAmt))
		cell = row.AddCell()
		cell.Value = d.AgentName
		cell.Merge(1, 0)
	}
}

func genHead(sheet *xlsx.Sheet, row *xlsx.Row, cell *xlsx.Cell) {
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
	cell.Merge(1, 0)
	row.AddCell()
	cell = row.AddCell()
	cell.Value = "支付宝"
	cell.Merge(1, 0)
	row.AddCell()
	cell = row.AddCell()
	cell.Value = "微信"
	cell.Merge(1, 0)
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
	cell.Value = "笔数"
	cell = row.AddCell()
	cell.Value = "金额"
	cell = row.AddCell()
	cell.Value = "笔数"
	cell = row.AddCell()
	cell.Value = "金额"
}
