package sysinfo

import (
	"context"
	"testing"
	"time"
)

func TestParseToolVersion(t *testing.T) {
	cases := []struct {
		tool   string
		output string
		want   string
	}{
		{"go", "go version go1.26.5 windows/arm64\n", "go1.26.5"},
		{"bombardier", "bombardier version v2.0.2 linux/amd64\n", "v2.0.2"},
		{"bombardier", "bombardier version unspecified linux/amd64\n", ""},
		{"node", "v24.4.1\n", "v24.4.1"},
		{"dotnet", "10.0.100\r\n", "10.0.100"},
		{"dotnet", "10.0.100\r\nextra line", "10.0.100"},
		{"x", "", ""},
	}

	for _, tc := range cases {
		if got := parseToolVersion(tc.tool, tc.output); got != tc.want {
			t.Errorf("parseToolVersion(%q, %q) = %q, want %q", tc.tool, tc.output, got, tc.want)
		}
	}
}

func TestFormatMemory(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, ""},
		{512, "512 B"},
		{1024, "1 KB"},
		{16 * 1024 * 1024 * 1024, "16 GB"},
		{16868593664, "15.71 GB"},
	}

	for _, tc := range cases {
		if got := formatMemory(tc.in); got != tc.want {
			t.Errorf("formatMemory(%f) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestCollectNeverPanics(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info := Collect(ctx, nil)

	// On any development or CI machine the Go toolchain exists (this test
	// is running under it) and basic hardware info must resolve.
	if info.Go == "" {
		t.Error("Go version empty; expected the running toolchain to be detected")
	}
	if info.Processor == "" && info.RAM == "" && info.OS == "" {
		t.Error("no hardware info collected at all")
	}
}
