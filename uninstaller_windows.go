package main

import (
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/TheJ0lly/GoChain/wallet"
	"log"
	"os"
	"strings"
)

func main() {

	fmt.Printf("Uninstalling BlockChain!\n")

	dir, err := os.Getwd()

	if err != nil {
		generalerrors.HandleError(err)
		return
	}

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}

	BC, err := blockchain.ImportChain()

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}

	err = osspecifics.ClearFolder(BC.GetDBLocation())

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("Deleting all wallets and their folder...\n")

	files, err := os.ReadDir(dir)

	for _, f := range files {
		if f.Name() != "bcs.json" && !strings.Contains(f.Name(), ".exe") {
			Wallet, err := wallet.ImportWallet(f.Name())

			if err != nil {
				generalerrors.HandleError(err)
				continue
			}

			err = osspecifics.ClearFolder(Wallet.GetDBLocation())

			if err != nil {
				generalerrors.HandleError(err)
				continue
			}

			WalletSavePath := osspecifics.CreatePath(dir, Wallet.GetUsername())

			err = os.Remove(WalletSavePath)

			if err != nil {
				generalerrors.HandleError(err)
				fmt.Printf("Error: Failed to remove the wallet save\n")
				continue
			}
		} else if f.Name() != "GoChain_Uninstaller.exe" {
			err = os.Remove(osspecifics.CreatePath(dir, f.Name()))

			if err != nil {
				generalerrors.HandleError(err)
			}
		}
	}

	fmt.Printf("Uninstall successful\n")
	fmt.Printf("For now you have to manually delete the uninstaller. Sorry for the inconvenience.\n")
}
