package cli

import (
	"fmt"
	"io"
	"runtime/debug"
)

// versionCmd implements `server-benchmarks version` using the build
// information stamped by the Go toolchain.
func versionCmd(stdout io.Writer) int {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprintln(stdout, "server-benchmarks (no build info)")
		return ExitOK
	}

	version := info.Main.Version
	if version == "" || version == "(devel)" {
		version = "devel"
	}

	var revision, modified string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			if setting.Value == "true" {
				modified = " (modified)"
			}
		}
	}

	fmt.Fprintf(stdout, "server-benchmarks %s", version)
	if len(revision) >= 12 {
		fmt.Fprintf(stdout, " (%s%s)", revision[:12], modified)
	}
	fmt.Fprintf(stdout, " %s\n", info.GoVersion)

	return ExitOK
}
