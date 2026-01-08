package mask

import (
	"testing"
)

func TestMask_Check(t *testing.T) {
	m := make(Mask, 1)
	m.Mark(0)
	if !m.Check(0) {
		t.Error("expected bit 0 to be set")
	}
	if m.Check(1) {
		t.Error("expected bit 1 to be unset")
	}
}

func TestMask_Mark(t *testing.T) {
	m := make(Mask, 1)
	m.Mark(3)
	if !m.Check(3) {
		t.Error("expected bit 3 to be set")
	}
}
