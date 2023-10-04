package main

import (
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/TheJ0lly/GoChain/wallet"
	"log"
	"os"
)

func main() {

	log.Printf("Uninstalling BlockChain!\n")

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

	log.Printf("Deleting all wallets and their folder...\n")

	files, err := os.ReadDir(dir)

	for _, f := range files {
		fi, err = os.Stat(osspecifics.CreatePath(dir, f.Name()))

		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			return
		}
		//Check what IsRegular does and if it's the correct behaviour
		if f.Name() != "bcs.json" && fi.Mode().IsRegular() {
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

			WalletSavePath := osspecifics.CreatePath(dir, Wallet.GetUsername()+".json")

			err = os.Remove(WalletSavePath)

			if err != nil {
				generalerrors.HandleError(err)
				log.Printf("Error: Failed to remove the wallet save\n")
				continue
			}
		} else if f.Name() != "GoChain_Uninstaller" {
			err = os.Remove(osspecifics.CreatePath(dir, f.Name()))

			if err != nil {
				generalerrors.HandleError(err)
			}
		}
	}

	log.Printf("Uninstall successful\n")
	log.Printf("For now you have to manually delete the uninstaller. Sorry for the inconvenience.\n")
}
