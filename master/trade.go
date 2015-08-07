package master

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	// "github.com/CardInfoLink/quickpay/mongo"
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

	var merId = params.Get("mchntid")
	req := &model.QueryCondition{
		Mchntid:     merId,
		Busicd:      params.Get("busicd"),
		StartTime:   params.Get("startTime"),
		EndTime:     params.Get("endTime"),
		Size:        maxReportRec,
		IsForReport: true,
		Page:        1,
	}

	// 查询
	ret := core.TransQuery(req)

	// 生成报表
	before := time.Now()
	genReport(merId, file, ret.Rec)
	after := time.Now()
	log.Debugf("gen trans report spent %s", after.Sub(before))

	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`,
		fmt.Sprintf(`attachment; filename="%s";  filename*=utf-8''%s`, filename, filename))
	file.Write(w)
}

// genReport 生成报表
func genReport(merId string, file *xlsx.File, trans []model.Trans) {

	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	// 可能有多个sheet
	sheet = file.AddSheet("商户交易报表")

	// 先空3行，最后写入汇总数据
	for i := 0; i < 3; i++ {
		sheet.AddRow()
	}

	// 生成title
	row = sheet.AddRow()
	headRow := &struct {
		MerId        string
		OrderNum     string
		OrigOrderNum string
		TerminalId   string
		Inscd        string
		RespCode     string
		TransAmt     string
		TransStatus  string
		Busicd       string
		ChanCode     string
		TransTime    string
		ErrorDetail  string
	}{"商户号", "订单号", "原订单号", "终端号", "机构号", "应答码", "交易金额（元）", "交易状态", "交易类型", "渠道", "交易时间", "详情"}
	row.WriteStruct(headRow, -1)

	// 设置列宽
	sheet.SetColWidth(0, 9, 18)

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
		// 终端号
		cell = row.AddCell()
		cell.Value = v.Terminalid
		// 机构号
		cell = row.AddCell()
		cell.Value = v.Inscd
		// 应答码
		cell = row.AddCell()
		cell.Value = v.RespCode
		// 交易金额
		cell = row.AddCell()
		cell.SetFloat(float64(v.TransAmt) / 100)
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
		switch v.ChanCode {
		case "WXP":
			cell.Value = "微信"
		case "ALP":
			cell.Value = "支付宝"
		default:
			cell.Value = "未知"
		}
		// 交易时间
		cell = row.AddCell()
		cell.Value = v.CreateTime
		// 详情
		cell = row.AddCell()
		cell.Value = v.ErrorDetail
	}

	// TODO 利用商户数据，完善报表数据
	// mer, err := mongo.MerchantColl.Find(merId)
	// if err != nil {

	// }

	// 写入汇总数据
	rows := sheet.Rows
	row = rows[0]
	row.WriteStruct(&summary{
		"名称：", "讯联测试报表", "支付宝交易金额：", "500", "支付宝退款金额：", "200", "支付宝手续费：", "50", "支付宝清算金额：", "50",
	}, -1)
	row = rows[1]
	row.WriteStruct(&summary{
		"", "", "微信交易金额：", "500", "微信退款金额：", "200", "微信手续费：", "50", "微信清算金额：", "50",
	}, -1)
	row = rows[2]
	row.WriteStruct(&summary{
		"总计：", "", "交易总额：", "1000", "退款总额：", "400", "手续费总额：", "100", "清算总额：", "100",
	}, -1)
}

type summary struct {
	Cell0 string
	Cell1 string
	Cell2 string
	Cell3 string
	Cell4 string
	Cell5 string
	Cell6 string
	Cell7 string
	Cell8 string
	Cell9 string
}
