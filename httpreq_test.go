package httpreq

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	joinURLTestCase struct {
		https    bool
		host     string
		port     string
		uri      string
		expected string
	}
)

func TestMakeBaseURLWithHTTP(t *testing.T) {
	expected := "http://localhost:8080"
	cfg := Config{
		HTTPS: false,
		Host:  "localhost",
		Port:  "8080",
	}

	url := makeBaseURL(cfg)
	assert.Equal(t, expected, url)
}

var (
	joinURLTestCases = []joinURLTestCase{
		{
			https:    false,
			host:     "localhost",
			port:     "",
			uri:      "/hello/world",
			expected: "http://localhost/hello/world",
		},
		{
			https:    true,
			host:     "localhost",
			port:     "",
			uri:      "/hello/world",
			expected: "https://localhost/hello/world",
		},
		{
			https:    false,
			host:     "localhost",
			port:     "8080",
			uri:      "/hello/world",
			expected: "http://localhost:8080/hello/world",
		},
		{
			https:    true,
			host:     "localhost",
			port:     "8080",
			uri:      "//hello/world/",
			expected: "https://localhost:8080/hello/world",
		},
		{
			https:    false,
			host:     "localhost",
			port:     "8080",
			uri:      "/hello//world",
			expected: "http://localhost:8080/hello/world",
		},
	}
)

func TestURLConstruction(t *testing.T) {
	for i, testCase := range joinURLTestCases {
		tc := testCase
		desc := fmt.Sprintf("%vth run", i)
		t.Run(desc, func(t *testing.T) {
			cfg := Config{
				HTTPS: tc.https,
				Host:  tc.host,
				Port:  tc.port,
			}
			base := makeBaseURL(cfg)
			res := joinURL(base, tc.uri)
			assert.Equal(t, tc.expected, res)
		})
	}
}
