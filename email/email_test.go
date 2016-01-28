package email

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/SKatiyar/qr"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
	openTemplate = `
	<html>
		<body>
			您好，您申请的商户参数信息如下：
			<h3>%s</h3>
			<h4>app登录信息：</h4>
			注册邮箱：%s<br>
			<h4>桌面版信息：</h4>
			商户号： 	%s<br>
			密钥：	%s
			<h4>网页版信息：</h4>
			<div id="code">
			<img src="data:image/png;base64,%s"/>
			</div>
		</body>
	</html>
	`
)

func TestSend(t *testing.T) {
	now := time.Now()
	e := Email{To: "379630413@qq.com", Title: title, Body: fmt.Sprintf(body, "https://www.baidu.com", "click me")}

	// file, err := os.Open("/Users/zhiruichen/Desktop/test.zip")

	var keys = []string{"sett/report/20151016/IC202_WXP_20151016.xlsx", "sett/report/20151019/IC202_WXP_20151019.xlsx"}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for _, k := range keys {
		resp, err := http.Get(qiniu.MakePrivateUrl(k))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		ebytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Log(len(ebytes))
		f, err := w.Create(k)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		_, err = f.Write(ebytes)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}

	err := w.Close()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// t.Log(buf.Len())
	e.Attach(buf, "测试.zip", "application/octet-stream; charset=UTF-8") //application/octet-stream; charset=UTF-8

	err = e.Send()
	if err != nil {
		t.Error(err)
	}
	after := time.Now()
	t.Log(after.Sub(now))
}

func TestSendOpen(t *testing.T) {
	code, err := qr.Encode("awdawdawdadawwaawdawdawdawda", qr.Q)
	png64 := base64.StdEncoding.EncodeToString(code.PNG())
	e := Email{To: "379630413@qq.com", Title: "【欢迎注册云收银】", Body: fmt.Sprintf(openTemplate, "金鸿洗餐厅", "379630413@qq.com", "999118880000001", "dajshdjasjbadasdasd", png64)}
	err = e.Send()
	if err != nil {
		t.Error(err)
	}
}

func TestSendTranscationFile(t *testing.T) {
	e := Email{
		To:    "379630413@qq.com",
		Title: "QR Payment Channel Files",
		Body: `
		<html>
		<body>
		Dear NTTDATA,<br>
		Please find the channel files in attachment.<br>
		Email is sent by system automatically, please do not reply this email.<br>
		Thanks
		</body>
	</html>
	`,
	}

	f, err := os.Open("/Users/zhiruichen/Desktop/201501_trans.xlsx")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	dst := bytes.NewBuffer([]byte{})
	io.Copy(dst, f)

	e.Attach(dst, f.Name(), "")
	err = e.Send()
	if err != nil {
		t.Error(err)
	}

	// t.Logf("%d", f.)
}
