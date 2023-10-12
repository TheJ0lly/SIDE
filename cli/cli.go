package cli

import (
	"flag"
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/TheJ0lly/GoChain/wallet"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	NoValuePassed = "NO_VALUE_PASSED"
	Success       = iota
	HelpCalled
	NoPassOrUser
	CreateNewWalletNoDBGiven
	WrongNumberOfArgsGivenToOp
	AddAssetFailed
	RemoveAssetFailed
	WrongPass
	FailedGetBC
	FailedDeleteWallet
	FailedGetWallet
	FailedToGetExeFolder
	UnknownOperation
)

type OPERATION int

const (
	AddAsset OPERATION = iota
	RemoveAsset
	ViewAssets
)

type FlagValues struct {
	Username         string
	Password         string
	NewWallet        bool
	WalletDir        string
	Operation        string
	DeleteBCSave     bool
	DeleteWalletSave bool
}

// displayHelp - will be used when the help flag is called, or when user fails to comply to execution requirements.
func displayHelp() {
	fmt.Printf("Usage: <exec> (-u <string> & -p <string>) [ACTIONS]\n\n")

	fmt.Print("  -h          \n      Display help menu.\n")
	fmt.Print("  -u <string> \n      Input the username of the wallet you want to log in.\n")
	fmt.Print("  -p <string> \n      Input the password of the wallet you want to log in.\n")
	fmt.Print("  -nw         \n      Creates a new instance of a wallet.\n")
	fmt.Print("  -dw <string>\n      Input the location of the database of the wallet. Effective only if `nw` flag is used.\n")
	fmt.Print("  -DW <string>\n      Delete the wallet of an user.\n")
	fmt.Print("  -op <string>\n      Input the name of the operation you want to perform:\n")
	fmt.Print("        AddAsset <New Asset Name:string> <Initial location on machine:string>\n")
	fmt.Print("        RemoveAsset <Asset Name:string>\n")
	fmt.Print("        ViewAssets\n")
}

func InitFlags() *FlagValues {
	H := flag.Bool("h", false, "")
	U := flag.String("u", NoValuePassed, "")
	P := flag.String("p", NoValuePassed, "")
	NewWallet := flag.Bool("nw", false, "")
	DbWallet := flag.String("dw", NoValuePassed, "")
	Operation := flag.String("op", NoValuePassed, "")
	DeleteWalletSave := flag.Bool("DW", false, "")

	flag.Usage = displayHelp

	flag.Parse()

	if *H {
		displayHelp()
		os.Exit(HelpCalled)
	}

	if *U == NoValuePassed || *P == NoValuePassed {
		displayHelp()
		os.Exit(NoPassOrUser)
	}

	if *DbWallet != NoValuePassed {
		if *NewWallet == false {
			fmt.Print("Warning: flag 'dw' is ineffective.\n")
		}
	} else {
		if *NewWallet {
			fmt.Print("Error: Cannot create a new Wallet without a folder for the database!\n\n")
			os.Exit(CreateNewWalletNoDBGiven)
		}
	}

	return &FlagValues{
		Username:         *U,
		Password:         *P,
		NewWallet:        *NewWallet,
		WalletDir:        *DbWallet,
		Operation:        *Operation,
		DeleteWalletSave: *DeleteWalletSave,
	}

}

func getBlockchain() (*blockchain.BlockChain, error) {
	var BC *blockchain.BlockChain
	var err error

	// Import blockchain
	BC, err = blockchain.ImportChain()

	if err != nil {
		return nil, err
	}

	return BC, nil
}

func getWallet(fv *FlagValues) (*wallet.Wallet, error) {
	var Wallet *wallet.Wallet
	var err error
	var files []fs.DirEntry

	if fv.NewWallet { // Create new wallet
		files, err = os.ReadDir(fv.WalletDir)

		if err != nil {
			return nil, &generalerrors.ReadDirFailed{Dir: fv.WalletDir}
		}

		if len(files) > 0 {
			return nil, &generalerrors.WalletDBHasItems{Dir: fv.WalletDir}
		}

		Wallet, err = wallet.CreateNewWallet(fv.Username, fv.Password, fv.WalletDir)

		if err != nil {
			return nil, err
		}

		log.Printf("Created a new Wallet\n")
	} else { // Import wallet
		Wallet, err = wallet.ImportWallet(fv.Username)

		if err != nil {
			return nil, err
		}
	}

	return Wallet, nil
}

func exportStates(Wallet *wallet.Wallet, BC *blockchain.BlockChain) {

	fmt.Print("\n")
	err := BC.ExportChain()

	if err != nil {
		generalerrors.HandleError(err, err)
	}

	err = Wallet.ExportWallet()

	if err != nil {
		generalerrors.HandleError(err, err)
	}
}

