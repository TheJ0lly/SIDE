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

func exportStates(Wallet *wallet.Wallet, BC *blockchain.BlockChain) error {

	fmt.Print("\n")
	err := BC.ExportChain()

	if err != nil {
		return err
	}

	err = Wallet.ExportWallet()

	if err != nil {
		return err
	}

	return nil
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

func RequestHandler(s network.Stream) {
	var newAdd = netutils.GetHostAddressFromConnection(s.Conn())
	log.Printf("INFO: received new stream - %s\n", newAdd)

	_, err := W.AddNode(newAdd)

	if err != nil {
		log.Printf("INFO: %s\n", err)
	} else {
		log.Printf("INFO: added new address\n")
	}

	defer func(s network.Stream) {
		log.Printf("INFO: closing stream - %s\n", s.ID())
		err := s.Close()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
		}
	}(s)

	var stor = make([]byte, 200)

	log.Printf("INFO: reading from stream\n")
	_, err = s.Read(stor)

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

	stor = make([]byte, 10)
	log.Printf("INFO: waiting for ready signal\n")
	_, err = s.Read(stor)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	resp := netutils.ConvertBytesToString(stor)

	log.Printf("INFO: received signal - %s\n", resp)
	if resp == "READY" {
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

func InitializeHandler(s network.Stream) {
	defer func(s network.Stream) {
		BC.Unlock()

		log.Printf("INFO: closing stream - %s\n", s.ID())
		err := s.Close()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
		}
	}(s)

	log.Printf("INFO: received new stream - %s\n", netutils.GetHostAddressFromConnection(s.Conn()))
	err := BC.Lock()

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("INFO: getting number of blocks\n")
	blocks := BC.GetBlocks()

	log.Printf("INFO: sending number of blocks\n")
	_, err = s.Write([]byte(strconv.Itoa(len(blocks))))

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	for _, b := range blocks {
		byt := blockchain.ExportBlockForConn(b)

		log.Printf("INFO: sending block size\n")
		_, err := s.Write([]byte(strconv.Itoa(len(byt))))

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		log.Printf("INFO: sending block bytes\n")
		_, err = s.Write(byt)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}
	}

	err = BC.ExportChain()

	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}

	BC.Unlock()
	log.Printf("INFO: updating the blockchain instance\n")
	BC, err = blockchain.ImportChain()

	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}
}

func StartListener(ctx context.Context) {
	fullAddr := W.GetHostAddress()
	log.Printf("INFO: listening on - %s\n", fullAddr)

	W.GetHost().SetStreamHandler("REQUEST", RequestHandler)
	W.GetHost().SetStreamHandler("INITIALIZE", InitializeHandler)
}
