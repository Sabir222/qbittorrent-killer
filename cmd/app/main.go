package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Sabir222/torrent-at-home/data/descriptor"
)

const (
	version = "1.0.0"
	name    = "torrent-at-home"
)

func banner() {
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("  %s v%s\n", name, version)
	fmt.Println(strings.Repeat("─", 50))
}

func entry() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <torrent-file> <output-path>\n", os.Args[0])
		os.Exit(1)
	}

	src := os.Args[1]
	dst := os.Args[2]

	banner()
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	meta, err := descriptor.Open(src)
	if err != nil {
		log.Fatalf("failed to load torrent: %v", err)
	}

	fmt.Printf("\n[info] torrent: %s\n", meta.Name)
	fmt.Printf("[info] size: %.2f MB\n", float64(meta.Length)/1024/1024)
	fmt.Printf("[info] pieces: %d\n\n", len(meta.PieceHashes))
	fmt.Println("[info] starting download...\n")

	if err := meta.DownloadToFile(dst); err != nil {
		log.Fatalf("transfer failed: %v", err)
	}

	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("[success] download complete!")
	fmt.Printf("[output] %s\n", dst)
	fmt.Println(strings.Repeat("─", 50))
}

func main() {
	entry()
}
