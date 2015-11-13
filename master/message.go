package master

import (
	"encoding/json"
	"fmt"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var avaiableLocale = make(map[string]*LocaleTemplate)

const (
	prefix        = "message_"
	DefaultLocale = "zh_CN"
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
	// 统计报表
	StatReport struct {
		Title      string
		Total      string
		StartDate  string
		EndDate    string
		Remark     string
		MerId      string
		MerName    string
		Summary    string
		ALP        string
		WXP        string
		AgentName  string
		TotalCount string
		TotalAmt   string
		TotalFee   string
		Count      string
		Amt        string
		Fee        string
	}
	// TODO...
}
