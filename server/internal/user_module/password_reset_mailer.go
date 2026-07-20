package user_module

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

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
<body style="margin:0;padding:0;background:#050505;font-family:'Helvetica Neue',Arial,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background:#050505;padding:56px 16px;">
    <tr>
      <td align="center">
        <table role="presentation" width="440" cellpadding="0" cellspacing="0" style="max-width:440px;width:100%%;">

          <!-- Wordmark -->
          <tr>
            <td align="center" style="padding-bottom:32px;">
              <span style="font-size:18px;font-weight:700;letter-spacing:4px;color:#ffffff;">REMOV</span>
              <div style="width:32px;height:2px;background:#e50914;margin:12px auto 0 auto;"></div>
            </td>
          </tr>

          <!-- Card -->
          <tr>
            <td style="background:#0f0f0f;border:1px solid #1c1c1c;border-radius:4px;">
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0">
                <tr>
                  <td style="padding:44px 40px 4px 40px;" align="center">
                    <p style="margin:0 0 12px 0;color:#e50914;font-size:11px;font-weight:700;letter-spacing:2.5px;
                              text-transform:uppercase;">
                      ความปลอดภัยของบัญชี
                    </p>
                    <p style="margin:0 0 32px 0;color:#8f8f8f;font-size:14px;line-height:1.8;max-width:300px;">
                      เราได้รับคำขอตั้งรหัสผ่านใหม่สำหรับบัญชีของคุณ กดปุ่มด้านล่างเพื่อดำเนินการต่อ
                    </p>
                  </td>
                </tr>
                <tr>
                  <td align="center" style="padding:4px 40px 0 40px;">
                    <table role="presentation" cellpadding="0" cellspacing="0">
                      <tr>
                        <td style="background:#e50914;border-radius:6px;">
                          <a href="%s"
                             style="display:block;padding:15px 44px;color:#ffffff;font-size:13px;
                                    font-weight:700;letter-spacing:1.5px;text-transform:uppercase;
                                    text-decoration:none;">
                            ตั้งรหัสผ่านใหม่
                          </a>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>
                <tr>
                  <td align="center" style="padding:22px 40px 8px 40px;">
                    <p style="margin:0;color:#666666;font-size:12px;">
                      ลิงก์มีอายุการใช้งาน 1 ชั่วโมง
                    </p>
                  </td>
                </tr>
                <tr>
                  <td style="padding:12px 40px 40px 40px;">
                    <p style="margin:0;color:#4a4a4a;font-size:11px;line-height:1.7;text-align:center;word-break:break-all;">
                      หรือคัดลอกลิงก์นี้ไปวางในเบราว์เซอร์<br>
                      <a href="%s" style="color:#e5504f;text-decoration:underline;">%s</a>
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Footnote -->
          <tr>
            <td style="padding:24px 8px 0 8px;">
              <p style="margin:0;color:#555555;font-size:11.5px;line-height:1.7;text-align:center;">
                หากคุณไม่ได้ขอรีเซ็ตรหัสผ่าน กรุณาเพิกเฉยต่ออีเมลฉบับนี้ รหัสผ่านเดิมของคุณจะยังคงใช้งานได้ตามปกติ
              </p>
            </td>
          </tr>
          <tr>
            <td align="center" style="padding:28px 0 0 0;">
              <p style="margin:0;color:#3a3a3a;font-size:10.5px;letter-spacing:0.5px;">
                REMOV — ระบบส่งอัตโนมัติ กรุณาอย่าตอบกลับอีเมลนี้
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
