package utils

import (
	"encoding/base64"
	cf "goBoss/config"
	"gopkg.in/gomail.v2"
	"log"
)

type Mail struct {
	Content string
	Subject string
	Attach  string
}

func (m *Mail) Send() {
	handle := gomail.NewMessage()
	handle.SetHeader("From", cf.Config.Sender)
	handle.SetHeader("To", cf.Config.Receiver)
	handle.SetHeader("Subject", "[Auto]:来自goBoss--"+m.Subject)
	handle.SetBody("text/html", m.Content)
	if m.Attach != "" {
		handle.Attach(m.Attach)
	}
	s := gomail.NewDialer(cf.Config.MailServer, cf.Config.MailPort, cf.Config.Sender, cf.Config.SenderPwd)
	if err := s.DialAndSend(handle); err != nil {
		log.Println("发送邮件失败!Error: ", err.Error())
	}
}

func Encode(enc *base64.Encoding, bt []byte) string {
	// 编码
	encStr := enc.EncodeToString(bt)
	return encStr
}
