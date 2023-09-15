package cli

import (
	"flag"
	"os"

	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/prettyfmt"
	"github.com/TheJ0lly/GoChain/wallet"
)

const (
	NO_VALUE_PASSED          = "NO_VALUE_PASSED"
	NO_PASS_OR_USER          = 1
	LOAD_BC_NO_DB_PASSED     = 2
	LOAD_WALLET_NO_DB_PASSED = 3
	FAILED_TO_LOAD_BC        = 4
	FAILED_TO_LOAD_WALLET    = 5
)

type Flag_Values struct {
	Username       string
	Password       string
	Blockchain_Dir string
	Wallet_Dir     string
	LoadBC         bool
	LoadWallet     bool
}

func Display_Help() {
	prettyfmt.ErrorF("Usage: <exec> (-u <string> & -p <string>) [ACTIONS]\n\n")

	prettyfmt.Print("  -h          \n      Display help menu.\n\n")
	prettyfmt.Print("  -u <string> \n      Input the username of the wallet you want to log in.\n\n")
	prettyfmt.Print("  -p <string> \n      Input the password of the wallet you want to log in.\n\n")
	prettyfmt.Print("  -lb         \n      Load the last blockchain save. Otherwise create new instance for every run.\n\n")
	prettyfmt.Print("  -lw         \n      Load the last wallet save. Otherwise create new instance for every run.\n\n")
	prettyfmt.Print("  -db <string>\n      Input the location of the database of the blockchain. Only use when creating new instance, otherwise ineffective.\n\n")
	prettyfmt.Print("  -da <string>\n      Input the location of the database of the wallet. Only use when creating new instance, otherwise ineffective.\n\n")

}

func Init_Flags() *Flag_Values {
	U := flag.String("u", NO_VALUE_PASSED, "")
	P := flag.String("p", NO_VALUE_PASSED, "")
	Load_BC := flag.Bool("lb", false, "")
	Load_Wallet := flag.Bool("lw", false, "")
	DB_BC := flag.String("db", NO_VALUE_PASSED, "")
	DB_ASSETS := flag.String("da", NO_VALUE_PASSED, "")

	flag.Usage = Display_Help

	flag.Parse()

	if *U == NO_VALUE_PASSED || *P == NO_VALUE_PASSED {
		Display_Help()
		os.Exit(NO_PASS_OR_USER)
	}

	if !*Load_BC && *DB_BC == NO_VALUE_PASSED {
		prettyfmt.ErrorF("Cannot start new instance of a Blockchain without a folder for the database!\n\n")
		Display_Help()
		os.Exit(LOAD_BC_NO_DB_PASSED)
	}

	if !*Load_Wallet && *DB_ASSETS == NO_VALUE_PASSED {
		prettyfmt.ErrorF("Cannot start new instance of a Wallet without a folder for the database!\n\n")
		Display_Help()
		os.Exit(LOAD_WALLET_NO_DB_PASSED)
	}

	return &Flag_Values{Username: *U, Password: *P, Blockchain_Dir: *DB_BC, Wallet_Dir: *DB_ASSETS, LoadBC: *Load_BC, LoadWallet: *Load_Wallet}

}

func Execute(fv *Flag_Values) {
	if fv.LoadBC {
		_, err := blockchain.Load_Blockchain()

		if err != nil {
			generalerrors.HandleError(err)
			prettyfmt.ErrorF("Failed to load blockchain!\n")
			os.Exit(FAILED_TO_LOAD_BC)
		}
	} else {
		prettyfmt.CPrintf("Starting creating new blockchain!\nDatabase location: %s\n\n", prettyfmt.BLUE, fv.Blockchain_Dir)
		blockchain.Initialize_BlockChain(fv.Blockchain_Dir)
	}

	if fv.LoadWallet {
		_, err := wallet.Load_Wallet()

		if err != nil {
			generalerrors.HandleError(err)
			prettyfmt.ErrorF("Failed to load wallet!\n")
			os.Exit(FAILED_TO_LOAD_WALLET)
		}
	} else {
		prettyfmt.CPrintf("Starting creating new wallet!\nDatabase location: %s\n\n", prettyfmt.BLUE, fv.Wallet_Dir)
		wallet.Initialize_Wallet(fv.Username, fv.Password, fv.Wallet_Dir)
	}
}
