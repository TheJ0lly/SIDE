package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/metadata"
	"github.com/TheJ0lly/GoChain/netutils"
	"github.com/TheJ0lly/GoChain/wallet"
	"github.com/libp2p/go-libp2p/core/network"
	"io"
	"log"
	"os"
	"strconv"
)

var W *wallet.Wallet
var BC *blockchain.BlockChain
var UserName string

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

func updateWallet(ctx context.Context, cancel func()) {
	for {
		_, err := os.ReadFile("importNotification")

		if err != nil {
			continue
		}

		log.Printf("INFO: received notification to update %s\n", UserName)
		log.Printf("INFO: updating wallet\n")
		W, err = wallet.ImportWallet(UserName)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: terminating the service\n")
			cancel()
			return
		}
		SetupStreams()

		log.Printf("INFO: updated wallet successfully\n")

		log.Printf("INFO: updating blockchain\n")
		BC, err = blockchain.ImportChain()

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: terminating the service\n")
			cancel()
			return
		}

		log.Printf("INFO: updated blockchain successfully\n")

		err = os.Remove("importNotification")

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: terminating the service\n")
			cancel()
			return
		}
	}
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

	UserName = *User

	log.Printf("INFO: getting the blockchain\n")
	BC, err = blockchain.ImportChain()

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		return
	}

	back, cancel := context.WithCancel(context.Background())
	fullAddr := W.GetHostAddress()
	log.Printf("INFO: listening on - %s\n", fullAddr)
	SetupStreams()

	go updateWallet(back, cancel)
	select {
	case <-back.Done():
	}

}

func RequestHandler(s network.Stream) {
	var newAdd = netutils.GetHostAddressFromConnection(s.Conn())
	log.Printf("INFO: received new stream on request - %s\n", newAdd)

	_, err := W.AddNode(newAdd)

	if err != nil {
		log.Printf("INFO: %s\n", err)
	} else {
		log.Printf("INFO: added new address\n")
	}

	defer func(s network.Stream) {
		log.Printf("INFO: finishing stream - %s\n", s.ID())
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
		log.Printf("INFO: finshing stream - %s\n", s.ID())
	}(s)

	log.Printf("INFO: received new stream on initialize - %s\n", netutils.GetHostAddressFromConnection(s.Conn()))
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
		return
	}

	BC.Unlock()
	log.Printf("INFO: updating the blockchain instance\n")
	BC, err = blockchain.ImportChain()

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
}

func FloodHandler(s network.Stream) {
	defer func(s network.Stream) {
		BC.Unlock()
		log.Printf("INFO: finshing stream - %s\n", s.ID())
	}(s)

	var newAdd = netutils.GetHostAddressFromConnection(s.Conn())
	log.Printf("INFO: received new stream on flood - %s\n", newAdd)
	err := BC.Lock()

	_, err = W.AddNode(newAdd)

	if err != nil {
		log.Printf("INFO: %s\n", err)
	} else {
		log.Printf("INFO: added new address\n")
	}

	log.Printf("INFO: reading from stream\n")
	stor := make([]byte, 10)

	_, err = s.Read(stor)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	resp := netutils.GetNumberFromResponse(stor)

	if resp == netutils.FailedConversion {
		log.Printf("ERROR: failed to convert bytes to response\n")
		return
	} else if resp == netutils.AssetNotFound {
		log.Printf("ERROR: stream failed to marshal metadata\n")
		return
	} else {
		log.Printf("INFO: sender has successfully marshalled the metadata\n")
		log.Printf("INFO: creating buffer of capacity - %d\n", resp)

		buff := make([]byte, resp)

		_, err = s.Write([]byte("READY"))

		if err != nil {
			log.Printf("ERROR: failed to send ready - %s\n", err)
			return
		}

		log.Printf("INFO: attempting to read %d bytes\n", resp)
		n, err := io.ReadFull(bufio.NewReader(s), buff)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		log.Printf("INFO: read %d bytes from the connection\n", n)
		if n != resp {
			log.Printf("ERROR: read a different amount of bytes than expected\n")
			return
		}

		var mie metadata.MetadataIE

		log.Printf("INFO: trying to unmarshal the bytes\n")
		err = json.Unmarshal(buff, &mie)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		log.Printf("INFO: unmarshalling successful\n")

		md := metadata.CreateNewMetaData(mie.Source, mie.Destination, mie.AssetName)

		if BC.GetLastMetaData().GetMetadataHash() == md.GetMetadataHash() {
			log.Printf("INFO: metadata already present\n")

			log.Printf("INFO: exporting new wallet data\n")
			err = W.ExportWallet()
			if err != nil {

				log.Printf("ERROR: %s\n", err)
			}

			log.Printf("INFO: importing new wallet data\n")
			W, err = wallet.ImportWallet(UserName)
			if err != nil {

				log.Printf("ERROR: %s\n", err)
			}
			SetupStreams()
			return
		}

		err = BC.AddData(mie.Source, mie.Destination, mie.AssetName)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		log.Printf("INFO: successfully added metadata from stream\n")

		netutils.ForwardMetadata(W.GetNodesAddresses(), W.GetHost(), buff, newAdd)
	}

	log.Printf("INFO: exporting new wallet data\n")
	err = W.ExportWallet()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("INFO: importing new wallet data\n")
	W, err = wallet.ImportWallet(UserName)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	SetupStreams()
}

func SetupStreams() {
	W.GetHost().SetStreamHandler("REQUEST", RequestHandler)
	W.GetHost().SetStreamHandler("INITIALIZE", InitializeHandler)
	W.GetHost().SetStreamHandler("FLOOD", FloodHandler)

}
