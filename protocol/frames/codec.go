package frames

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	TypeChoke         = iota
	TypeUnchoke
	TypeInterested
	TypeNotInterested
	TypeHave
	TypeBitfield
	TypeRequest
	TypePiece
	TypeCancel
)

var (
	ErrInvalidType     = errors.New("invalid message type")
	ErrPayloadTooShort = errors.New("payload too short")
	ErrIndexMismatch   = errors.New("index mismatch")
	ErrOffsetTooHigh   = errors.New("offset too high")
	ErrDataTooLong     = errors.New("data too long")
)

type Frame struct {
	Type uint8
	Data []byte
}

type Codec struct{}

func NewCodec() *Codec {
	return &Codec{}
}

func (c *Codec) Request(pieceIdx, offset, size int) *Frame {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf[0:4], uint32(pieceIdx))
	binary.BigEndian.PutUint32(buf[4:8], uint32(offset))
	binary.BigEndian.PutUint32(buf[8:12], uint32(size))
	return &Frame{Type: TypeRequest, Data: buf}
}

func (c *Codec) Have(pieceIdx int) *Frame {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(pieceIdx))
	return &Frame{Type: TypeHave, Data: buf}
}

func ReadPieceData(target []byte, pieceIdx int, frm *Frame) (int, error) {
	if frm.Type != TypePiece {
		return 0, ErrInvalidType
	}
	if len(frm.Data) < 8 {
		return 0, ErrPayloadTooShort
	}

	idx := int(binary.BigEndian.Uint32(frm.Data[0:4]))
	if idx != pieceIdx {
		return 0, ErrIndexMismatch
	}

	offset := int(binary.BigEndian.Uint32(frm.Data[4:8]))
	if offset >= len(target) {
		return 0, ErrOffsetTooHigh
	}

	payload := frm.Data[8:]
	if offset+len(payload) > len(target) {
		return 0, ErrDataTooLong
	}

	copy(target[offset:], payload)
	return len(payload), nil
}

func ReadHave(frm *Frame) (int, error) {
	if frm.Type != TypeHave {
		return 0, ErrInvalidType
	}
	if len(frm.Data) != 4 {
		return 0, ErrPayloadTooShort
	}
	return int(binary.BigEndian.Uint32(frm.Data)), nil
}

func (f *Frame) Pack() []byte {
	if f == nil {
		return []byte{0, 0, 0, 0}
	}

	size := uint32(len(f.Data) + 1)
	buf := make([]byte, 4+size)
	binary.BigEndian.PutUint32(buf, size)
	buf[4] = f.Type
	copy(buf[5:], f.Data)
	return buf
}

func Unpack(r io.Reader) (*Frame, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}

	msgLen := binary.BigEndian.Uint32(lenBuf[:])
	if msgLen == 0 {
		return nil, nil
	}

	data := make([]byte, msgLen)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	return &Frame{
		Type: data[0],
		Data: data[1:],
	}, nil
}

func (f *Frame) Label() string {
	if f == nil {
		return "KeepAlive"
	}

	names := []string{
		"Choke", "Unchoke", "Interested", "NotInterested",
		"Have", "Bitfield", "Request", "Piece", "Cancel",
	}

	if int(f.Type) < len(names) {
		return names[f.Type]
	}
	return "Unknown"
}

func (f *Frame) String() string {
	if f == nil {
		return "KeepAlive"
	}
	return f.Label()
}
