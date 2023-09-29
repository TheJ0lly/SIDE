package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/wallet"
)

const (
	NoValuePassed        = "NO_VALUE_PASSED"
	HelpCalled           = 0
	NoPassOrUser         = 1
	LoadBcNoDbPassed     = 2
	LoadWalletNoDbPassed = 3
)

type FlagValues struct {
	Username      string
	Password      string
	BlockchainDir string
	WalletDir     string
}

// displayHelp - will be used when the help flag is called, or when user fails to comply to execution requirements.
func displayHelp() {
	fmt.Printf("Usage: <exec> (-u <string> & -p <string>) [ACTIONS]\n\n")

	fmt.Print("  -h          \n      Display help menu.\n\n")
	fmt.Print("  -u <string> \n      Input the username of the wallet you want to log in.\n\n")
	fmt.Print("  -p <string> \n      Input the password of the wallet you want to log in.\n\n")
	fmt.Print("  -db <string>\n      Input the location of the database of the blockchain. Only use when creating new instance, otherwise ineffective.\n\n")
	fmt.Print("  -da <string>\n      Input the location of the database of the wallet. Only use when creating new instance, otherwise ineffective.\n\n")

}

func InitFlags() *FlagValues {
	H := flag.Bool("h", false, "")
	U := flag.String("u", NoValuePassed, "")
	P := flag.String("p", NoValuePassed, "")
	DbBc := flag.String("db", NoValuePassed, "")
	DbAssets := flag.String("da", NoValuePassed, "")

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

	if *DbBc == NoValuePassed {
		fmt.Print("Error: Cannot start new instance of a Blockchain without a folder for the database!\n\n")
		displayHelp()
		os.Exit(LoadBcNoDbPassed)
	}

	if *DbAssets == NoValuePassed {
		fmt.Print("Error: Cannot start new instance of a Wallet without a folder for the database!\n\n")
		displayHelp()
		os.Exit(LoadWalletNoDbPassed)
	}

	return &FlagValues{Username: *U, Password: *P, BlockchainDir: *DbBc, WalletDir: *DbAssets}

}

// Execute - will execute the action chosen by the user on the blockchain, local and remote.
func Execute(fv *FlagValues) {

	fmt.Printf("Starting creating new blockchain!\nDatabase location: %s\n\n", fv.BlockchainDir)
	blockchain.CreateNewBlockchain(fv.BlockchainDir)

	fmt.Printf("Starting creating new wallet!\nDatabase location: %s\n\n", fv.WalletDir)
	wallet.CreateNewWallet(fv.Username, fv.Password, fv.WalletDir)

}
