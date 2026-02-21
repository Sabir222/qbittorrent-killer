package connector

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/Sabir222/torrent-at-home/data/mask"
	"github.com/Sabir222/torrent-at-home/protocol/greeting"
	"github.com/Sabir222/torrent-at-home/protocol/frames"
	"github.com/Sabir222/torrent-at-home/network/endpoints"
)

type PeerConn struct {
	Conn     net.Conn
	Choked   bool
	Bitfield mask.Mask
	peer     endpoints.Endpoint
	infoHash [20]byte
	peerID   [20]byte
}

func doHandshake(conn net.Conn, infohash, peerID [20]byte) (*greeting.Greeting, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	outbound := greeting.Build(infohash, peerID)
	if _, err := conn.Write(outbound.Pack()); err != nil {
		return nil, err
	}

	incoming, err := greeting.Unpack(conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(incoming.Hash[:], infohash[:]) {
		return nil, fmt.Errorf("infohash does not match")
	}

	return incoming, nil
}

func readBitfield(conn net.Conn) (mask.Mask, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	frm, err := frames.Unpack(conn)
	if err != nil {
		return nil, err
	}

	if frm == nil || frm.Type != frames.TypeBitfield {
		return nil, fmt.Errorf("bitfield frame expected")
	}

	return frm.Data, nil
}

func Connect(p endpoints.Endpoint, peerID, infoHash [20]byte) (*PeerConn, error) {
	conn, err := net.DialTimeout("tcp", p.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	_, err = doHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := readBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &PeerConn{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
		peer:     p,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}

func (p *PeerConn) Read() (*frames.Frame, error) {
	return frames.Unpack(p.Conn)
}

func (p *PeerConn) SendRequest(index, begin, length int) error {
	codec := frames.NewCodec()
	frm := codec.Request(index, begin, length)
	_, err := p.Conn.Write(frm.Pack())
	return err
}

func (p *PeerConn) SendInterested() error {
	frm := &frames.Frame{Type: frames.TypeInterested}
	_, err := p.Conn.Write(frm.Pack())
	return err
}

func (p *PeerConn) SendNotInterested() error {
	frm := &frames.Frame{Type: frames.TypeNotInterested}
	_, err := p.Conn.Write(frm.Pack())
	return err
}

func (p *PeerConn) SendUnchoke() error {
	frm := &frames.Frame{Type: frames.TypeUnchoke}
	_, err := p.Conn.Write(frm.Pack())
	return err
}

func (p *PeerConn) SendHave(index int) error {
	codec := frames.NewCodec()
	frm := codec.Have(index)
	_, err := p.Conn.Write(frm.Pack())
	return err
}
