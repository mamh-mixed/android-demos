package settle

import (
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
	"strings"
)

//  reportType
const (
	TransferReport       = 1 // 划款报表
	ReconciliationReport = 2 // 对账报表
	InsFlowReport        = 3 // 机构流水报表
	ChanMerReport        = 4 // 渠道商户报表
	ChanLessReport       = 5 // 对账不平报表-渠道少清
	ChanMoreReport       = 6 // 对账不平报表-渠道多清
	// 分润报表
)

const filePrefix = "sett/report/%s/" // 文件名：sett/report/20151012/IC202_99911888_20151012.xlsx

// 最外部key为代理号，接下来是渠道，接下来是清算角色，value是角色下数据
type reconciliationMap map[string]map[string]map[string]*reconciliatReportData

func getRsRecord(reportType int, date string) *model.RoleSett {
	rs := &model.RoleSett{
		ReportName: filePrefix,
		ReportType: reportType,
		SettDate:   date,
	}
	fd := strings.Replace(date, "-", "", -1)
	prefix := fmt.Sprintf(filePrefix, fd)
	switch reportType {
	case TransferReport:
		// multi role
		rs.ReportName = prefix + "IC202_%s_" + fd + ".xlsx"
	case ReconciliationReport:
		rs.ReportName = prefix + "IC002_" + fd + ".xlsx"
		rs.SettRole = "ALL"
	case InsFlowReport:
		// TODO
		rs.ReportName = ""
	case ChanMerReport:
		// TODO
		rs.ReportName = ""
	case ChanLessReport:
		rs.ReportName = prefix + "IC401_" + fd + ".xlsx"
		rs.SettRole = "ALL"
	case ChanMoreReport:
		rs.ReportName = prefix + "IC402_" + fd + ".xlsx"
		rs.SettRole = "ALL"
	}
	return rs
}

// SpSettReport 扫码清算报表
func SpSettReport(settDate string) error {

	data, err := mongo.SpTransSettColl.GroupBySettRole(settDate)
	if err != nil {
		log.Errorf("fail to find trans group by settRole: %s", err)
		return err
	}

	// 报表日期显示格式
	sd := strings.Replace(settDate, "-", "", -1)

	// 遍历数据
	for _, sg := range data {
		// 记录
		rs := getRsRecord(TransferReport, settDate)
		rs.ReportName = fmt.Sprintf(rs.ReportName, sg.SettRole)
		rs.SettRole = sg.SettRole

		rpData := settDataHandle(sg)
		// 有数据才生成报表
		if len(rpData) != 0 {
			// 生成报表上传
			if err = upload(rs.ReportName, genSpTransferReportExcel(rpData, sd)); err != nil {
				continue
			}
			if err = mongo.RoleSettCol.Upsert(rs); err != nil {
				log.Errorf("roleSett upsert error: %s", err)
			}
		}
	}

	return nil
}

