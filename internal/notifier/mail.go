package notifier

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	host string
	addr string
	auth smtp.Auth
	from string
	to   string
}

func NewMailer(host, port, user, pass, from, to string) *Mailer {
	return &Mailer{
		host: host,
		addr: host + ":" + port,
		auth: smtp.PlainAuth("", user, pass, host),
		from: from,
		to:   to,
	}
}

func (m *Mailer) Send(subject, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, m.to, subject, body)

	return smtp.SendMail(m.addr, m.auth, m.from, []string{m.to}, []byte(msg))
}
