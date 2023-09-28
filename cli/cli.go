package cli

import (
	"flag"
	"os"

	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/prettyfmt"
	"github.com/TheJ0lly/GoChain/wallet"
)

const (
	NoValuePassed        = "NO_VALUE_PASSED"
	NoPassOrUser         = 1
	LoadBcNoDbPassed     = 2
	LoadWalletNoDbPassed = 3
	FailedToLoadBc       = 4
	FailedToLoadWallet   = 5
)

type FlagValues struct {
	Username      string
	Password      string
	BlockchainDir string
	WalletDir     string
	LoadBC        bool
	LoadWallet    bool
}

func DisplayHelp() {
	prettyfmt.ErrorF("Usage: <exec> (-u <string> & -p <string>) [ACTIONS]\n\n")

	prettyfmt.Print("  -h          \n      Display help menu.\n\n")
	prettyfmt.Print("  -u <string> \n      Input the username of the wallet you want to log in.\n\n")
	prettyfmt.Print("  -p <string> \n      Input the password of the wallet you want to log in.\n\n")
	prettyfmt.Print("  -lb         \n      Load the last blockchain save. Otherwise create new instance for every run.\n\n")
	prettyfmt.Print("  -lw         \n      Load the last wallet save. Otherwise create new instance for every run.\n\n")
	prettyfmt.Print("  -db <string>\n      Input the location of the database of the blockchain. Only use when creating new instance, otherwise ineffective.\n\n")
	prettyfmt.Print("  -da <string>\n      Input the location of the database of the wallet. Only use when creating new instance, otherwise ineffective.\n\n")

}

func InitFlags() *FlagValues {
	U := flag.String("u", NoValuePassed, "")
	P := flag.String("p", NoValuePassed, "")
	LoadBc := flag.Bool("lb", false, "")
	LoadWallet := flag.Bool("lw", false, "")
	DbBc := flag.String("db", NoValuePassed, "")
	DbAssets := flag.String("da", NoValuePassed, "")

	flag.Usage = DisplayHelp

	flag.Parse()

	if *U == NoValuePassed || *P == NoValuePassed {
		DisplayHelp()
		os.Exit(NoPassOrUser)
	}

	if !*LoadBc && *DbBc == NoValuePassed {
		prettyfmt.ErrorF("Cannot start new instance of a Blockchain without a folder for the database!\n\n")
		DisplayHelp()
		os.Exit(LoadBcNoDbPassed)
	}

	if !*LoadWallet && *DbAssets == NoValuePassed {
		prettyfmt.ErrorF("Cannot start new instance of a Wallet without a folder for the database!\n\n")
		DisplayHelp()
		os.Exit(LoadWalletNoDbPassed)
	}

	return &FlagValues{Username: *U, Password: *P, BlockchainDir: *DbBc, WalletDir: *DbAssets, LoadBC: *LoadBc, LoadWallet: *LoadWallet}

}

func Execute(fv *FlagValues) {
	if fv.LoadBC {
		// _, err := blockchain.Load_Blockchain()

		// if err != nil {
		// 	generalerrors.HandleError(err)
		// 	prettyfmt.ErrorF("Failed to load blockchain!\n")
		// 	os.Exit(FAILED_TO_LOAD_BC)
		// }
	} else {
		prettyfmt.CPrintf("Starting creating new blockchain!\nDatabase location: %s\n\n", prettyfmt.BLUE, fv.BlockchainDir)
		blockchain.CreateNewBlockchain(fv.BlockchainDir)
	}

	if fv.LoadWallet {
		// _, err := wallet.Load_Wallet()

		// if err != nil {
		// 	generalerrors.HandleError(err)
		// 	prettyfmt.ErrorF("Failed to load wallet!\n")
		// 	os.Exit(FAILED_TO_LOAD_WALLET)
		// }
	} else {
		prettyfmt.CPrintf("Starting creating new wallet!\nDatabase location: %s\n\n", prettyfmt.BLUE, fv.WalletDir)
		wallet.CreateNewWallet(fv.Username, fv.Password, fv.WalletDir)
	}
}
