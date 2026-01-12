package greeting

import (
	"errors"
	"io"
)

const reservedBytes = 8
const hashSize = 20

var ErrInvalidProtocolLen = errors.New("invalid protocol length")

type Greeting struct {
	Protocol string
	Hash     [hashSize]byte
	ID       [hashSize]byte
}

func Build(hash, id [hashSize]byte) *Greeting {
	return &Greeting{
		Protocol: "BitTorrent protocol",
		Hash:     hash,
		ID:       id,
	}
}

func (g *Greeting) Pack() []byte {
	pstrLen := byte(len(g.Protocol))
	out := make([]byte, 1+len(g.Protocol)+reservedBytes+hashSize+hashSize)

	out[0] = pstrLen
	pos := 1

	pos += copy(out[pos:], g.Protocol)
	pos += copy(out[pos:], make([]byte, reservedBytes))
	copy(out[pos:], append(g.Hash[:], g.ID[:]...))

	return out
}

func Unpack(stream io.Reader) (*Greeting, error) {
	var lenBuf [1]byte
	if _, err := io.ReadFull(stream, lenBuf[:]); err != nil {
		return nil, err
	}

	pstrLen := int(lenBuf[0])
	if pstrLen == 0 {
		return nil, ErrInvalidProtocolLen
	}

	totalLen := pstrLen + reservedBytes + hashSize + hashSize
	raw := make([]byte, totalLen)
	if _, err := io.ReadFull(stream, raw); err != nil {
		return nil, err
	}

	var result Greeting
	result.Protocol = string(raw[:pstrLen])

	hashStart := pstrLen + reservedBytes
	copy(result.Hash[:], raw[hashStart:hashStart+hashSize])
	copy(result.ID[:], raw[hashStart+hashSize:])

	return &result, nil
}
