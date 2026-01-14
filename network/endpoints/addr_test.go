package endpoints

import (
	"net"
	"testing"
)

func TestParse(t *testing.T) {
	raw := []byte{127, 0, 0, 1, 0x1f, 0x90}
	endpoints, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(endpoints))
	}
	if !endpoints[0].Addr.Equal(net.IPv4(127, 0, 0, 1)) {
		t.Errorf("expected 127.0.0.1, got %v", endpoints[0].Addr)
	}
	if endpoints[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", endpoints[0].Port)
	}
}

func TestParse_Malformed(t *testing.T) {
	raw := []byte{127, 0, 0, 1, 0x1f}
	_, err := Parse(raw)
	if err != ErrMalformedPeerData {
		t.Errorf("expected ErrMalformedPeerData, got %v", err)
	}
}
