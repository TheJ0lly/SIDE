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
	fmt.Printf("  -n\n      Clears the folder in which you are trying to set the database.\n")
}

func main() {
	Database := flag.String("d", "", "")
	ClearFolder := flag.Bool("n", false, "")
	flag.Usage = InstallerHelp

	flag.Parse()

	if *Database == "" {
		log.Printf("Error: No directory to install into. Set the value with \"d\".")
		return
	}

	*Database = osspecifics.GetFullPathFromArg(*Database)

	if *ClearFolder {
		log.Printf("Clearing selected folder...\n")
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
		log.Printf("Error: Folder already contains some files!\n To clear the folder use the flag: -n\n")
		return
	}

	log.Printf("Initializing BlockChain...\n")
	BC, err := blockchain.CreateNewBlockchain(*Database)

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		return
	}
	log.Printf("BlockChain intialized!\n\n")

	err = BC.ExportChain()

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		return
	}

	log.Printf("SIDE executable ready to use!\n")

}
