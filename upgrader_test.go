package main

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_updatePackages(t *testing.T) {
	dir, err := setupTestProject(t, "a")
	if err != nil {
		t.Fatal(err)
	}

	err = updatePackages(dir, "github.com/google/go-github", "v20", false)
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(filepath.Clean(dir), "main.go"))
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != sampleMain20 {
		t.Errorf("got diffs:\n%s", quickDiff(string(content), sampleMain20))
	}
}

func Test_createNewImport(t *testing.T) {
	testCases := []struct {
		desc       string
		parts      []string
		newVersion string
		expected   string
	}{
		{
			desc:       "no version",
			parts:      []string{"github.com", "ldez", "foobar"},
			newVersion: "v2",
			expected:   "github.com/ldez/foobar/v2",
		},
		{
			desc:       "version",
			parts:      []string{"github.com", "ldez", "foobar", "v1"},
			newVersion: "v2",
			expected:   "github.com/ldez/foobar/v2",
		},
		{
			desc:       "no version and subpackage",
			parts:      []string{"github.com", "ldez", "foobar", "foo"},
			newVersion: "v2",
			expected:   "github.com/ldez/foobar/v2/foo",
		},
		{
			desc:       "version and subpackage",
			parts:      []string{"github.com", "ldez", "foobar", "v1", "foo"},
			newVersion: "v2",
			expected:   "github.com/ldez/foobar/v2/foo",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			newImport, err := createNewImport(test.parts, test.newVersion)
			if err != nil {
				t.Fatal(err)
			}

			if newImport != test.expected {
				t.Errorf("got %s, want %s", newImport, test.expected)
			}
		})
	}
}

func Test_updateModFile(t *testing.T) {
	testCases := []struct {
		desc     string
		dir      string
		expected string
	}{
		{
			desc:     "one block",
			dir:      "a",
			expected: sampleGoMod20,
		},
		{
			desc:     "two blocks",
			dir:      "b",
			expected: sampleGoMod20Blocks,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			dir, err := setupTestProject(t, test.dir)
			if err != nil {
				t.Fatal(err)
			}

			err = updateModFile(dir, "github.com/google/go-github", "v20.0.0", "v20")
			if err != nil {
				t.Fatal(err)
			}

			mod, err := os.ReadFile(filepath.Join(filepath.Clean(dir), "go.mod"))
			if err != nil {
				t.Fatal(err)
			}

			t.Log(string(mod))

			if string(mod) != test.expected {
				t.Errorf("got diffs:\n%s", quickDiff(string(mod), test.expected))
			}
		})
	}
}

func Test_guessVersion(t *testing.T) {
	type expected struct {
		Major string
		Full  string
	}

	testCases := []struct {
		desc           string
		baseModuleName string
		raw            string
		expected       expected
	}{
		{
			desc:           "only number",
			baseModuleName: "github.com/google/go-github",
			raw:            "28",
			expected: expected{
				Major: "v28",
				Full:  "v28.1.1",
			},
		},
		{
			desc:           "version prefixed by v",
			baseModuleName: "github.com/google/go-github",
			raw:            "v30",
			expected: expected{
				Major: "v30",
				Full:  "v30.1.0",
			},
		},
		{
			desc:           "latest major (go-github)",
			baseModuleName: "github.com/google/go-github",
			raw:            "latest",
			expected: expected{
				Major: "v71",
				Full:  "v71.0.0",
			},
		},
		{
			desc:           "latest major (backoff)",
			baseModuleName: "github.com/cenkalti/backoff",
			raw:            "latest",
			expected: expected{
				Major: "v5",
				Full:  "v5.0.2",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			full, major, err := guessVersion(test.baseModuleName, false, test.raw)
			if err != nil {
				t.Fatal(err)
			}

			if full != test.expected.Full {
				t.Errorf("Got: %s, want: %s", full, test.expected.Full)
			}

			if major != test.expected.Major {
				t.Errorf("Got: %s, want: %s", major, test.expected.Major)
			}
		})
	}
}
