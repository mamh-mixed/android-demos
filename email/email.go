package email

import (
	"github.com/omigo/log"
	"net/smtp"
)

const (
	ADDR     = "smtp.qiye.163.com:25"
	HOST     = "smtp.qiye.163.com"
	PASSWORD = "cil1009"
	USER     = "it.support@cardinfolink.com"
)

var auth = smtp.PlainAuth("", USER, PASSWORD, HOST)
var contentType = "Content-Type: text/html; charset=UTF-8"

type Email struct {
	To    string
	Title string
	Body  string
}

// Send 发送邮件
func (e *Email) Send() error {

	// var iv [20]byte
	// if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
	// 	log.Errorf("io.ReadFull error: %s", err)
	// }
	// body := fmt.Sprintf(template, e.ActUrl, b64Encoding.EncodeToString(iv[:]))
	msg := "To: " + e.To + "\r\nFrom: " + USER + ">\r\nSubject: " + e.Title + "\r\n" + contentType + "\r\n\r\n" + e.Body
	log.Infof("SendBody: %s", msg)
	// send
	err := smtp.SendMail(ADDR, auth, USER, []string{e.To}, []byte(msg))
	if err != nil {
		log.Errorf("TO %s, send email fail: %s", e.To, err)
	}
	return err
}
