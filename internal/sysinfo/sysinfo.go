// Package sysinfo collects a best-effort fingerprint of the machine and
// the toolchains a benchmark ran on. Collection never fails the caller:
// anything that cannot be determined stays empty and is only logged at
// debug level, and the report templates omit empty entries.
package sysinfo

import (
	"context"
	"log/slog"
	"math"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

// Info is the machine and toolchain fingerprint shown in reports.
type Info struct {
	Processor string `json:"processor,omitempty"`
	RAM       string `json:"ram,omitempty"`
	OS        string `json:"os,omitempty"`

	Bombardier string `json:"bombardier,omitempty"` // e.g. "v2.0.2"
	Go         string `json:"go,omitempty"`         // e.g. "go1.26.5"
	Dotnet     string `json:"dotnet,omitempty"`     // e.g. "10.0.100"
	Node       string `json:"node,omitempty"`       // e.g. "v24.4.1"
}

// Collect gathers hardware information and tool versions. It always
// returns, with empty fields for whatever could not be determined.
// A nil logger discards diagnostics.
func Collect(ctx context.Context, logger *slog.Logger) Info {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	var info Info

	if stats, err := cpu.InfoWithContext(ctx); err == nil && len(stats) > 0 {
		info.Processor = strings.TrimSpace(stats[0].ModelName)
	} else {
		logger.Debug("sysinfo: cpu info unavailable", "error", err)
	}

	if vm, err := mem.VirtualMemoryWithContext(ctx); err == nil {
		info.RAM = formatMemory(float64(vm.Total))
	} else {
		logger.Debug("sysinfo: memory info unavailable", "error", err)
	}

	if h, err := host.InfoWithContext(ctx); err == nil {
		info.OS = h.Platform
	} else {
		logger.Debug("sysinfo: host info unavailable", "error", err)
	}

	info.Go = toolVersion(ctx, logger, "go", "version")
	info.Bombardier = toolVersion(ctx, logger, "bombardier", "--version")
	info.Dotnet = toolVersion(ctx, logger, "dotnet", "--version")
	info.Node = toolVersion(ctx, logger, "node", "--version")

	return info
}

// toolVersion runs a tool's version command, returning an empty string
// when the tool is missing or misbehaving.
func toolVersion(ctx context.Context, logger *slog.Logger, name string, args ...string) string {
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil {
		logger.Debug("sysinfo: tool version unavailable", "tool", name, "error", err)
		return ""
	}

	return parseToolVersion(name, string(out))
}

// parseToolVersion extracts a compact version from a version command's
// output, e.g. "go version go1.26.5 windows/arm64\n" -> "go1.26.5".
func parseToolVersion(tool, output string) string {
	v := strings.TrimSpace(output)
	v = strings.TrimPrefix(v, tool+" version ")

	if i := strings.IndexAny(v, "\r\n"); i >= 0 {
		v = v[:i]
	}

	// Drop a trailing platform suffix ("go1.26.5 windows/arm64"); the OS
	// already has its own row in the report.
	if i := strings.IndexByte(v, ' '); i > 0 {
		v = v[:i]
	}

	// Tools built from source without release metadata report this
	// placeholder; treat it as unknown.
	if v == "unspecified" {
		return ""
	}

	return v
}

// formatMemory renders a byte count as a human-readable size, e.g.
// 16868593664 -> "15.71 GB".
func formatMemory(size float64) string {
	if size <= 0 {
		return ""
	}

	suffixes := [...]string{"B", "KB", "MB", "GB", "TB"}

	base := math.Log(size) / math.Log(1024)
	idx := min(int(math.Floor(base)), len(suffixes)-1)
	value := round(math.Pow(1024, base-math.Floor(base)), 0.5, 2)

	return strconv.FormatFloat(value, 'f', -1, 64) + " " + suffixes[idx]
}

// round rounds val to the given number of decimal places, rounding the
// final digit up when its remainder reaches roundOn.
func round(val, roundOn float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	digit := pow * val

	_, div := math.Modf(digit)
	if div >= roundOn {
		return math.Ceil(digit) / pow
	}

	return math.Floor(digit) / pow
}
