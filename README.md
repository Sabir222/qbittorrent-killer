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

## License

MIT License
