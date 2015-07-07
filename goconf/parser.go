package goconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
)

func init() {
	fileName := fmt.Sprintf("%s/config/config_%s.js", util.WorkDir, util.Env)
	fmt.Printf("config file:\t %s\n", fileName)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("read config file %s error: %s\n", fileName, err)
		os.Exit(4)
	}

	// 删除注释
	re := regexp.MustCompile(`\s*//\s.+|/\*.+?\*/`)
	content = re.ReplaceAll(content, []byte(""))
	log.Debugf("config content: %s", content)

	err = json.Unmarshal(content, &Config)
	if err != nil {
		fmt.Printf("config file %s parser error: %s\n", fileName, err)
		os.Exit(5)
	}

	Config.postProcess()
}
