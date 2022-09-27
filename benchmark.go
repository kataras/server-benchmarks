package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const defaultLanguage = "Go"

type (
	// Test contains information about test will be performed.
	Test struct {
		Name        string     `yaml:"Name"` // e.g. "Parameterized"
		Description string     `yaml:"Description"`
		Envs        []*TestEnv `yaml:"Envs"`

		NumberOfConnections uint64            `yaml:"NumberOfConnections"` // defaults to 125.
		NumberOfRequests    uint64            `yaml:"NumberOfRequests"`
		Duration            time.Duration     `yaml:"Duration"`
		Timeout             time.Duration     `yaml:"Timeout"` // defaults to 2s
		Headers             map[string]string `yaml:"Headers"`
		Method              string            `yaml:"Method"`
		URL                 string            `yaml:"URL"`
		BodyFile            string            `yaml:"BodyFile"`
	}

	// TestEnv is the place at which stress-test code should be able to located, per framework.
	TestEnv struct {
		Name              string `yaml:"Name"` // can be empty and retrieved by repo.
		Link              string `yaml:"Link"` // the link behind the name. Read GetLink for more.
		Repo              string `yaml:"Repo"` // e.g. kataras/iris
		Dir               string `yaml:"Dir"`  // e.g. ./benchmarks/iris
		absDir            string //  Dir converted to absolute path.
		Exec              string `yaml:"Exec"`     // e.g. go run main.go, can be multiline.
		Language          string `yaml:"Language"` // defaults to Go.
		NotYetImplemented bool   `yaml:"NotYetImplemented"`
		NotSupported      bool   `yaml:"NotSupported"`

		Result *TestResult `json:"result" yaml:"-"`
	}

	// TestResult holds results of the test.
	TestResult struct {
		BytesRead        int64   `json:"bytesRead"`
		BytesWritten     int64   `json:"bytesWritten"`
		TimeTakenSeconds float64 `json:"timeTakenSeconds"`

		Req1XX uint64 `json:"req1xx"`
		Req2XX uint64 `json:"req2xx"`
		Req3XX uint64 `json:"req3xx"`
		Req4XX uint64 `json:"req4xx"`
		Req5XX uint64 `json:"req5xx"`
		Others uint64 `json:"others"`
		Errors []struct {
			Description string `json:"description"`
			Count       uint64 `json:"count"`
		} `json:"errors"`

		Latency           Stats `json:"latency"`
		RequestsPerSecond Stats `json:"rps"`
	}

	// Stats holds the Latency and RPS numbers.
	Stats struct {
		Mean   float64 `json:"mean"`
		Stddev float64 `json:"stddev"`
		Max    float64 `json:"max"`
	}
)

// ParseDescription returns the template-parsed description text.
func (t *Test) ParseDescription(tmplData interface{}) string {
	description := t.Description
	if strings.Contains(description, "{{") {
		// it's a template.
		tmpl, err := template.New("description").Parse(t.Description)
		if err != nil {
			return err.Error()
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, tmplData)
		if err != nil {
			return err.Error()
		}

		description = buf.String()
	}

	return strings.TrimSuffix(description, "\n")
}

