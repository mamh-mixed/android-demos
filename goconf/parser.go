package goconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/CardInfoLink/quickpay/tools"
)

func init() {
	fileName := fmt.Sprintf("%s/config/config_%s.json", tools.WorkDir, tools.Env)
	fmt.Printf("config file:\t %s\n", fileName)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("read config file %s error: %s\n", fileName, err)
		os.Exit(4)
	}

	err = json.Unmarshal(content, &Config)
	if err != nil {
		fmt.Printf("config file %s parser error: %s\n", fileName, err)
		os.Exit(5)
	}

	Config.postProcess()
}
