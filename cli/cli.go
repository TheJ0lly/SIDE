package cli

import (
	"os"

	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/prettyfmt"
	"github.com/TheJ0lly/GoChain/wallet"
)

// This function will handle the login/signup proccess of the user. If either the Wallet or the Blockchain could not be initialized it exits the program.
// Otherwise, it returns the blockchain and wallet pointers.
func Login_Or_Signup() (*blockchain.BlockChain, *wallet.Wallet) {
	bc, err := blockchain.Load_Blockchain()

	if err != nil {
		generalerrors.HandleError(err)

		var new_blockchain_db_loc string

		prettyfmt.CPrint("Hello, there! You will soon become a new node of the Toy Blockchain! Where do you want your new database to be?\n", prettyfmt.GREEN)

		prettyfmt.Scanln(&new_blockchain_db_loc)

		bc, err = blockchain.Initialize_BlockChain(new_blockchain_db_loc)

		if err != nil {
			generalerrors.HandleError(err, &generalerrors.All_Errors_Exit{Exit_Code: 1})
		}

		prettyfmt.CPrint("Fantastic! You have become a node of the Toy Blockchain!\n", prettyfmt.GREEN)
	} else {
		prettyfmt.CPrint("Found save file!\n", prettyfmt.GREEN)
	}

	w, err := wallet.Load_Wallet()

	if err != nil {
		generalerrors.HandleError(err)

		var new_username string
		var new_password string
		var new_db_loc string

		prettyfmt.CPrint("Where do you want to store your new wallet?\n", prettyfmt.GREEN)
		prettyfmt.Scanln(&new_db_loc)

		prettyfmt.CPrint("What is your new username?\n", prettyfmt.GREEN)
		prettyfmt.Scanln(&new_username)

		prettyfmt.CPrint("What is your new password?\n", prettyfmt.GREEN)
		prettyfmt.Scanln(&new_password)

		w = wallet.Initialize_Wallet(new_username, new_password, new_db_loc)

		if w == nil {
			generalerrors.HandleError(err, &generalerrors.All_Errors_Exit{Exit_Code: 2})
		}
	} else {
		prettyfmt.CPrint("Found save file!\n", prettyfmt.GREEN)

		prettyfmt.CPrintf("Found wallet for user: %s!\nEnter password to confirm identity:\n", prettyfmt.BLUE, w.Get_Username())

		var password string
		prettyfmt.Scanln(&password)

		if !w.Confirm_Password(password) {
			prettyfmt.CPrint("Wrong password! Goodbye!\n", prettyfmt.RED)
			os.Exit(3)
		}
	}

	prettyfmt.CPrintf("Welcome, %s!\n", prettyfmt.GREEN, w.Get_Username())

	return bc, w
}
