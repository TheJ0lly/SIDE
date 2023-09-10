package generalerrors

import (
	"errors"
	"os"

	"github.com/TheJ0lly/GoChain/prettyfmt"
)

type BlockCapacityReached struct {
	Capacity uint8
}

type ReadDirFailed struct {
	Dir string
}

type ReadFileFailed struct {
	File string
}

type JSONMarshalFailed struct {
	Object string
}

type JSONUnMarshalFailed struct {
	Object string
}

type RemoveFileFailed struct {
	File string
}

type DataTooBig struct {
	Data        string
	Data_Length uint8
}

type WriteFileFailed struct {
	File string
}

type BlockChainDBEmpty struct {
	Dir string
}

type BlockMissing struct {
	Block_Hash string
}

// ======== ERROR FUNCTIONS TO IMPLEMENT THE ERROR INTERFACE ========

func (bcr *BlockCapacityReached) Error() string {
	return prettyfmt.Sprintf("Block capacity of %d has been reached! Need new block!", bcr.Capacity)
}

func (rdf *ReadDirFailed) Error() string {
	return prettyfmt.Sprintf("Failed to read directory: %s", rdf.Dir)
}

func (rff *ReadFileFailed) Error() string {
	if rff.File == "./bcs" {
		return "There is no save file of the blockchain!"
	} else if rff.File == "./ws" {
		return "There is no save file of the blockchain!"
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

type All_Errors_Exit struct {
	Exit_Code int
}

func (aee *All_Errors_Exit) Error() string {
	return "Any error will exit the program"
}

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