func getOpArgs(op OPERATION) []string {
	args := os.Args
	var opArgs []string

	var operation string

	switch op {
	case AddAsset:
		operation = "AddAsset"
		for i := 0; i < len(args); i++ {
			if args[i] == operation && i < len(args)-2 {
				opArgs = append(opArgs, args[i+1])
				opArgs = append(opArgs, args[i+2])
				break
			}
		}

	case RemoveAsset:
		operation = "RemoveAsset"
		for i := 0; i < len(args); i++ {
			if args[i] == operation && i < len(args)-1 {
				opArgs = append(opArgs, args[i+1])
				break
			}
		}
	case ViewAssets:
		//There is nothing to gather
	}

	return opArgs
}

func walletExists(username string, files []fs.DirEntry) bool {
	for _, f := range files {
		if strings.Contains(f.Name(), username) {
			return true
		}
	}

	return false
}

func performOperation(fv *FlagValues, Wallet *wallet.Wallet, BC *blockchain.BlockChain) int {

	if fv.Operation == NoValuePassed {
		return Success
	}

	switch fv.Operation {
	case "AddAsset":
		args := getOpArgs(AddAsset)

		if len(args) != 2 {
			log.Printf("Error: Operation AddAsset did not receive the right amount of arguments\n")
			return WrongNumberOfArgsGivenToOp
		}

		asset, err := Wallet.AddAsset(args[0], args[1])

		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			log.Printf("Failed to add asset: %s\n", args[0])
			return AddAssetFailed
		}

		err = BC.AddData("ADDED", Wallet.GetUsername(), asset)

		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			log.Printf("Failed to add metadata: %s\n", asset.GetName())
			return AddAssetFailed
		}

		log.Printf("Added Asset \"%s\" successfully!\n", asset.GetName())
		return Success
	case "RemoveAsset":
		args := getOpArgs(RemoveAsset)

		if len(args) != 1 {
			log.Printf("Error: Operation RemoveAsset did not receive the right amount of arguments\n")
			return WrongNumberOfArgsGivenToOp
		}

		asset, err := Wallet.RemoveAsset(args[0])

		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			log.Printf("Failed to remove asset: %s\n", args[0])
			return RemoveAssetFailed
		}

		err = BC.AddData(Wallet.GetUsername(), "REMOVED", asset)

		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			log.Printf("Failed to add metadata: %s\n", asset.GetName())
			return AddAssetFailed
		}

		log.Printf("Removed Asset \"%s\" successfully!\n", asset.GetName())
		return Success
	case "ViewAssets":
		assetSlice := Wallet.ViewAssets()

		if assetSlice == nil {
			log.Printf("There are no assets to show\n")
			return Success
		}

		for _, a := range assetSlice {
			a.PrintInfo()
			fmt.Print("\n")
		}

		return Success

	default:
		return UnknownOperation

	}
}

// Execute - will execute the action chosen by the user on the blockchain, local and remote.
func Execute(fv *FlagValues) {

	var BC *blockchain.BlockChain

	var Wallet *wallet.Wallet

	exePath, err := os.Executable()

	if err != nil {
		log.Printf("Error: %s\n", err)
		os.Exit(FailedToGetExeFolder)
	}

	dir := filepath.Dir(exePath)

	//Blockchain handling
	BC, err = getBlockchain()

	if err != nil {
		generalerrors.HandleError(err)
		os.Exit(FailedGetBC)
	}

	//Wallet handling
	Wallet, err = getWallet(fv)

	if err != nil {
		generalerrors.HandleError(err)
		os.Exit(FailedGetWallet)
	}

	if fv.DeleteWalletSave {

		files, err := os.ReadDir(dir)

		if err != nil {
			generalerrors.HandleError(err)
			os.Exit(FailedDeleteWallet)
		}

		if !walletExists(Wallet.GetUsername(), files) {
			log.Printf("Error: Username \"%s\" does not exist!\n", Wallet.GetUsername())
			os.Exit(FailedDeleteWallet)
		}

		err = osspecifics.ClearFolder(Wallet.GetDBLocation())

		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			log.Printf("Could not delete Wallet save and Assets folder\n")
			os.Exit(FailedDeleteWallet)
		}

		WalletSavePath := osspecifics.CreatePath(dir, Wallet.GetUsername())

		err = os.Remove(WalletSavePath)

		if err != nil {
			generalerrors.HandleError(err)
			log.Printf("Error: Failed to remove the wallet save\n")
			os.Exit(FailedDeleteWallet)
		}

		log.Printf("Successfully deleted Wallet save and Assets folder!\n")
		os.Exit(Success)

	}

	if !Wallet.ConfirmPassword(fv.Password) {
		log.Printf("Wrong password for user: %s\n", fv.Username)
		os.Exit(WrongPass)
	}
	log.Printf("Logged in successfully as: %s\n", Wallet.GetUsername())

	//Perform actions based on Flag Values
	retVal := performOperation(fv, Wallet, BC)

	if retVal != Success {
		if retVal == UnknownOperation {
			log.Printf("Unknown operation: %s\n", fv.Operation)
		}
		os.Exit(retVal)
	}

	if fv.Operation == "ViewAssets" {
		os.Exit(Success)
	}

	//Export states
	exportStates(Wallet, BC)
}
