package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_updateCmd(t *testing.T) {
	testCases := []struct {
		desc     string
		version  string
		expected string
	}{
		{
			desc:     "only number",
			version:  "20",
			expected: sampleGoMod20,
		},
		{
			desc:     "version prefixed by v",
			version:  "v20",
			expected: sampleGoMod20,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {

			dir, cleanUp, err := setupTestProject()
			defer cleanUp()
			if err != nil {
				t.Fatal(err)
			}

			if os.Chdir(dir) != nil {
				t.Fatal(err)
			}

			updateCmd(false, false, "github.com/google/go-github", test.version)

			content, err := ioutil.ReadFile(filepath.Join(dir, "main.go"))
			if err != nil {
				t.Fatal(err)
			}

			if string(content) != test.expected {
				t.Errorf("got diffs:\n%s", quickDiff(string(content), test.expected))
			}
		})
	}
}
