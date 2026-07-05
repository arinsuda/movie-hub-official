package user_module

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
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

	subject := "รีเซ็ตรหัสผ่านบัญชีของคุณ — REMOV"
	htmlBody := buildResetPasswordEmailBody(resetURL)

	headers := strings.Join([]string{
		fmt.Sprintf("From: REMOV <%s>", m.from),
		fmt.Sprintf("To: %s", toEmail),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		`Content-Type: text/html; charset="UTF-8"`,
	}, "\r\n")

	msg := []byte(headers + "\r\n\r\n" + htmlBody)
	addr := fmt.Sprintf("%s:%s", m.host, m.port)

	if err := smtp.SendMail(addr, auth, m.from, []string{toEmail}, msg); err != nil {
		return fmt.Errorf("mailer: send reset link to %s: %w", toEmail, err)
	}
	return nil
}

func buildResetPasswordEmailBody(resetURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin:0;padding:0;background:#0a0a0a;font-family:'Helvetica Neue',Arial,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background:#0a0a0a;padding:48px 16px;">
    <tr>
      <td align="center">
        <table role="presentation" width="480" cellpadding="0" cellspacing="0"
               style="max-width:480px;width:100%%;background:linear-gradient(160deg,#1a1a1a 0%%,#141414 100%%);
                      border-radius:16px;overflow:hidden;border:1px solid #2a2a2a;
                      box-shadow:0 20px 60px rgba(229,9,20,0.15);">
          <tr>
            <td style="background:linear-gradient(90deg,#e50914,#b0060f);padding:28px 32px;">
              <h1 style="margin:0;color:#fff;font-size:22px;letter-spacing:1px;">🎬 REMOV</h1>
            </td>
          </tr>
          <tr>
            <td style="padding:40px 36px 8px 36px;">
              <p style="margin:0 0 4px 0;color:#8a8a8a;font-size:13px;letter-spacing:2px;text-transform:uppercase;">
                ความปลอดภัยของบัญชี
              </p>
              <h2 style="margin:0 0 16px 0;color:#fff;font-size:24px;line-height:1.4;">
                รีเซ็ตรหัสผ่านของคุณ
              </h2>
              <p style="margin:0 0 32px 0;color:#b3b3b3;font-size:15px;line-height:1.7;">
                เราได้รับคำขอรีเซ็ตรหัสผ่านสำหรับบัญชี REMOV ของคุณ
                กดปุ่มด้านล่างเพื่อตั้งรหัสผ่านใหม่
              </p>
            </td>
          </tr>
          <tr>
            <td align="center" style="padding:0 36px 8px 36px;">
              <a href="%s"
                 style="display:inline-block;width:100%%;max-width:340px;box-sizing:border-box;
                        padding:16px 24px;background:linear-gradient(90deg,#e50914,#f6121d);
                        color:#ffffff;border-radius:10px;text-decoration:none;font-weight:700;
                        font-size:15px;letter-spacing:0.4px;text-align:center;
                        box-shadow:0 8px 24px rgba(229,9,20,0.35);">
                ตั้งรหัสผ่านใหม่
              </a>
            </td>
          </tr>
          <tr>
            <td style="padding:24px 36px 8px 36px;">
              <p style="margin:0;color:#8a8a8a;font-size:13px;line-height:1.6;">
                ⏱️ ลิงก์นี้จะหมดอายุใน <strong style="color:#e5e5e5;">1 ชั่วโมง</strong>
              </p>
            </td>
          </tr>
          <tr>
            <td style="padding:16px 36px 8px 36px;">
              <p style="margin:0;color:#5c5c5c;font-size:12px;line-height:1.6;word-break:break-all;">
                หากปุ่มด้านบนใช้งานไม่ได้ ให้คัดลอกลิงก์นี้ไปวางในเบราว์เซอร์:<br>
                <a href="%s" style="color:#e5504f;text-decoration:underline;">%s</a>
              </p>
            </td>
          </tr>
          <tr>
            <td style="padding:20px 36px 36px 36px;">
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0"
                     style="background:#1f1a1a;border-left:3px solid #e50914;border-radius:6px;">
                <tr>
                  <td style="padding:14px 16px;">
                    <p style="margin:0;color:#c9a3a3;font-size:12.5px;line-height:1.6;">
                      หากคุณไม่ได้ขอรีเซ็ตรหัสผ่าน กรุณาเพิกเฉยต่ออีเมลฉบับนี้
                      รหัสผ่านของคุณจะยังคงเดิมและไม่มีการเปลี่ยนแปลงใดๆ
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:20px 36px 32px 36px;border-top:1px solid #232323;">
              <p style="margin:0;color:#5c5c5c;font-size:11.5px;text-align:center;">
                © REMOV — อีเมลนี้ถูกส่งโดยระบบอัตโนมัติ กรุณาอย่าตอบกลับ
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, resetURL, resetURL, resetURL)
}
