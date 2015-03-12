package tools

import (
	u "github.com/nu7hatch/gouuid"
	"github.com/omigo/g"
)

// UUID 产生 UUID
func UUID() string {
	u4, err := u.NewV4()
	if err != nil {
		g.Error("error: ", err)
		return ""
	}
	return u4.String()
}
