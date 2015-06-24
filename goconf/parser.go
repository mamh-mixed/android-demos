package goconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/omigo/log"
)

const (
	sysEnv = "QUICKPAY_ENV"
	pkg    = "github.com/CardInfoLink/quickpay"
)

func init() {
	env := os.Getenv(sysEnv)
	if env == "" {
		fmt.Printf("system environment variable `%s` not set, must set to `develop` or `testing` or `product`\n", sysEnv)
		os.Exit(1)
	}
	fmt.Printf("environment:\t %s\n", env)

	var err error
	workDir, err = os.Getwd()
	if err != nil {
		fmt.Printf("can not get work directory: %s\n", err)
		os.Exit(2)
	}
	if pos := strings.Index(workDir, pkg); pos >= 0 {
		workDir = workDir[:(pos + len(pkg))]
	}

	fmt.Printf("work directory:\t %s\n", workDir)
	fileName := fmt.Sprintf("%s/config/config_%s.json", workDir, env)

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

}

var workDir string

// GetWorkDir 获取程序启动目录
func GetWorkDir() string {
	return workDir
}

// Hostname 取主机名，如果没取到，返回 `unknown`
func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		log.Errorf("get hostname error: %s", name)
		return "unknown"
	}

	return name
}
