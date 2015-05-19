// 渠道csv文件字段统一的格式顺序如下
// code,   msg,   respCode,     respMsg
// 渠道代码,返回信息,对应系统代码,对应系统返回信息
package data

import (
	"encoding/csv"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"os"
	"regexp"
	"strconv"
)

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

// AddCardBinFromCsv 从csv里导入卡bin
// rebuild: true 删除集合再重建
// rebuild: false 做更新操作，即存在的更新，不存在的增加
func AddCardBinFromCsv(path string, rebuild bool) error {

	cardBins, err := ReadCardBinCsv(path)
	if err != nil {
		return err
	}
	// 重建
	if rebuild {
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
	}

	return nil

}

// AddFromCsv 从csv文件里读取应答码表
func AddSysCodeFromCsv(path string) error {

	data, err := ReadQuickpayCsv(path)
	if err != nil {
		return err
	}
	// 添加到mongodb，若存在的跳过
	// 若新增的便添加
	for _, v := range data {
		_, err := mongo.RespCodeColl.FindOne(v.RespCode)
		if err != nil {
			fmt.Println(v)
			mongo.RespCodeColl.Add(v)
		}
	}
	return nil
}

func AddChanCodeFromScv(channel, path string) error {

	data, err := ReadChanCsv(path)
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
			q.Cfca = append(q.Cfca, v)
		case channel == "cil":
			q.Cil = append(q.Cil, v)
		default:
			// ...更多渠道
		}
		mongo.RespCodeColl.Update(q)
	}
	return nil
}

// ReadQuickpayCsv 读取系统应答码csv文件
// 并持久化
func ReadQuickpayCsv(path string) ([]*model.QuickpayCsv, error) {

	data, err := readCsv(path)
	if err != nil {
		return nil, err
	}
	qs := make([]*model.QuickpayCsv, 0, len(data))

	// 根据数据规则遍历
	for i, each := range data {
		if i == 0 {
			continue
		}
		q := &model.QuickpayCsv{RespCode: each[0], RespMsg: each[1]}
		// fmt.Printf("%+v \n", q)
		qs = append(qs, q)
	}
	fmt.Println(len(qs))
	return qs, nil
}

// ReadChanCsv 读取渠道应答码文件
func ReadChanCsv(path string) ([]*model.ChanCsv, error) {
	data, err := readCsv(path)
	if err != nil {
		return nil, err
	}
	qs := make([]*model.ChanCsv, 0, len(data))

	// 根据渠道应答码文件规则遍历
	for i, each := range data {
		if i == 0 {
			continue
		}
		q := &model.ChanCsv{each[0], each[1], each[2], each[3]}
		// fmt.Printf("%+v \n", q)
		qs = append(qs, q)
	}
	return qs, nil
}

// ReadCardBinCsv 从csv读取卡bin转为对象
func ReadCardBinCsv(path string) ([]*model.CardBin, error) {

	data, err := readCsv(path)
	if err != nil {
		return nil, err
	}
	cs := make([]*model.CardBin, 0, len(data))

	for i, each := range data {
		// 跳过第一条
		if i == 0 {
			continue
		}
		// 判断该记录的长度是否为5
		if len(each) != 5 {
			return nil, fmt.Errorf("%d行格式错误，检测到有%d个字段", i+1, len(each))
		}

		if matched, _ := regexp.MatchString(`^\d+$`, each[0]); !matched {
			return nil, fmt.Errorf("%d行，bin应为数字，实际为：%s", i+1, each[0])
		}

		binLen, err := strconv.Atoi(each[1])
		if err != nil {
			return nil, fmt.Errorf("%d行，binLen应为数字，实际为：%s", i+1, each[1])
		}

		if matched, _ := regexp.MatchString(`^\d+$`, each[2]); !matched {
			return nil, fmt.Errorf("%d行，insCode应为数字，实际为：%s", i+1, each[2])
		}

		cardLen, err := strconv.Atoi(each[3])
		if err != nil {
			return nil, fmt.Errorf("%d行，cardLen应为数字，实际为：%s", i+1, each[3])
		}

		if matched, _ := regexp.MatchString(`^[A-Z]+$`, each[4]); !matched {
			return nil, fmt.Errorf("%d行，cardBrand应为大写字母，实际为：%s", i+1, each[4])
		}

		c := &model.CardBin{Bin: each[0], BinLen: binLen,
			InsCode: each[2], CardLen: cardLen, CardBrand: each[4]}
		cs = append(cs, c)
	}
	return cs, nil
}

// readCsv 读取文件返回数据
func readCsv(path string) ([][]string, error) {

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
