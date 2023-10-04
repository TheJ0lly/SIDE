package wallet

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"os"
)

type walletIE struct {
	Username    string         `json:"Username"`
	Password    [32]byte       `json:"Password"`
	PrivateKey  rsa.PrivateKey `json:"PrivateKey"`
	DatabaseDir string         `json:"DatabaseDir"`
	Assets      []string       `json:"Assets"`
}

func ImportWallet(username string) (*Wallet, error) {
	dir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	path := osspecifics.CreatePath(dir, username)
	allBytes, err := os.ReadFile(path)

	if err != nil {
		//return nil, &generalerrors.ReadFileFailed{File: path}
		return nil, errors.New("No user \"" + username + "\" has been found.")
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

		ft := asset.DetermineAssetType(bytesFromFile)

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
	dir, err := os.Getwd()

	if err != nil {
		return err
	}

	fmt.Printf("Exporting Wallet state...\n")

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

	fmt.Printf("Wallet state exported successfully!\n")
	return nil
}
