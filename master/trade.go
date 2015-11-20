package master

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/query"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
)

var maxReportRec = 10000

func getTradeMsg(q *model.QueryCondition, msgType int) (ret *model.ResultBody) {
	ls, total, err := query.GetSpTransLogs(q, msgType)
	if err != nil {
		return model.NewResultBody(1, "查询数据库失败")
	}

	paging := model.Pagination{
		Page:  q.Page,
		Total: total,
		Size:  q.Size,
		Data:  ls,
	}

	return &model.ResultBody{
		Status: 0,
		Data:   paging,
	}

}

// tradeSettleReportQuery 清算报表查询
func tradeSettleReportQuery(role, date string, size, page int) (result *model.ResultBody) {
	log.Debugf("role=%s; date=%s", role, date)

	if page <= 0 {
		return model.NewResultBody(400, "page 参数错误")
	}

	if size == 0 {
		size = 10
	}

	results, total, err := mongo.RoleSettCol.PaginationFind(role, date, size, page)
	if err != nil {
		log.Errorf("分页查询出错%s", err)
		return model.NewResultBody(1, "查询失败")
	}

	// 分页信息
	pagination := &model.Pagination{
		Page:  page,
		Total: total,
		Size:  size,
		Count: len(results),
		Data:  results,
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    pagination,
	}

	return result
}

// tradeQuery 交易查询
func tradeQuery(q *model.QueryCondition) (ret *model.ResultBody) {

	switch {
	case q.Col == "bp":
		return query.BpTransQuery(q)
	case q.Col == "coupon":
		return query.CouponTransQuery(q)
	default:
		return query.SpTransQuery(q)
	}
}

// tradeQuery 交易查询
func tradeFindOne(q *model.QueryCondition) (ret *model.ResultBody) {
	return query.SpTransFindOne(q)
}

