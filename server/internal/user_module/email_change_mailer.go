package user_module

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
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
	// เว้นวรรคระหว่างตัวเลขให้ดูเป็นรหัสที่อ่านง่ายและดูมีระดับ
	spaced := strings.Join(strings.Split(otp, ""), " ")

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
                ยืนยันตัวตน
              </p>
              <h2 style="margin:0 0 16px 0;color:#fff;font-size:24px;line-height:1.4;">
                คำขอเปลี่ยนอีเมล
              </h2>
              <p style="margin:0 0 28px 0;color:#b3b3b3;font-size:15px;line-height:1.7;">
                เราได้รับคำขอเปลี่ยนอีเมลของบัญชี REMOV ของคุณ
                กรุณากรอกรหัส OTP ด้านล่างเพื่อยืนยันการเปลี่ยนแปลง
              </p>
            </td>
          </tr>
          <tr>
            <td style="padding:0 36px;">
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0"
                     style="background:#0f0f0f;border:1px dashed #e50914;border-radius:12px;">
                <tr>
                  <td align="center" style="padding:28px 16px;">
                    <span style="display:inline-block;font-size:38px;font-weight:700;letter-spacing:10px;
                                 color:#ffffff;font-family:'Courier New',monospace;">
                      %s
                    </span>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:24px 36px 8px 36px;">
              <p style="margin:0;color:#8a8a8a;font-size:13px;line-height:1.6;">
                ⏱️ รหัสนี้จะหมดอายุใน <strong style="color:#e5e5e5;">15 นาที</strong>
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
                      หากคุณไม่ได้ทำรายการนี้ กรุณาเพิกเฉยต่ออีเมลฉบับนี้
                      และพิจารณาตรวจสอบความปลอดภัยของบัญชีคุณ
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
</html>`, spaced)
}
