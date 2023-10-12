package blockchain

import (
	"encoding/json"
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
		return nil, &generalerrors.JSONUnMarshalFailed{Object: "BlockChain"}
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

	return bc, nil

}

func (bc *BlockChain) ExportChain() error {
	exePath, err := os.Executable()

	if err != nil {
		log.Printf("Error: %s\n", err)
		return err
	}

	dir := filepath.Dir(exePath)

	log.Printf("Exporting BlockChain state...\n")
	bcIE := blockchainIE{
		DatabaseDir:   bc.mDatabaseDir,
		LastBlockHash: fmt.Sprintf("%X", bc.mLastBlock.mCurrHash),
	}

	bytesToWrite, err := json.MarshalIndent(bcIE, "", "    ")

	if err != nil {
		return &generalerrors.JSONMarshalFailed{Object: "BlockChain"}
	}
	path := osspecifics.CreatePath(dir, "bcs.json")

	err = os.WriteFile(path, bytesToWrite, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: path}
	}

	for _, b := range bc.mBlocks {
		err = ExportBlock(bc.GetDBLocation(), b)

		if err != nil {
			log.Printf("Error: Failed to export block!\n")
			return &generalerrors.FailedExport{Object: "Block"}
		}

	}

	log.Print("Blockchain state exported successfully!\n")
	return nil
}