// SpReconciliatReport 扫码对账报表
// date should be yyyy-mm-dd
func SpReconciliatReport(date string, transSetts ...model.TransSett) error {

	// 判断数据源
	if len(transSetts) == 0 {
		tss, err := mongo.SpTransSettColl.Find(&model.QueryCondition{Date: date})
		if err != nil {
			return err
		}
		if len(tss) == 0 {
			return nil
		}
		transSetts = tss
	}

	// 外key为机构号，内map的key为渠道编号
	reconciliatMMap := make(reconciliationMap)

	// 处理数据
	for _, transSett := range transSetts {
		t := transSett.Trans
		// 没有机构号或者没有清算角色
		if t.AgentCode == "" || t.SettRole == "" {
			log.Errorf("gen reconciliation report, but params empty, merId=%s, orderNum=%s", t.MerId, t.OrderNum)
			continue
		}

		// 机构
		if chanDataMap, ok := reconciliatMMap[t.AgentCode]; ok {
			// 渠道
			if roleDataMap, ok := chanDataMap[t.ChanCode]; ok {
				// 角色
				if d, found := roleDataMap[t.SettRole]; found {
					//下单和预下单相加，退款、取消和撤销相减
					if t.TransType == model.PayTrans {
						d.transAmt += t.TransAmt
						d.MerFee += transSett.MerFee
						d.MerSettAmt += transSett.MerSettAmt
						d.AcqFee += transSett.AcqFee
						d.AcqSettAmt += transSett.AcqSettAmt
						d.transNum++
					} else {
						d.transAmt -= t.TransAmt
						d.MerFee -= transSett.MerFee
						d.MerSettAmt -= transSett.MerSettAmt
						d.AcqFee -= transSett.AcqFee
						d.AcqSettAmt -= transSett.AcqSettAmt
						// d.transNum++
					}
				} else {
					// 没找到角色
					roleDataMap[t.SettRole] = NewReconciliatData(transSett)
				}
			} else {
				// 没有渠道
				roleDataMap := make(map[string]*reconciliatReportData)
				roleDataMap[t.SettRole] = NewReconciliatData(transSett)
				chanDataMap[t.ChanCode] = roleDataMap
			}
		} else {
			// 还没存在该机构下记录
			roleDataMap := make(map[string]*reconciliatReportData)
			roleDataMap[t.SettRole] = NewReconciliatData(transSett)
			chanDataMap := make(map[string]map[string]*reconciliatReportData)
			chanDataMap[t.ChanCode] = roleDataMap
			reconciliatMMap[t.AgentCode] = chanDataMap
		}
	}

	// for _, v := range reconciliatMMap {
	// 	for _, v1 := range v {
	// 		for k, v2 := range v1 {
	// 			log.Debugf("key=%s,value=%+v", k, v2)
	// 		}
	// 	}
	// }

	if len(reconciliatMMap) != 0 {
		// 上传并记录
		rs := getRsRecord(ReconciliationReport, date)
		if err := upload(rs.ReportName, genReconciliatReportExcel(reconciliatMMap, date)); err == nil {
			if err = mongo.RoleSettCol.Upsert(rs); err != nil {
				log.Errorf("roleSett upsert error: %s", err)
			}
		}
	}

	return nil
}

// settDataHandle 清算数据处理
func settDataHandle(sg model.SettRoleGroup) []reportData {

	var rds []reportData
	for _, mg := range sg.MerGroups {
		// 交易金额为0的过滤掉
		if mg.TransAmt == 0 {
			continue
		}
		m, err := mongo.MerchantColl.Find(mg.MerId)
		if err != nil {
			// cmMap[mg.MerId] = 0 // 标识不成功
			// continue
			m = &model.Merchant{MerId: mg.MerId} // 兼容老系统数据，可能商户没同步到新系统
		}
		// 补充开户银行和支行
		if m.Detail.OpenBankName == "" {
			m.Detail.OpenBankName = m.Detail.BankName
		}
		if m.Detail.BankName == "" {
			m.Detail.BankName = m.Detail.OpenBankName
		}
		rds = append(rds, reportData{mg: mg, m: *m})
	}

	return rds
}

type reportData struct {
	mg model.MerGroup
	m  model.Merchant
}

type reconciliatReportData struct {
	insCode    string //客户代码
	insName    string //客户名称
	chcd       string //渠道名称
	role       string //清算角色
	transNum   int    //交易笔数
	transAmt   int64  //交易金额
	MerFee     int64  //商户手续费
	MerSettAmt int64  //商户应收金额
	AcqFee     int64  //讯联手续费
	AcqSettAmt int64  //讯联应收金额
}

func NewReconciliatData(ts model.TransSett) *reconciliatReportData {
	return &reconciliatReportData{
		insCode:    ts.Trans.AgentCode,
		insName:    ts.Trans.AgentName,
		chcd:       ts.Trans.ChanCode,
		role:       ts.Trans.SettRole,
		transNum:   1,
		transAmt:   ts.Trans.TransAmt,
		MerFee:     ts.MerFee,
		MerSettAmt: ts.MerSettAmt,
		AcqFee:     ts.AcqFee,
		AcqSettAmt: ts.AcqSettAmt,
	}
}

type errorReportData struct {
	chanCode    string //渠道编号
	chanName    string //渠道名称
	merCode     string //商户编号
	chanMerCode string //虚拟商户号
	merName     string //商户名称
	orderId     string //订单号
	transTime   string //交易时间
	transType   string //交易类型
	transAmt    int64  //交易金额
	compareType string //对账标识
}
