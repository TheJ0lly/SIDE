package cli

import (
	"os"

	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/prettyfmt"
	"github.com/TheJ0lly/GoChain/wallet"
)

var BC *blockchain.BlockChain
var W *wallet.Wallet

func Login_Or_Signup() {
	bc := blockchain.Load_Blockchain()

	if bc == nil {
		var new_blockchain_db_loc string

		prettyfmt.Print("Hello, there! You will soon become a new node of the Toy Blockchain! Where do you want your new database to be?\n", prettyfmt.GREEN)

		prettyfmt.Scanln(&new_blockchain_db_loc)

		bc = blockchain.Initialize_BlockChain(new_blockchain_db_loc)

		if bc == nil {
			prettyfmt.Print("Failed to create new instance of the Toy Blockchain!\nExiting...\n", prettyfmt.RED)
			os.Exit(1)
		}
		prettyfmt.Print("Fantastic! You have become a node of the Toy Blockchain!\n", prettyfmt.GREEN)
	} else {
		prettyfmt.Print("Found save file! Restoring Toy Blockchain...\n", prettyfmt.GREEN)
	}

	w := wallet.Load_Wallet()

	if w == nil {

		var new_username string
		var new_password string
		var new_db_loc string

		prettyfmt.Print("Where do you want to store your new wallet?\n", prettyfmt.GREEN)
		prettyfmt.Scanln(&new_db_loc)

		prettyfmt.Print("What is your new username?\n", prettyfmt.GREEN)
		prettyfmt.Scanln(&new_username)

		prettyfmt.Print("What is your new password?\n", prettyfmt.GREEN)
		prettyfmt.Scanln(&new_password)

		w = wallet.Initialize_Wallet(new_username, new_password, new_db_loc)

		if w == nil {
			prettyfmt.Print("Failed to create new wallet!\nExiting...\n", prettyfmt.RED)
			os.Exit(2)
		}
	} else {
		prettyfmt.Print("Found save file! Restoring Wallet...\n", prettyfmt.GREEN)

		prettyfmt.Printf("Found wallet for user: %s!\nEnter password to confirm identity:\n", prettyfmt.BLUE, w.Get_Username())

		var password string
		prettyfmt.Scanln(&password)

		w.Confirm_Password(password)

		if !w.Confirm_Password(password) {
			prettyfmt.Print("Wrong password! Goodbye!\n", prettyfmt.RED)
			os.Exit(3)
		}
	}

	prettyfmt.Printf("Welcome, %s!\n", prettyfmt.GREEN, w.Get_Username())

	BC = bc
	W = w
}
