package wallet

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"os"

	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
)

const (
	keyBitSize = 256
)

type Wallet struct {
	username    string
	password    string
	publicKey   rsa.PublicKey
	privateKey  rsa.PrivateKey
	databaseDir string
	assets      []*asset.Asset
}

// CreateNewWallet - This function will create a wallet.
func CreateNewWallet(username string, password string, dbLoc string) (*Wallet, error) {
	files, err := os.ReadDir(dbLoc)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil, &generalerrors.ReadDirFailed{Dir: dbLoc}
	}

	if len(files) > 0 {
		fmt.Printf("Warning: Folder %s contains files! Do you want to delete them?[y\\n]\n", dbLoc)
		var choice string
		_, err := fmt.Scanln(&choice)

		if err != nil {
			return nil, err
		}

		if choice != "y" {
			return nil, &generalerrors.WalletDBHasItems{Dir: dbLoc}
		}

		err = clearFolder(dbLoc, files)

		if err != nil {
			return nil, err
		}
	}

	privKey, err := rsa.GenerateKey(rand.Reader, keyBitSize)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil, err
	}

	passBytes := sha256.Sum256([]byte(password))

	passBytesStr := fmt.Sprintf("%X", passBytes)

	w := &Wallet{username: username, password: passBytesStr, privateKey: *privKey, publicKey: privKey.PublicKey, assets: nil, databaseDir: dbLoc}

	return w, nil
}

// AddAsset - This function will add an asset to the wallet.
func (w *Wallet) AddAsset(assetName string, fileLocation string) bool {

	if w.checkAssetExists(assetName) {
		fmt.Printf("Error: There is already an asset with this name - \"%s\"\n", assetName)
		return false
	}

	temp, err := os.Stat(fileLocation)

	if os.IsNotExist(err) {
		fmt.Printf("Error: Asset location does not exist! - \"%s\"\n", fileLocation)
		return false
	}

	if temp.IsDir() {
		fmt.Printf("Error: Asset is a folder! - \"%s\"\n", fileLocation)
		return false
	}

	fileData, err := os.ReadFile(fileLocation)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return false
	}

	switch asset.DetermineAssetType(fileData) {
	case asset.JPEG:
		assetToAdd := asset.CreateNewAsset(assetName, asset.JPEG, fileData)
		w.assets = append(w.assets, assetToAdd)
		err = os.WriteFile(osspecifics.CreatePath(w.databaseDir, assetName), fileData, 0444)

		if err != nil {
			fmt.Printf("Error: Failed to add \"%s\" as an asset.\nError: ", assetName)
			generalerrors.HandleError(err)
			return false
		}

		fmt.Printf("Success: Added \"%s\" as an asset.\nFormat: JPEG\nSize: %d bytes\n", assetName, len(fileData))
		return true
	case asset.PDF:
		assetToAdd := asset.CreateNewAsset(assetName, asset.PDF, fileData)
		w.assets = append(w.assets, assetToAdd)
		err = os.WriteFile(osspecifics.CreatePath(w.databaseDir, assetName), fileData, 0444)

		if err != nil {
			fmt.Printf("Error: Failed to add \"%s\" as an asset.\nError: ", assetName)
			generalerrors.HandleError(err)
			return false
		}

		fmt.Printf("Succes: Added \"%s\" as an asset.\nFormat: PDF\nSize: %d bytes\n", assetName, len(fileData))
		return true
	default:
		fmt.Printf("Error: Failed to add \"%s\" as an asset.\nError: Unknown format!\n", assetName)
		return false
	}
}

// RemoveAsset - This function will remove an asset from the wallet.
func (w *Wallet) RemoveAsset(assetName string) bool {
	a := w.getAsset(assetName)

	if a != nil {
		err := os.Remove(osspecifics.CreatePath(w.databaseDir, assetName))

		if err != nil {
			generalerrors.HandleError(err)
			return false
		}
	} else {
		return false
	}

	// This removes the asset from the slice of assets
	for i, as := range w.assets {
		if as == a {
			w.assets = append(w.assets[:i], w.assets[i+1:]...)
			break
		}
	}

	return true
}

func (w *Wallet) GetUsername() string {
	return w.username
}

func (w *Wallet) ConfirmPassword(pass string) bool {
	passBytes := sha256.Sum256([]byte(pass))

	passBytesStr := fmt.Sprintf("%X", passBytes)

	if w.password == passBytesStr {
		w.password = passBytesStr
		return true
	}

	return false
}
