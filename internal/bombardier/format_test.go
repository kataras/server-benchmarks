package bombardier

import "testing"

// The expectations below match bombardier's own console output for the
// same values; do not "fix" them without checking bombardier first.

func TestFormatBinary(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0.00"},
		{512, "512.00"},
		{2048, "2.00KB"},
		{11189196, "10.67MB"},
		{3 * 1024 * 1024 * 1024, "3.00GB"},
	}

	for _, tc := range cases {
		if got := FormatBinary(tc.in); got != tc.want {
			t.Errorf("FormatBinary(%f) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatTimeUs(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0.00us"},
		{438.94, "438.94us"},
		{1560.13, "1.56ms"},
		{1000000, "1.00s"},
		{90000000, "1.50m"},
	}

	for _, tc := range cases {
		if got := FormatTimeUs(tc.in); got != tc.want {
			t.Errorf("FormatTimeUs(%f) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