func (t *Test) buildArgs() (args []string) {
	// default concurrent connections (this can be omitted, as it's the bombardier's default).
	if t.NumberOfConnections == 0 {
		t.NumberOfConnections = 125
	}

	// default timeout to 2 seconds (this can be omitted, as it's the bombardier's default).
	if t.Timeout == 0 {
		t.Timeout = 2 * time.Second
	}

	// if not number of requests to fire defined and not a test duration,
	// then default test duration to 5 seconds.
	if t.NumberOfRequests == 0 && t.Duration == 0 {
		t.Duration = 5 * time.Second
	}

	if v := t.NumberOfConnections; v > 0 {
		args = append(args, []string{"-c", fmt.Sprintf("%d", v)}...)
	}

	if v := t.NumberOfRequests; v > 0 {
		args = append(args, []string{"-n", fmt.Sprintf("%d", v)}...)
	}

	if v := t.Duration; v > 0 {
		args = append(args, []string{"-d", v.String()}...)
	}

	if v := t.Timeout; v > 0 {
		args = append(args, []string{"-t", v.String()}...)
	}

	if v := t.Method; v != "" {
		args = append(args, []string{"-m", v}...)
	}

	if v := t.BodyFile; v != "" {
		args = append(args, []string{"-f", v}...)

		// try to fill the content type if missing.
		if t.Method == "POST" || t.Method == "PUT" {
			contentTypeHeaderKey := "Content-Type"
			cType := "application/x-www-form-urlencoded"
			if _, exists := t.Headers[contentTypeHeaderKey]; !exists {
				if t.Headers == nil {
					t.Headers = make(map[string]string)
				}

				switch filepath.Ext(v) {
				case ".bin":
					cType = "application/octet-stream"
				case ".json":
					cType = "application/json"
				case ".xml":
					cType = "text/xml"
				case ".txt":
					cType = "text/plain"
				}

				t.Headers[contentTypeHeaderKey] = cType
			}
		}
	}

	if headers := t.Headers; len(headers) > 0 {
		for k, v := range headers {
			args = append(args, []string{"-H", fmt.Sprintf(`%s: %s`, k, v)}...)
		}
	}

	// "--http2", "--insecure",
	args = append(args, []string{"--format=json", "--print=result", t.URL}...)
	return args
}

// GetName returns the last segment of the Repo.
func (e *TestEnv) GetName() string {
	if e.Name != "" {
		return e.Name
	}

	if e.Repo == "" {
		panic("missing Repo field")
	}

	e.Repo = strings.TrimSuffix(e.Repo, "/")

	idx := strings.LastIndexByte(e.Repo, '/')
	if idx < 1 {
		panic("invalid repo <" + e.Repo + ">")
	}

	// name := strings.Title(e.Repo[idx+1:])
	caser := cases.Title(language.English)
	name := caser.String(e.Repo[idx+1:])

	e.Name = name
	return name
}

// GetLink returns the link under the name.
// Should target the source code of the test.
// If it's empty then the repo or the dir is the link.
// Dir: useful as link for local-only projects.
// Repo: useful as link when comparing different projects.
// Link: useful as a relative link to published benchmark,
// which the code exist in the same or other repository.
//
// If Link is empty and Repo is empty too but Dir exists
// then this Link is automatically set to Dir (as given by the user).
func (e *TestEnv) GetLink() string {
	if e.Link != "" {
		return e.Link
	}

	if e.Repo != "" {
		return "https://github.com/" + e.Repo
	}

	return filepath.ToSlash(e.Dir)
}

// CanBenchmark reports whether this test can run on this env.
func (e *TestEnv) CanBenchmark() bool {
	return !e.NotSupported && !e.NotYetImplemented
}

// Throughput returns total throughput (read + write) in bytes per
// second.
func (r TestResult) Throughput() float64 {
	return float64(r.BytesRead+r.BytesWritten) / r.TimeTakenSeconds
}

const defaultCodeDir = "./_code"

