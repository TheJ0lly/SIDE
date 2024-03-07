package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/wallet"
	"github.com/libp2p/go-libp2p/core/network"
	"log"
)

var W *wallet.Wallet
var BC *blockchain.BlockChain

func displayHelp() {
	fmt.Printf("Usage: <exec> -u <string>\n\n")
	fmt.Printf("  -u  \n      Choose the user for which the service will run.\n")
}

func main() {
	User := flag.String("u", "", "")

	flag.Usage = displayHelp

	flag.Parse()

	if *User == "" {
		flag.Usage()
		return
	}

	var err error

	log.Printf("INFO: getting wallet %s\n", *User)
	W, err = wallet.ImportWallet(*User)

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		return
	}

	log.Printf("INFO: getting the blockchain\n")
	BC, err = blockchain.ImportChain()

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		return
	}
	//
	//cancel := make(chan os.Signal, 1)
	//
	//signal.Notify(cancel)

	back := context.Background()
	go StartListener(back)

	//select {
	//case s := <-cancel:
	//	fmt.Printf("\nService stopped: %v\n", s.String())
	//}

	select {
	case <-back.Done():
	}

}

func ListenHandler(s network.Stream) {
	log.Printf("INFO: received new stream - %s\n", s.ID())

	var stor = make([]byte, 200)

	log.Printf("INFO: reading from stream\n")
	_, err := s.Read(stor)

	//_, err := io.ReadFull(bufio.NewReader(s), stor)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("INFO: received - %s\n", stor)

	_, err = s.Write([]byte("DONE"))

	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}

	log.Printf("INFO: closing stream - %s\n", s.ID())
	err = s.Close()

	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}

	return
}

func StartListener(ctx context.Context) {
	fullAddr := W.GetHostAddress()
	log.Printf("INFO: listening on - %s\n", fullAddr)

	W.GetHost().SetStreamHandler("LISTEN", ListenHandler)
}
