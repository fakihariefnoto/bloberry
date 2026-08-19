package mailer

import (
	"fmt"
	"html"
)

// Render builds a branded HTML email. kind selects the content block. All
// user-supplied values pass through html.EscapeString.
func Render(kind string, v map[string]string) string {
	title := ""
	body := ""
	switch kind {
	case "otp":
		title = "Your one-time login code"
		body = fmt.Sprintf(`
      <p style="margin:0 0 16px;font-size:15px;line-height:1.6;color:#374151;">
        A sign-in to your Bloberry account used this one-time code. Enter it on
        the login screen to finish signing in. It expires in a few minutes.
      </p>
      <div style="margin:0 0 20px;padding:16px;border-radius:12px;background:#0f0b2e;text-align:center;">
        <span style="font-size:34px;font-weight:800;letter-spacing:10px;color:#8b7deb;">%s</span>
      </div>
      <p style="margin:0;font-size:13px;line-height:1.6;color:#6b7280;">
        If you didn't request this code, you can safely ignore this email —
        your account stays protected.
      </p>`, html.EscapeString(v["code"]))
	case "reset":
		title = "Reset your password"
		body = fmt.Sprintf(`
      <p style="margin:0 0 16px;font-size:15px;line-height:1.6;color:#374151;">
        We received a request to reset the password for
        <strong>%s</strong>. The link below is valid for 30 minutes.
      </p>
      <p style="margin:0 0 20px;">
        <a href="%s" style="display:inline-block;padding:12px 24px;border-radius:10px;background:#8b7deb;color:#ffffff;text-decoration:none;font-weight:700;font-size:14px;">
          Reset password
        </a>
      </p>
      <p style="margin:0 0 16px;font-size:13px;color:#6b7280;">
        If the button doesn't work, copy this link into your browser:
      </p>
      <p style="margin:0 0 20px;font-size:12px;color:#6b7280;word-break:break-all;">%s</p>
      <p style="margin:0;font-size:13px;line-height:1.6;color:#6b7280;">
        If you didn't ask to reset your password, ignore this email.
      </p>`,
			html.EscapeString(v["email"]),
			html.EscapeString(v["url"]),
			html.EscapeString(v["url"]))
	case "invite":
		title = "You've been invited"
		body = fmt.Sprintf(`
      <p style="margin:0 0 16px;font-size:15px;line-height:1.6;color:#374151;">
        You've been invited to join <strong>%s</strong> on Bloberry.
        Create your password to get started.
      </p>
      <p style="margin:0 0 20px;">
        <a href="%s" style="display:inline-block;padding:12px 24px;border-radius:10px;background:#8b7deb;color:#ffffff;text-decoration:none;font-weight:700;font-size:14px;">
          Accept invitation
        </a>
      </p>
      <p style="margin:0 0 16px;font-size:13px;color:#6b7280;">
        Or copy this link into your browser:
      </p>
      <p style="margin:0 0 20px;font-size:12px;color:#6b7280;word-break:break-all;">%s</p>
      <p style="margin:0;font-size:13px;line-height:1.6;color:#6b7280;">
        The invitation expires in 7 days.
      </p>`,
			html.EscapeString(v["tenant"]),
			html.EscapeString(v["url"]),
			html.EscapeString(v["url"]))
	case "activation":
		title = "Activate your account"
		body = fmt.Sprintf(`
      <p style="margin:0 0 16px;font-size:15px;line-height:1.6;color:#374151;">
        You've been added to <strong>%s</strong> on Bloberry. Set your password
        once to activate your account.
      </p>
      <p style="margin:0 0 20px;">
        <a href="%s" style="display:inline-block;padding:12px 24px;border-radius:10px;background:#8b7deb;color:#ffffff;text-decoration:none;font-weight:700;font-size:14px;">
          Activate account
        </a>
      </p>
      <p style="margin:0 0 16px;font-size:13px;color:#6b7280;">
        Enter your email <strong>%s</strong> at the activation page. This only
        works once.
      </p>
      <p style="margin:0 0 20px;font-size:12px;color:#6b7280;word-break:break-all;">%s</p>`,
			html.EscapeString(v["tenant"]),
			html.EscapeString(v["url"]),
			html.EscapeString(v["email"]),
			html.EscapeString(v["url"]))
	}

	// Outer shell shared by every template.
	return fmt.Sprintf(`<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body style="margin:0;padding:0;background:#f3f4f6;">
    <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background:#f3f4f6;padding:32px 16px;">
      <tr>
        <td align="center">
          <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="max-width:520px;background:#ffffff;border-radius:16px;overflow:hidden;border:1px solid #e5e7eb;">
            <tr>
              <td style="background:#0f0b2e;padding:24px 32px;">
                <span style="font-size:17px;font-weight:800;color:#8b7deb;">Bloberry</span>
              </td>
            </tr>
            <tr>
              <td style="padding:32px;">
                <h1 style="margin:0 0 16px;font-size:21px;font-weight:800;color:#111827;">%s</h1>
                %s
              </td>
            </tr>
            <tr>
              <td style="padding:20px 32px;background:#f9fafb;border-top:1px solid #e5e7eb;">
                <p style="margin:0;font-size:12px;color:#9ca3af;">
                  Bloberry · self-hosted object storage ·
                  <a href="https://bloberry.app" style="color:#9ca3af;text-decoration:none;">bloberry.app</a>
                </p>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`, title, body)
}
