package notification

import (
	"strings"
	"testing"
)

// v1 sliced the endpoint with sub.Endpoint[:50], which panics on any endpoint
// shorter than 50 characters.
func TestTruncateEndpoint(t *testing.T) {
	cases := map[string]struct {
		in            string
		wantSuffix    bool
		wantUnchanged bool
	}{
		"short endpoint is left alone":   {"https://push.example/abc", false, true},
		"empty string does not panic":    {"", false, true},
		"exactly at the limit is intact": {strings.Repeat("x", 50), false, true},
		"longer is truncated":            {strings.Repeat("x", 200), true, false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := truncateEndpoint(tc.in)

			if tc.wantUnchanged && got != tc.in {
				t.Errorf("expected %q unchanged, got %q", tc.in, got)
			}
			if tc.wantSuffix && !strings.HasSuffix(got, "...") {
				t.Errorf("expected a truncation marker, got %q", got)
			}
			if tc.wantSuffix && len(got) != 53 {
				t.Errorf("expected 50 chars plus the marker, got %d", len(got))
			}
		})
	}
}
