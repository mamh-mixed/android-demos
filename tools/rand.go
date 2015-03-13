package tools

import (
	"fmt"
	u "github.com/nu7hatch/gouuid"
	"github.com/omigo/g"
)

// serialNumber 生成序列号，也就是UUID
func SerialNumber() string {
	u4, err := u.NewV4()
	if err != nil {
		g.Error("error: ", err)
		return ""
	}
	return fmt.Sprintf("%x", u4[:])
}
