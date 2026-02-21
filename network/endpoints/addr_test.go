package endpoints

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := map[string]struct {
		raw   string
		want  []Endpoint
		fail  bool
	}{
		"valid peer list": {
			raw: string([]byte{127, 0, 0, 1, 0x00, 0x50, 1, 1, 1, 1, 0x01, 0xbb}),
			want: []Endpoint{
				{Addr: net.IP{127, 0, 0, 1}, Port: 80},
				{Addr: net.IP{1, 1, 1, 1}, Port: 443},
			},
		},
		"incomplete peer entry": {
			raw:  string([]byte{127, 0, 0, 1, 0x00}),
			want: nil,
			fail: true,
		},
	}

	for _, c := range cases {
		got, err := Parse([]byte(c.raw))
		if c.fail {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, c.want, got)
	}
}

func TestEndpointString(t *testing.T) {
	cases := []struct {
		ep   Endpoint
		want string
	}{
		{
			ep:   Endpoint{Addr: net.IP{127, 0, 0, 1}, Port: 8080},
			want: "127.0.0.1:8080",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, c.ep.String())
	}
}
