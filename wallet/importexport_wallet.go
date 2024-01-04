package wallet

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/multiformats/go-multiaddr"
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
	Addresses   []string       `json:"Addresses"`
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
		return nil, err
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

	var MASlice []multiaddr.Multiaddr

	for _, a := range wie.Addresses {
		ma, err := multiaddr.NewMultiaddr(a)

		if err != nil {
			return nil, err
		}

		MASlice = append(MASlice, ma)
	}

	w := &Wallet{
		mUsername:    wie.Username,
		mPassword:    wie.Password,
		mPublicKey:   wie.PrivateKey.PublicKey,
		mPrivateKey:  wie.PrivateKey,
		mDatabaseDir: wie.DatabaseDir,
		mAssets:      assetSlice,
		mAddresses:   MASlice,
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

	var MAString []string

	for _, a := range w.mAddresses {
		MAString = append(MAString, a.String())
	}

	wie := walletIE{
		Username:    w.mUsername,
		Password:    w.mPassword,
		PrivateKey:  w.mPrivateKey,
		DatabaseDir: w.mDatabaseDir,
		Assets:      walletAssets,
		Addresses:   MAString,
	}

	bytesToWrite, err := json.MarshalIndent(wie, "", "    ")

	if err != nil {
		return err
	}

	path := osspecifics.CreatePath(dir, wie.Username)

	err = os.WriteFile(path, bytesToWrite, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: path}
	}

	log.Printf("Wallet state exported successfully!\n")
	return nil
}
