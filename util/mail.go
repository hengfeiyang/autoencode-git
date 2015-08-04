package util

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

type MailT struct {
	Addr  string // smtp.163.com:25
	User  string // robot@163.com
	Pass  string // 12345
	From  string // 发件人
	To    string // 收件人
	Title string // 邮件标题
	Body  string // 邮件正文
	Type  string // 内容类型，纯文本plain或网页html
}

// 发送邮件
func SendMail(conf *MailT) error {
	encode := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	host := strings.Split(conf.Addr, ":")
	auth := smtp.PlainAuth("", conf.User, conf.Pass, host[0])
	var contentType string
	if conf.Type == "html" {
		contentType = "text/html; charset=UTF-8"
	} else {
		contentType = "text/plain; charset=UTF-8"
	}
	from := mail.Address{"9466代码中心", conf.From}
	to := mail.Address{"", conf.To}
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = conf.Title
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = contentType
	header["Content-Transfer-Encoding"] = "base64"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + encode.EncodeToString([]byte(conf.Body))
	err := smtp.SendMail(
		conf.Addr,
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
	return err
}
