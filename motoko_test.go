package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_updateCmd(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		// Because "GetSizesGolist" doesn't work well on Travis.
		// https://github.com/golang/tools/blob/16909d206f00da7d0d5ba28cd9dc7fb223648ecf/go/internal/packagesdriver/sizes.go#L80
		t.Skipf("TRAVIS=true")
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
		test := test
		t.Run(test.desc, func(t *testing.T) {
			dir, err := setupTestProject(t)
			if err != nil {
				t.Fatal(err)
			}

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

			content, err := ioutil.ReadFile(filepath.Join(filepath.Clean(dir), "main.go"))
			if err != nil {
				t.Fatal(err)
			}

			if string(content) != test.expected.code {
				t.Errorf("got diffs:\n%s", quickDiff(string(content), test.expected.code))
			}

			mod, err := ioutil.ReadFile(filepath.Join(filepath.Clean(dir), "go.mod"))
			if err != nil {
				t.Fatal(err)
			}

			if string(mod) != test.expected.mod {
				t.Errorf("got diffs:\n%s", quickDiff(string(mod), test.expected.mod))
			}
		})
	}
}
