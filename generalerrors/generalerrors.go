package generalerrors

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// BlockCapacityReached - This error means that a block's data/transaction capacity has been reached.
type BlockCapacityReached struct {
	Capacity uint8
}

// ReadDirFailed - This error means that a call to os.ReadDir has failed.
type ReadDirFailed struct {
	Dir string
}

// ReadFileFailed - This error means that a call to os.ReadFile has failed.
type ReadFileFailed struct {
	File string
}

// RemoveFileFailed -This error means that a call to os.Remove has failed.
type RemoveFileFailed struct {
	File string
}

// DataTooBig - This error means that some data that is trying to be added is too big.
// SOON WILL BECOME OBSOLETE WHEN REPLACING STRINGS WITH TRANSACTIONS
type DataTooBig struct {
	Data       string
	DataLength uint8
}

// WriteFileFailed - This error means that a call to os.WriteFile has failed.
type WriteFileFailed struct {
	File string
}

// BlockChainDBEmpty - This error means that the database directory where the blocks have been stored, is empty.
// Is used only when rebuilding the blockchain, because no files means that there has been 100% some tampering going on.
type BlockChainDBEmpty struct {
	Dir string
}

// BlockMissing - This error means that a block file is missing from the reconstruction of the blockchain.
type BlockMissing struct {
	BlockHash string
}

type WalletDBHasItems struct {
	Dir string
}

type BlockHashDifferent struct {
	BlockHash    string
	ComputedHash string
}

type AssetAlreadyExists struct {
	AssetName string
}

type AssetInitialLocationDoesNotExist struct {
	Location string
}

type AssetIsAFolder struct {
	Location string
}

type UsernameTooLong struct {
	Length int
}

type UnknownFormat struct {
	FileExt string
}

type FailedExport struct {
	Object string
}

type AssetDoesNotExist struct {
	AssetName string
}

type UserNotFound struct {
	UserName string
}

type FailedToCreateUserDir struct {
	UserName string
}

// ======== ERROR FUNCTIONS TO IMPLEMENT THE ERROR INTERFACE ========

func (bcr *BlockCapacityReached) Error() string {
	return fmt.Sprintf("block capacity of %d has been reached - need new block!", bcr.Capacity)
}

func (bcr *BlockCapacityReached) Is(target error) bool {
	var blockCapacityReached *BlockCapacityReached
	ok := errors.As(target, &blockCapacityReached)
	return ok
}

func (rdf *ReadDirFailed) Error() string {
	return fmt.Sprintf("failed to read directory: %s", rdf.Dir)
}

func (rff *ReadFileFailed) Error() string {
	if strings.Contains(rff.File, "bcs.json") {
		return "there is no save file of the blockchain!"
	} else if strings.Contains(rff.File, "ws.json") {
		return "there is no save file of the wallet!"
	}
	return fmt.Sprintf("failed to read file: %s", rff.File)
}

func (rff *RemoveFileFailed) Error() string {
	return fmt.Sprintf("failed to remove file: %s", rff.File)
}

func (dtb *DataTooBig) Error() string {
	return fmt.Sprintf("\"%s\" - is too big - maximum length allowed: %d!", dtb.Data, dtb.DataLength)
}

func (wff *WriteFileFailed) Error() string {
	return fmt.Sprintf("failed to write file: %s", wff.File)
}

func (bcdbe *BlockChainDBEmpty) Error() string {
	return fmt.Sprintf("there are no files in the blockchain database directory: %s!", bcdbe.Dir)
}

func (bm *BlockMissing) Error() string {
	return fmt.Sprintf("there is no block with the hash: %s!", bm.BlockHash)
}

func (w *WalletDBHasItems) Error() string {
	return fmt.Sprintf("folder selected for storing wallet assets contains files - directory: %s", w.Dir)
}

func (bhd *BlockHashDifferent) Error() string {
	return fmt.Sprintf("block hash does not match the computed root hash!\ncomputed hash: %s\nblock hash: %s",
		bhd.ComputedHash, bhd.BlockHash)
}

func (aae *AssetAlreadyExists) Error() string {
	return fmt.Sprintf("there already is an asset with this name: %s", aae.AssetName)
}

func (ailne *AssetInitialLocationDoesNotExist) Error() string {
	return fmt.Sprintf("the file path is invalid: %s", ailne.Location)
}

func (aif *AssetIsAFolder) Error() string {
	return fmt.Sprintf("the location is a folder: %s", aif.Location)
}

func (utl *UsernameTooLong) Error() string {
	return fmt.Sprintf("the entered username passes the maximum length of: %d", utl.Length)
}

func (ufa *UnknownFormat) Error() string {
	return fmt.Sprintf("unknown file format - file extension: %s", ufa.FileExt)
}

func (fe *FailedExport) Error() string {
	return fmt.Sprintf("failed to export: %s", fe.Object)
}

func (ane *AssetDoesNotExist) Error() string {
	return fmt.Sprintf("asset does not exist: %s", ane.AssetName)
}

func (unf *UserNotFound) Error() string {
	return fmt.Sprintf("no user \"%s\" has been found.", unf.UserName)
}

func (fcud *FailedToCreateUserDir) Error() string {
	return fmt.Sprintf("cannot create the dir for the user %s", fcud.UserName)
}

// ======== HANDLE ERROR FUNCTION ========

type LogLevel int

const (
	ERROR LogLevel = iota
	WARNING
	INFO
)

// AllErrorsExit - This error means, that if any error has occurred, just exit with the ExitCode value.
type AllErrorsExit struct {
	ExitCode int
}

func (aee *AllErrorsExit) Error() string {
	return "any error will exit the program"
}

// HandleError - This function will print the errors given.
// If you want to exit on a specific error, just add it after the initial error, and if it matches, the program will exit with 1.
// If you want to exit on all errors, doesn't matter which specific one, just use AllErrorsExit, and pass the error code you want to exit with.
func HandleError(loglevel LogLevel, err error, errorsToFail ...error) {
	switch loglevel {
	case ERROR:
		log.Printf("ERROR: %s\n", err.Error())
	case WARNING:
		log.Printf("WARNING: %s\n", err.Error())
	case INFO:
		log.Printf("INFO: %s\n", err.Error())
	}

	if len(errorsToFail) != 0 {
		var aee *AllErrorsExit
		ok := errors.As(errorsToFail[0], &aee)

		if ok {
			os.Exit(aee.ExitCode)
		}

		for _, e := range errorsToFail {
			if errors.Is(err, e) {
				os.Exit(1)
			}
		}
	}
}
