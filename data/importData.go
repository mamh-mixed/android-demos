// Package data 读取csv文件持久化到数据库
// 具体的字段顺序在csv文件里的第一行
package data

import (
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// readUserCSV
func readUserCSV(path string) ([]model.User, error) {
	data, err := readCSV(path)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0, len(data))
	// 门店类型	商户编号	门店名称	手机号	对账平台账号	代理编号	公司编号	商户编号
	for _, v := range data {
		u := model.User{}
		u.UserType = model.UserTypeShop
		u.MerId = strings.TrimSpace(v[1])
		u.PhoneNum = strings.TrimSpace(v[3])
		u.UserName = strings.TrimSpace(v[4])
		u.AgentCode = strings.TrimSpace(v[6])
		u.SubAgentCode = strings.TrimSpace(v[7])
		u.GroupCode = strings.TrimSpace(v[8])
		passData := []byte(model.RAND_PWD + "{" + u.UserName + "}" + model.DEFAULT_PWD)
		u.Password = fmt.Sprintf("%x", sha1.Sum(passData))
		users = append(users, u)
	}

	return users, nil
}

// InitTestMer 初始化测试商户
// start:从哪个数值开始 end:结束 cardBrand:卡品牌
// 初始化后格式:test000001
func InitTestMer(start, end int, cardBrand string) error {
	if start > end {
		return fmt.Errorf("%s", "end must large than start")
	}
	if end > 999999 {
		return fmt.Errorf("%s", "end must smaller than 999999")
	}
	for i := start; i <= end; i++ {
		rp := &model.RouterPolicy{
			MerId:     fmt.Sprintf("test%06d", i),
			CardBrand: cardBrand,
			ChanCode:  "CFCA",
			ChanMerId: "001405",
		}
		mongo.RouterPolicyColl.Insert(rp)
		m := &model.Merchant{
			MerId:      fmt.Sprintf("test%06d", i),
			MerStatus:  "Normal",
			TransCurr:  "156",
			SignKey:    "0123456789",
			EncryptKey: "AAECAwQFBgcICQoLDA0ODwABAgMEBQYHCAkKCwwNDg8=",
		}
		mongo.MerchantColl.Insert(m)
	}
	return nil
}

