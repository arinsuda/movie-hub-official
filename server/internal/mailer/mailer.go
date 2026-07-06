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
	subject := "ยืนยันอีเมลของคุณ — REMOV"
	body := buildVerifyEmailBody(username, verifyURL)
	return m.send(toEmail, subject, body)
}

func (m *Mailer) send(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)

	headers := strings.Join([]string{
		fmt.Sprintf("From: REMOV <%s>", m.cfg.From),
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

// buildVerifyEmailBody: ปรับให้ใช้โครง card เดียวกับ email_change_mailer.go
// (wordmark REMOV + เส้นแดงบาง, dark card #0f0f0f border #1c1c1c) เพื่อความสม่ำเสมอของทุกอีเมลในระบบ
func buildVerifyEmailBody(username, verifyURL string) string {
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
                      ยินดีต้อนรับ
                    </p>
                    <p style="margin:0 0 20px 0;color:#ffffff;font-size:18px;font-weight:600;line-height:1.4;">
                      สวัสดี %s
                    </p>
                    <p style="margin:0 0 32px 0;color:#8f8f8f;font-size:14px;line-height:1.8;max-width:300px;">
                      ขอบคุณที่สมัครสมาชิก REMOV กดปุ่มด้านล่างเพื่อยืนยันอีเมล
                      และเริ่มต้นใช้งานบัญชีของคุณได้ทันที
                    </p>
                  </td>
                </tr>
                <tr>
                  <td align="center" style="padding:4px 40px 0 40px;">
                    <table role="presentation" cellpadding="0" cellspacing="0" width="100%%">
                      <tr>
                        <td style="background:#e50914;border-radius:6px;">
                          <a href="%s"
                             style="display:block;padding:15px 44px;color:#ffffff;font-size:13px;
                                    font-weight:700;letter-spacing:1.5px;text-transform:uppercase;
                                    text-decoration:none;text-align:center;">
                            ยืนยันอีเมลของฉัน
                          </a>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>
                <tr>
                  <td align="center" style="padding:22px 40px 8px 40px;">
                    <p style="margin:0;color:#666666;font-size:12px;">
                      ลิงก์มีอายุการใช้งาน 24 ชั่วโมง
                    </p>
                  </td>
                </tr>
                <tr>
                  <td style="padding:12px 40px 40px 40px;">
                    <p style="margin:0;color:#4a4a4a;font-size:11px;line-height:1.7;text-align:center;word-break:break-all;">
                      หากปุ่มด้านบนใช้งานไม่ได้ ให้คัดลอกลิงก์นี้ไปวางในเบราว์เซอร์<br>
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
                หากคุณไม่ได้สมัครสมาชิก REMOV กรุณาเพิกเฉยต่ออีเมลฉบับนี้
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
</html>`, username, verifyURL, verifyURL, verifyURL)
}
