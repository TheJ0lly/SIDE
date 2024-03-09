package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"log"
	"os"
	"path/filepath"
	"slices"
)

type blockchainIE struct {
	DatabaseDir   string `json:"DatabaseDir"`
	LastBlockHash string `json:"LastBlockHash"`
}

// Lock - will create a file that signals that the blockchain is currently in use.
func (bc *BlockChain) Lock() error {
	log.Printf("INFO: locking the blockchain save file.\n")
	err := osspecifics.LockFile("bcs.json")

	if err != nil {
		return err
	}

	return nil
}

// Unlock - will remove the lock file, thus signaling that the blockchain is ready to use.
func (bc *BlockChain) Unlock() {
	osspecifics.UnlockFile("bcs.json")
	log.Printf("INFO: the blockchain save file has been unlocked.\n")
}

func ImportChain() (*BlockChain, error) {
	exePath, err := os.Executable()

	if err != nil {
		log.Printf("Error: %s\n", err)
		return nil, err
	}

	dir := filepath.Dir(exePath)

	if err != nil {
		return nil, err
	}

	path := osspecifics.CreatePath(dir, "bcs.json")
	allBytes, err := os.ReadFile(path)

	if err != nil {
		return nil, &generalerrors.ReadFileFailed{File: path}
	}

	var bcIE blockchainIE

	err = json.Unmarshal(allBytes, &bcIE)

	if err != nil {
		return nil, err
	}

	bc := &BlockChain{
		mBlocks:      nil,
		mDatabaseDir: bcIE.DatabaseDir,
		mLastBlock:   nil,
	}

	var LastBlockHash = bcIE.LastBlockHash

	for {
		b, err := ImportBlock(osspecifics.CreatePath(bc.mDatabaseDir, LastBlockHash))

		if err != nil {
			return nil, err
		}

		bc.mBlocks = append(bc.mBlocks, b)

		if b.mPrevHash == nil {
			break
		}

		LastBlockHash = fmt.Sprintf("%X", b.mPrevHash)
	}

	slices.Reverse(bc.mBlocks)

	bc.mLastBlock = bc.mBlocks[len(bc.mBlocks)-1]

	if osspecifics.IsLocked("bcs.json") {
		return nil, errors.New("the blockchain save file is locked")
	}

	return bc, nil

}

// ExportChain - will export the whole chain
func (bc *BlockChain) ExportChain() error {
	exePath, err := os.Executable()

	if err != nil {
		log.Printf("Error: %s\n", err)
		return err
	}

	dir := filepath.Dir(exePath)

	log.Printf("INFO: exporting blockchain state...\n")
	bcIE := blockchainIE{
		DatabaseDir:   bc.mDatabaseDir,
		LastBlockHash: fmt.Sprintf("%X", bc.mLastBlock.mCurrHash),
	}

	bytesToWrite, err := json.MarshalIndent(bcIE, "", "    ")

	if err != nil {
		return err
	}
	path := osspecifics.CreatePath(dir, "bcs.json")

	err = os.WriteFile(path, bytesToWrite, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: path}
	}

	for _, b := range bc.mBlocks {
		err = ExportBlock(bc.GetDBLocation(), b)

		if err != nil {
			return &generalerrors.FailedExport{Object: "Block"}
		}

	}

	log.Print("INFO: blockchain state exported successfully!\n")

	return nil
}
