package master

import (
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/CardInfoLink/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var avaiableLocale = make(map[string]*LocaleTemplate)

const (
	prefix        = "message_"
	DefaultLocale = "zh-CN"
)

func init() {

	var mp []string
	md := util.WorkDir + "/config/message/"
	filepath.Walk(md, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fn := f.Name()
		if strings.HasPrefix(fn, prefix) {
			mp = append(mp, fn)
		}
		return nil
	})
	// load
	for _, fn := range mp {
		filePath := md + fn
		bs, err := ioutil.ReadFile(md + fn)
		if err != nil {
			fmt.Printf("read config file %s error: %s\n", filePath, err)
			os.Exit(4)
		}
		lt := &LocaleTemplate{}
		err = json.Unmarshal(bs, lt)
		if err != nil {
			fmt.Printf("config file %s parser error: %s\n", filePath, err)
			os.Exit(5)
		}
		avaiableLocale[fn[:len(fn)-3]] = lt
	}

	// 检查是否有默认语言-中文
	if !IsLocaleExist(DefaultLocale) {
		log.Errorf("fail to find default locale: %s", DefaultLocale)
		os.Exit(6)
	}

	// log.Debugf("%+v", GetLocale(DefaultLocale))
}

// IsLocaleExist 是否有该语言模板
func IsLocaleExist(l string) bool {
	if _, ok := avaiableLocale[prefix+l]; ok {
		return ok
	}
	return false
}

// GetLocale 获得一份语言模板
func GetLocale(l string) *LocaleTemplate {
	if locale, ok := avaiableLocale[prefix+l]; ok {
		return locale
	}
	// 默认中文
	return avaiableLocale[prefix+DefaultLocale]
}

// LocaleTemplate 语言模板
type LocaleTemplate struct {
	Currency string
	// 渠道代号
	ChanCode struct {
		ALP     string
		WXP     string
		Unknown string
	}

	ReportName struct {
		SettleJournal string
		SettleSummary string
	}

	// 角色
	Role struct {
		Company string
		Agent   string
		Group   string
		Mer     string
	}

	// 交易类型
	BusicdType struct {
		Purc    string
		Paut    string
		Inqy    string
		Refd    string
		Void    string
		Canc    string
		Qyzf    string
		Jszf    string
		Veri    string
		Unknown string
	}

	// 交易状态
	TransStatus struct {
		TransHandling string
		TransFail     string
		TransSuccess  string
		TransClosed   string
		Unknown       string
	}

	// 统计报表
	StatReport    StatReport
	ImportMessage ImportMessage
	TransReport   TransReport
}

// TransReport 交易明细报表
type TransReport struct {
	Yes            string
	No             string
	SheetName      string
	MerId          string
	MerName        string
	OrderNum       string
	TransAmt       string
	TransCurr      string
	MerFee         string
	ChanCode       string
	TransTime      string
	PayTime        string
	TransStatus    string
	ChanMerId      string
	AgentCode      string
	TerminalId     string
	Busicd         string
	OrigOrderNum   string
	Remark         string
	IsSettled      string
	RefundAmt      string
	Fee            string
	SettAmt        string
	TotalTransAmt  string
	TotalRefundAmt string
	TotalFee       string
	TotalSettAmt   string
}

// StatReport 统计报表
type StatReport struct {
	Title       string
	Total       string
	StartDate   string
	EndDate     string
	Remark      string
	MerId       string
	MerName     string
	Summary     string
	AgentName   string
	CompanyName string
	GroupName   string
	TotalCount  string
	TotalAmt    string
	TotalFee    string
	Count       string
	Amt         string
	Fee         string
}

// ImportMessage 批导信息
type ImportMessage struct {
	Yes           string
	No            string
	SysErr        string
	EmptyErr      string
	FileErr       string
	CellMapErr    string
	ColNumErr     string
	ImportSuccess string
	MerIdRepeat   string
	DataHandleErr struct {
		NotSupportOperation string
		NoMerId             string
		MerIdFormatErr      string
		MerIdExist          string
		MerIdNotExist       string
		ALPMerchantErr      string
		WXPMerchantErr      string
		UsernameExist       string
		UsernameNotExist    string
		AgentNotExist       string
		CompanyNotExist     string
		CompanyBelongsErr   string
		GroupNotExist       string
		GroupBelongsErr     string
		NoALPKey            string
		NoALPRouteToUdpSf   string
		NoWXPRouteToUdpSf   string
		WXPNotAgentMode     string
		SysConfigErr        string
		AgentMerInfoErr     string
		AgentModeNotMatch   string
		NoSuchAgentMer      string
		NoWXPKey            string
		CILFeeErr           string
		MerFeeErr           string
		CILFeeOverMax       string
		MerFeeOverMax       string
	}
	ValidateErr struct {
		NoMerName         string
		NoSignKey         string
		NoAgentCode       string
		OpenSignValueErr  string
		AddAcctValueErr   string
		UNOrPWDEmptyErr   string
		SignLengthErr     string
		NoCommodityName   string
		IsAgentStrErr     string
		NoWXPMer          string
		WXPSettFlagErr    string
		ALPSettFlagErr    string
		NoOverseasChanMer string
		NoSchemeType      string
	}
}
