package wallet

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"log"
	"os"
	"path/filepath"
)

type walletIE struct {
	Username    string         `json:"Username"`
	Password    [32]byte       `json:"Password"`
	PrivateKey  rsa.PrivateKey `json:"PrivateKey"`
	DatabaseDir string         `json:"DatabaseDir"`
	Assets      []string       `json:"Assets"`
}

func ImportWallet(username string) (*Wallet, error) {
	exePath, err := os.Executable()

	if err != nil {
		log.Printf("Error: %s\n", err)
		return nil, err
	}

	dir := filepath.Dir(exePath)

	path := osspecifics.CreatePath(dir, username)
	allBytes, err := os.ReadFile(path)

	if err != nil {
		return nil, &generalerrors.UserNotFound{UserName: username}
	}

	var wie walletIE

	err = json.Unmarshal(allBytes, &wie)

	if err != nil {
		return nil, &generalerrors.JSONUnMarshalFailed{Object: "Wallet"}
	}

	var assetSlice []*asset.Asset

	for _, a := range wie.Assets {
		path := osspecifics.CreatePath(wie.DatabaseDir, a)

		bytesFromFile, err := os.ReadFile(path)

		if err != nil {
			return nil, &generalerrors.ReadFileFailed{File: path}
		}

		ft := asset.DetermineType(bytesFromFile)

		if ft == asset.UNKNOWN {
			return nil, &generalerrors.UnknownFormat{FileExt: path}
		}

		newAsset := asset.CreateNewAsset(a, ft, bytesFromFile)

		assetSlice = append(assetSlice, newAsset)
	}

	w := &Wallet{
		mUsername:    wie.Username,
		mPassword:    wie.Password,
		mPublicKey:   wie.PrivateKey.PublicKey,
		mPrivateKey:  wie.PrivateKey,
		mDatabaseDir: wie.DatabaseDir,
		mAssets:      assetSlice,
	}

	return w, nil
}

func (w *Wallet) ExportWallet() error {
	exePath, err := os.Executable()

	if err != nil {
		log.Printf("Error: %s\n", err)
		return err
	}

	dir := filepath.Dir(exePath)

	log.Printf("Exporting Wallet state...\n")

	var walletAssets []string

	for _, a := range w.mAssets {
		walletAssets = append(walletAssets, a.GetName())
	}

	wie := walletIE{
		Username:    w.mUsername,
		Password:    w.mPassword,
		PrivateKey:  w.mPrivateKey,
		DatabaseDir: w.mDatabaseDir,
		Assets:      walletAssets,
	}

	bytesToWrite, err := json.MarshalIndent(wie, "", "    ")

	if err != nil {
		return &generalerrors.JSONMarshalFailed{Object: "Wallet"}
	}

	path := osspecifics.CreatePath(dir, wie.Username)

	err = os.WriteFile(path, bytesToWrite, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: path}
	}

	log.Printf("Wallet state exported successfully!\n")
	return nil
}
