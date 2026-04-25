package mailer

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/arinsuda/movie-hub/config"
)

type Mailer struct {
	cfg config.SMTPConfig
}

func New(cfg config.SMTPConfig) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) SendVerificationEmail(toEmail, username, verifyURL string) error {
	subject := "ยืนยันอีเมลของคุณ — Movie Hub"
	body := buildVerifyEmailBody(username, verifyURL)
	return m.send(toEmail, subject, body)
}

func (m *Mailer) send(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)

	headers := strings.Join([]string{
		fmt.Sprintf("From: Movie Hub <%s>", m.cfg.From),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		`Content-Type: text/html; charset="UTF-8"`,
	}, "\r\n")

	msg := []byte(headers + "\r\n\r\n" + htmlBody)
	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)

	if err := smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg); err != nil {
		return fmt.Errorf("mailer: send to %s: %w", to, err)
	}
	return nil
}

func buildVerifyEmailBody(username, verifyURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family:sans-serif;background:#0f0f0f;color:#fff;padding:40px">
  <div style="max-width:480px;margin:auto;background:#1a1a1a;border-radius:12px;padding:32px">
    <h2 style="color:#e50914;margin-bottom:8px">🎬 Movie Hub</h2>
    <p>สวัสดี <strong>%s</strong>,</p>
    <p>คลิกปุ่มด้านล่างเพื่อยืนยันอีเมลของคุณ</p>
    <a href="%s"
       style="display:inline-block;margin:24px 0;padding:12px 28px;
              background:#e50914;color:#fff;border-radius:8px;text-decoration:none;font-weight:bold">
      ยืนยันอีเมล
    </a>
    <p style="color:#888;font-size:12px">
      ลิงก์นี้จะหมดอายุใน 24 ชั่วโมง<br>
      หากคุณไม่ได้สมัครสมาชิก กรุณาเพิกเฉยต่ออีเมลนี้
    </p>
  </div>
</body>
</html>`, username, verifyURL)
}
