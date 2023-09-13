package generalerrors

import (
	"errors"
	"os"

	"github.com/TheJ0lly/GoChain/prettyfmt"
)

// This error means that a block's data/transaction capacity has been reached.
type BlockCapacityReached struct {
	Capacity uint8
}

// This error means that a os.ReadDir has failed.
type ReadDirFailed struct {
	Dir string
}

// This error means that a os.ReadDir has failed.
type ReadFileFailed struct {
	File string
}

// This error means that a json.Marshal has failed.
type JSONMarshalFailed struct {
	Object string
}

// This error means that a json.Unmarshal has failed.
type JSONUnMarshalFailed struct {
	Object string
}

// This error means that a os.Remove has failed.
type RemoveFileFailed struct {
	File string
}

// This error means that some data that is trying to be added is too big.
// SOON WILL BECOME OBSOLETE WHEN REPLACING STRINGS WITH TRANSACTIONS
type DataTooBig struct {
	Data        string
	Data_Length uint8
}

// This error means that a os.WriteFile has failed.
type WriteFileFailed struct {
	File string
}

// This error means that the database directory where the blocks have been stored, is empty.
// Is used only when rebuilding the blockchain, because no files means that there has been 100% some tampering going on.
type BlockChainDBEmpty struct {
	Dir string
}

// This error means that a block file is missing from the reconstruction of the blockchain.
type BlockMissing struct {
	Block_Hash string
}

// ======== ERROR FUNCTIONS TO IMPLEMENT THE ERROR INTERFACE ========

func (bcr *BlockCapacityReached) Error() string {
	return prettyfmt.Sprintf("Block capacity of %d has been reached! Need new block!", bcr.Capacity)
}

func (bcr *BlockCapacityReached) Is(target error) bool {
	_, ok := target.(*BlockCapacityReached)
	return ok
}

func (rdf *ReadDirFailed) Error() string {
	return prettyfmt.Sprintf("Failed to read directory: %s", rdf.Dir)
}

func (rff *ReadFileFailed) Error() string {
	if rff.File == "./bcs" {
		return "There is no save file of the blockchain!"
	} else if rff.File == "./ws" {
		return "There is no save file of the wallet!"
	}
	return prettyfmt.Sprintf("Failed to read file: %s", rff.File)
}

func (jmf *JSONMarshalFailed) Error() string {
	return prettyfmt.Sprintf("Failed to marshal object of type: %s", jmf.Object)
}

func (jumf *JSONUnMarshalFailed) Error() string {
	return prettyfmt.Sprintf("Failed to unmarshal object of type: %s", jumf.Object)
}

func (rff *RemoveFileFailed) Error() string {
	return prettyfmt.Sprintf("Failed to remove file: %s", rff.File)
}

func (dtb *DataTooBig) Error() string {
	return prettyfmt.Sprintf("\"%s\" - is too big! Maximum length allowed: %d!", dtb.Data, dtb.Data_Length)
}

func (wff *WriteFileFailed) Error() string {
	return prettyfmt.Sprintf("Failed to write file: %s", wff.File)
}

func (bcdbe *BlockChainDBEmpty) Error() string {
	return prettyfmt.Sprintf("There are no files in the BlockChain Database directory: %s!", bcdbe.Dir)
}

func (bm *BlockMissing) Error() string {
	return prettyfmt.Sprintf("There is no block with the hash: %s!", bm.Block_Hash)
}

// ======== HANDLE ERROR FUNCTION ========

// This error means, that if any error has occured, just exit with the Exit_Code value.
type All_Errors_Exit struct {
	Exit_Code int
}

func (aee *All_Errors_Exit) Error() string {
	return "Any error will exit the program"
}

// This function will print the errors given.
// If you want to exit on a specific error, just add it after the initial error, and if it matches, the program will exit with 1.
// If you want to exit on all errors, doesn't matter which specific one, just use All_Errors_Exit, and pass the error code you want to exit with.
func HandleError(err error, errors_to_fail ...error) {
	prettyfmt.ErrorF("%s\n", err.Error())

	if len(errors_to_fail) != 0 {
		aee, ok := errors_to_fail[0].(*All_Errors_Exit)

		if ok {
			os.Exit(aee.Exit_Code)
		}

		for _, e := range errors_to_fail {
			if errors.Is(err, e) {
				os.Exit(1)
			}
		}
	}
}
