package cli

import (
	"flag"
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/TheJ0lly/GoChain/wallet"
	"io/fs"
	"os"
	"strings"
)

const (
	NoValuePassed              = "NO_VALUE_PASSED"
	HelpCalled                 = 0
	NoPassOrUser               = 1
	CreateNewBCNoDBGiven       = 2
	CreateNewWalletNoDBGiven   = 3
	BlockchainFolderIsNotEmpty = 4
	WalletFolderIsNotEmpty     = 5
	WrongNumberOfArgsGivenToOp = 6
	AddAssetFailed             = 7
	WrongPass                  = 8
	FailedDeleteBC             = 9
	FailedDeleteWallet         = 10
	ExitAfter                  = 11
	FailedToGetCWD             = 12
)

type OPERATION int

const (
	AddAsset OPERATION = iota
)

type FlagValues struct {
	Username         string
	Password         string
	NewBlockchain    bool
	NewWallet        bool
	BlockchainDir    string
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
	fmt.Print("  -nb         \n      Creates a new instance of the blockchain.\n")
	fmt.Print("  -nw         \n      Creates a new instance of a wallet.\n")
	fmt.Print("  -db <string>\n      Input the location of the database of the blockchain. Effective only if `nb` flag is used.\n")
	fmt.Print("  -dw <string>\n      Input the location of the database of the wallet. Effective only if `nw` flag is used.\n")
	fmt.Print("  -DB         \n      Delete the blockchain from the machine.\n")
	fmt.Print("  -DW <string>\n      Delete the wallet of an user.\n")
	fmt.Print("  -op <string>\n      Input the name of the operation you want to perform:\n")
	fmt.Print("        AddAsset <New Asset Name:string> <Initial location on machine:string>\n")
}

func InitFlags() *FlagValues {
	H := flag.Bool("h", false, "")
	U := flag.String("u", NoValuePassed, "")
	P := flag.String("p", NoValuePassed, "")
	NewBlockchain := flag.Bool("nb", false, "")
	NewWallet := flag.Bool("nw", false, "")
	DbBc := flag.String("db", NoValuePassed, "")
	DbWallet := flag.String("dw", NoValuePassed, "")
	Operation := flag.String("op", NoValuePassed, "")
	DeleteBlockchainSave := flag.Bool("DB", false, "")
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

	if *DbBc != NoValuePassed {
		if *NewBlockchain == false {
			fmt.Print("Warning: flag 'db' is ineffective.\n")
		}
	} else {
		if *NewBlockchain {
			fmt.Print("Error: Cannot create a new Blockchain instance without a folder for the database!\n\n")
			os.Exit(CreateNewBCNoDBGiven)
		}
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
		NewBlockchain:    *NewBlockchain,
		NewWallet:        *NewWallet,
		BlockchainDir:    *DbBc,
		WalletDir:        *DbWallet,
		Operation:        *Operation,
		DeleteBCSave:     *DeleteBlockchainSave,
		DeleteWalletSave: *DeleteWalletSave,
	}

}

func getBlockchain(fv *FlagValues) *blockchain.BlockChain {
	var BC *blockchain.BlockChain
	var err error
	var files []fs.DirEntry

	if fv.NewBlockchain { // Create new blockchain
		files, err = os.ReadDir(fv.BlockchainDir)

		if err != nil {
			generalerrors.HandleError(err, err)
		}

		if len(files) > 0 {
			err := &generalerrors.BlockchainDBHasItems{Dir: fv.BlockchainDir}
			AEE := &generalerrors.AllErrorsExit{ExitCode: BlockchainFolderIsNotEmpty}
			generalerrors.HandleError(err, AEE)
		}

		BC, err = blockchain.CreateNewBlockchain(fv.BlockchainDir)

		if err != nil {
			generalerrors.HandleError(err, err)
		}

		fmt.Print("Created a new Blockchain instance.\n")
	} else { // Import blockchain
		BC, err = blockchain.ImportChain()

		if err != nil {
			generalerrors.HandleError(err, err)
		}
	}

	return BC
}

