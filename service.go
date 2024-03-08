package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/netutils"
	"github.com/TheJ0lly/GoChain/wallet"
	"github.com/libp2p/go-libp2p/core/network"
	"log"
	"strconv"
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

	log.Printf("INFO: getting wallet %s address\n", *User)
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

	back := context.Background()
	go StartListener(back)

	select {
	case <-back.Done():
	}

}

func ListenHandler(s network.Stream) {
	log.Printf("INFO: received new stream - %s\n", wallet.GetHostAddressFromConnection(s.Conn()))

	defer func(s network.Stream) {
		log.Printf("INFO: closing stream - %s\n", s.ID())
		err := s.Close()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
		}
	}(s)

	var stor = make([]byte, 200)

	log.Printf("INFO: reading from stream\n")
	_, err := s.Read(stor)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("INFO: received - %s\n", stor)
	log.Printf("INFO: searching if %s has %s", W.GetUsername(), stor)

	assetName := netutils.ConvertBytesToString(stor)

	as, err := W.GetAsset(assetName)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		_, err := s.Write([]byte("-1"))

		if err != nil {
			log.Printf("ERROR: %s\n", err)
		}
		return
	}

	log.Printf("INFO: asset found - sending size\n")
	_, err = s.Write([]byte(strconv.Itoa(as.GetAssetSize())))

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("INFO: waiting for ready signal\n")
	_, err = s.Read(stor)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	if netutils.ConvertBytesToString(stor) == "READY" {
		log.Printf("INFO: ready signal recieved - sending bytes\n")
		_, err = s.Write(as.GetAssetBytes())

		if err != nil {
			log.Printf("ERROR: failed to send bytes - %s\n", err)
			return
		}

		log.Printf("INFO: bytes sent successfully\n")
	}

	return
}

func StartListener(ctx context.Context) {
	fullAddr := W.GetHostAddress()
	log.Printf("INFO: listening on - %s\n", fullAddr)

	W.GetHost().SetStreamHandler("REQUEST", ListenHandler)
}
