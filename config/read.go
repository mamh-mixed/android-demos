package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/omigo/log"
)

func init() {
	env := "develop"

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

	c, err := goconfig.LoadConfigFile(fileName)
	if err != nil {
		log.Fatalf("can not load config file(%s): %s", fileName, err)
	}

	mongoHost, err := c.GetValue("mongo", "host")
	if err != nil {
		log.Errorf("can not get value from selection `%s` on key `%s`", "mongo", "host")
	}
	log.Debugf("%s.%s = %s", "mongo", "host", mongoHost)

}
