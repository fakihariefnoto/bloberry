package usecase

import "testing"
func TestRejectUnsafeName(t *testing.T) {
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
			err := rejectUnsafeName(tt.name)
			if tt.block && err == nil {
				t.Fatalf("%s: expected blocked, got nil", tt.name)
			}
			if !tt.block && err != nil {
				t.Fatalf("%s: expected allowed, got %v", tt.name, err)
			}
		})
	}
}

func TestRejectUnsafeNameCaseInsensitive(t *testing.T) {
	// .PHP must block just like .php — uploaders try case tricks.
	if err := rejectUnsafeName("shell.PHP"); err == nil {
		t.Fatal("expected uppercase .PHP to be blocked")
	}
	if err := rejectUnsafeName("SH.SH"); err == nil {
		t.Fatal("expected uppercase .SH to be blocked")
	}
}
