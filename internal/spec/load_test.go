package spec

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadDefaultsAndNormalization(t *testing.T) {
	const doc = `
- Name: Static
  Description: Fires {{.NumberOfRequests}} requests.
  NumberOfRequests: 100000
  Method: GET
  URL: http://localhost:5000
  Envs:
    - Repo: kataras/iris
    - Repo: expressjs/express
      Language: Javascript
    - Name: ASP.NET Core
      Repo: dotnet/aspnetcore
      Language: "C#"
      Dir: ./_code/static/aspnetcore
`
	base := t.TempDir()
	tests, err := Load(strings.NewReader(doc), base)
	if err != nil {
		t.Fatal(err)
	}

	if len(tests) != 1 {
		t.Fatalf("expected 1 test, got %d", len(tests))
	}

	test := tests[0]
	if got, want := test.Description, "Fires 100000 requests."; got != want {
		t.Errorf("Description = %q, want %q (templates must render at load)", got, want)
	}
	if got, want := test.NumberOfConnections, uint64(DefaultConnections); got != want {
		t.Errorf("NumberOfConnections = %d, want %d", got, want)
	}
	if got, want := test.Timeout.Std(), DefaultTimeout; got != want {
		t.Errorf("Timeout = %s, want %s", got, want)
	}
	if test.Duration != 0 {
		t.Errorf("Duration = %s, want 0 (NumberOfRequests is set)", test.Duration)
	}

	envs := test.Envs
	if len(envs) != 3 {
		t.Fatalf("expected 3 envs, got %d", len(envs))
	}

	iris := envs[0]
	if got, want := iris.Name, "Iris"; got != want {
		t.Errorf("derived name = %q, want %q", got, want)
	}
	if got, want := iris.Language, DefaultLanguage; got != want {
		t.Errorf("language = %q, want %q", got, want)
	}
	if got, want := iris.Exec, "go run ."; got != want {
		t.Errorf("exec = %q, want %q", got, want)
	}
	if got, want := iris.Dir, filepath.Join(base, DefaultCodeDir, "static", "iris"); got != want {
		t.Errorf("dir = %q, want %q", got, want)
	}
	if got, want := iris.Link, "https://github.com/kataras/iris"; got != want {
		t.Errorf("link = %q, want %q", got, want)
	}

	express := envs[1]
	if got, want := express.Exec, "npm install\nnode ."; got != want {
		t.Errorf("express exec = %q, want %q", got, want)
	}

	aspnet := envs[2]
	if got, want := aspnet.Exec, "dotnet run -c Release"; got != want {
		t.Errorf("aspnet exec = %q, want %q", got, want)
	}
	if got, want := aspnet.Dir, filepath.Join(base, "_code", "static", "aspnetcore"); got != want {
		t.Errorf("aspnet dir = %q, want %q (explicit relative Dir)", got, want)
	}
}

