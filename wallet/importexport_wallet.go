package wallet

import (
	"encoding/json"
	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
	"log"
	"os"
	"path/filepath"
)

type walletIE struct {
	Username    string   `json:"Username"`
	Password    [32]byte `json:"Password"`
	PrivateKey  []byte   `json:"PrivateKey"`
	DatabaseDir string   `json:"DatabaseDir"`
	Assets      []string `json:"Assets"`
	Addresses   []string `json:"Addresses"`
	KnownHosts  []string `json:"KnownHosts"`
}

func ImportWallet(username string) (*Wallet, error) {
	exePath, err := os.Executable()

	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(exePath)

	path := osspecifics.CreatePath(dir, username, "config")

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

		newAsset := asset.CreateNewAsset(a, bytesFromFile)

		assetSlice = append(assetSlice, newAsset)
	}

	PrivKey, err := crypto.UnmarshalPrivateKey(wie.PrivateKey)

	if err != nil {
		return nil, err
	}

	host, err := libp2p.New(
		libp2p.ListenAddrStrings(wie.Addresses...),
		libp2p.Identity(PrivKey),
	)

	if err != nil {
		log.Printf("ERROR: could not initialize node\n")
		return nil, err
	}

	var KnownHosts []core.Multiaddr

	for _, kh := range wie.KnownHosts {
		ma, err := multiaddr.NewMultiaddr(kh)

		if err != nil {
			return nil, err
		}

		KnownHosts = append(KnownHosts, ma)
	}

	w := &Wallet{
		mUsername:    wie.Username,
		mPassword:    wie.Password,
		mPublicKey:   PrivKey.GetPublic(),
		mPrivateKey:  PrivKey,
		mDatabaseDir: wie.DatabaseDir,
		mAssets:      assetSlice,
		mHost:        host,
		mKnownHosts:  KnownHosts,
	}

	return w, nil
}

func (w *Wallet) ExportWallet() error {
	exePath, err := os.Executable()

	if err != nil {
		return err
	}

	dir := filepath.Dir(exePath)

	log.Printf("INFO: exporting Wallet state...\n")

	var walletAssets []string

	for _, a := range w.mAssets {
		walletAssets = append(walletAssets, a.GetName())
	}

	var Addresses []string

	for _, a := range w.mHost.Addrs() {
		Addresses = append(Addresses, a.String())
	}

	var KnownHosts []string

	for _, kh := range w.mKnownHosts {
		KnownHosts = append(KnownHosts, kh.String())
	}

	b, err := crypto.MarshalPrivateKey(w.mPrivateKey)

	if err != nil {
		return err
	}

	wie := walletIE{
		Username:    w.mUsername,
		Password:    w.mPassword,
		PrivateKey:  b,
		DatabaseDir: w.mDatabaseDir,
		Assets:      walletAssets,
		Addresses:   Addresses,
		KnownHosts:  KnownHosts,
	}

	bytesToWrite, err := json.MarshalIndent(wie, "", "    ")

	if err != nil {
		return err
	}

	path := osspecifics.CreatePath(dir, wie.Username)

	_, err = os.Stat(path)

	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0777)
		if err != nil {
			return &generalerrors.FailedToCreateUserDir{UserName: wie.Username}
		}
	}

	path = osspecifics.CreatePath(path, "config")

	err = os.WriteFile(path, bytesToWrite, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: path}
	}

	log.Printf("INFO: wallet state exported successfully!\n")
	return nil
}
