package rsync

import (
	"testing"
)

func TestParseProgress(t *testing.T) {
	cases := []struct {
		line    string
		wantPct int
		ok      bool
	}{
		{"123.45M  99%   12.34MB/s    0:00:05 (xfr#495, to-chk=0/500)", 99, true},
		{"  123,456,789  45%   1.23MB/s    0:01:23", 45, true},
		{"building file list ... done", 0, false},
	}

	for _, tc := range cases {
		p, ok := ParseProgress(tc.line)
		if ok != tc.ok {
			t.Fatalf("ParseProgress(%q) ok=%v, want %v", tc.line, ok, tc.ok)
		}
		if !ok {
			continue
		}
		if p.Percent != tc.wantPct {
			t.Fatalf("ParseProgress(%q) percent=%d, want %d", tc.line, p.Percent, tc.wantPct)
		}
	}
}
