package bombardier

import (
	"maps"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/kataras/server-benchmarks/internal/spec"
)

// Args builds the bombardier command-line arguments for t. It is a pure
// function: t is never mutated and the result is deterministic (headers
// are emitted in sorted order).
//
// bodyFilePath, when non-empty, overrides t.BodyFile; the runner passes
// the local path of a downloaded BodyFile URL here.
func Args(t spec.Test, bodyFilePath string) []string {
	// Sized for the common case: 5 flag pairs + headers + trailing 3.
	args := make([]string, 0, 16+2*len(t.Headers))

	if v := t.NumberOfConnections; v > 0 {
		args = append(args, "-c", strconv.FormatUint(v, 10))
	}

	if v := t.NumberOfRequests; v > 0 {
		args = append(args, "-n", strconv.FormatUint(v, 10))
	}

	if v := t.Duration; v > 0 {
		args = append(args, "-d", v.String())
	}

	if v := t.Timeout; v > 0 {
		args = append(args, "-t", v.String())
	}

	if v := t.Method; v != "" {
		args = append(args, "-m", v)
	}

	headers := t.Headers
	if bodyFile := firstNonEmpty(bodyFilePath, t.BodyFile); bodyFile != "" {
		args = append(args, "-f", bodyFile)
		headers = withInferredContentType(t, bodyFile)
	}

	for _, k := range slices.Sorted(maps.Keys(headers)) {
		args = append(args, "-H", k+": "+headers[k])
	}

	return append(args, "--format=json", "--print=result", t.URL)
}

// firstNonEmpty returns a when non-empty, otherwise b.
func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}

	return b
}

// withInferredContentType returns the test's headers, extended with a
// Content-Type inferred from the body file extension when the request
// carries a body (POST/PUT) and no Content-Type was declared.
func withInferredContentType(t spec.Test, bodyFile string) map[string]string {
	if t.Method != "POST" && t.Method != "PUT" {
		return t.Headers
	}

	if _, exists := t.Headers["Content-Type"]; exists {
		return t.Headers
	}

	contentType := "application/x-www-form-urlencoded"
	switch filepath.Ext(bodyFile) {
	case ".bin":
		contentType = "application/octet-stream"
	case ".json":
		contentType = "application/json"
	case ".xml":
		contentType = "text/xml"
	case ".txt":
		contentType = "text/plain"
	}

	headers := make(map[string]string, len(t.Headers)+1)
	maps.Copy(headers, t.Headers)
	headers["Content-Type"] = contentType

	return headers
}