func getWallet(fv *FlagValues) *wallet.Wallet {
	var Wallet *wallet.Wallet
	var err error
	var files []fs.DirEntry

	if fv.NewWallet { // Create new wallet
		files, err = os.ReadDir(fv.WalletDir)

		if err != nil {
			generalerrors.HandleError(err, err)
		}

		if len(files) > 0 {
			err := &generalerrors.WalletDBHasItems{Dir: fv.WalletDir}
			AEE := &generalerrors.AllErrorsExit{ExitCode: WalletFolderIsNotEmpty}
			generalerrors.HandleError(err, AEE)
		}

		Wallet, err = wallet.CreateNewWallet(fv.Username, fv.Password, fv.WalletDir)

		if err != nil {
			generalerrors.HandleError(err, err)
		}

		fmt.Printf("Created a new Wallet\n")
	} else { // Import wallet
		Wallet, err = wallet.ImportWallet(fv.Username)

		if err != nil {
			generalerrors.HandleError(err, err)
		}
	}

	return Wallet
}

func exportStates(Wallet *wallet.Wallet, BC *blockchain.BlockChain) {
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
	}

	for i := 0; i < len(args); i++ {
		if args[i] == operation {
			opArgs = append(opArgs, args[i+1])
			opArgs = append(opArgs, args[i+2])
			break
		}
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

// Execute - will execute the action chosen by the user on the blockchain, local and remote.
func Execute(fv *FlagValues) {

	var BC *blockchain.BlockChain
	var Wallet *wallet.Wallet
	var exitAfter = false

	dir, err := os.Getwd()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(FailedToGetCWD)
	}

	//Blockchain handling
	BC = getBlockchain(fv)
	//Wallet handling
	Wallet = getWallet(fv)

	if fv.DeleteBCSave {
		err := osspecifics.ClearFolder(BC.GetDBLocation())

		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			fmt.Printf("Could not delete BlockChain save and folder\n")
			os.Exit(FailedDeleteBC)
		}

		bcSavePath := osspecifics.CreatePath(dir, "bcs.json")

		err = os.Remove(bcSavePath)

		if err != nil {
			generalerrors.HandleError(err)
			fmt.Printf("Error: Failed to remove the blockchain save\n")
			os.Exit(FailedDeleteBC)
		}

		fmt.Printf("Successfully deleted blockchain save and folder!\n")
		exitAfter = true
	}

	if fv.DeleteWalletSave {

		files, err := os.ReadDir(dir)

		if err != nil {
			generalerrors.HandleError(err)
			os.Exit(FailedDeleteWallet)
		}

		if !walletExists(Wallet.GetUsername(), files) {
			fmt.Printf("Error: Username \"%s\" does not exist!\n", Wallet.GetUsername())
			os.Exit(FailedDeleteWallet)
		}

		err = osspecifics.ClearFolder(Wallet.GetDBLocation())

		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			fmt.Printf("Could not delete Wallet save and Assets folder\n")
			os.Exit(FailedDeleteBC)
		}

		WalletSavePath := osspecifics.CreatePath(dir, Wallet.GetUsername()+".json")

		err = os.Remove(WalletSavePath)

		if err != nil {
			generalerrors.HandleError(err)
			fmt.Printf("Error: Failed to remove the wallet save\n")
			os.Exit(FailedDeleteBC)
		}

		fmt.Printf("Successfully deleted Wallet save and Assets folder!\n")
		exitAfter = true
	}

	if exitAfter {
		os.Exit(ExitAfter)
	}

	if !Wallet.ConfirmPassword(fv.Password) {
		fmt.Printf("Wrong password for user: %s\n", fv.Username)
		os.Exit(WrongPass)
	}
	fmt.Printf("Logged in successfully as: %s\n", Wallet.GetUsername())

	//Perform actions based on Flag Values

	switch fv.Operation {
	case "AddAsset":
		args := getOpArgs(AddAsset)

		if len(args) != 2 {
			fmt.Printf("Error: Operation AddAsset did not receive the right amount of arguments\n")
			os.Exit(WrongNumberOfArgsGivenToOp)
		}

		asset, err := Wallet.AddAsset(args[0], args[1])

		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			fmt.Printf("Failed to add asset: %s\n", args[0])
			os.Exit(AddAssetFailed)
		}

		err = BC.AddData(Wallet.GetUsername(), Wallet.GetUsername(), asset)

		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			fmt.Printf("Failed to add asset: %s\n", asset.GetName())
			os.Exit(AddAssetFailed)
		}

		fmt.Printf("Added Asset \"%s\" successfully!\n", asset.GetName())
	}

	//Export states
	exportStates(Wallet, BC)
}
