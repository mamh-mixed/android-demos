package master

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/currency"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/query"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
	"net/http"
	"time"
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
		trans, total := query.SpTransQuery(q)
		return &model.ResultBody{
			Status:  0,
			Message: "查询成功",
			Data: &model.Pagination{
				Page:  q.Page,
				Total: total,
				Size:  q.Size,
				Count: len(trans),
				Data:  trans,
			},
		}
	}
}

// tradeQuery 交易查询
func tradeFindOne(q *model.QueryCondition) (ret *model.ResultBody) {
	return query.SpTransFindOne(q)
}

// tradeReport 处理查找所有商户的请求
func tradeReport(w http.ResponseWriter, cond *model.QueryCondition, filename string) {
	var file = xlsx.NewFile()

	// 语言模板
	rl := GetLocale(cond.Locale)

	// 查询
	trans, _ := query.SpTransQuery(cond)

	// 生成报表
	genReport(file, trans, rl, &Zone{cond.UtcOffset, time.Local})

	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`, fmt.Sprintf(`attachment; filename="%s"`, filename))
	file.Write(w)
}

// genReport 生成报表
func genReport(file *xlsx.File, trans []*model.Trans, locale *LocaleTemplate, z *Zone) {
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	// 语言
	m := locale.TransReport
	lALP, lWXP := locale.ChanCode.ALP, locale.ChanCode.WXP

	// 币种单位
	cur := currency.Get(locale.Currency)

	// format
	floatFormat := locale.ExportF64Format

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
		MerFee       string
		ChanCode     string
		TransTime    string
		PayTime      string
		TransStatus  string
		AgentCode    string
		TerminalId   string
		Busicd       string
		OrigOrderNum string
		Remark       string
	}{m.MerId, m.MerName, m.OrderNum, m.TransAmt, m.MerFee, m.ChanCode, m.TransTime, m.PayTime, m.TransStatus, m.AgentCode, m.TerminalId, m.Busicd, m.OrigOrderNum, m.Remark}
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
			amt = cur.F64(v.TransAmt)
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
			amt = -cur.F64(v.TransAmt)
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
		cell.SetFloatWithFormat(amt, floatFormat)
		// 商户手续费
		cell = row.AddCell()
		cell.SetFloatWithFormat(cur.F64(v.NetFee), floatFormat)
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
		cell.Value = z.GetTime(v.CreateTime)
		// 支付时间，维持北京时间
		cell = row.AddCell()
		cell.Value = v.PayTime + " +0800"
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
		// 备注
		cell = row.AddCell()
		cell.Value = v.TicketNum
	}

	// 总金额
	transAmt = wxpTransAmt + alpTransAmt
	refundAmt = wxpRefundAmt + alpRefundAmt
	fee = alpFee + wxpFee

	// 写入汇总数据
	rows := sheet.Rows
	row = rows[0]
	row.WriteStruct(&summary{
		lALP + m.TransAmt + "：", cur.F64(alpTransAmt),
		lALP + m.RefundAmt + "：", -cur.F64(alpRefundAmt),
		lALP + m.Fee + "：", cur.F64(alpFee),
		lALP + m.SettAmt + "：", cur.F64(alpTransAmt - alpRefundAmt - alpFee),
	}, -1)
	row = rows[1]
	row.WriteStruct(&summary{
		lWXP + m.TransAmt + "：", cur.F64(wxpTransAmt),
		lWXP + m.RefundAmt + "：", -cur.F64(wxpRefundAmt),
		lWXP + m.Fee + "：", cur.F64(wxpFee),
		lWXP + m.SettAmt + "：", cur.F64(wxpTransAmt - wxpRefundAmt - wxpFee),
	}, -1)
	row = rows[2]
	row.WriteStruct(&summary{
		m.TotalTransAmt + "：", cur.F64(transAmt),
		m.TotalRefundAmt + "：", -cur.F64(refundAmt),
		m.TotalFee + "：", cur.F64(fee),
		m.TotalSettAmt + "：", cur.F64(transAmt - refundAmt - fee),
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

// Zone 代表时区
type Zone struct {
	UtcOffset     int            // UTC 的偏移量
	ParseLocation *time.Location // 原时间时区
}

// GetTime 获得某个地区时间
func (z *Zone) GetTime(ctime string) string {

	// 东八区时间偏移量
	cstOffset := 60 * 60 * 8

	// 假如是东八区北京时间，直接返回
	if z.UtcOffset == cstOffset {
		return ctime
	}

	// 以服务器时区为准，即东八区
	t, err := time.ParseInLocation("2006-01-02 15:04:05", ctime, z.ParseLocation)
	if err != nil {
		log.Errorf("fail to parse time in Local: %s", ctime)
		return ctime
	}

	if loc, ok := locationsMap[z.UtcOffset]; ok {
		t = t.In(loc)
		return t.Format(layout)
	}

	loc := time.FixedZone("CUR", z.UtcOffset)
	locationsMap[z.UtcOffset] = loc

	t = t.In(loc)
	return t.Format(layout)
}

var layout = "2006-01-02 15:04:05 -0700"
var locationsMap = make(map[int]*time.Location)
