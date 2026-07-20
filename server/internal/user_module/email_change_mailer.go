package user_module

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

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

	subject := "รหัสยืนยันการเปลี่ยนอีเมล — REMOV"
	htmlBody := buildOTPEmailBody(otp)

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
		return fmt.Errorf("mailer: send otp to %s: %w", toEmail, err)
	}
	return nil
}

func buildOTPEmailBody(otp string) string {
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
                      คำขอเปลี่ยนอีเมล
                    </p>
                    <p style="margin:0 0 28px 0;color:#8f8f8f;font-size:14px;line-height:1.8;max-width:300px;">
                      กรอกรหัสยืนยันด้านล่างเพื่อดำเนินการเปลี่ยนอีเมลของบัญชีคุณ
                    </p>
                  </td>
                </tr>

                <!-- OTP code: copy-friendly, ระยะห่างทำด้วย letter-spacing ล้วน text สะอาด -->
                <tr>
                  <td align="center" style="padding:0 40px;">
                    <table role="presentation" cellpadding="0" cellspacing="0"
                           style="background:#050505;border:1px solid #e50914;border-radius:6px;">
                      <tr>
                        <td style="padding:20px 32px;">
                          <span style="font-size:30px;font-weight:700;letter-spacing:10px;color:#ffffff;
                                       font-family:'Courier New',monospace;user-select:all;">%s</span>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>

                <tr>
                  <td align="center" style="padding:14px 40px 4px 40px;">
                    <p style="margin:0;color:#4a4a4a;font-size:11px;line-height:1.6;">
                      แตะ/คลิกที่รหัสด้านบนค้างไว้เพื่อเลือกทั้งหมด แล้วคัดลอกได้ทันที
                    </p>
                  </td>
                </tr>
                <tr>
                  <td align="center" style="padding:16px 40px 40px 40px;">
                    <p style="margin:0;color:#666666;font-size:12px;">
                      รหัสมีอายุการใช้งาน 15 นาที
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
                หากคุณไม่ได้เป็นผู้ทำรายการนี้ กรุณาเพิกเฉยต่ออีเมลฉบับนี้
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
</html>`, otp)
}
