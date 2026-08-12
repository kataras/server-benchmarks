package spec

import (
	"reflect"
	"testing"
)

func filterFixture() []Test {
	return []Test{
		{
			Name: "Static",
			Envs: []Env{{Name: "Iris"}, {Name: "Gin"}, {Name: "ASP.NET Core", Dir: "/repo/_code/static/aspnetcore"}},
		},
		{
			Name: "REST",
			Envs: []Env{{Name: "Iris"}, {Name: "Fastify"}},
		},
	}
}

func envNames(t Test) []string {
	names := make([]string, len(t.Envs))
	for i, e := range t.Envs {
		names[i] = e.Name
	}
	return names
}

func TestFilter(t *testing.T) {
	cases := []struct {
		name          string
		selectors     []string
		wantTests     []string
		wantEnvs      map[string][]string
		wantUnmatched []string
	}{
		{
			name:      "no selectors keeps everything",
			selectors: nil,
			wantTests: []string{"Static", "REST"},
			wantEnvs:  map[string][]string{"Static": {"Iris", "Gin", "ASP.NET Core"}, "REST": {"Iris", "Fastify"}},
		},
		{
			name:      "whole test by name",
			selectors: []string{"rest"},
			wantTests: []string{"REST"},
			wantEnvs:  map[string][]string{"REST": {"Iris", "Fastify"}},
		},
		{
			name:      "single env",
			selectors: []string{"REST.Iris"},
			wantTests: []string{"REST"},
			wantEnvs:  map[string][]string{"REST": {"Iris"}},
		},
		{
			name:      "env name containing dots",
			selectors: []string{"static.asp.net core"},
			wantTests: []string{"Static"},
			wantEnvs:  map[string][]string{"Static": {"ASP.NET Core"}},
		},
		{
			name:      "env matched by directory slug",
			selectors: []string{"static.aspnetcore"},
			wantTests: []string{"Static"},
			wantEnvs:  map[string][]string{"Static": {"ASP.NET Core"}},
		},
		{
			name:      "multiple selectors across tests",
			selectors: []string{"static.gin", "rest.fastify"},
			wantTests: []string{"Static", "REST"},
			wantEnvs:  map[string][]string{"Static": {"Gin"}, "REST": {"Fastify"}},
		},
		{
			name:      "whole test wins over env selector",
			selectors: []string{"static", "static.gin"},
			wantTests: []string{"Static"},
			wantEnvs:  map[string][]string{"Static": {"Iris", "Gin", "ASP.NET Core"}},
		},
		{
			name:          "unmatched selector is reported",
			selectors:     []string{"static.gin", "nope", "rest.vapor"},
			wantTests:     []string{"Static"},
			wantEnvs:      map[string][]string{"Static": {"Gin"}},
			wantUnmatched: []string{"nope", "rest.vapor"},
		},
		{
			name:          "nothing matches",
			selectors:     []string{"bogus"},
			wantTests:     nil,
			wantUnmatched: []string{"bogus"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := filterFixture()
			got, unmatched := Filter(input, tc.selectors)

			var gotNames []string
			for _, tt := range got {
				gotNames = append(gotNames, tt.Name)
			}
			if !reflect.DeepEqual(gotNames, tc.wantTests) {
				t.Fatalf("tests = %v, want %v", gotNames, tc.wantTests)
			}

			for _, tt := range got {
				if want, ok := tc.wantEnvs[tt.Name]; ok {
					if got := envNames(tt); !reflect.DeepEqual(got, want) {
						t.Errorf("%s envs = %v, want %v", tt.Name, got, want)
					}
				}
			}

			if !reflect.DeepEqual(unmatched, tc.wantUnmatched) {
				t.Errorf("unmatched = %v, want %v", unmatched, tc.wantUnmatched)
			}

			// The input must never be mutated.
			if !reflect.DeepEqual(input, filterFixture()) {
				t.Error("Filter mutated its input")
			}
		})
	}
}
