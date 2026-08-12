package spec

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/goccy/go-yaml"
)

// Load reads a YAML specification document from r and returns the fully
// normalized and validated tests. Unknown or duplicate YAML fields are
// rejected. Relative Dir and BodyFile paths are resolved against baseDir;
// pass an empty baseDir to resolve them against the current directory.
func Load(r io.Reader, baseDir string) ([]Test, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read spec: %w", err)
	}

	var tests []Test
	if err := yaml.UnmarshalWithOptions(data, &tests, yaml.Strict()); err != nil {
		return nil, fmt.Errorf("parse spec: %w", err)
	}

	if len(tests) == 0 {
		return nil, errors.New("spec contains no tests")
	}

	seen := make(map[string]struct{}, len(tests))
	for i := range tests {
		if err := normalizeTest(&tests[i], baseDir); err != nil {
			return nil, err
		}

		key := strings.ToLower(tests[i].Name)
		if _, dup := seen[key]; dup {
			return nil, fmt.Errorf("duplicate test name %q", tests[i].Name)
		}
		seen[key] = struct{}{}
	}

	return tests, nil
}

// LoadFile reads the specification from the YAML file at path.
// Relative Dir and BodyFile entries are resolved against the file's directory.
func LoadFile(path string) ([]Test, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open spec: %w", err)
	}
	defer f.Close()

	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve spec path: %w", err)
	}

	tests, err := Load(f, filepath.Dir(abs))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	return tests, nil
}

func normalizeTest(t *Test, baseDir string) error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return errors.New("test with empty Name")
	}

	if t.URL == "" {
		return fmt.Errorf("test %q: missing URL", t.Name)
	}

	if t.NumberOfConnections == 0 {
		t.NumberOfConnections = DefaultConnections
	}

	if t.Timeout == 0 {
		t.Timeout = Duration(DefaultTimeout)
	}

	if t.NumberOfRequests == 0 && t.Duration == 0 {
		t.Duration = Duration(DefaultDuration)
	}

	if t.BodyFile != "" && !IsURL(t.BodyFile) && !filepath.IsAbs(t.BodyFile) {
		t.BodyFile = filepath.Join(baseDir, t.BodyFile)
	}

	// Render the description template now, after defaults are in place,
	// so a broken description fails at load time instead of at report
	// time — long after the benchmarks have run.
	description, err := renderDescription(t)
	if err != nil {
		return fmt.Errorf("test %q: Description: %w", t.Name, err)
	}
	t.Description = description

	seen := make(map[string]struct{}, len(t.Envs))
	for i := range t.Envs {
		if err := normalizeEnv(&t.Envs[i], t.Name, baseDir); err != nil {
			return fmt.Errorf("test %q: env #%d: %w", t.Name, i+1, err)
		}

		key := strings.ToLower(t.Envs[i].Name)
		if _, dup := seen[key]; dup {
			return fmt.Errorf("test %q: duplicate env name %q", t.Name, t.Envs[i].Name)
		}
		seen[key] = struct{}{}
	}

	return nil
}

func normalizeEnv(e *Env, testName, baseDir string) error {
	e.Repo = strings.TrimSuffix(strings.TrimSpace(e.Repo), "/")
	e.Name = strings.TrimSpace(e.Name)

	if e.Name == "" {
		if e.Repo == "" {
			return errors.New("missing both Name and Repo fields")
		}

		segment := e.Repo[strings.LastIndexByte(e.Repo, '/')+1:]
		if segment == "" {
			return fmt.Errorf("invalid Repo %q", e.Repo)
		}

		e.Name = upperFirst(segment)
	}

	if e.Language == "" {
		e.Language = DefaultLanguage
	}

	switch {
	case e.Dir == "":
		e.Dir = filepath.Join(baseDir, DefaultCodeDir, strings.ToLower(testName), strings.ToLower(e.Name))
	case !filepath.IsAbs(e.Dir):
		e.Dir = filepath.Join(baseDir, e.Dir)
	default:
		e.Dir = filepath.Clean(e.Dir)
	}

	if e.Exec == "" {
		execCommand, ok := defaultExec(e.Language)
		if !ok {
			return fmt.Errorf("unsupported language %q (set Exec to override)", e.Language)
		}
		e.Exec = execCommand
	}

	if e.Link == "" && e.Repo != "" {
		e.Link = "https://github.com/" + e.Repo
	}

	return nil
}

// defaultExec returns the default Exec commands for the given language.
func defaultExec(language string) (string, bool) {
	switch strings.ToLower(language) {
	case "go", "golang":
		return "go run .", true
	case "c#", "csharp", "net", ".net", "dotnet", "netcore", "net.core", ".net core", "aspnetcore", "kestrel":
		return "dotnet run -c Release", true
	case "node", "nodejs", "javascript", "js":
		return "npm install\nnode .", true
	default:
		return "", false
	}
}

// renderDescription renders t.Description, which may itself be a
// text/template over the test (e.g. "Fires {{.NumberOfRequests}} requests"),
// and trims surrounding whitespace.
func renderDescription(t *Test) (string, error) {
	description := strings.TrimSpace(t.Description)
	if !strings.Contains(description, "{{") {
		return description, nil
	}

	tmpl, err := template.New("description").Parse(description)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, t); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

// IsURL reports whether s is an http or https URL.
func IsURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// upperFirst returns s with its first rune upper-cased, e.g. "iris" -> "Iris".
func upperFirst(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}

	return string(unicode.ToUpper(r)) + s[size:]
}
