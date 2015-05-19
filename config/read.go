// Package config 用户集中配置管理。使用前，需要在系统中配置 QUICKPAY_ENV 环境变量，
// 变量值为 develop/testing/product，如果不配置，系统会报错
// export QUICKPAY_ENV=develop
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/omigo/log"
)

var conf *goconfig.ConfigFile
var workDir string

func init() {
	sysEnv := "QUICKPAY_ENV"
	env := os.Getenv(sysEnv)
	if env == "" {
		fmt.Printf("system environment variable `%s` not set, must set to `develop/testing/product`\n", sysEnv)
		os.Exit(1)
	}
	fmt.Printf("quickpay environment: %s\n", env)

	var err error
	workDir, err = os.Getwd()
	if err != nil {
		log.Fatalf("can not get work directory: %s", err)
	}
	pkg := "github.com/CardInfoLink/quickpay"
	if pos := strings.Index(workDir, pkg); pos >= 0 {
		workDir = workDir[:(pos + len(pkg))]
	}

	fmt.Printf("work directory: %s\n", workDir)
	fileName := fmt.Sprintf("%s/config/config_%s.ini", workDir, env)

	conf, err = goconfig.LoadConfigFile(fileName)
	if err != nil {
		log.Fatalf("can not load config file(%s): %s", fileName, err)
	}

	// print all configurations
	sects := conf.GetSectionList()
	for _, sect := range sects {
		fmt.Printf("\n[%s]\n", sect)
		for _, key := range conf.GetKeyList(sect) {
			value, err := conf.GetValue(sect, key)
			if err != nil {
				log.Errorf("read config error in section %s and key %s: %s", sect, key, err)
			}
			fmt.Printf("%-20s = %s\n", key, value)
		}
	}
	fmt.Println()
}

// GetWorkDir 获取程序启动目录
func GetWorkDir() string {
	return workDir
}

// GetFile 从配置文件中读取文件全名，包含绝对路径
func GetFile(section, key string) (filename string, err error) {
	v := GetValue(section, key)
	if v == "" {
		return "", errors.New("selection %s key %s not configurated")
	}

	// 如果配置的是绝对路径，直接返回
	if v[0] == '/' {
		return v, nil
	}

	// 如果配置的是相对路径，那么就是以启动程序的目录为相对路径
	return GetWorkDir() + "/" + v, nil
}

// GetValue 从配置文件中取值
func GetValue(section, key string) (v string) {
	v, err := conf.GetValue(section, key)
	if err != nil {
		log.Errorf("can not get value from selection `%s` on key `%s`", section, key)
	}
	// log.Debugf("%s.%s = %s", section, key, v)
	if v == "" {
		log.Warnf("selection `%s` key `%s` value is blank", section, key)
	}

	return v
}

// Int 从配置文件中整数值
func Int(section, key string) (v int) {
	v, err := conf.Int(section, key)
	if err != nil {
		log.Errorf("can not get value from selection `%s` on key `%s`", section, key)
	}
	// log.Debugf("%s.%s = %d", section, key, v)

	return v
}
