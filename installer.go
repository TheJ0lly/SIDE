package main

import (
	"flag"
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"log"
	"os"
)

func InstallerHelp() {
	fmt.Printf("Usage: <exec> -db <string>\n")
	fmt.Printf("  -n\n      Clears the folder in which you are trying to set the database.\n")
}

func main() {
	Database := flag.String("db", "", "")
	ClearFolder := flag.Bool("n", false, "")
	flag.Usage = InstallerHelp

	flag.Parse()

	if *Database == "" {
		log.Printf("Error: Blockchain Database is empty. Set the value with \"db\".")
		return
	}

	if *ClearFolder {
		err := osspecifics.ClearFolder(*Database)

		if err != nil {
			generalerrors.HandleError(err)
			return
		}
	}

	files, err := os.ReadDir(*Database)

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}

	if len(files) > 0 {
		log.Printf("Error: Folder already contains some files!\n To delete the file use the flag: -n\n")
		return
	}

	log.Printf("Initializing BlockChain...\n")
	BC, err := blockchain.CreateNewBlockchain(*Database)

	if err != nil {
		generalerrors.HandleError(err)
		return
	}
	log.Printf("BlockChain intialized!\n\n")

	log.Printf("Exporting BlockChain...\n")
	err = BC.ExportChain()

	if err != nil {
		generalerrors.HandleError(err)
		return
	}

	log.Printf("BlockChain exported!\n")
	log.Printf("GoChain executable ready to use!\n")

}
