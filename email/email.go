package email

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/omigo/log"
	"io"
	"net/smtp"
)

const (
	ADDR     = "smtp.qiye.163.com:25"
	HOST     = "smtp.qiye.163.com"
	PASSWORD = "cil1009"
	USER     = "it.support@cardinfolink.com"
	AndyLi   = "andy.li@cardinfolink.com"
)

var auth = smtp.PlainAuth("", USER, PASSWORD, HOST)
var b64Encoding = base64.StdEncoding

type Email struct {
	To     string
	ActUrl string
}

func (e *Email) Send() error {

	var iv [20]byte
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		log.Errorf("io.ReadFull error: %s", err)
	}

	body := fmt.Sprintf(template, e.ActUrl, b64Encoding.EncodeToString(iv[:]))
	msg := "To: " + e.To + "\r\nFrom: " + USER + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body
	log.Infof("SendBody: %s", msg)

	// send
	err := smtp.SendMail(ADDR, auth, USER, []string{e.To}, []byte(msg))
	if err != nil {
		log.Errorf("TO %s, send email fail: %s", e.To, err)
	}
	return err
}

var (
	subject     = "云收银帐号激活"
	contentType = "Content-Type: text/html; charset=UTF-8"

	actTemplate = `	
	<html>
		<body>
			<h3>
				点击以下链接以激活账户</br>
				<a href="%s">%s</a>
			</h3>
		</body>
	</html>
	`

	promoteTemplate = `
	<html>
		<body>
			<h3>
				Hello,Andy.li:
			</h3>
		</body>
	</html>
	`
)
