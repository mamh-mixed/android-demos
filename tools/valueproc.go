package tools

import (
	"os"
	"strings"

	"github.com/omigo/log"
)

// FirstExistValue 返回第一个不为空的值。values 是多个以 “|” 分隔的值，这些值可以是
// 字符串，也可以是环境变量，环境变量以 “$” 开头。系统会依次判断这些环境变量是否存在，取
// 第一个存在的。如果都不存在，返回 空字符串
func FirstExistValue(values string) string {
	vs := strings.Split(values, "|")
	for _, v := range vs {
		v := strings.TrimSpace(v)
		// 如果以 $ 开始，表示系统环境变量
		if v[0] == '$' {
			v = os.Getenv(v[1:])
			if v != "" {
				return v
			}
		}

		if v != "" {
			return v
		}
	}

	log.Errorf("config value %s not correct", values)
	return ""
}
