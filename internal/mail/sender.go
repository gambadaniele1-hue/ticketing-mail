package mail

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/gambadaniele1-hue/ticketing-mail/internal/queue"
)

type Sender interface {
	Send(job queue.MailJob) error
}

type SMTPConfing struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type SMTPSender struct {
	config SMTPConfing
}

func NewSMTPSender(cfg SMTPConfing) *SMTPSender {
	return &SMTPSender{
		config: cfg,
	}
}

func (s *SMTPSender) Send(job queue.MailJob) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	msg := buildMessage(s.config.From, job)

	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	err := smtp.SendMail(addr, auth, s.config.From, []string{job.To}, []byte(msg))

	if err != nil {
		return fmt.Errorf("invio fallito per %s: %w", job.To, err)
	}

	return nil
}

func buildMessage(from string, job queue.MailJob) string {
	boundary := "==ticketing=="

	var b strings.Builder

	b.WriteString(fmt.Sprintf("From: %s\r\n", from))
	b.WriteString(fmt.Sprintf("To: %s\r\n", job.To))
	b.WriteString(fmt.Sprintf("Subject: %s\r\n", job.Subject))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%q\r\n", boundary))
	b.WriteString("\r\n")

	b.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	b.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	b.WriteString("\r\n")
	b.WriteString(job.Text)
	b.WriteString("\r\n")

	b.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	b.WriteString("Content-Type: text/html; charset=utf-8\r\n")
	b.WriteString("\r\n")
	b.WriteString(job.HTML)
	b.WriteString("\r\n")

	b.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return b.String()
}
