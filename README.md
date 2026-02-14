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

- No DHT support
- No PEX (Peer Exchange)
- No magnet links
- No protocol encryption

Best results with popular torrents (Linux ISOs) with 50+ peers.

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Disclaimer

This project is for educational purposes. Please ensure you comply with your local laws and regulations when downloading torrents.
