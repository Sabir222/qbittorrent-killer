package engine

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/Sabir222/torrent-at-home/network/connector"
	"github.com/Sabir222/torrent-at-home/protocol/frames"
	"github.com/Sabir222/torrent-at-home/network/endpoints"
)

const (
	DefaultChunkSize = 16384
	MaxPending       = 5
)

type Session struct {
	Peers       []endpoints.Endpoint
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type job struct {
	index  int
	hash   [20]byte
	length int
}

type result struct {
	index int
	buf   []byte
}

type transferState struct {
	index     int
	conn      *connector.PeerConn
	buf       []byte
	received  int
	requested int
	pending   int
}

func (s *transferState) processMsg() error {
	msg, err := s.conn.Read()
	if err != nil {
		return err
	}

	if msg == nil {
		return nil
	}

	switch msg.Type {
	case frames.TypeUnchoke:
		s.conn.Choked = false
	case frames.TypeChoke:
		s.conn.Choked = true
	case frames.TypeHave:
		idx, err := frames.ReadHave(msg)
		if err != nil {
			return err
		}
		s.conn.Bitfield.Mark(idx)
	case frames.TypePiece:
		n, err := frames.ReadPieceData(s.buf, s.index, msg)
		if err != nil {
			return err
		}
		s.received += n
		s.pending--
	}
	return nil
}

func fetchPiece(conn *connector.PeerConn, j *job) ([]byte, error) {
	state := transferState{
		index: j.index,
		conn:  conn,
		buf:   make([]byte, j.length),
	}

	conn.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer conn.Conn.SetDeadline(time.Time{})

	// Pipelined request loop - keep multiple outstanding requests
	for state.received < j.length {
		if !state.conn.Choked {
			for state.pending < MaxPending && state.requested < j.length {
				chunkSize := DefaultChunkSize
				if j.length-state.requested < chunkSize {
					chunkSize = j.length - state.requested
				}

				err := conn.SendRequest(j.index, state.requested, chunkSize)
				if err != nil {
					return nil, err
				}
				state.pending++
				state.requested += chunkSize
			}
		}

		err := state.processMsg()
		if err != nil {
			return nil, err
		}
	}

	return state.buf, nil
}

func verifyPiece(j *job, buf []byte) error {
	// SHA-1 hash verification for data integrity
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], j.hash[:]) {
		return fmt.Errorf("piece %d failed validation", j.index)
	}
	return nil
}

// spawnWorker creates a goroutine that connects to a peer and processes piece jobs
func (s *Session) spawnWorker(peer endpoints.Endpoint, jobs chan *job, results chan *result) {
	conn, err := connector.Connect(peer, s.PeerID, s.InfoHash)
	if err != nil {
		log.Printf("[peer] ✗ %s unreachable\n", peer.Addr)
		return
	}
	defer conn.Conn.Close()
	log.Printf("[peer] ✓ connected to %s\n", peer.Addr)

	conn.SendUnchoke()
	conn.SendInterested()

	for j := range jobs {
		if !conn.Bitfield.Check(j.index) {
			jobs <- j
			continue
		}

		buf, err := fetchPiece(conn, j)
		if err != nil {
			log.Printf("[peer] fetch error: %v\n", err)
			jobs <- j
			return
		}

		err = verifyPiece(j, buf)
		if err != nil {
			log.Printf("piece %d corrupted\n", j.index)
			jobs <- j
			continue
		}

		conn.SendHave(j.index)
		results <- &result{j.index, buf}
	}
}

func (s *Session) pieceRange(index int) (begin int, end int) {
	begin = index * s.PieceLength
	end = begin + s.PieceLength
	if end > s.Length {
		end = s.Length
	}
	return begin, end
}

func (s *Session) pieceSize(index int) int {
	begin, end := s.pieceRange(index)
	return end - begin
}

func (s *Session) Download() ([]byte, error) {
	log.Printf("[session] starting download: %s\n", s.Name)
	log.Printf("[session] %d piece(s), %d peer(s) available\n", len(s.PieceHashes), len(s.Peers))

	jobs := make(chan *job, len(s.PieceHashes))
	results := make(chan *result)

	for index, hash := range s.PieceHashes {
		length := s.pieceSize(index)
		jobs <- &job{index, hash, length}
	}

	for _, peer := range s.Peers {
		go s.spawnWorker(peer, jobs, results)
	}

	buf := make([]byte, s.Length)
	completed := 0
	for completed < len(s.PieceHashes) {
		res := <-results
		begin, end := s.pieceRange(res.index)
		copy(buf[begin:end], res.buf)
		completed++

		pct := float64(completed) / float64(len(s.PieceHashes)) * 100
		workers := runtime.NumGoroutine() - 1
		log.Printf("[progress] [%5.1f%%] ✓ piece %d (%d worker(s))\n", pct, res.index, workers)
	}
	close(jobs)

	log.Printf("[session] ✓ download complete: %d piece(s)\n", completed)
	return buf, nil
}
