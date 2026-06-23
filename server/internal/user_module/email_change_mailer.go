package user_module

import (
	"fmt"
	"net/smtp"
	"os"
)

// Mailer interface ทำให้ test ง่าย (mock ได้)
type Mailer interface {
	SendOTP(toEmail, otp string) error
}

type smtpMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPMailer() Mailer {
	return &smtpMailer{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

func (m *smtpMailer) SendOTP(toEmail, otp string) error {
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	subject := "รหัส OTP สำหรับเปลี่ยน Email - MovieHub"
	body := fmt.Sprintf(`สวัสดี,

คุณได้ขอเปลี่ยน Email บัญชี MovieHub ของคุณ

รหัส OTP ของคุณคือ: %s

รหัสนี้จะหมดอายุใน 15 นาที
หากคุณไม่ได้ขอเปลี่ยน Email กรุณาเพิกเฉยต่ออีเมลนี้

MovieHub Team`, otp)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, toEmail, subject, body,
	)

	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	return smtp.SendMail(addr, auth, m.from, []string{toEmail}, []byte(msg))
}
