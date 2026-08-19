package mailer

import (
	"strings"
	"testing"
)

func TestRenderOTP(t *testing.T) {
	h := Render("otp", map[string]string{"code": "123456"})
	if !strings.Contains(h, "123456") { t.Fatal("code missing") }
	if !strings.Contains(h, "Bloberry") { t.Fatal("brand missing") }
	if !strings.Contains(h, "</html>") { t.Fatal("not html") }
}

func TestRenderEscapes(t *testing.T) {
	h := Render("reset", map[string]string{"email": "<script>", "url": "https://x/?a=1&b=2"})
	if strings.Contains(h, "<script>") { t.Fatal("email not escaped") }
	if !strings.Contains(h, "&amp;") { t.Fatal("url not escaped") }
}

func TestSMTPBuildMultipart(t *testing.T) {
	s := &SMTP{From: "noreply@x.com"}
	msg, err := s.buildMessage("a@b.com", "Subj", "plain", "<p>html</p>")
	if err != nil { t.Fatal(err) }
	if !strings.Contains(msg, "multipart/alternative") { t.Fatal("not multipart") }
	if !strings.Contains(msg, "text/html") { t.Fatal("no html part") }
	msg2, _ := s.buildMessage("a@b.com", "Subj", "plain", "")
	if strings.Contains(msg2, "multipart") { t.Fatal("plain should not be multipart") }
}
