package user_module

import (
	"fmt"
	"net/smtp"
	"os"
)

// PasswordResetMailer ส่ง reset-password link ทาง email
type PasswordResetMailer interface {
	SendResetLink(toEmail, resetURL string) error
}

type smtpPasswordResetMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPPasswordResetMailer() PasswordResetMailer {
	return &smtpPasswordResetMailer{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

func (m *smtpPasswordResetMailer) SendResetLink(toEmail, resetURL string) error {
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	subject := "รีเซ็ตรหัสผ่าน - MovieHub"
	body := fmt.Sprintf(`สวัสดี,

เราได้รับคำขอรีเซ็ตรหัสผ่านบัญชี MovieHub ของคุณ

คลิกลิงก์ด้านล่างเพื่อตั้งรหัสผ่านใหม่:
%s

ลิงก์นี้จะหมดอายุใน 1 ชั่วโมง
หากคุณไม่ได้ขอรีเซ็ตรหัสผ่าน กรุณาเพิกเฉยต่ออีเมลนี้

MovieHub Team`, resetURL)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, toEmail, subject, body,
	)

	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	return smtp.SendMail(addr, auth, m.from, []string{toEmail}, []byte(msg))
}
