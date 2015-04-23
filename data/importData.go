package data

import (
	"encoding/csv"
	"fmt"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"os"
)

// AddFromCsv 从csv文件里读取应答码表
func AddSysCodeFromCsv(path string) error {

	data, err := ReadQuickpayCsv(path)
	if err != nil {
		return err
	}

	// 添加到mongodb，若存在的跳过
	// 若新增的便添加
	for _, v := range data {
		q, _ := mongo.RespCodeColl.FindOne(v.RespCode)
		if q == nil {
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
	qs := make([]*model.QuickpayCsv, len(data))

	// 根据数据规则遍历
	for i, each := range data {
		if i == 0 {
			continue
		}
		q := &model.QuickpayCsv{RespCode: each[0], RespMsg: each[1]}
		// fmt.Printf("%+v \n", q)
		qs = append(qs, q)
	}

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
