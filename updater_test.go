package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_update(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		// Because "GetSizesGolist" doesn't work well on Travis.
		// https://github.com/golang/tools/blob/16909d206f00da7d0d5ba28cd9dc7fb223648ecf/go/internal/packagesdriver/sizes.go#L80
		t.Skipf("TRAVIS=true")
	}

	dir, err := setupTestProject(t)
	if err != nil {
		t.Fatal(err)
	}

	err = update(dir, "github.com/google/go-github", "v20", false)
	if err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadFile(filepath.Join(filepath.Clean(dir), "main.go"))
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != sampleGoMod20 {
		t.Errorf("got diffs:\n%s", quickDiff(string(content), sampleGoMod20))
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
