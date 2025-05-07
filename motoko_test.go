package main

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"
)

func Test_updateCmd(t *testing.T) {
	if os.Getenv("CI") == "true" && build.Default.GOOS == "windows" {
		// Because the cleanup doesn't work well on Windows in GitHub Actions
		// The process cannot access the file because it is being used by another process.
		t.Skipf("Windows and GitHub Actions")
	}

	type expected struct {
		code string
		mod  string
	}

	testCases := []struct {
		desc     string
		version  string
		expected expected
	}{
		{
			desc:    "only number",
			version: "20",
			expected: expected{
				code: sampleMain20,
				mod:  sampleGoMod20,
			},
		},
		{
			desc:    "version prefixed by v",
			version: "v20",
			expected: expected{
				code: sampleMain20,
				mod:  sampleGoMod20,
			},
		},
		{
			desc:    "full version",
			version: "v20.0.0",
			expected: expected{
				code: sampleMain20,
				mod:  sampleGoMod20,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			dir, err := setupTestProject(t, "a")
			if err != nil {
				t.Fatal(err)
			}

			wd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = os.Chdir(wd) })

			if os.Chdir(dir) != nil {
				t.Fatal(err)
			}

			cfg := config{
				lib:     "github.com/google/go-github",
				version: test.version,
			}

			err = updateCmd(cfg)
			if err != nil {
				t.Fatal(err)
			}

			content, err := os.ReadFile(filepath.Join(filepath.Clean(dir), "main.go"))
			if err != nil {
				t.Fatal(err)
			}

			if string(content) != test.expected.code {
				t.Errorf("got diffs:\n%s", quickDiff(string(content), test.expected.code))
			}

			mod, err := os.ReadFile(filepath.Join(filepath.Clean(dir), "go.mod"))
			if err != nil {
				t.Fatal(err)
			}

			if string(mod) != test.expected.mod {
				t.Log(string(mod))
				t.Errorf("got diffs:\n%s", quickDiff(string(mod), test.expected.mod))
			}
		})
	}
}
