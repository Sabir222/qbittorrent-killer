package mask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	m := Mask{0b01010100, 0b01010100}
	expected := []bool{false, true, false, true, false, true, false, false, false, true, false, true, false, true, false, false, false, false, false, false}
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], m.Check(i))
	}
}

func TestMark(t *testing.T) {
	cases := []struct {
		data  Mask
		idx   int
		want  Mask
	}{
		{
			data: Mask{0b01010100, 0b01010100},
			idx:  4,
			want: Mask{0b01011100, 0b01010100},
		},
		{
			data: Mask{0b01010100, 0b01010100},
			idx:  9,
			want: Mask{0b01010100, 0b01010100},
		},
		{
			data: Mask{0b01010100, 0b01010100},
			idx:  15,
			want: Mask{0b01010100, 0b01010101},
		},
		{
			data: Mask{0b01010100, 0b01010100},
			idx:  19,
			want: Mask{0b01010100, 0b01010100},
		},
	}
	for _, c := range cases {
		m := c.data
		m.Mark(c.idx)
		assert.Equal(t, c.want, m)
	}
}
