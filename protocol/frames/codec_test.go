package frames

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPack(t *testing.T) {
	tests := map[string]struct {
		input  *Frame
		output []byte
	}{
		"normal": {
			input:  &Frame{Type: TypeHave, Data: []byte{1, 2, 3, 4}},
			output: []byte{0, 0, 0, 5, 4, 1, 2, 3, 4},
		},
		"keepalive": {
			input:  nil,
			output: []byte{0, 0, 0, 0},
		},
	}

	for _, test := range tests {
		buf := test.input.Pack()
		assert.Equal(t, test.output, buf)
	}
}

func TestUnpack(t *testing.T) {
	tests := map[string]struct {
		input  []byte
		output *Frame
		fails  bool
	}{
		"normal": {
			input:  []byte{0, 0, 0, 5, 4, 1, 2, 3, 4},
			output: &Frame{Type: TypeHave, Data: []byte{1, 2, 3, 4}},
			fails:  false,
		},
		"keepalive": {
			input:  []byte{0, 0, 0, 0},
			output: nil,
			fails:  false,
		},
		"too short": {
			input:  []byte{1, 2, 3},
			output: nil,
			fails:  true,
		},
		"incomplete": {
			input:  []byte{0, 0, 0, 5, 4, 1, 2},
			output: nil,
			fails:  true,
		},
	}

	for _, test := range tests {
		r := bytes.NewReader(test.input)
		m, err := Unpack(r)
		if test.fails {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.output, m)
	}
}

func TestReadPieceData(t *testing.T) {
	buf := make([]byte, 10)
	frm := &Frame{
		Type: TypePiece,
		Data: []byte{
			0x00, 0x00, 0x00, 0x04,
			0x00, 0x00, 0x00, 0x02,
			0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
		},
	}

	n, err := ReadPieceData(buf, 4, frm)
	assert.NoError(t, err)
	assert.Equal(t, 6, n)
	assert.Equal(t, []byte{0x00, 0x00, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x00}, buf)
}
