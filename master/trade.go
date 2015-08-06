package master

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
	"net/http"
	"time"
)

var maxReportRec = 10000

// tradeReport 处理查找所有商户的请求
func tradeReport(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	filename := params.Get("filename")

	var file = xlsx.NewFile()

	req := &model.QueryCondition{
		Mchntid:   params.Get("mchntid"),
		Busicd:    params.Get("busicd"),
		StartTime: params.Get("startTime"),
		EndTime:   params.Get("endTime"),
		Size:      maxReportRec,
		Page:      1,
	}

	// 查询
	ret := core.TransQuery(req)

	// 生成报表
	before := time.Now()
	genReport(file, ret.Rec)
	after := time.Now()
	log.Debugf("gen trans report spent %s", after.Sub(before))

	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`,
		fmt.Sprintf(`attachment; filename="%s";  filename*=utf-8''%s`, filename, filename))
	file.Write(w)
}

// genReport 生成报表
func genReport(file *xlsx.File, trans []model.Trans) {

	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	// 可能有多个sheet
	sheet = file.AddSheet("商户交易报表")

	// 生成title
	// TODO 字段待补充
	row = sheet.AddRow()
	headRow := &struct {
		MerId        string
		OrderNum     string
		OrigOrderNum string
		Inscd        string
		RespCode     string
		TransAmt     string
		TransStatus  string
		Busicd       string
		ChanCode     string
		TransTime    string
		ErrorDetail  string
	}{"商户号", "订单号", "原订单号", "机构号", "应答码", "交易金额（分）", "交易状态", "交易类型", "渠道", "交易时间", "详情"}
	row.WriteStruct(headRow, -1)

	// 设置列宽
	sheet.SetColWidth(0, 9, 18)
	// col := sheet.Col(0)
	// newStyle := xlsx.NewStyle()
	// newStyle.Alignment = xlsx.Alignment{"align", "align"}
	// newStyle.ApplyAlignment = true
	// col.SetStyle(newStyle)

	// 生成数据
	for _, v := range trans {
		row = sheet.AddRow()
		// 商户号
		cell = row.AddCell()
		cell.Value = v.MerId
		// 订单号
		cell = row.AddCell()
		cell.Value = v.OrderNum
		// 原订单号
		cell = row.AddCell()
		cell.Value = v.OrigOrderNum
		// 机构号
		cell = row.AddCell()
		cell.Value = v.Inscd
		// 应答码
		cell = row.AddCell()
		cell.Value = v.RespCode
		// 交易金额
		cell = row.AddCell()
		cell.SetInt64(v.TransAmt)
		// 交易状态
		cell = row.AddCell()
		switch v.TransStatus {
		case model.TransSuccess:
			cell.Value = "交易成功"
		case model.TransFail:
			cell.Value = "交易失败"
		case model.TransHandling:
			cell.Value = "交易处理中"
		case model.TransClosed:
			cell.Value = "交易已关闭"
		default:
			cell.Value = "未知"
		}
		// 交易类型
		cell = row.AddCell()
		switch v.Busicd {
		case model.Purc:
			cell.Value = "下单并支付"
		case model.Paut:
			cell.Value = "预下单"
		case model.Refd:
			cell.Value = "退款"
		case model.Void:
			cell.Value = "撤销"
		case model.Canc:
			cell.Value = "取消"
		case model.Qyfk:
			cell.Value = "企业付款"
		default:
			cell.Value = "未知"
		}
		// 渠道
		cell = row.AddCell()
		cell.Value = v.ChanCode
		// 交易时间
		cell = row.AddCell()
		cell.Value = v.CreateTime
		// 详情
		cell = row.AddCell()
		cell.Value = v.ErrorDetail
		// TODO 字段待补充 ...
	}
}
