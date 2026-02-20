package descriptor

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/Sabir222/torrent-at-home/engine"
	"github.com/jackpal/bencode-go"
)

const Port uint16 = 6881

type TorrentFile struct {
	Announce     string
	AnnounceList [][]string
	InfoHash     [20]byte
	PieceHashes  [][20]byte
	PieceLength  int
	Length       int
	Name         string
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce     string      `bencode:"announce"`
	AnnounceList [][]string  `bencode:"announce-list"`
	Info         bencodeInfo `bencode:"info"`
}

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()

	bto := bencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return TorrentFile{}, err
	}
	return bto.toTorrentFile()
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	raw := []byte(i.Pieces)
	if len(raw)%hashLen != 0 {
		return nil, fmt.Errorf("invalid piece hash data")
	}
	count := len(raw) / hashLen
	result := make([][20]byte, count)

	for i := 0; i < count; i++ {
		copy(result[i][:], raw[i*hashLen:(i+1)*hashLen])
	}
	return result, nil
}

func (bto *bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	t := TorrentFile{
		Announce:     bto.Announce,
		AnnounceList: bto.AnnounceList,
		InfoHash:     infoHash,
		PieceHashes:  pieceHashes,
		PieceLength:  bto.Info.PieceLength,
		Length:       bto.Info.Length,
		Name:         bto.Info.Name,
	}
	return t, nil
}

// DownloadToFile downloads a torrent and writes it to a file
func (t *TorrentFile) DownloadToFile(path string) error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil {
		return err
	}

	peers, err := t.announce(peerID, Port)
	if err != nil {
		return err
	}

	session := engine.Session{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    t.InfoHash,
		PieceHashes: t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.Length,
		Name:        t.Name,
	}
	buf, err := session.Download()
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
