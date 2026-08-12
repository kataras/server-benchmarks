package bombardier

import (
	"encoding/json"
	"fmt"
)

// envelope mirrors the top level of bombardier's --format=json output.
// Decoding into it (never into caller-owned types) guarantees that no
// unrelated JSON key can overwrite caller state.
type envelope struct {
	Result *Result `json:"result"`
}

// Parse decodes bombardier's --format=json --print=result output.
// It fails when the document has no "result" object, so callers can never
// observe a half-filled Result.
func Parse(data []byte) (Result, error) {
	var doc envelope
	if err := json.Unmarshal(data, &doc); err != nil {
		return Result{}, fmt.Errorf("parse bombardier output: %w (output: %s)", err, snippet(data))
	}

	if doc.Result == nil {
		return Result{}, fmt.Errorf(`parse bombardier output: missing "result" object (output: %s)`, snippet(data))
	}

	return *doc.Result, nil
}

// snippet returns a short, single-purpose preview of raw output for error
// messages.
func snippet(data []byte) string {
	const max = 256

	s := string(data)
	if len(s) > max {
		s = s[:max] + "..."
	}

	if s == "" {
		s = "<empty>"
	}

	return s
}
