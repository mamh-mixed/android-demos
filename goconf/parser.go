package goconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/CardInfoLink/quickpay/util"
)

func init() {
	fileName := fmt.Sprintf("%s/config/config_%s.js", util.WorkDir, util.Env)
	fmt.Printf("config file:\t %s\n", fileName)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("read config file %s error: %s\n", fileName, err)
		os.Exit(4)
	}

	// 删除注释，为避免删除非注释文本，配置文件有如下规定:
	// 1. 单行注释文本前必须加一个空字符，与 http://example.com 区分开
	// 2. 多行注释文本前后都必须加一个空格
	re := regexp.MustCompile(`\s*//\s.+|\s*/\*\s[\S\s]+?\s\*/`)
	content = re.ReplaceAll(content, []byte(""))
	// fmt.Println("config content: %s\n", content)

	err = json.Unmarshal(content, &Config)
	if err != nil {
		fmt.Printf("config file %s parser error: %s\n", fileName, err)
		os.Exit(5)
	}

	Config.postProcess()
}
