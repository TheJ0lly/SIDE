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
	fmt.Printf("Usage: <exec> -d <string>\n")
	fmt.Printf("  -d\n      Input the directory in which you want to hold the blockchain data.\n")
	fmt.Printf("  -n\n      Clears the folder in which you are trying to set the database.\n")
}

func main() {
	Database := flag.String("d", "", "")
	ClearFolder := flag.Bool("n", false, "")
	flag.Usage = InstallerHelp

	flag.Parse()

	if *Database == "" {
		log.Printf("ERROR: No directory to install into. Set the value with the \"d\" flag.")
		return
	}

	*Database = osspecifics.GetFullPathFromArg(*Database)

	if *ClearFolder {
		log.Printf("INFO: clearing selected folder...\n")
		err := osspecifics.ClearFolder(*Database)

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
	BC, err := blockchain.CreateNewBlockchain(*Database)

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
