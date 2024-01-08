package cli

import (
	"errors"
	"flag"
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/network"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/TheJ0lly/GoChain/wallet"
	"github.com/multiformats/go-multiaddr"
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
	WrongNumberOfArgsGivenToOp
	AddAssetFailed
	RemoveAssetFailed
	SendFailed
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
	Send
)

type FlagValues struct {
	Username         string
	Password         string
	NewWallet        string
	Operation        string
	DeleteWalletSave bool
	IP6Def           bool
	IP4Def           bool
	Addrs            []string
}

// displayHelp - will be used when the help flag is called, or when user fails to comply to execution requirements.
func displayHelp() {
	fmt.Printf("Usage: <exec> (-u <string> & -p <string>) [ACTIONS]\n\n")

	fmt.Print("  -h             \n      Display help menu.\n")
	fmt.Print("  -u <string>    \n      Input the username of the wallet you want to log in.\n")
	fmt.Print("  -p <string>    \n      Input the password of the wallet you want to log in.\n")
	fmt.Print("  -n             \n      Creates a new instance of a wallet.\n")
	fmt.Print("  -d             \n      Delete the wallet.\n")
	fmt.Print("  -ip6           \n      Allow the auto-search for available IPv6 addresses when creating the wallet.\n")
	fmt.Print("  -ip4           \n      Allow the auto-search for available IPv4 addresses when creating the wallet.\n")
	fmt.Print("  -a  <string(s)>\n      Give the address(es) to listen to when creating a new wallet.\n")
	fmt.Print("        Example: (/ip4/192.168.1.1/tcp/8080)\n")
	fmt.Print("        Use address 192.168.1.1(IPv4) on port 8080 to handle a TCP connection\n")
	fmt.Print("  -op <string>   \n      Input the name of the op you want to perform:\n")
	fmt.Print("        AddAsset <New Asset Name:string> <Initial location on machine:string>\n")
	fmt.Print("        RemoveAsset <Asset Name:string>\n")
	fmt.Print("        ViewAssets\n")
	fmt.Print("        Send <Asset Name:string> <Peer address:string>")
}

func InitFlags() *FlagValues {

	H := flag.Bool("h", false, "")
	U := flag.String("u", NoValuePassed, "")
	P := flag.String("p", NoValuePassed, "")
	NewWallet := flag.String("n", NoValuePassed, "")
	Operation := flag.String("op", NoValuePassed, "")
	DeleteWalletSave := flag.Bool("d", false, "")
	IP6 := flag.Bool("ip6", false, "")
	IP4 := flag.Bool("ip4", false, "")
	ADDRS := flag.String("a", NoValuePassed, "")

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

	var Addresses = make([]string, 0)

	if *ADDRS != NoValuePassed {
		Addresses = strings.Split(*ADDRS, " ")
	}

	return &FlagValues{
		Username:         *U,
		Password:         *P,
		NewWallet:        *NewWallet,
		Operation:        *Operation,
		DeleteWalletSave: *DeleteWalletSave,
		IP6Def:           *IP6,
		IP4Def:           *IP4,
		Addrs:            Addresses,
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

	if fv.NewWallet != NoValuePassed { // Create new wallet

		if walletExists(fv.Username) {
			return nil, errors.New(fmt.Sprintf("The user %s already exists!", fv.Username))
		}

		fv.NewWallet = osspecifics.GetFullPathFromArg(fv.NewWallet)

		files, err = os.ReadDir(fv.NewWallet)

		if err != nil {
			return nil, &generalerrors.ReadDirFailed{Dir: fv.NewWallet}
		}

		if len(files) > 0 {
			return nil, &generalerrors.WalletDBHasItems{Dir: fv.NewWallet}
		}

		Wallet, err = wallet.CreateNewWallet(fv.Username, fv.Password, fv.NewWallet, fv.IP4Def, fv.IP6Def, fv.Addrs...)

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
		generalerrors.HandleError(generalerrors.ERROR, err, err)
	}

	err = Wallet.ExportWallet()

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err, err)
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
	case ViewAssets: //There is nothing to gather

	case Send:
		operation = "Send"
		for i := 0; i < len(args); i++ {
			if args[i] == operation && i < len(args)-2 {
				opArgs = append(opArgs, args[i+1])
				opArgs = append(opArgs, args[i+2])
				break
			}
		}
	}

	return opArgs
}

func walletExists(username string) bool {

	exePath, err := os.Executable()

	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	exeDir := filepath.Dir(exePath)

	files, err := os.ReadDir(exeDir)

	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}

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
			log.Printf("ERROR: Operation AddAsset did not receive the right amount of arguments\n")
			return WrongNumberOfArgsGivenToOp
		}

		asset, err := Wallet.AddAsset(args[0], args[1])

		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			log.Printf("Failed to add asset: %s\n", args[0])
			return AddAssetFailed
		}

		err = BC.AddData("ADDED", Wallet.GetUsername(), asset)

		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			log.Printf("Failed to add metadata: %s\n", asset.GetName())
			return AddAssetFailed
		}

		log.Printf("Added Asset \"%s\" successfully!\n", asset.GetName())
		return Success
	case "RemoveAsset":
		args := getOpArgs(RemoveAsset)

		if len(args) != 1 {
			log.Printf("ERROR: Operation RemoveAsset did not receive the right amount of arguments\n")
			return WrongNumberOfArgsGivenToOp
		}

		asset, err := Wallet.RemoveAsset(args[0])

		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			log.Printf("Failed to remove asset: %s\n", args[0])
			return RemoveAssetFailed
		}

		err = BC.AddData(Wallet.GetUsername(), "REMOVED", asset)

		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
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
	case "Send":
		args := getOpArgs(Send)

		if len(args) != 2 {
			log.Printf("ERROR: Operation Send did not receive the right amount of arguments\n")
			return WrongNumberOfArgsGivenToOp
		}

		asset, err := Wallet.GetAsset(args[0])

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return SendFailed
		}

		ma, err := multiaddr.NewMultiaddr(args[1])

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return SendFailed
		}

		err = network.SendTo(Wallet.GetHost(), asset.GetAssetBytes(), ma)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return SendFailed
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

	//Blockchain handling
	BC, err := getBlockchain()

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		os.Exit(FailedGetBC)
	}

	//Wallet handling
	Wallet, err = getWallet(fv)

	if err != nil {
		generalerrors.HandleError(generalerrors.ERROR, err)
		os.Exit(FailedGetWallet)
	}

	if fv.DeleteWalletSave {
		if !walletExists(Wallet.GetUsername()) {
			log.Printf("Error: Username \"%s\" does not exist!\n", Wallet.GetUsername())
			os.Exit(FailedDeleteWallet)
		}

		exePath, err := os.Executable()

		if err != nil {
			log.Printf("Error: %s\n", err)
			os.Exit(FailedToGetExeFolder)
		}

		dir := filepath.Dir(exePath)

		err = osspecifics.ClearFolder(Wallet.GetDBLocation())

		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			log.Printf("Could not delete Wallet save and Assets folder\n")
			os.Exit(FailedDeleteWallet)
		}

		WalletSavePath := osspecifics.CreatePath(dir, Wallet.GetUsername())

		err = os.Remove(WalletSavePath)

		if err != nil {
			generalerrors.HandleError(generalerrors.ERROR, err)
			log.Printf("Error: Failed to remove the wallet save\n")
			os.Exit(FailedDeleteWallet)
		}

		log.Printf("Successfully deleted wallet save and and cleared the assets folder!\n")
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
