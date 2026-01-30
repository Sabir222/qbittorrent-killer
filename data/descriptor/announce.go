package descriptor

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Sabir222/torrent-at-home/network/endpoints"
	"github.com/jackpal/bencode-go"
)

const (
	trackerTimeout  = 15 * time.Second
	defaultPort     = 6881
)

var ErrTrackerResponse = errors.New("invalid tracker response")

type trackerReply struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (t *TorrentFile) announce(peerID [20]byte, port uint16) ([]endpoints.Endpoint, error) {
	var allPeers []endpoints.Endpoint

	trackers := t.getTrackerList()
	log.Printf("[tracker] querying %d tracker(s)\n", len(trackers))

	for i, trackerURL := range trackers {
		t.Announce = trackerURL
		protocol := "HTTP"
		if t.isUDPTracker() {
			protocol = "UDP"
		}

		peers, err := t.announceSingleTracker(peerID, port)
		if err != nil {
			log.Printf("[tracker] [%d/%d] %s %s → failed: %v\n", i+1, len(trackers), protocol, truncateURL(trackerURL), err)
			continue
		}
		if len(peers) > 0 {
			log.Printf("[tracker] [%d/%d] %s %s → %d peer(s)\n", i+1, len(trackers), protocol, truncateURL(trackerURL), len(peers))
			allPeers = append(allPeers, peers...)
		} else {
			log.Printf("[tracker] [%d/%d] %s %s → no peers\n", i+1, len(trackers), protocol, truncateURL(trackerURL))
		}
	}

	if len(allPeers) == 0 {
		return nil, errors.New("no peers received from any tracker")
	}

	// Remove duplicate peers
	seen := make(map[string]bool)
	uniquePeers := make([]endpoints.Endpoint, 0, len(allPeers))
	for _, p := range allPeers {
		key := p.String()
		if !seen[key] {
			seen[key] = true
			uniquePeers = append(uniquePeers, p)
		}
	}

	log.Printf("[tracker] total: %d unique peer(s) after deduplication\n", len(uniquePeers))
	return uniquePeers, nil
}

func truncateURL(u string) string {
	if len(u) <= 40 {
		return u
	}
	return u[:37] + "..."
}

func (t *TorrentFile) getTrackerList() []string {
	var trackers []string

	// First add all tiers from announce-list
	for _, tier := range t.AnnounceList {
		for _, tracker := range tier {
			trackers = append(trackers, tracker)
		}
	}

	// Add primary announce if not already in list
	if t.Announce != "" {
		found := false
		for _, tracker := range trackers {
			if tracker == t.Announce {
				found = true
				break
			}
		}
		if !found {
			trackers = append(trackers, t.Announce)
		}
	}

	return trackers
}

func (t *TorrentFile) announceSingleTracker(peerID [20]byte, port uint16) ([]endpoints.Endpoint, error) {
	if t.isUDPTracker() {
		peers, err := t.announceUDP(peerID, port)
		if err == nil {
			return peers, nil
		}
	}

	target, err := t.assembleURL(peerID, port)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: trackerTimeout}
	reply, err := client.Get(target)
	if err != nil {
		return nil, err
	}
	defer reply.Body.Close()

	var decoded trackerReply
	if err := bencode.Unmarshal(reply.Body, &decoded); err != nil {
		return nil, err
	}

	return endpoints.Parse([]byte(decoded.Peers))
}

func (t *TorrentFile) assembleURL(id [20]byte, p uint16) (string, error) {
	base, parseErr := url.Parse(t.Announce)
	if parseErr != nil {
		return "", parseErr
	}

	query := url.Values{}
	query.Set("info_hash", string(t.InfoHash[:]))
	query.Set("peer_id", string(id[:]))
	query.Set("port", strconv.Itoa(int(p)))
	query.Set("uploaded", "0")
	query.Set("downloaded", "0")
	query.Set("compact", "1")
	query.Set("left", strconv.Itoa(t.Length))
	query.Set("event", "started")

	base.RawQuery = query.Encode()
	return base.String(), nil
}