func runBenchmark(t *Test, env *TestEnv) (err error) {
	if env.Dir == "" {
		env.absDir, err = filepath.Abs(filepath.Join(defaultCodeDir, strings.ToLower(t.Name), strings.ToLower(env.GetName())))
		if err != nil {
			return
		}
	} else if !filepath.IsAbs(env.Dir) {
		env.absDir, err = filepath.Abs(env.Dir)
		if err != nil {
			return
		}
	}

	if env.Language == "" {
		env.Language = defaultLanguage
	}

	if env.Exec == "" {
		if env.absDir == "" {
			return fmt.Errorf("%s:%s missing Exec and Dir fields", t.Name, env.Repo)
		}

		var execCommand string
		switch lang := strings.ToLower(env.Language); lang {
		case "go", "golang":
			execCommand = "go run main.go"
		case "csharp", "c#", "net", ".net", "aspnetcore", "kestrel", "netcore", "net.core", ".net core":
			execCommand = "dotnet run -c Release"
		case "node", "nodejs", "javascript", "js":
			execCommand = "npm install\nnode ."
		default:
			return fmt.Errorf("%s:%s unsupported language: %s", t.Name, env.Repo, lang)
		}

		env.Exec = execCommand
	}

	// build the bench command before server ran.
	args := t.buildArgs()
	bombardierCommand := "bombardier " + strings.Join(args, " ") // for logging.

	benchCmd := exec.Command("bombardier", args...)
	buf := new(bytes.Buffer)
	benchCmd.Stdout = buf

	fmt.Fprintf(os.Stdout, "[%s]\n", env.absDir)

	commandsToRun := strings.Split(env.Exec, "\n")
	for i, commandToRun := range commandsToRun {
		if commandToRun == "" {
			continue
		}

		fmt.Fprintf(os.Stdout, "$ %s\n", commandToRun)

		cmd := newCmd(commandToRun)
		cmd.Dir = env.absDir
		watchCmd(cmd)

		// last command should be the server, which blocks, so don't wait for it.
		shouldWait := i < len(commandsToRun)-1
		if shouldWait {
			err = cmd.Run()
		} else {
			err = cmd.Start()
			defer func() {
				if killErr := killCmd(cmd); killErr != nil {
					fmt.Fprintf(os.Stdout, "kill: error: %v\n", killErr)
				}
			}()
		}

		if err != nil {
			return err
		}
	}

	// time.Sleep(*waitServerDur)
	waitServer(t.URL)

	fmt.Fprintf(os.Stdout, "$ %s\n", bombardierCommand)

	err = benchCmd.Run()
	if err != nil {
		return fmt.Errorf("%s\n%s", buf.String(), err)
	}

	if err = json.NewDecoder(buf).Decode(env); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	return nil
}

func benchmark(t *Test) error {
	for _, env := range t.Envs {
		if !env.CanBenchmark() {
			continue
		}

		if err := runBenchmark(t, env); err != nil {
			return fmt.Errorf("run benchmark :%w", err)
		}

		var httpErrors []string
		for _, httpErr := range env.Result.Errors {
			httpErrors = append(httpErrors, httpErr.Description)
		}

		if len(httpErrors) > 0 {
			return fmt.Errorf("%s:%s errors: %s", t.Name, env.Repo, strings.Join(httpErrors, ", "))
		}

		if expected, got := uint64(0), env.Result.Req4XX; expected != got {
			return fmt.Errorf("%s:%s failed with %d 404 Not Found codes", t.Name, env.Repo, got)
		}

		if expected, got := t.NumberOfRequests, env.Result.Req2XX; expected > 0 && expected != got {
			return fmt.Errorf("%s:%s expected successful requests: %d but got: %d [3xx: %d, 4xx: %d, 5xx: %d, others: %d, errors: %s]",
				t.Name, env.Repo, expected, got,
				env.Result.Req3XX,
				env.Result.Req4XX,
				env.Result.Req5XX,
				env.Result.Others,
				strings.Join(httpErrors, ", "),
			)
		}

		// fmt.Fprintf(os.Stdout, "> sleeping for %s...\n", waitRunDur.String())
		time.Sleep(*waitRunDur)
	}

	sort.Slice(t.Envs, func(i, j int) bool {
		if !t.Envs[i].CanBenchmark() {
			return false
		}

		return t.Envs[i].Result.RequestsPerSecond.Mean > t.Envs[j].Result.RequestsPerSecond.Mean
	})

	fmt.Fprintf(os.Stdout, "[%s] Winner: %s!\n\n", t.Name, t.Envs[0].GetName())

	return nil
}

func newCmd(command string) *exec.Cmd {
	full := strings.Split(command, " ")
	name := full[0]
	var args []string
	if len(full) > 1 {
		args = full[1:]
	}

	cmd := exec.Command(name, args...)
	wrapCmd(cmd)
	return cmd
}

func watchCmd(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
}

func waitServer(rawURL string) (int, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return -1, err
	}

	maxTries := 10
	tries := 1
	timeout := time.Duration(2 * time.Second)

	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for range ticker.C {
		if tries >= maxTries {
			break
		}

		conn, err := net.DialTimeout("tcp4", u.Host, timeout)
		if err == nil {
			conn.Close()
			fmt.Fprintf(os.Stdout, "> Connected to %s\n", u.Host)
			return tries, nil
		}
		tries++
	}

	return maxTries, fmt.Errorf("unreachable %s", rawURL)
}
