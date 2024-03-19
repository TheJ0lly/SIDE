package cli

import (
	"errors"
	"flag"
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/netutils"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/TheJ0lly/GoChain/wallet"
	"github.com/howeyc/gopass"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	NoValuePassed = ""
	Success       = iota
	HelpCalled
	NoPassOrUser
	WrongNumberOfArgsGivenToOp
	AddAssetFailed
	RemoveAssetFailed
	AddNodeFailed
	RequestAssetFailed
	WrongPass
	FailedGetBC
	FailedDeleteWallet
	FailedGetWallet
	FailedToGetExeFolder
	UnknownOperation
	FailedExporting
)

type OPERATION int

const (
	AddAsset OPERATION = iota
	RemoveAsset
	ViewAssets
	AddNode
	ViewNodes
	RequestAsset
)

type FlagValues struct {
	Username         string
	Password         string
	NewWallet        string
	Operation        string
	DeleteWalletSave bool
	IP               string
	Port             string
}

// displayHelp - will be used when the help flag is called, or when user fails to comply to execution requirements.
func displayHelp() {
	fmt.Printf("Usage: <exec> (-u <string> & -p <string>) [ACTIONS]\n\n")

	fmt.Print("  -h             \n      Display help menu.\n")
	fmt.Print("  -u <string>    \n      Input the username of the wallet you want to log in.\n")
	fmt.Print("  -n <string>    \n      Creates a new instance of a wallet.\n")
	fmt.Print("  -d             \n      Delete the wallet.\n")
	fmt.Print("  -ip            \n      Bind an IP address to this node.\n")
	fmt.Print("  -p             \n      Bind a port to this nodes address\n")
	fmt.Print("  -op <string>   \n      Input the name of the op you want to perform:\n")
	fmt.Print("        AddAsset <New Asset Name:string> <Initial location on machine:string>\n")
	fmt.Print("        RemoveAsset <Asset Name:string>\n")
	fmt.Print("        ViewAssets\n")
	fmt.Print("        AddNode <Address in MultiAddress format:string>\n")
	fmt.Print("        ViewNodes\n")
	fmt.Print("        Request <AssetName : string>\n")
}

// InitFlags - will initialize the flags that will be used to execute the client.
func InitFlags() *FlagValues {
	H := flag.Bool("h", false, "")
	U := flag.String("u", NoValuePassed, "")
	NewWallet := flag.String("n", NoValuePassed, "")
	Operation := flag.String("op", NoValuePassed, "")
	DeleteWalletSave := flag.Bool("d", false, "")
	IP := flag.String("ip", NoValuePassed, "")
	PORT := flag.String("p", NoValuePassed, "")

	flag.Usage = displayHelp

	flag.Parse()

	if *H {
		flag.Usage()
		os.Exit(HelpCalled)
	}

	if *U == NoValuePassed {
		os.Exit(NoPassOrUser)
	}

	if *NewWallet != NoValuePassed {
		fmt.Printf("Enter password for new user %s:", *U)
	} else {
		fmt.Printf("Enter password for user %s:", *U)
	}

	password, err := gopass.GetPasswd()

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return nil
	}

	return &FlagValues{
		Username:         *U,
		Password:         string(password),
		NewWallet:        *NewWallet,
		Operation:        *Operation,
		DeleteWalletSave: *DeleteWalletSave,
		IP:               *IP,
		Port:             *PORT,
	}

}

// getBlockchain - is a helper function that will import the blockchain if possible, otherwise return error
func getBlockchain() (*blockchain.BlockChain, error) {
	var BC *blockchain.BlockChain
	var err error

	// Import blockchain
	BC, err = blockchain.ImportChain()

	if err != nil {
		return nil, err
	}

	err = BC.Lock()

	if err != nil {
		return nil, err
	}

	return BC, nil
}

// getWallet - is a helper function that will import a wallet if possible, otherwise return error
func getWallet(fv *FlagValues) (*wallet.Wallet, error) {
	var Wallet *wallet.Wallet
	var err error
	var files []fs.DirEntry

	if fv.NewWallet != NoValuePassed { // Create new wallet

		if bv, err := walletExists(fv.Username); bv == true {
			return nil, errors.New(fmt.Sprintf("the user %s already exists!", fv.Username))
		} else if err != nil {
			return nil, err
		}

		fv.NewWallet = osspecifics.GetFullPathFromArg(fv.NewWallet)

		files, err = os.ReadDir(fv.NewWallet)

		if err != nil {
			return nil, &generalerrors.ReadDirFailed{Dir: fv.NewWallet}
		}

		if len(files) > 0 {
			return nil, &generalerrors.WalletDBHasItems{Dir: fv.NewWallet}
		}

		Wallet, err = wallet.CreateNewWallet(fv.Username, fv.Password, fv.NewWallet, fv.IP, fv.Port)

		if err != nil {
			return nil, err
		}

		log.Printf("INFO: created a new wallet\n")
	} else { // Import wallet
		Wallet, err = wallet.ImportWallet(fv.Username)

		if err != nil {
			return nil, err
		}
	}

	return Wallet, nil
}

