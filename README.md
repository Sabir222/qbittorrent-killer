# Torrent Client

A lightweight BitTorrent client written in Go, designed for efficient file downloads through the BitTorrent protocol.

## Features

- **Parallel Downloads** - Connects to multiple peers simultaneously for maximum throughput
- **Piece Validation** - SHA-1 hash verification ensures data integrity
- **Peer Management** - Automatic peer discovery and connection handling
- **Progress Tracking** - Real-time download progress reporting
- **UDP Tracker Support** - Fast, efficient tracker communication with HTTP fallback

## Project Structure

```
torrent-client/
├── cmd/app/              # CLI entry point
├── engine/               # Download engine and worker management
├── protocol/             # BitTorrent protocol implementation
│   ├── greeting/         # Peer handshake protocol
│   └── frames/           # Message encoding and decoding
├── network/              # Network layer
│   ├── connector/        # Peer connection handling
│   └── endpoints/        # Peer address parsing
├── data/                 # Data structures and metadata
│   ├── mask/             # Bitfield operations for piece tracking
│   └── descriptor/       # Torrent file parsing and tracker communication
└── tools/                # Utilities and scripts
```

## Installation

### Build from Source

```bash
git clone https://github.com/Sabir222/torrent-at-home.git
cd torrent-at-home
go build -o torrent-at-home ./cmd/app
```

### Requirements

- Go 1.21 or later

## Usage

```bash
./torrent-at-home <torrent-file> <output-path>
```

### Example

```bash
# Download a Linux ISO
./torrent-at-home kali-linux-2025.4-installer-amd64.iso.torrent ./kali.iso
```

## How It Works

### 1. Torrent File Parsing

The client reads the `.torrent` file to extract:

- Tracker URLs for peer discovery
- Piece hashes for validation
- File metadata (name, size, piece length)

### 2. Peer Discovery

Connects to the tracker to obtain a list of peers currently seeding or downloading the same torrent.

### 3. Handshake Protocol

Establishes connections with peers using the BitTorrent handshake protocol, exchanging info hashes and peer IDs.

### 4. Download Process

- Requests pieces from multiple peers in parallel
- Validates each piece using SHA-1 hash verification
- Tracks completed pieces using bitfields
- Re-requests corrupted pieces automatically

### 5. Completion

Once all pieces are downloaded and validated, they're assembled into the final file.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI (cmd/app)                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Descriptor (data/descriptor)               │
│              Torrent parsing & Tracker communication         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Engine (engine/transfer)                  │
│           Download coordination & Worker management          │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
    ┌─────────────────┐ ┌─────────────┐ ┌─────────────┐
    │   Connector     │ │   Protocol  │ │    Mask     │
    │  (connections)  │ │  (messages) │ │  (bitfield) │
    └─────────────────┘ └─────────────┘ └─────────────┘
```

**Data Flow:**
1. CLI loads torrent file → Descriptor parses metadata
2. Descriptor queries trackers → Returns peer list
3. Engine spawns workers → Each connects to a peer
4. Workers request pieces → Validate and return results
5. Engine assembles pieces → Writes final file

## Testing

### Unit Tests

Run the test suite:

```bash
go test ./...
```

### Live Testing

Use the included Kali Linux torrent:

```bash
./torrent-at-home kali-linux-2025.4-installer-amd64.iso.torrent ./kali.iso
```

## Performance

The client optimizes download speed through:

- **Concurrent peer connections** - Multiple workers download pieces in parallel
- **Pipelined requests** - Keeps multiple outstanding requests per peer
- **Efficient memory usage** - Pre-allocated buffers reduce GC pressure

## Limitations

This is a minimal implementation focusing on core BitTorrent functionality:

- **No DHT support** - Requires tracker for peer discovery
- **No PEX (Peer Exchange)** - Cannot learn about peers from connected peers
- **No magnet links** - Requires `.torrent` file
- **No protocol encryption** - May be blocked by some ISPs
- **Single-file torrents only** - Multi-file torrents not supported

**Recommendation:** Best results with popular torrents (Linux ISOs) with 50+ peers.

## Future Improvements

- [ ] DHT support for trackerless downloads
- [ ] PEX (Peer Exchange) implementation
- [ ] Magnet link support
- [ ] Protocol encryption (PE/MSE)
- [ ] Multi-file torrent support

## Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please make sure your code follows Go best practices and includes tests where appropriate.

## Disclaimer

This project is for educational purposes only. The author(s) are not responsible for any misuse of this software. Please ensure you comply with your local laws and regulations when downloading torrents. Only download content you have the legal right to access.

## License

MIT License - See [LICENSE](LICENSE) for details.
