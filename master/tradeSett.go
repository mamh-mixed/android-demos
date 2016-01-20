package master

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/currency"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/tealeg/xlsx"
	"net/http"
	"time"
)

// tradeSettJournalReport 对账流水下载
func tradeSettJournalReport(w http.ResponseWriter, cond *model.QueryCondition) {
	// 语言模板
	rl := GetLocale(cond.Locale)

	var filename string
	reportName := rl.ReportName.SettleJournal
	switch cond.UserType {
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

	// 设置返回content
	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`, fmt.Sprintf(`attachment; filename="%s"`, filename))

	// 查询
	transSetts, _ := mongo.SpTransSettColl.Find(cond)

	// 生成报表
	file := settJornalReport(transSetts, rl, &Zone{cond.UtcOffset, time.Local})

	file.Write(w)
}

// genReport 生成报表
func settJornalReport(transSetts []model.TransSett, locale *LocaleTemplate, z *Zone) (file *xlsx.File) {

	file = xlsx.NewFile()
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
		TransCurr    string
		MerFee       string
		ChanCode     string
		TransTime    string
		PayTime      string
		TransStatus  string
		AgentCode    string
		CompanyName  string
		GroupName    string
		TerminalId   string
		Busicd       string
		OrigOrderNum string
		Remark       string
		IsSettled    string
	}{m.MerId, m.MerName, m.OrderNum, m.TransAmt, m.TransCurr, m.MerFee, m.ChanCode, m.TransTime, m.PayTime, m.TransStatus, m.AgentCode, locale.Role.Company, locale.Role.Group, m.TerminalId, m.Busicd, m.OrigOrderNum, m.Remark, m.IsSettled}
	row.WriteStruct(headRow, -1)

	// 设置列宽
	sheet.SetColWidth(0, 9, 18)

	// 支付宝交易金额、退款金额
	var alpTransAmt, alpRefundAmt, alpFee int64 = 0, 0, 0
	// 微信交易金额、退款金额
	var wxpTransAmt, wxpRefundAmt, wxpFee int64 = 0, 0, 0
	// 总交易金额、退款金额
	var transAmt, refundAmt, fee int64 = 0, 0, 0
	// 参与清算总额
	var alpSettAmt, wxpSettAmt int64 = 0, 0

	var cur currency.Cur
	// 生成数据
	if len(transSetts) != 0 {
		// TODO 先随机取一条交易的币种确定单位
		transCurr := transSetts[0].Trans.Currency

		// 币种单位
		cur = currency.Get(transCurr)

		// 金额显示格式
		var floatFormat = "#,##0"
		for i := 0; i < cur.Precision; i++ {
			if i == 0 {
				floatFormat += "."
			}
			floatFormat += "0"
		}
		for _, ts := range transSetts {

			v := ts.Trans
			var amt float64

			// 交易金额 = 成功的交易金额
			// 手续费 = 支付交易的手续费-（退款、撤销、取消）手续费
			switch v.TransType {
			case model.PayTrans:
				amt = cur.F64(v.TransAmt)
				if v.ChanCode == channel.ChanCodeAlipay {
					alpTransAmt += v.TransAmt
					alpFee += ts.MerFee
					if ts.BlendType == 0 {
						alpSettAmt += (v.TransAmt - ts.MerFee)
					}
				}
				if v.ChanCode == channel.ChanCodeWeixin {
					wxpTransAmt += v.TransAmt
					wxpFee += ts.MerFee
					if ts.BlendType == 0 {
						wxpSettAmt += (v.TransAmt - ts.MerFee)
					}
				}
			// 退款、撤销、取消
			default:
				amt = -cur.F64(v.TransAmt)
				if v.ChanCode == channel.ChanCodeAlipay {
					alpRefundAmt += v.TransAmt
					alpFee -= ts.MerFee
					if ts.BlendType == 0 {
						alpSettAmt -= (v.TransAmt - ts.MerFee)
					}
				}
				if v.ChanCode == channel.ChanCodeWeixin {
					wxpRefundAmt += v.TransAmt
					wxpFee -= ts.MerFee
					if ts.BlendType == 0 {
						wxpSettAmt -= (v.TransAmt - ts.MerFee)
					}
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
			// 交易币种
			cell = row.AddCell()
			cell.Value = v.Currency
			// 商户手续费
			cell = row.AddCell()
			cell.SetFloatWithFormat(cur.F64(ts.MerFee), floatFormat)
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
			if v.PayTime == "" {
				v.PayTime = v.CreateTime
			}
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
			// 公司
			cell = row.AddCell()
			cell.Value = v.SubAgentCode
			// 商户
			cell = row.AddCell()
			cell.Value = v.GroupCode
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
			// 是否已清算
			cell = row.AddCell()
			if ts.BlendType == 0 {
				cell.Value = "Y"
			} else {
				cell.Value = "N"
			}
		}
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
		lALP + m.SettAmt + "：", cur.F64(alpSettAmt),
	}, -1)
	row = rows[1]
	row.WriteStruct(&summary{
		lWXP + m.TransAmt + "：", cur.F64(wxpTransAmt),
		lWXP + m.RefundAmt + "：", -cur.F64(wxpRefundAmt),
		lWXP + m.Fee + "：", cur.F64(wxpFee),
		lWXP + m.SettAmt + "：", cur.F64(wxpSettAmt),
	}, -1)
	row = rows[2]
	row.WriteStruct(&summary{
		m.TotalTransAmt + "：", cur.F64(transAmt),
		m.TotalRefundAmt + "：", -cur.F64(refundAmt),
		m.TotalFee + "：", cur.F64(fee),
		m.TotalSettAmt + "：", cur.F64(alpSettAmt + wxpSettAmt),
	}, -1)

	return file
}

// settJornalReport2 生成交易明细
func settJornalReport2(transSetts []model.TransSett, locale *LocaleTemplate, z *Zone) (file *xlsx.File) {

	file = xlsx.NewFile()
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	// 语言
	m := locale.TransReport
	lALP, lWXP := locale.ChanCode.ALP, locale.ChanCode.WXP

	// 可能有多个sheet
	sheet, _ = file.AddSheet(m.SheetName)

	// 生成title
	row = sheet.AddRow()
	headRow := &struct {
		MerId        string
		MerName      string
		OrderNum     string
		TransAmt     string
		TransCurr    string
		MerFee       string
		ChanCode     string
		TransTime    string
		PayTime      string
		TransStatus  string
		ChanMerId    string
		AgentCode    string
		TerminalId   string
		Busicd       string
		OrigOrderNum string
		Remark       string
		IsSettled    string
	}{m.MerId, m.MerName, m.OrderNum, m.TransAmt, m.TransCurr, m.MerFee, m.ChanCode, m.TransTime, m.PayTime, m.TransStatus, m.ChanMerId, m.AgentCode, m.TerminalId, m.Busicd, m.OrigOrderNum, m.Remark, m.IsSettled}
	row.WriteStruct(headRow, -1)

	// 设置列宽
	sheet.SetColWidth(0, 9, 18)
	var cur currency.Cur
	// 生成数据
	if len(transSetts) != 0 {
		// TODO 先随机取一条交易的币种确定单位
		transCurr := transSetts[0].Trans.Currency

		// 币种单位
		cur = currency.Get(transCurr)

		// 金额显示格式
		var floatFormat = "#,##0"
		for i := 0; i < cur.Precision; i++ {
			if i == 0 {
				floatFormat += "."
			}
			floatFormat += "0"
		}
		for _, ts := range transSetts {

			v := ts.Trans
			var amt float64

			// 交易金额 = 成功的交易金额
			// 手续费 = 支付交易的手续费-（退款、撤销、取消）手续费
			switch v.TransType {
			case model.PayTrans:
				amt = cur.F64(v.TransAmt)
			// 退款、撤销、取消
			default:
				amt = -cur.F64(v.TransAmt)
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
			// 交易币种
			cell = row.AddCell()
			cell.Value = v.Currency
			// 商户手续费
			cell = row.AddCell()
			if v.TransType == model.PayTrans {
				cell.SetFloatWithFormat(cur.F64(ts.MerFee), floatFormat)
			} else {
				cell.SetFloatWithFormat(cur.F64(-ts.MerFee), floatFormat)
			}
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
			if v.PayTime == "" {
				v.PayTime = v.CreateTime
			}
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
			// 渠道商户号
			cell = row.AddCell()
			cell.Value = v.ChanMerId
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
			// 是否已清算
			cell = row.AddCell()
			if ts.BlendType == 0 {
				cell.Value = m.Yes
			} else {
				cell.Value = m.No
			}
		}
	}
	return file
}
