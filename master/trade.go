package master

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CardInfoLink/quickpay/channel"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
)

var maxReportRec = 10000

// tradeQuery 交易查询
func tradeQuery(w http.ResponseWriter, data []byte) {

	// // 允许跨域
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "*")

	// 交易查询
	q := &model.QueryCondition{}
	err := json.Unmarshal(data, q)
	if err != nil {
		log.Errorf("unmarshal json(%s) error: %s", data, err)
		http.Error(w, "json format error: "+err.Error(), http.StatusPreconditionFailed)
		return
	}
	if q.EndTime != "" {
		q.EndTime += " 23:59:59"
	}
	ret := core.TransQuery(q)
	retBytes, err := json.Marshal(ret)
	if err != nil {
		log.Error(err)
		http.Error(w, "system error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(retBytes)
}

// tradeReport 处理查找所有商户的请求
func tradeReport(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	filename := params.Get("filename")

	var file = xlsx.NewFile()

	var merId = params.Get("merId")
	req := &model.QueryCondition{
		MerId:        merId,
		Busicd:       params.Get("busicd"),
		StartTime:    params.Get("startTime"),
		EndTime:      params.Get("endTime"),
		OrderNum:     params.Get("orderNum"),
		OrigOrderNum: params.Get("origOrderNum"),
		Size:         maxReportRec,
		IsForReport:  true,
		Page:         1,
		RefundStatus: model.TransRefunded,
		TransStatus:  model.TransSuccess,
	}

	// 查询
	ret := core.TransQuery(req)

	// 类型转换
	if trans, ok := ret.Rec.([]model.Trans); ok {
		// 生成报表
		before := time.Now()
		genReport(merId, file, trans)
		after := time.Now()
		log.Debugf("gen trans report spent %s", after.Sub(before))
	}

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
		// ErrorDetail  string
	}{"商户号", "订单号", "原订单号", "终端号", "机构号", "应答码", "交易金额（元）", "交易状态", "交易类型", "渠道", "交易时间"}
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
		cell.Value = v.AgentCode
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
			// 针对退款的交易
			cell.Value = "交易已退款"
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
		case model.Jszf:
			cell.Value = "公众号支付"
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

		// 金额
		switch v.TransType {
		case model.PayTrans:
			if v.ChanCode == channel.ChanCodeAlipay {
				alpTransAmt += v.TransAmt - v.RefundAmt
				alpFee += v.Fee
			}
			if v.ChanCode == channel.ChanCodeWeixin {
				wxpTransAmt += v.TransAmt - v.RefundAmt
				wxpFee += v.Fee
			}
		// 退款、撤销、取消
		default:
			if v.ChanCode == channel.ChanCodeAlipay {
				alpRefundAmt += v.TransAmt
			}
			if v.ChanCode == channel.ChanCodeWeixin {
				wxpRefundAmt += v.TransAmt
			}
		}
	}

	// 利用商户数据，完善报表数据
	var merName string
	if merId != "" {
		mer, err := mongo.MerchantColl.Find(merId)
		if err == nil {
			merName = mer.Detail.MerName
		}
	}

	// 总金额
	transAmt = wxpTransAmt + alpTransAmt
	refundAmt = wxpRefundAmt + alpRefundAmt
	fee = alpFee + wxpFee

	// 写入汇总数据
	// TODO 手续费计算，在记录交易时计算
	rows := sheet.Rows
	row = rows[0]
	row.WriteStruct(&summary{
		"名称：", merName,
		"支付宝交易金额：", float64(alpTransAmt) / 100,
		"支付宝退款金额：", float64(alpRefundAmt) / 100,
		"支付宝手续费：", float64(alpFee) / 100,
		// "支付宝清算金额：", 50.00,
	}, -1)
	row = rows[1]
	row.WriteStruct(&summary{
		"", "",
		"微信交易金额：", float64(wxpTransAmt) / 100,
		"微信退款金额：", float64(wxpRefundAmt) / 100,
		"微信手续费：", float64(wxpFee) / 100,
		// "微信清算金额：", 50.00,
	}, -1)
	row = rows[2]
	row.WriteStruct(&summary{
		"总计：", "",
		"交易总额：", float64(transAmt) / 100,
		"退款总额：", float64(refundAmt) / 100,
		"手续费总额：", float64(fee) / 100,
		// "清算总额：", 100.00,
	}, -1)
}

type summary struct {
	Cell0 string
	Cell1 string
	Cell2 string
	Cell3 float64
	Cell4 string
	Cell5 float64
	Cell6 string
	Cell7 float64
	// Cell8 string
	// Cell9 float64
}
