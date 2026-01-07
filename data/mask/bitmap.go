package mask

type Mask []byte

func (m Mask) Check(idx int) bool {
	if idx < 0 {
		return false
	}
	byteIdx := idx >> 3
	if byteIdx >= len(m) {
		return false
	}
	bitPos := 7 - (idx & 7)
	return (m[byteIdx]>>bitPos)&1 == 1
}

func (m Mask) Mark(idx int) {
	if idx < 0 {
		return
	}
	byteIdx := idx >> 3
	if byteIdx >= len(m) {
		return
	}
	bitPos := 7 - (idx & 7)
	m[byteIdx] |= 1 << bitPos
}
