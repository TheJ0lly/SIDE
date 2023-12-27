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
		generalerrors.HandleError(generalerrors.INFO, err)
	} else {
		err = osspecifics.ClearFolder(BC.GetDBLocation())

		if err != nil {
			generalerrors.HandleError(generalerrors.ERROR, err, err)
		}

		log.Printf("Deleting all wallets and their folder...\n")
	}

	files, err := os.ReadDir(dir)

	for _, f := range files {
		if !osspecifics.IsExecutable(f.Name()) && f.Name() != "bcs.json" {
			Wallet, err := wallet.ImportWallet(f.Name())

			if err != nil {
				generalerrors.HandleError(generalerrors.WARNING, err)
				continue
			}

			log.Printf("Clearing wallet folder...\n")
			err = osspecifics.ClearFolder(Wallet.GetDBLocation())

			if err != nil {
				generalerrors.HandleError(generalerrors.WARNING, err)
				continue
			}

			WalletSavePath := osspecifics.CreatePath(dir, Wallet.GetUsername())

			err = os.Remove(WalletSavePath)

			if err != nil {
				generalerrors.HandleError(generalerrors.WARNING, err)
				continue
			}
		} else if f.Name() != "GoChain_Uninstaller.exe" {
			err = os.Remove(osspecifics.CreatePath(dir, f.Name()))

			if err != nil {
				generalerrors.HandleError(generalerrors.WARNING, err)
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
		generalerrors.HandleError(generalerrors.ERROR, err)
	}

	log.Printf("Uninstall successful\n")
}
