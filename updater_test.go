package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_updatePackages(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		// Because "GetSizesGolist" doesn't work well on Travis.
		// https://github.com/golang/tools/blob/16909d206f00da7d0d5ba28cd9dc7fb223648ecf/go/internal/packagesdriver/sizes.go#L80
		t.Skipf("TRAVIS=true")
	}

	dir, err := setupTestProject(t)
	if err != nil {
		t.Fatal(err)
	}

	err = updatePackages(dir, "github.com/google/go-github", "v20", false)
	if err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadFile(filepath.Join(filepath.Clean(dir), "main.go"))
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
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			newImport := createNewImport(test.parts, test.newVersion)

			if newImport != test.expected {
				t.Errorf("got %s, want %s", newImport, test.expected)
			}
		})
	}
}

func Test_updateModFile(t *testing.T) {
	dir, err := setupTestProject(t)
	if err != nil {
		t.Fatal(err)
	}

	err = updateModFile(dir, "github.com/google/go-github", "v20.0.0", "v20")
	if err != nil {
		t.Fatal(err)
	}

	mod, err := ioutil.ReadFile(filepath.Join(filepath.Clean(dir), "go.mod"))
	if err != nil {
		t.Fatal(err)
	}

	if string(mod) != sampleGoMod20 {
		t.Errorf("got diffs:\n%s", quickDiff(string(mod), sampleGoMod20))
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
			desc:           "",
			baseModuleName: "github.com/google/go-github",
			raw:            "latest",
			expected: expected{
				Major: "v32",
				Full:  "v32.1.0",
			},
		},
		{
			desc:           "",
			baseModuleName: "github.com/cenkalti/backoff",
			raw:            "latest",
			expected: expected{
				Major: "v4",
				Full:  "v4.0.2",
			},
		},
	}

	for _, test := range testCases {
		test := test
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