// AddSettSchemeCdFromCSV 导入计费方案代码
func AddSettSchemeCdFromCSV(path string) error {

	schemes, err := readSettSchemeCdCSV(path)
	// fmt.Println(len(schemes), err)
	if err != nil {
		return nil
	}

	for _, v := range schemes {
		err := mongo.SettSchemeCdCol.Upsert(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddCardBinFromCSV 从 csv 里导入卡 bin
// rebuild: true 删除集合再重建
// rebuild: false 做更新操作，即存在的更新，不存在的增加
func AddCardBinFromCSV(path string, rebuild bool) error {

	cardBins, err := readCardBinCSV(path)
	if err != nil {
		return err
	}
	fmt.Printf("%+v	", cardBins[2])
	// return nil
	// fmt.Println(len(cardBins))
	// 重建
	if rebuild {
		fmt.Println("drop cardBin collection...")
		err = mongo.CardBinColl.Drop()
		if err != nil {
			return err
		}
	}
	for _, v := range cardBins {
		err = mongo.CardBinColl.Upsert(v)
		if err != nil {
			return err
		}
		fmt.Print(".")
	}
	fmt.Printf("\nImported cardBin %d records\n", len(cardBins))

	return nil
}

// AddScanPayFromCSV 从csv文件里读取应答码表，若存在的跳过，若新增的便添加
func AddScanPayFromCSV(path string) error {

	data, err := readScanPayCSV(path)
	if err != nil {
		return err
	}

	for _, v := range data {
		sp, err := mongo.ScanPayRespCol.FindOne(v.ISO8583Code)
		if err != nil {
			// fmt.Printf("New Add: %+v \n", v)
			mongo.ScanPayRespCol.Add(v)
			continue
		}
		// update
		v.Alp = sp.Alp
		v.Wxp = sp.Wxp
		mongo.ScanPayRespCol.Update(v)
	}
	fmt.Printf("\nImported ScanPayRespCode %d records\n", len(data))
	return nil
}

// AddSysCodeFromCSV 从csv文件里读取应答码表，若存在的跳过，若新增的便添加
func AddSysCodeFromCSV(path string) error {

	data, err := readQuickpayCSV(path)
	if err != nil {
		return err
	}

	for _, v := range data {
		_, err := mongo.RespCodeColl.FindOne(v.RespCode)
		if err != nil {
			// fmt.Printf("New Add: %+v \n", v)
			fmt.Print(".")
			mongo.RespCodeColl.Add(v)
		}
	}
	fmt.Printf("\nImported RespCode %d records\n", len(data))
	return nil
}

// AddSpChanCodeFromScv 增加扫码渠道应答码
func AddSpChanCodeFromScv(channel, path string) error {

	data, err := readSpChanCSV(path)
	if err != nil {
		return err
	}

	for _, v := range data {
		// fmt.Printf("%+v \n", v)
		q, err := mongo.ScanPayRespCol.FindOne(v.ISOCode)
		if err != nil {
			fmt.Printf("error: %s, code: %s, errorData: %+v \n", err, q.ISO8583Code, v)
			continue
		}
		// 不保存 ISOCode,ISOMsg 两个字段
		v.ISOCode = ""
		v.ISOMsg = ""
		switch {
		case channel == "wxp":
			for i, wxp := range q.Wxp {
				if wxp.Code == v.Code && wxp.Busicd == v.Busicd {
					// delete
					q.Wxp = append(q.Wxp[:i], q.Wxp[i+1:]...)
				}
			}
			q.Wxp = append(q.Wxp, v)
		case channel == "alp":
			// 过滤重复的
			for i, alp := range q.Alp {
				if alp.Code == v.Code && alp.Busicd == v.Busicd {
					// delete
					q.Alp = append(q.Alp[:i], q.Alp[i+1:]...)
				}
			}
			q.Alp = append(q.Alp, v)
		default:
			// ...更多渠道
		}
		// fmt.Printf("%+v \n", q)
		err = mongo.ScanPayRespCol.Update(q)
		if err != nil {
			fmt.Printf("error: %s, code: %s", err, q.ISO8583Code)
		}
	}
	return nil
}

// AddChanCodeFromScv 增加渠道应答码
func AddChanCodeFromScv(channel, path string) error {

	data, err := readChanCSV(path)
	if err != nil {
		return err
	}

	for _, v := range data {
		q, err := mongo.RespCodeColl.FindOne(v.RespCode)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// 不保存respCode,respMsg 两个字段
		v.RespCode = ""
		v.RespMsg = ""
		switch {
		case channel == "cfca":
			for i, cfca := range q.Cil {
				if cfca.Code == v.Code {
					// delete
					q.Cfca = append(q.Cfca[:i], q.Cfca[i+1:]...)
				}
			}
			q.Cfca = append(q.Cfca, v)
		case channel == "cil":
			// 过滤重复的
			for i, cil := range q.Cil {
				if cil.Code == v.Code {
					// delete
					q.Cil = append(q.Cil[:i], q.Cil[i+1:]...)
				}
			}
			q.Cil = append(q.Cil, v)
		default:
			// ...更多渠道
		}
		mongo.RespCodeColl.Update(q)
	}
	return nil
}

func readSettSchemeCdCSV(path string) ([]*model.SettSchemeCd, error) {

	data, err := readCSV(path)
	if err != nil {
		return nil, err
	}
	fmt.Println(data[2])
	qs := make([]*model.SettSchemeCd, 0, len(data))

	for i, each := range data {
		if i == 0 {
			continue
		}
		openIn, _ := strconv.Atoi(each[4])
		eventId, _ := strconv.Atoi(each[5])
		recId, _ := strconv.Atoi(each[6])

		q := &model.SettSchemeCd{SchemeCd: each[0], FitBitMap: each[1], Nm: each[2], Descs: each[3],
			OperIn: openIn, EventId: eventId, RecId: recId, RecUpdTs: each[7], RecCrtTs: each[8]}
		qs = append(qs, q)
	}
	return qs, nil
}

// readScanPayCSV 读取扫码应答码csv文件
// 并持久化
func readScanPayCSV(path string) ([]*model.ScanPayCSV, error) {

	data, err := readCSV(path)
	if err != nil {
		return nil, err
	}
	qs := make([]*model.ScanPayCSV, 0, len(data))

	// 根据数据规则遍历
	for i, each := range data {
		if i == 0 {
			continue
		}
		isUseIso, _ := strconv.ParseBool(each[3])
		q := &model.ScanPayCSV{ISO8583Code: each[0], ISO8583Msg: each[1], ErrorCode: each[2], IsUseISO: isUseIso}
		// fmt.Printf("%+v \n", q)
		qs = append(qs, q)
	}
	// fmt.Println(len(qs))
	return qs, nil
}

// ReadQuickpayCSV 读取系统应答码csv文件
// 并持久化
func readQuickpayCSV(path string) ([]*model.QuickpayCSV, error) {

	data, err := readCSV(path)
	if err != nil {
		return nil, err
	}
	qs := make([]*model.QuickpayCSV, 0, len(data))

	// 根据数据规则遍历
	for i, each := range data {
		if i == 0 {
			continue
		}
		q := &model.QuickpayCSV{RespCode: each[0], RespMsg: each[1]}
		// fmt.Printf("%+v \n", q)
		qs = append(qs, q)
	}
	// fmt.Println(len(qs))
	return qs, nil
}

// readSpChanCSV 读取扫码渠道应答码文件
func readSpChanCSV(path string) ([]*model.SpChanCSV, error) {
	data, err := readCSV(path)
	if err != nil {
		return nil, err
	}
	qs := make([]*model.SpChanCSV, 0, len(data))

	// 根据渠道应答码文件规则遍历
	for i, each := range data {
		if i == 0 {
			continue
		}
		q := &model.SpChanCSV{each[0], each[1], each[2], each[3], ""}
		// fmt.Printf("%+v \n", q)
		qs = append(qs, q)
	}
	return qs, nil
}

// ReadChanCSV 读取渠道应答码文件
func readChanCSV(path string) ([]*model.ChanCSV, error) {
	data, err := readCSV(path)
	if err != nil {
		return nil, err
	}
	qs := make([]*model.ChanCSV, 0, len(data))

	// 根据渠道应答码文件规则遍历
	for i, each := range data {
		if i == 0 {
			continue
		}
		q := &model.ChanCSV{each[0], each[1], each[2], each[3]}
		// fmt.Printf("%+v \n", q)
		qs = append(qs, q)
	}
	return qs, nil
}

// ReadCardBinCSV 从csv读取卡bin转为对象
func readCardBinCSV(path string) ([]*model.CardBin, error) {

	data, err := readCSV(path)
	if err != nil {
		return nil, err
	}
	cs := make([]*model.CardBin, 0, len(data))

	for i, each := range data {
		// 跳过第一条
		if i == 0 {
			continue
		}
		// BIN,BIN_LEN,INS_CODE,CARD_LEN,ACCT_TYPE,INS_NAME,CARD_NAME,CARD_BRAND
		// 判断该记录的长度是否为8
		if len(each) != 8 {
			return nil, fmt.Errorf("%d 行格式错误，检测到有%d 个字段", i+1, len(each))
		}

		if matched, _ := regexp.MatchString(`^\d+$`, each[0]); !matched {
			return nil, fmt.Errorf("%d 行，bin 应为数字，实际为：%s", i+1, each[0])
		}

		binLen, err := strconv.Atoi(each[1])
		if err != nil {
			return nil, fmt.Errorf("%d 行，binLen 应为数字，实际为：%s", i+1, each[1])
		}

		if matched, _ := regexp.MatchString(`^\d+`, each[2]); !matched {
			return nil, fmt.Errorf("%d 行，insCode 应为数字，实际为：%s", i+1, each[2])
		}

		cardLen, err := strconv.Atoi(each[3])
		if err != nil {
			return nil, fmt.Errorf("%d 行，cardLen 应为数字，实际为：%s", i+1, each[3])
		}
		acctType := ""
		switch each[4] {
		case "1":
			// 表示借记卡
			acctType = "10"
		case "2", "3":
			// 贷记卡\准贷记卡
			acctType = "20"
		case "4":
			// 预付卡
			acctType = "30"
		default:
			return nil, fmt.Errorf("%d 行，账户类型 应为【1，2，3，4】，实际为：%s", i+1, each[4])
		}

		if matched, _ := regexp.MatchString(`^[A-Z]+$`, each[7]); !matched {
			return nil, fmt.Errorf("%d  行，cardBrand  应为大写字母，实际为：%s", i+1, each[7])
		}

		c := &model.CardBin{Bin: each[0], BinLen: binLen, InsName: each[5], CardName: each[6],
			InsCode: strings.TrimSpace(each[2]), CardLen: cardLen, CardBrand: each[7], AcctType: acctType}
		cs = append(cs, c)
	}
	return cs, nil
}

// readCSV 读取文件返回数据
func readCSV(path string) ([][]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	reader.TrimLeadingSpace = true
	reader.Comment = '\n'

	rawCSVdata, err := reader.ReadAll()

	return rawCSVdata, err
}
