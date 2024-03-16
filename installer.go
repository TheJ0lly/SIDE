package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/netutils"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
	"log"
	"os"
	"strings"
)

func InstallerHelp() {
	fmt.Printf("Usage: <exec> -d <string>\n")
	fmt.Printf("  -d\n      Input the directory in which you want to hold the blockchain data.\n")
	fmt.Printf("  -n\n      Clears the folder in which you are trying to set the database.\n")
	fmt.Printf("  -a\n      The address of the node you want to initialize from.\n")
}

func main() {
	Database := flag.String("d", "", "")
	ClearFolder := flag.Bool("n", false, "")
	Address := flag.String("a", "", "")
	flag.Usage = InstallerHelp

	flag.Parse()

	if *Database == "" {
		log.Printf("ERROR: No directory to install into. Set the value with the \"d\" flag.")
		return
	}
	var BC *blockchain.BlockChain
	var err error
	var privateKey crypto.PrivKey
	var ha core.Host
	var ma multiaddr.Multiaddr

	*Database = osspecifics.GetFullPathFromArg(*Database)

	if *ClearFolder {
		log.Printf("INFO: clearing selected folder...\n")
		err = osspecifics.ClearFolder(*Database)

		if err != nil {
			generalerrors.HandleError(generalerrors.ERROR, err)
			return
		}
	}

	files, err := os.ReadDir(*Database)

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}

	if len(files) > 0 {
		log.Printf("ERROR: folder already contains some files!\n To clear the folder use the flag: -n\n")
		return
	}

	log.Printf("INFO: initializing blockchain...\n")

	if *Address != "" {
		ma, err = multiaddr.NewMultiaddr(*Address)

		if err != nil {
			log.Printf("ERROR: cannot convert given address - %s\n", err)
			return
		}

		privateKey, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

		if err != nil {
			log.Printf("ERROR: failed to create temporary key - %s\n", err)
			return
		}

		ha, err = netutils.CreateNewNode(netutils.CreateNodeOptions(privateKey, "", "0"))

		if err != nil {
			log.Printf("ERROR: failed to create temporary host - %s\n", err)
			return
		}

		BC, err = netutils.CreateNewBlockchainFromConn(ha, *Database, ma)
	} else {

		log.Printf("INFO: creating blockchain from scratch\n")
		fmt.Printf("Are you sure you don't want to try to connect to other users? y/n: ")

		var choice string

		_, err := fmt.Scanln(&choice)

		if err != nil || strings.ToLower(choice) == "n" {
			log.Printf("INFO: aborting installation\n")
			return
		}

		BC, err = blockchain.CreateNewBlockchain(*Database)
	}

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		return
	}
	log.Printf("INFO: blockchain intialized!\n\n")

	err = BC.ExportChain()

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		return
	}

	log.Printf("INFO: SIDE executable ready to use!\n")
}
