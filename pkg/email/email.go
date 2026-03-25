package email

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendVerificationEmail(toEmail, callbackUrl, token string) error {
	url := fmt.Sprintf("%s?token=%s", callbackUrl, token)
	html := verifyEmailTemplate(url)
	return send(toEmail, "Verifikasi Akun Anda - SkillGap", html)
}

func SendResetPasswordEmail(toEmail, callbackUrl, token string) error {
	url := fmt.Sprintf("%s?token=%s", callbackUrl, token)
	html := resetPasswordTemplate(url)
	return send(toEmail, "Reset Password - SkillGap", html)
}

func send(toEmail, subject, html string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "norepy@undev.fun",
		To:      []string{toEmail},
		Subject: subject,
		Html:    html,
	}

	resp, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println("RESEND ERROR: ", err)
		return err
	}
	fmt.Println("RESEND SUCCESS: ", resp.Id)
	return err
}

func verifyEmailTemplate(url string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"/><title>Verify Your Account</title></head>
<body style="margin:0;padding:0;background-color:#f7f6f3;font-family:Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0">
    <tr>
      <td align="center" style="padding:40px 16px;">
        <table style="max-width:520px;background:#ffffff;border-radius:10px;padding:32px;">
          <tr>
            <td style="text-align:center;padding-bottom:24px;">
              <h2 style="margin:0;color:#171815;">Verifikasi Akun Anda</h2>
            </td>
          </tr>
          <tr>
            <td style="color:#544337;font-size:14px;line-height:1.7;padding-bottom:20px;">
              Terima kasih telah mendaftar di <strong>SkillGap</strong>.
              Klik tombol di bawah untuk mengaktifkan akun Anda.
            </td>
          </tr>
          <tr>
            <td align="center" style="padding:24px 0;">
              <a href="%s" style="background-color:#4F46E5;color:#ffffff;text-decoration:none;
                padding:14px 28px;border-radius:8px;font-size:14px;font-weight:600;display:inline-block;">
                Verifikasi Email
              </a>
            </td>
          </tr>
          <tr>
            <td style="color:#544337;font-size:13px;">
              Link ini akan kedaluwarsa dalam <strong>15 menit</strong>.
            </td>
          </tr>
          <tr>
            <td style="padding-top:24px;border-top:1px solid #e6e3dc;color:#544337;font-size:12px;">
              Jika tombol tidak berfungsi, salin link ini:<br/>
              <span style="word-break:break-all;">%s</span>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, url, url)
}

func resetPasswordTemplate(url string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"/><title>Reset Password</title></head>
<body style="margin:0;padding:0;background-color:#f7f6f3;font-family:Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0">
    <tr>
      <td align="center" style="padding:40px 16px;">
        <table style="max-width:520px;background:#ffffff;border-radius:10px;padding:32px;">
          <tr>
            <td style="text-align:center;padding-bottom:24px;">
              <h2 style="margin:0;color:#171815;">Reset Password</h2>
            </td>
          </tr>
          <tr>
            <td style="color:#544337;font-size:14px;line-height:1.7;padding-bottom:20px;">
              Kami menerima permintaan reset password untuk akun <strong>SkillGap</strong> Anda.
            </td>
          </tr>
          <tr>
            <td align="center" style="padding:24px 0;">
              <a href="%s" style="background-color:#4F46E5;color:#ffffff;text-decoration:none;
                padding:14px 28px;border-radius:8px;font-size:14px;font-weight:600;display:inline-block;">
                Reset Password
              </a>
            </td>
          </tr>
          <tr>
            <td style="color:#544337;font-size:13px;">
              Link ini akan kedaluwarsa dalam <strong>15 menit</strong>.
            </td>
          </tr>
          <tr>
            <td style="padding-top:24px;border-top:1px solid #e6e3dc;color:#544337;font-size:12px;">
              Jika tombol tidak berfungsi, salin link ini:<br/>
              <span style="word-break:break-all;">%s</span>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, url, url)
}
