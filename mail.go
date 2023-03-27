package main

import "net/smtp"

type mailSender struct {
	id, acc, pwd, host string
}

type list struct {
}

func newMailSender(id, acc, pwd, host string) *mailSender {
	return &mailSender{
		id:   id,
		acc:  acc,
		pwd:  pwd,
		host: host,
	}
}

func (m *mailSender) send(msg string) error {
	auth := smtp.PlainAuth(m.id, m.acc, m.pwd, m.host)

	return smtp.SendMail(m.host+":587", auth,
		m.acc, conf.MailList, []byte(msg))
}
