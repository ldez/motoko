package internal

import (
	"context"
	"testing"
)

func TestFindHighestFromDepsDev(t *testing.T) {
	testCases := []struct {
		desc string

		expected string
	}{
		{
			desc:     "github.com/akamai/AkamaiOPEN-edgegrid-golang",
			expected: "v13.3.0",
		},
		{
			desc:     "github.com/vultr/govultr",
			expected: "v3.31.2",
		},
		{
			desc:     "gopkg.in/yaml.v2",
			expected: "v3.0.1",
		},
		{
			desc:     "github.com/namedotcom/go/v4",
			expected: "v4.0.2",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			highest, err := FindHighestFromDepsDev(context.Background(), test.desc)
			if err != nil {
				t.Fatal(err)
			}

			if highest != test.expected {
				t.Errorf("got %s, want %s", highest, test.expected)
			}
		})
	}
}

func TestFindHighestFromDepsDev_notFound(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "github.com/go-viper/mapstructure/v6",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			_, err := FindHighestFromDepsDev(context.Background(), test.desc)
			if err == nil {
				t.Fatal("No error returned")
			}
		})
	}
}
