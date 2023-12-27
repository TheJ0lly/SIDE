package main

import (
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/TheJ0lly/GoChain/wallet"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	log.Printf("Uninstalling BlockChain!\n")

	exePath, err := os.Executable()

	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}

	dir := filepath.Dir(exePath)

	BC, err := blockchain.ImportChain()

	if err != nil {
		generalerrors.HandleError(err)
	} else {
		err = osspecifics.ClearFolder(BC.GetDBLocation())

		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			return
		}

		log.Printf("Deleting all wallets and their folder...\n")
	}

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
				log.Printf("Error: Failed to remove the wallet save\n")
				continue
			}
		} else if f.Name() != "GoChain_Uninstaller.exe" {
			err = os.Remove(osspecifics.CreatePath(dir, f.Name()))

			if err != nil {
				generalerrors.HandleError(err)
			}
		}
	}

	uninstallerPath := osspecifics.CreatePath(dir, "GoChain_Uninstaller.exe")
	removeUninstaller := "/c start timeout /t 1 /NOBREAK > NUL && del " + uninstallerPath

	s := exec.Command("cmd.exe", removeUninstaller)
	s.Stdout = os.Stdout
	s.Stdin = os.Stdin
	s.Stderr = os.Stderr

	err = s.Start()

	if err != nil {
		generalerrors.HandleError(err)
	}

	log.Printf("Uninstall successful\n")
}
