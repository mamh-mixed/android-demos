// Package config 用户集中配置管理。使用前，需要在系统中配置 QUICKPAY_ENV 环境变量，
// 变量值为 develop/testing/product，如果不配置，默认 develop
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/omigo/log"
)

var configFile *goconfig.ConfigFile

func loadConfigFile() {
	sysEnv := "QUICKPAY_ENV"
	env := os.Getenv(sysEnv)
	if env == "" {
		log.Warnf("system env `%s` not set, use `develop` instead", sysEnv)
		env = "develop"
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("can not get work directory: %s", err)
	}
	pkg := "github.com/CardInfoLink/quickpay"
	if pos := strings.Index(wd, pkg); pos >= 0 {
		wd = wd[:(pos + len(pkg))]
	}

	log.Debugf("work directory: %s", wd)
	fileName := fmt.Sprintf("%s/config/config_%s.ini", wd, env)

	configFile, err = goconfig.LoadConfigFile(fileName)
	if err != nil {
		log.Fatalf("can not load config file(%s): %s", fileName, err)
	}
}

// GetValue 从配置文件中取值
func GetValue(section, key string) (v string) {
	if configFile == nil {
		loadConfigFile()
	}

	v, err := configFile.GetValue(section, key)
	if err != nil {
		log.Errorf("can not get value from selection `%s` on key `%s`", section, key)
	}
	log.Debugf("%s.%s = %s", section, key, v)

	return v
}

// Int 从配置文件中整数值
func Int(section, key string) (v int) {
	if configFile == nil {
		loadConfigFile()
	}

	v, err := configFile.Int(section, key)
	if err != nil {
		log.Errorf("can not get value from selection `%s` on key `%s`", section, key)
	}
	log.Debugf("%s.%s = %d", section, key, v)

	return v
}
