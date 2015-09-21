package email

import (
	"fmt"
	"testing"
	"time"
)

var (
	title = "测试"
	body  = `	
	<html>
		<body>
			<h3>
				点击以下链接以激活账户</br>
				<a href="%s">%s</a>
			</h3>
		</body>
	</html>
	`
)

func TestSend(t *testing.T) {
	now := time.Now()
	e := Email{To: "andy.li@cardinfolink.com", Title: title, Body: fmt.Sprintf(body, "https://www.baidu.com", "click me")}
	err := e.Send()
	if err != nil {
		t.Error(err)
	}
	after := time.Now()
	t.Log(after.Sub(now))
}
