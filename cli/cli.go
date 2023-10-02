package cli

import (
	"errors"
	"flag"
	"fmt"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"io/fs"
	"os"
)

const (
	NoValuePassed              = "NO_VALUE_PASSED"
	HelpCalled                 = 0
	NoPassOrUser               = 1
	CreateNewBCNoDBGiven       = 2
	CreateNewWalletNoDBGiven   = 3
	BlockchainFolderIsNotEmpty = 4
	WalletFolderIsNotEmpty     = 5
)

type FlagValues struct {
	Username      string
	Password      string
	NewBlockchain bool
	NewWallet     bool
	BlockchainDir string
	WalletDir     string
}

// displayHelp - will be used when the help flag is called, or when user fails to comply to execution requirements.
func displayHelp() {
	fmt.Printf("Usage: <exec> (-u <string> & -p <string>) [ACTIONS]\n\n")

	fmt.Print("  -h          \n      Display help menu.\n\n")
	fmt.Print("  -u <string> \n      Input the username of the wallet you want to log in.\n\n")
	fmt.Print("  -p <string> \n      Input the password of the wallet you want to log in.\n\n")
	fmt.Print("  -nb         \n      Creates a new instance of the blockchain.\n\n")
	fmt.Print("  -nw         \n      Creates a new instance of a wallet.\n\n")
	fmt.Print("  -db <string>\n      Input the location of the database of the blockchain. Should only be used when `nb` flag is set, otherwise ineffective.\n\n")
	fmt.Print("  -dw <string>\n      Input the location of the database of the wallet.  Should only be used when `nw` flag is set, otherwise ineffective.\n\n")
}

func InitFlags() *FlagValues {
	H := flag.Bool("h", false, "")
	U := flag.String("u", NoValuePassed, "")
	P := flag.String("p", NoValuePassed, "")
	NewBlockchain := flag.Bool("nb", false, "")
	NewWallet := flag.Bool("nw", false, "")
	DbBc := flag.String("db", NoValuePassed, "")
	DbWallet := flag.String("dw", NoValuePassed, "")

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

	if *NewBlockchain && *DbBc == NoValuePassed {
		fmt.Print("Error: Cannot create a new Blockchain instance without a folder for the database!\n\n")
		//displayHelp()
		os.Exit(CreateNewBCNoDBGiven)
	}

	if *NewWallet && *DbWallet == NoValuePassed {
		fmt.Print("Error: Cannot create a new Wallet without a folder for the database!\n\n")
		//displayHelp()
		os.Exit(CreateNewWalletNoDBGiven)
	}

	return &FlagValues{
		Username:      *U,
		Password:      *P,
		NewBlockchain: *NewBlockchain,
		NewWallet:     *NewWallet,
		BlockchainDir: *DbBc,
		WalletDir:     *DbWallet,
	}

}

// Execute - will execute the action chosen by the user on the blockchain, local and remote.
func Execute(fv *FlagValues) {
	var BC *blockchain.BlockChain
	//var Wallet *wallet.Wallet
	var err error
	var files []fs.DirEntry

	if fv.NewBlockchain {
		files, err = os.ReadDir(fv.BlockchainDir)

		if err != nil {
			generalerrors.HandleError(err, err)
		}

		if len(files) > 0 {
			err := errors.New("Folder contains files. " +
				"Before you select a folder as a database, clear it or make a new one.\n")
			AEE := &generalerrors.AllErrorsExit{ExitCode: BlockchainFolderIsNotEmpty}
			generalerrors.HandleError(err, AEE)
		}

		BC, err = blockchain.CreateNewBlockchain(fv.BlockchainDir)

		if err != nil {
			generalerrors.HandleError(err, err)
		}

		fmt.Print("Created a new Blockchain instance.\n")
	} else {
		BC, err = blockchain.ImportChain()

		if err != nil {
			generalerrors.HandleError(err, err)
		}
	}

	err = BC.ExportChain()

	if err != nil {
		generalerrors.HandleError(err, err)
	}

	//To add later for wallet
}
