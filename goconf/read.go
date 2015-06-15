// Package goconf 配置集中管理, 使用前，需要在系统中配置 QUICKPAY_ENV 环境变量，
// 变量值为 develop/testing/product，如果不配置，系统会报错
// export ANGRYCARD_ENV=develop
package goconf

import (
	"fmt"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/omigo/log"
)

const (
	sysEnv = "QUICKPAY_ENV"
	pkg    = "github.com/CardInfoLink/quickpay"
)

var conf *goconfig.ConfigFile
var workDir string

func init() {
	env := os.Getenv(sysEnv)
	if env == "" {
		fmt.Printf("system environment variable `%s` not set, must set to `develop` or `testing` or `product`\n", sysEnv)
		os.Exit(1)
	}
	fmt.Printf("environment: %s\n", env)

	var err error
	workDir, err = os.Getwd()
	if err != nil {
		log.Fatalf("can not get work directory: %s", err)
	}
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

// Hostname 取主机名，如果没取到，返回 `unknown`
func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		log.Errorf("get hostname error: %s", name)
		return "unknown"
	}

	return name
}

// GetFile 从配置文件中读取文件全名，包含绝对路径
func GetFile(section, key string) (filename string) {
	v := GetValue(section, key)
	if v == "" {
		return ""
	}

	// 如果配置的是绝对路径，直接返回
	if v[0] == '/' {
		return v
	}

	// 如果配置的是相对路径，那么就是以启动程序的目录为相对路径
	return GetWorkDir() + "/" + v
}

// GetValue 从配置文件中取值
func GetValue(section, key string) (v string) {
	v, err := conf.GetValue(section, key)
	if err != nil {
		log.Errorf("can not get value from selection `%s` on key `%s`", section, key)
	}
	// log.Debugf("%s.%s = %s", section, key, v)

	return v
}

// LogLevel 从配置文件中取日志级别
func LogLevel() (v int) {
	l, err := conf.GetValue("app", "logLevel")
	if err != nil {
		log.Error("can not get value from selection `app` on key `logLevel`")
	}
	if l == "" {
		return log.Linfo
	}

	switch strings.ToLower(l) {
	case "trace":
		return log.Ltrace
	case "debug":
		return log.Ldebug
	case "info":
		return log.Linfo
	case "warn":
		return log.Lwarn
	case "error":
		return log.Lerror
	case "fatal":
		return log.Lfatal
	default:
		return log.Linfo
	}
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