func TestLoadDurationDefaultAndParsing(t *testing.T) {
	const doc = `
- Name: Timed
  Duration: 8s
  Timeout: 30s
  URL: http://localhost:5000
  Envs:
    - Repo: kataras/iris
- Name: Defaulted
  URL: http://localhost:5000
  Envs:
    - Repo: kataras/iris
`
	tests, err := Load(strings.NewReader(doc), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if got, want := tests[0].Duration.Std(), 8*time.Second; got != want {
		t.Errorf("Duration = %s, want %s", got, want)
	}
	if got, want := tests[0].Timeout.Std(), 30*time.Second; got != want {
		t.Errorf("Timeout = %s, want %s", got, want)
	}
	if got, want := tests[1].Duration.Std(), DefaultDuration; got != want {
		t.Errorf("defaulted Duration = %s, want %s", got, want)
	}
}

func TestLoadLocalPrivateEnv(t *testing.T) {
	// A private/local framework: no public Repo, absolute Dir.
	dir := t.TempDir()
	doc := `
- Name: Static
  URL: http://localhost:5000
  Envs:
    - Name: Iris (next)
      Dir: ` + filepath.ToSlash(dir) + `
`
	tests, err := Load(strings.NewReader(doc), "ignored-base")
	if err != nil {
		t.Fatal(err)
	}

	env := tests[0].Envs[0]
	if got, want := env.Name, "Iris (next)"; got != want {
		t.Errorf("name = %q, want %q", got, want)
	}
	if got, want := env.Dir, filepath.Clean(dir); got != want {
		t.Errorf("dir = %q, want %q (absolute Dir must be kept)", got, want)
	}
	if env.Link != "" {
		t.Errorf("link = %q, want empty (no Repo)", env.Link)
	}
	if got, want := env.Exec, "go run ."; got != want {
		t.Errorf("exec = %q, want %q", got, want)
	}
}

func TestLoadBodyFileResolution(t *testing.T) {
	base := t.TempDir()
	doc := `
- Name: REST
  Method: POST
  BodyFile: ./payload.json
  URL: http://localhost:5000/42
  Envs:
    - Repo: kataras/iris
- Name: REST2
  Method: POST
  BodyFile: "https://example.com/1MB.json"
  URL: http://localhost:5000/42
  Envs:
    - Repo: kataras/iris
`
	tests, err := Load(strings.NewReader(doc), base)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := tests[0].BodyFile, filepath.Join(base, "payload.json"); got != want {
		t.Errorf("local BodyFile = %q, want %q", got, want)
	}
	if got, want := tests[1].BodyFile, "https://example.com/1MB.json"; got != want {
		t.Errorf("URL BodyFile = %q, want %q (must stay untouched)", got, want)
	}
}

func TestLoadErrors(t *testing.T) {
	cases := []struct {
		name    string
		doc     string
		wantSub string
	}{
		{"empty spec", `[]`, "no tests"},
		{"missing url", "- Name: X\n  Envs:\n    - Repo: a/b", "missing URL"},
		{"missing name and repo", "- Name: X\n  URL: http://localhost:5000\n  Envs:\n    - Language: Go", "missing both Name and Repo"},
		{"invalid repo", "- Name: X\n  URL: http://localhost:5000\n  Envs:\n    - Repo: a//", "invalid Repo"},
		{"unsupported language", "- Name: X\n  URL: http://localhost:5000\n  Envs:\n    - Repo: a/b\n      Language: Rust", "unsupported language"},
		{"unknown field", "- Name: X\n  URL: http://localhost:5000\n  Bogus: yes\n  Envs:\n    - Repo: a/b", "unknown field"},
		{"duplicate test names", "- Name: X\n  URL: http://localhost:5000\n  Envs:\n    - Repo: a/b\n- Name: x\n  URL: http://localhost:5000\n  Envs:\n    - Repo: a/b", "duplicate test"},
		{"duplicate env names", "- Name: X\n  URL: http://localhost:5000\n  Envs:\n    - Repo: a/b\n    - Repo: c/b", "duplicate env"},
		{"invalid duration", "- Name: X\n  Duration: fast\n  URL: http://localhost:5000\n  Envs:\n    - Repo: a/b", "invalid duration"},
		{"broken description template", "- Name: X\n  Description: '{{.Bogus}}'\n  URL: http://localhost:5000\n  Envs:\n    - Repo: a/b", "Description"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Load(strings.NewReader(tc.doc), t.TempDir())
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.wantSub)
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Fatalf("error = %q, want it to contain %q", err, tc.wantSub)
			}
		})
	}
}

func TestLoadUnsupportedLanguageAllowedWithExec(t *testing.T) {
	const doc = `
- Name: X
  URL: http://localhost:5000
  Envs:
    - Repo: a/b
      Language: Rust
      Exec: |-
        cargo build --release
        ./target/release/app
`
	tests, err := Load(strings.NewReader(doc), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if got := tests[0].Envs[0].Exec; !strings.Contains(got, "cargo build") {
		t.Errorf("exec = %q, want the explicit commands kept", got)
	}
}

func TestUpperFirst(t *testing.T) {
	for in, want := range map[string]string{
		"iris": "Iris", "gin": "Gin", "chi": "Chi", "": "",
		"aspnetcore": "Aspnetcore", "Fiber": "Fiber",
	} {
		if got := upperFirst(in); got != want {
			t.Errorf("upperFirst(%q) = %q, want %q", in, got, want)
		}
	}
}
