package spec

import (
	"path/filepath"
	"strings"
)

// Filter returns the tests matched by the given selectors, without mutating
// the input. A selector is either a test name ("rest") or a test name and an
// env separated by a dot ("rest.iris"); matching is case-insensitive. An env
// matches by its display name or by its directory base name, so both
// "static.net/http" and "static.nethttp" select the net/http environment.
// Because both test and env names may themselves contain dots (e.g. an env
// named "asp.net core"), every dot position is tried when splitting a
// selector.
//
// With no selectors the input is returned as-is. Selectors that match
// nothing are reported in unmatched so callers can reject typos.
func Filter(tests []Test, selectors []string) (matched []Test, unmatched []string) {
	if len(selectors) == 0 {
		return tests, nil
	}

	used := make([]bool, len(selectors))

	for _, t := range tests {
		testName := strings.ToLower(t.Name)
		keepWhole := false
		var envKeep map[int]struct{}

		for si, raw := range selectors {
			sel := strings.ToLower(strings.TrimSpace(raw))
			if sel == "" {
				used[si] = true // ignore empty selectors silently.
				continue
			}

			if sel == testName {
				keepWhole = true
				used[si] = true
				continue
			}

			// Try "test.env" at every dot position, so names
			// containing dots still match.
			for i := range len(sel) {
				if sel[i] != '.' || sel[:i] != testName {
					continue
				}

				envName := sel[i+1:]
				for ei := range t.Envs {
					if envMatches(t.Envs[ei], envName) {
						if envKeep == nil {
							envKeep = make(map[int]struct{})
						}
						envKeep[ei] = struct{}{}
						used[si] = true
					}
				}
			}
		}

		switch {
		case keepWhole:
			matched = append(matched, t)
		case len(envKeep) > 0:
			filtered := t
			filtered.Envs = make([]Env, 0, len(envKeep))
			for ei := range t.Envs {
				if _, keep := envKeep[ei]; keep {
					filtered.Envs = append(filtered.Envs, t.Envs[ei])
				}
			}
			matched = append(matched, filtered)
		}
	}

	for i, ok := range used {
		if !ok {
			unmatched = append(unmatched, selectors[i])
		}
	}

	return matched, unmatched
}

// envMatches reports whether the (lowercase) selector fragment names the
// env, either by display name or by directory base name.
func envMatches(e Env, name string) bool {
	if strings.ToLower(e.Name) == name {
		return true
	}

	return e.Dir != "" && strings.ToLower(filepath.Base(e.Dir)) == name
}