// exportStates - is a utility function that exports the states of the current wallet and the blockchain
func exportStates(Wallet *wallet.Wallet, BC *blockchain.BlockChain) error {

	fmt.Print("\n")
	err := BC.ExportChain()

	if err != nil {
		return err
	}

	err = Wallet.ExportWallet()

	if err != nil {
		return err
	}

	return nil
}

// getOpArgs - this function will return the arguments to a client operation
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
	case ViewAssets: // there is nothing to gather

	case AddNode:
		operation = "AddNode"
		for i := 0; i < len(args); i++ {
			if args[i] == operation && i < len(args)-1 {
				opArgs = append(opArgs, args[i+1])
				break
			}
		}

	case ViewNodes: // there is nothing to gather

	case RequestAsset:
		operation = "Request"
		for i := 0; i < len(args); i++ {
			if args[i] == operation && i < len(args)-1 {
				opArgs = append(opArgs, args[i+1])
				break
			}
		}
	}

	return opArgs
}

// walletExists - checks if a wallet exists
func walletExists(username string) (bool, error) {

	exePath, err := os.Executable()

	if err != nil {
		return false, err
	}
	exeDir := filepath.Dir(exePath)

	files, err := os.ReadDir(exeDir)

	if err != nil {
		return false, err
	}

	for _, f := range files {
		if strings.Contains(f.Name(), username) {
			return true, nil
		}
	}

	return false, nil
}

// performOperation - will perform the operation chosen by the user, along with the specified flags
func performOperation(fv *FlagValues, Wallet *wallet.Wallet, BC *blockchain.BlockChain) int {

	if fv.Operation == NoValuePassed {
		return Success
	}

	switch fv.Operation {
	case "AddAsset":
		args := getOpArgs(AddAsset)

		if len(args) != 2 {
			log.Printf("ERROR: operation AddAsset did not receive the right amount of arguments\n")
			return WrongNumberOfArgsGivenToOp
		}

		asset, err := Wallet.AddAssetFromLocal(args[0], args[1])

		if err != nil {
			log.Printf("INFO: failed to add asset: %s\n", args[0])
			log.Printf("ERROR: %s\n", err.Error())
			return AddAssetFailed
		}

		err = BC.AddData("ADDED", Wallet.GetUsername(), asset.GetName())

		if err != nil {
			log.Printf("INFO: failed to add metadata: %s\n", asset.GetName())
			log.Printf("ERROR: %s\n", err.Error())
			return AddAssetFailed
		}

		log.Printf("INFO: added asset \"%s\" successfully!\n", asset.GetName())
		netutils.FloodProtocol(Wallet.GetNodesAddresses(), Wallet.GetHost(), BC.GetLastMetaData())
		return Success
	case "RemoveAsset":
		args := getOpArgs(RemoveAsset)

		if len(args) != 1 {
			log.Printf("ERROR: operation RemoveAsset did not receive the right amount of arguments\n")
			return WrongNumberOfArgsGivenToOp
		}

		asset, err := Wallet.RemoveAsset(args[0])

		if err != nil {
			log.Printf("INFO: failed to remove asset: %s\n", args[0])
			log.Printf("ERROR: %s\n", err.Error())
			return RemoveAssetFailed
		}

		err = BC.AddData(Wallet.GetUsername(), "REMOVED", asset.GetName())

		if err != nil {
			log.Printf("INFO: failed to add metadata: %s\n", asset.GetName())
			log.Printf("ERROR: %s\n", err.Error())
			return AddAssetFailed
		}

		log.Printf("INFO: removed asset \"%s\" successfully!\n", asset.GetName())
		netutils.FloodProtocol(Wallet.GetNodesAddresses(), Wallet.GetHost(), BC.GetLastMetaData())
		return Success
	case "ViewAssets":
		assetSlice := Wallet.ViewAssets()

		if assetSlice == nil {
			log.Printf("INFO: there are no assets to show\n")
			return Success
		}

		for _, a := range assetSlice {
			a.PrintInfo()
			fmt.Print("\n")
		}

		return Success
	case "AddNode":
		args := getOpArgs(AddNode)

		if len(args) != 1 {
			log.Printf("ERROR: operation AddNode did not receive the right amount of arguments\n")
			return WrongNumberOfArgsGivenToOp
		}

		ma, err := Wallet.AddNode(args[0])

		if err != nil {
			generalerrors.HandleError(generalerrors.ERROR, err)
			return AddNodeFailed
		}

		log.Printf(fmt.Sprintf("INFO: address %s has been successfully added\n", ma.String()))
		return Success

	case "ViewNodes":
		log.Printf("INFO: host address - %s\n", Wallet.GetHostAddress())

		addresses := Wallet.GetNodesAddresses()

		if addresses == nil {
			log.Printf("INFO: there are no other addresses known\n")
			return Success
		}

		log.Printf("INFO: known addresses:")
		for _, a := range addresses {
			fmt.Printf("  %s\n", a.String())
		}

		return Success

	case "Request":
		args := getOpArgs(RequestAsset)

		if len(args) != 1 {
			log.Printf("ERROR: operation Request did not receive the right amount of arguments\n")
			return WrongNumberOfArgsGivenToOp
		}

		ok, as := netutils.MakeRequest(Wallet.GetNodesAddresses(), Wallet.GetHost(), args[0])

		if ok {
			_, err := Wallet.AddAssetFromNode(as)

			if err != nil {
				log.Printf("ERROR: failed to add asset - %s\n", err)
				return RequestAssetFailed
			}

			return Success
		} else {
			return RequestAssetFailed
		}

	default:
		return UnknownOperation
	}
}

