package email

import (
	"github.com/omigo/log"
	"gopkg.in/jordan-wright/email.v1"
	"io"
	"net/smtp"
)

const (
	ADDR     = "smtp.qiye.163.com:25"
	HOST     = "smtp.qiye.163.com"
	PASSWORD = "cil1009"
	USER     = "it.support@cardinfolink.com"
)

var auth = smtp.PlainAuth("", USER, PASSWORD, HOST)

type Email struct {
	To      string
	Title   string
	Body    string
	attachs []*attachment
}

type attachment struct {
	r        io.Reader
	filename string
	c        string
}

func (e *Email) Attach(r io.Reader, filename, contentType string) {
	e.attachs = append(e.attachs, &attachment{
		r:        r,
		filename: filename,
		c:        contentType,
	})
}

// Send 发送邮件
func (e *Email) Send() error {

	em := email.NewEmail()

	// basic
	em.To = []string{e.To}
	em.From = USER
	em.Subject = e.Title
	// em.Cc

	em.HTML = []byte(e.Body) // Content-Type: text/html

	// add attachment
	for _, a := range e.attachs {
		em.Attach(a.r, a.filename, a.c)
	}

	err := em.Send(ADDR, auth)
	if err != nil {
		log.Errorf("TO %s, send email fail: %s", e.To, err)
	}

	return err
}
