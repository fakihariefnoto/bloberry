package usecase

import (
	"testing"

	"github.com/fakihariefnoto/bloberry/internal/domain"
)

func TestCheckUploadPolicyDefault(t *testing.T) {
	p := &domain.UploadPolicy{Mode: "default"}
	tests := []struct {
		name  string
		block bool
	}{
		{"photo.jpg", false},
		{"report.pdf", false},
		{"archive.tar.gz", false},
		{"image.png", false},
		{"shell.php", true},
		{"backdoor.php5", true},
		{"x.phtml", true},
		{"evil.aspx", true},
		{"script.js", true},
		{"module.mjs", true},
		{"page.html", true},
		{"x.htm", true},
		{"vector.svg", true},
		{"notes.xml", true},
		{"run.sh", true},
		{"payload.exe", true},
		{"malware.apk", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := checkUploadPolicy(tt.name, p)
			if tt.block && err == nil {
				t.Fatalf("%s: expected blocked, got nil", tt.name)
			}
			if !tt.block && err != nil {
				t.Fatalf("%s: expected allowed, got %v", tt.name, err)
			}
		})
	}
}

func TestCheckUploadPolicyCaseInsensitive(t *testing.T) {
	p := &domain.UploadPolicy{Mode: "default"}
	if err := checkUploadPolicy("shell.PHP", p); err == nil {
		t.Fatal("expected uppercase .PHP to be blocked")
	}
	if err := checkUploadPolicy("SH.SH", p); err == nil {
		t.Fatal("expected uppercase .SH to be blocked")
	}
}

func TestCheckUploadPolicyAllow(t *testing.T) {
	p := &domain.UploadPolicy{Mode: "allow", Extensions: []string{"png", "jpg", "pdf"}}
	if err := checkUploadPolicy("a.png", p); err != nil {
		t.Fatalf("png should be allowed: %v", err)
	}
	if err := checkUploadPolicy("a.PNG", p); err != nil {
		t.Fatalf("upper PNG should be allowed: %v", err)
	}
	if err := checkUploadPolicy("a.pdf", p); err != nil {
		t.Fatalf("pdf should be allowed: %v", err)
	}
	if err := checkUploadPolicy("a.exe", p); err == nil {
		t.Fatal("exe should be blocked by allowlist")
	}
	if err := checkUploadPolicy("a.php", p); err == nil {
		t.Fatal("php should be blocked by allowlist")
	}
}

func TestCheckUploadPolicyBlock(t *testing.T) {
	p := &domain.UploadPolicy{Mode: "block", Extensions: []string{"zip", "tar.gz"}}
	if err := checkUploadPolicy("a.jpg", p); err != nil {
		t.Fatalf("jpg should be allowed: %v", err)
	}
	if err := checkUploadPolicy("a.zip", p); err == nil {
		t.Fatal("zip should be blocked")
	}
	// built-in blocklist still applies
	if err := checkUploadPolicy("a.php", p); err == nil {
		t.Fatal("php should be blocked even in block mode")
	}
}