// tradeReport 处理查找所有商户的请求
func tradeReport(w http.ResponseWriter, cond *model.QueryCondition, filename string) {
	var file = xlsx.NewFile()

	// 查询
	ret := query.SpTransQuery(cond)

	// 类型转换
	if pagination, ok := ret.Data.(*model.Pagination); ok {
		if trans, ok := pagination.Data.([]*model.Trans); ok {
			// 生成报表
			before := time.Now()
			genReport(cond.MerId, file, trans, GetLocale(cond.Locale))
			after := time.Now()
			log.Debugf("gen trans report spent %s", after.Sub(before))
		}
	}

	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`, fmt.Sprintf(`attachment; filename="%s"`, filename))
	file.Write(w)
}

// TODO:genReport 生成报表
func genReport(merId string, file *xlsx.File, trans []*model.Trans, locale *LocaleTemplate) {
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	// 语言
	m := locale.TransReport
	lALP, lWXP := locale.ChanCode.ALP, locale.ChanCode.WXP

	// 可能有多个sheet
	sheet, _ = file.AddSheet(m.SheetName)

	// 先空3行，最后写入汇总数据
	for i := 0; i < 3; i++ {
		sheet.AddRow()
	}

	// 生成title
	row = sheet.AddRow()
	headRow := &struct {
		MerId        string
		MerName      string
		OrderNum     string
		TransAmt     string
		ChanCode     string
		TransTime    string
		PayTime      string
		TransStatus  string
		AgentCode    string
		TerminalId   string
		Busicd       string
		OrigOrderNum string
	}{m.MerId, m.MerName, m.OrderNum, m.TransAmt, m.ChanCode, m.TransTime, m.PayTime, m.TransStatus, m.AgentCode, m.TerminalId, m.Busicd, m.OrigOrderNum}
	row.WriteStruct(headRow, -1)

	// 设置列宽
	sheet.SetColWidth(0, 9, 18)

	// 支付宝交易金额、退款金额
	var alpTransAmt, alpRefundAmt, alpFee int64 = 0, 0, 0
	// 微信交易金额、退款金额
	var wxpTransAmt, wxpRefundAmt, wxpFee int64 = 0, 0, 0
	// 总交易金额、退款金额
	var transAmt, refundAmt, fee int64 = 0, 0, 0

	// 生成数据
	for _, v := range trans {

		var amt float64

		// 交易金额 = 成功的交易金额
		// 手续费 = 支付交易的手续费-（退款、撤销、取消）手续费
		switch v.TransType {
		case model.PayTrans:
			amt = float64(v.TransAmt) / 100
			if v.ChanCode == channel.ChanCodeAlipay {
				alpTransAmt += v.TransAmt
				alpFee += v.Fee
			}
			if v.ChanCode == channel.ChanCodeWeixin {
				wxpTransAmt += v.TransAmt
				wxpFee += v.Fee
			}
		// 退款、撤销、取消
		default:
			amt = -float64(v.TransAmt) / 100
			if v.ChanCode == channel.ChanCodeAlipay {
				alpRefundAmt += v.TransAmt
				alpFee -= v.Fee
			}
			if v.ChanCode == channel.ChanCodeWeixin {
				wxpRefundAmt += v.TransAmt
				wxpFee -= v.Fee
			}
		}

		//商户号，商户名称，订单号，金额，渠道，交易时间，交易状态，终端号，交易类型，原订单号
		row = sheet.AddRow()
		// 商户号
		cell = row.AddCell()
		cell.Value = v.MerId
		// 商户名称
		cell = row.AddCell()
		cell.Value = v.MerName
		// 订单号
		cell = row.AddCell()
		cell.Value = v.OrderNum
		// 交易金额
		cell = row.AddCell()
		cell.SetFloatWithFormat(float64(amt), floatFormat)
		// 渠道
		cell = row.AddCell()
		switch v.ChanCode {
		case "WXP":
			cell.Value = lWXP
		case "ALP":
			cell.Value = lALP
		default:
			cell.Value = locale.ChanCode.Unknown
		}
		// 交易时间
		cell = row.AddCell()
		cell.Value = v.CreateTime
		// 支付时间
		cell = row.AddCell()
		cell.Value = v.PayTime
		// 交易状态
		cell = row.AddCell()
		switch v.TransStatus {
		case model.TransSuccess:
			cell.Value = locale.TransStatus.TransSuccess
		case model.TransFail:
			cell.Value = locale.TransStatus.TransFail
		case model.TransHandling:
			cell.Value = locale.TransStatus.TransHandling
		case model.TransClosed:
			// 针对退款的交易
			cell.Value = locale.TransStatus.TransClosed
		default:
			cell.Value = locale.TransStatus.Unknown
		}
		// 机构号
		cell = row.AddCell()
		cell.Value = v.AgentCode
		// 终端号
		cell = row.AddCell()
		cell.Value = v.Terminalid
		// 交易类型
		cell = row.AddCell()
		switch v.Busicd {
		case model.Purc:
			cell.Value = locale.BusicdType.Purc
		case model.Paut:
			cell.Value = locale.BusicdType.Paut
		case model.Refd:
			cell.Value = locale.BusicdType.Refd
		case model.Void:
			cell.Value = locale.BusicdType.Void
		case model.Canc:
			cell.Value = locale.BusicdType.Canc
		case model.Qyzf:
			cell.Value = locale.BusicdType.Qyzf
		case model.Jszf:
			cell.Value = locale.BusicdType.Jszf
		default:
			cell.Value = locale.BusicdType.Unknown
		}

		// 原订单号
		cell = row.AddCell()
		cell.Value = v.OrigOrderNum

	}

	// 总金额
	transAmt = wxpTransAmt + alpTransAmt
	refundAmt = wxpRefundAmt + alpRefundAmt
	fee = alpFee + wxpFee

	// 写入汇总数据
	rows := sheet.Rows
	row = rows[0]
	row.WriteStruct(&summary{
		lALP + m.TransAmt + "：", float64(alpTransAmt) / 100,
		lALP + m.RefundAmt + "：", -float64(alpRefundAmt) / 100,
		lALP + m.Fee + "：", float64(alpFee) / 100,
		lALP + m.SettAmt + "：", float64(alpTransAmt-alpRefundAmt-alpFee) / 100,
	}, -1)
	row = rows[1]
	row.WriteStruct(&summary{
		lWXP + m.TransAmt + "：", float64(wxpTransAmt) / 100,
		lWXP + m.RefundAmt + "：", -float64(wxpRefundAmt) / 100,
		lWXP + m.Fee + "：", float64(wxpFee) / 100,
		lWXP + m.SettAmt + "：", float64(wxpTransAmt-wxpRefundAmt-wxpFee) / 100,
	}, -1)
	row = rows[2]
	row.WriteStruct(&summary{
		m.TotalTransAmt + "：", float64(transAmt) / 100,
		m.TotalRefundAmt + "：", -float64(refundAmt) / 100,
		m.TotalFee + "：", float64(fee) / 100,
		m.TotalSettAmt + "：", float64(transAmt-refundAmt-fee) / 100,
	}, -1)
}

type summary struct {
	Cell2 string
	Cell3 float64
	Cell4 string
	Cell5 float64
	Cell6 string
	Cell7 float64
	Cell8 string
	Cell9 float64
}
