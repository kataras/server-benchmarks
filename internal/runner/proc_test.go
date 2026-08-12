package runner

import (
	"reflect"
	"testing"
)

func TestSplitCommand(t *testing.T) {
	cases := []struct {
		in       string
		wantName string
		wantArgs []string
		wantErr  bool
	}{
		{in: "go run .", wantName: "go", wantArgs: []string{"run", "."}},
		{in: "node .", wantName: "node", wantArgs: []string{"."}},
		{in: `"C:\Program Files\dotnet\dotnet" run -c Release`, wantName: `C:\Program Files\dotnet\dotnet`, wantArgs: []string{"run", "-c", "Release"}},
		{in: `app --flag 'a value with spaces'`, wantName: "app", wantArgs: []string{"--flag", "a value with spaces"}},
		{in: `app --json '{"a": 1}'`, wantName: "app", wantArgs: []string{"--json", `{"a": 1}`}},
		{in: `app ""`, wantName: "app", wantArgs: []string{""}},
		{in: "  spaced   out  ", wantName: "spaced", wantArgs: []string{"out"}},
		{in: `broken "quote`, wantErr: true},
		{in: "", wantErr: true},
		{in: "   ", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			name, args, err := splitCommand(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("splitCommand(%q): expected error, got %q %q", tc.in, name, args)
				}
				return
			}

			if err != nil {
				t.Fatalf("splitCommand(%q): %v", tc.in, err)
			}
			if name != tc.wantName {
				t.Errorf("name = %q, want %q", name, tc.wantName)
			}
			if !reflect.DeepEqual(args, tc.wantArgs) {
				t.Errorf("args = %q, want %q", args, tc.wantArgs)
			}
		})
	}
}

func TestExecLines(t *testing.T) {
	got := execLines("npm install\r\n\n  node .  \n")
	want := []string{"npm install", "node ."}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("execLines = %q, want %q", got, want)
	}

	if got := execLines(""); got != nil {
		t.Errorf("execLines(\"\") = %q, want nil", got)
	}
}

func TestTail(t *testing.T) {
	tl := newTail(8)
	tl.Write([]byte("0123456789"))
	if got, want := tl.String(), "23456789"; got != want {
		t.Errorf("tail = %q, want %q", got, want)
	}

	tl.Write([]byte("ab"))
	if got, want := tl.String(), "456789ab"; got != want {
		t.Errorf("tail after second write = %q, want %q", got, want)
	}
}

func TestTargetAddr(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{in: "http://localhost:5000", want: "localhost:5000"},
		{in: "http://localhost:5000/hello/world", want: "localhost:5000"},
		{in: "http://127.0.0.1:8080/x", want: "127.0.0.1:8080"},
		{in: "http://example.com/x", want: "example.com:80"},
		{in: "https://example.com", want: "example.com:443"},
		{in: "not a url at all", wantErr: true},
	}

	for _, tc := range cases {
		got, err := targetAddr(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("targetAddr(%q): expected error, got %q", tc.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("targetAddr(%q): %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("targetAddr(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
