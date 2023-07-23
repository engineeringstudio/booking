package internel

import "net/smtp"

type mailSender struct {
	id, acc, pwd, host string
	mailList           []string
}

type list struct {
}

func NewMailSender(mailList []string, id, acc, pwd, host string) *mailSender {
	return &mailSender{
		id:       id,
		acc:      acc,
		pwd:      pwd,
		host:     host,
		mailList: mailList,
	}
}

func (m *mailSender) send(msg string) error {
	auth := smtp.PlainAuth(m.id, m.acc, m.pwd, m.host)

	return smtp.SendMail(m.host+":587", auth,
		m.acc, m.mailList, []byte(msg))
}