// Execute - will execute the action chosen by the user on the blockchain, local and remote.
func Execute(fv *FlagValues) int {

	var BC *blockchain.BlockChain

	var Wallet *wallet.Wallet

	//Blockchain handling
	BC, err := getBlockchain()

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		return FailedGetBC
	}

	defer BC.Unlock()

	//Wallet handling
	Wallet, err = getWallet(fv)

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		return FailedGetWallet
	}

	if !Wallet.ConfirmPassword(fv.Password) {
		log.Printf("ERROR: wrong password for user: %s\n", fv.Username)
		return WrongPass
	}

	if fv.DeleteWalletSave {
		exePath, err := os.Executable()

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return FailedToGetExeFolder
		}

		dir := filepath.Dir(exePath)

		err = osspecifics.ClearFolder(Wallet.GetDBLocation())

		if err != nil {
			log.Printf("INFO: could not clear assets folder\n")
			log.Printf("ERROR: %s\n", err.Error())
			return FailedDeleteWallet
		}

		WalletSavePath := osspecifics.CreatePath(dir, Wallet.GetUsername())

		err = osspecifics.ClearFolder(WalletSavePath)

		if err != nil {
			log.Printf("INFO: could not delete wallet folder\n")
			log.Printf("ERROR: %s\n", err.Error())
			return FailedDeleteWallet
		}

		err = os.Remove(WalletSavePath)

		if err != nil {
			log.Printf("ERROR: failed to remove the wallet save\n")
			generalerrors.HandleError(generalerrors.ERROR, err)
			return FailedDeleteWallet
		}

		log.Printf("INFO: successfully deleted wallet save and and cleared the assets folder\n")
		return Success

	}

	log.Printf("INFO: logged in successfully as: %s\n", Wallet.GetUsername())

	//Perform actions based on Flag Values
	retVal := performOperation(fv, Wallet, BC)

	if retVal != Success {
		if retVal == UnknownOperation {
			log.Printf("WARNING: unknown operation: %s\n", fv.Operation)
		}
		return retVal
	}

	switch fv.Operation {
	case "ViewAssets":
	case "ViewNodes":
		return Success
	default:
		//Export states
		err = exportStates(Wallet, BC)

		if err != nil {
			generalerrors.HandleError(generalerrors.ERROR, err)
			return FailedExporting
		}
	}

	log.Printf("INFO: writing notification\n")
	err = os.WriteFile("importNotification", nil, 0666)

	if err != nil {
		log.Printf("ERROR: failed to write notification to the service\n")
		return RequestAssetFailed
	}
	log.Printf("INFO: notification wrote\n")
	return Success
}
