package endpoints

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

const peerEntrySize = 6

var ErrMalformedPeerData = errors.New("malformed peer data")

type Endpoint struct {
	Addr net.IP
	Port uint16
}

func Parse(raw []byte) ([]Endpoint, error) {
	if len(raw)%peerEntrySize != 0 {
		return nil, ErrMalformedPeerData
	}

	count := len(raw) / peerEntrySize
	result := make([]Endpoint, count)

	for i := 0; i < count; i++ {
		pos := i * peerEntrySize
		result[i].Addr = net.IP(raw[pos : pos+4])
		result[i].Port = binary.BigEndian.Uint16(raw[pos+4 : pos+6])
	}

	return result, nil
}

func (e Endpoint) String() string {
	return net.JoinHostPort(e.Addr.String(), strconv.Itoa(int(e.Port)))
}
