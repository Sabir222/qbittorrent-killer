package descriptor

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"net/url"
	"time"

	"github.com/Sabir222/torrent-at-home/network/endpoints"
)

const (
	udpProtocolID   = 0x41727101980
	udpActionConnect = 0
	udpActionAnnounce = 1
	udpActionError   = 3
	udpConnectTimeout = 15 * time.Second
)

var (
	ErrUDPTransactionMismatch = errors.New("UDP transaction ID mismatch")
	ErrUDPUnexpectedAction    = errors.New("unexpected UDP action in response")
	ErrUDPConnectionExpired   = errors.New("UDP connection ID expired")
)

func (t *TorrentFile) announceUDP(peerID [20]byte, port uint16) ([]endpoints.Endpoint, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", t.trackerHost())
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	log.Printf("[udp] obtaining connection ID from %s\n", t.trackerHost())
	connID, err := t.getUDPConnectionID(conn)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	log.Printf("[udp] ✓ connection established\n")

	return t.sendUDPAnnounce(conn, connID, peerID, port)
}

func (t *TorrentFile) getUDPConnectionID(conn *net.UDPConn) (int64, error) {
	txID := randomTransactionID()

	req := make([]byte, 16)
	binary.BigEndian.PutUint64(req[0:8], udpProtocolID)
	binary.BigEndian.PutUint32(req[8:12], udpActionConnect)
	binary.BigEndian.PutUint32(req[12:16], txID)

	resp, err := sendUDPRequest(conn, req, udpConnectTimeout)
	if err != nil {
		return 0, err
	}

	if len(resp) < 16 {
		return 0, errors.New("UDP connect response too short")
	}

	action := binary.BigEndian.Uint32(resp[0:4])
	respTxID := binary.BigEndian.Uint32(resp[4:8])

	if action != udpActionConnect {
		return 0, ErrUDPUnexpectedAction
	}
	if respTxID != txID {
		return 0, ErrUDPTransactionMismatch
	}

	connID := int64(binary.BigEndian.Uint64(resp[8:16]))
	return connID, nil
}

func sendUDPRequest(conn *net.UDPConn, data []byte, timeout time.Duration) ([]byte, error) {
	maxRetries := 8

	for attempt := 0; attempt <= maxRetries; attempt++ {
		conn.SetDeadline(time.Now().Add(timeout))

		_, err := conn.Write(data)
		if err != nil {
			return nil, err
		}

		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err == nil {
			return buf[:n], nil
		}

		// Exponential backoff: 15 × 2^n seconds
		backoff := 15 * time.Second * time.Duration(1<<attempt)
		timeout = backoff
	}

	return nil, errors.New("UDP tracker timeout after all retries")
}

func randomTransactionID() uint32 {
	var b [4]byte
	rand.Read(b[:])
	return binary.BigEndian.Uint32(b[:])
}

func (t *TorrentFile) trackerHost() string {
	u, err := url.Parse(t.Announce)
	if err != nil {
		return ""
	}
	return u.Host
}

func (t *TorrentFile) isUDPTracker() bool {
	u, err := url.Parse(t.Announce)
	if err != nil {
		return false
	}
	return u.Scheme == "udp"
}
