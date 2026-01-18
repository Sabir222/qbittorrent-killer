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

	return &PeerConn{
		Conn:     conn,
		Choked:   true,
		peer:     p,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}
