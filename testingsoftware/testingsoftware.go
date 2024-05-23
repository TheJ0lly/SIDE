package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

func getHostAddress(h core.Host) string {
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", h.ID()))

	return h.Addrs()[1].Encapsulate(hostAddr).String()
}

var isRunning = true

func wait(ctx context.Context, cancel func()) {
	for isRunning {
	}

	cancel()
}

func main() {
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
	)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	address := flag.String("a", "", "")
	file := flag.String("s", "", "")
	receive := flag.Bool("r", false, "")

	flag.Parse()

	if *receive {
		log.Printf("The address to connect to: %s\n", getHostAddress(h))

		ctx, cancel := context.WithCancel(context.Background())

		h.SetStreamHandler("TRANSFER", func(s network.Stream) {
			log.Printf("Making a buffer of 32600\n")
			start := time.Now()
			b := make([]byte, 32600)

			log.Printf("Starting to read the file\n")
			n, err := io.ReadFull(bufio.NewReader(s), b)
			end := time.Since(start)
			if err != nil {
				log.Printf("ERROR: %s\n", err)
				isRunning = false
				return
			}
			log.Printf("Finished reading %d - %v\n", n, end)
			s.Write([]byte("Done"))
			isRunning = false
		})

		go wait(ctx, cancel)

		select {
		case <-ctx.Done():
		}

	} else {
		start := time.Now()
		b, err := os.ReadFile(*file)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		ma, err := multiaddr.NewMultiaddr(*address)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		pinfo, err := peer.AddrInfoFromP2pAddr(ma)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		h.Peerstore().AddAddrs(pinfo.ID, pinfo.Addrs, peerstore.TempAddrTTL)

		s, err := h.NewStream(context.Background(), pinfo.ID, "TRANSFER")

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		log.Printf("Established connection with %s\n", ma.String())

		log.Printf("Sending file bytes - %d\n", len(b))

		_, err = s.Write(b)
		end := time.Since(start)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		log.Printf("Finished sending bytes in %v ms\n", end.Milliseconds())

		b = make([]byte, 4)
		s.Read(b)

		s.Close()
	}
}
