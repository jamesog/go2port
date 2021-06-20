package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

func TestReverseDepByName(t *testing.T) {
	deps := []Dependency{
		{"foo.com/C/bb", ""},
		{"foo.com/a/b", ""},
		{"foo.com/a/bb", ""},
	}
	sort.Sort(reverseDepByName(deps))

	if deps[0].Name != "foo.com/a/bb" {
		t.Error("expected foo.com/a/bb in index 0")
	}
}

func TestModuleDependencies(t *testing.T) {
	tests := map[string]struct {
		modDir   string
		wantDeps int
	}{
		"8_deps":  {"module/8deps", 8},
		"no_deps": {"module/nodeps", 0},
	}

	pkg, _ := newPackage("github.com/owner/project", "v0.1.0")

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gomod, _ := os.Open(filepath.Join("testdata", tt.modDir, "go.mod"))
			gosum, _ := os.Open(filepath.Join("testdata", tt.modDir, "go.sum"))
			gock.New("https://raw.githubusercontent.com/owner/project/v0.1.0").
				Get("/go.mod").
				Reply(200).
				Body(gomod)
			gock.New("https://raw.githubusercontent.com/owner/project/v0.1.0").
				Get("/go.sum").
				Reply(200).
				Body(gosum)
			defer gock.Off()

			deps, err := moduleDependenciesNew(pkg, "/")
			if err != nil {
				t.Fatalf("expected no error; got %v", err)
			}
			if len(deps) != tt.wantDeps {
				t.Errorf("expected %d dependencies; got %d", tt.wantDeps, len(deps))
			}
		})
	}

	t.Run("no go.mod", func(t *testing.T) {
		gock.New("https://raw.githubusercontent.com/owner/project/v0.1.0").
			Get("/go.mod").
			Reply(404)
		defer gock.Off()

		_, err := moduleDependenciesNew(pkg, "/")
		if err == nil {
			t.Errorf("expected an error; got none")
		}
	})
}
