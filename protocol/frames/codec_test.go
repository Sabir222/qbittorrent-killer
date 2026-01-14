package frames

import (
	"bytes"
	"testing"
)

func TestFrame_Pack(t *testing.T) {
	f := &Frame{Type: TypeHave, Data: []byte{0, 0, 0, 1}}
	packed := f.Pack()
	if len(packed) != 9 {
		t.Errorf("expected 9 bytes, got %d", len(packed))
	}
}

func TestFrame_Unpack(t *testing.T) {
	data := []byte{0, 0, 0, 5, 4, 0, 0, 0, 1}
	r := bytes.NewReader(data)
	f, err := Unpack(r)
	if err != nil {
		t.Fatal(err)
	}
	if f.Type != TypeHave {
		t.Errorf("expected TypeHave, got %d", f.Type)
	}
}

func TestCodec_Request(t *testing.T) {
	codec := NewCodec()
	f := codec.Request(1, 0, 16384)
	if f.Type != TypeRequest {
		t.Errorf("expected TypeRequest, got %d", f.Type)
	}
}
