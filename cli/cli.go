package cli

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strconv"

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

func display_title() {
	Clear_Screen()

	prettyfmt.Print("#########    #####   ######  ######       #######         #######\n")
	prettyfmt.Print("#########  ##     ##   ####  ####         ###   ###      ##      \n")
	prettyfmt.Print("   ###     ##     ##    ###  ###          ###    ###    ##       \n")
	prettyfmt.Print("   ###     ##     ##     ######           #######      ##        \n")
	prettyfmt.Print("   ###     ##     ##      ####            ###    ###   ##        \n")
	prettyfmt.Print("   ###     ##     ##      ####            ###     ###   ##       \n")
	prettyfmt.Print("   ###     ##     ##      ####            ###    ###     ##      \n")
	prettyfmt.Print("   ###       #####        ####            #######         #######\n")

	prettyfmt.Print("\n\n\n")
}

func Clear_Screen() {
	if runtime.GOOS == "windows" {
		clear_screen_err := exec.Command("powershell", "clear").Run()

		if clear_screen_err != nil {
			prettyfmt.ErrorF("Error in displaying the menu: %s\n", clear_screen_err.Error())
			return
		}
	} else if runtime.GOOS == "linux" {
		clear_screen_err := exec.Command("clear").Run()

		if clear_screen_err != nil {
			prettyfmt.ErrorF("Error in displaying the menu: %s\n", clear_screen_err.Error())
			return
		}
	}
}

func ScanChoice() int {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()

	val, err := strconv.Atoi(s.Text())

	prettyfmt.CPrintf("%s\n", prettyfmt.BLUE, s.Text())

	if err != nil && s.Text() != "" {
		prettyfmt.ErrorF("Error in parsing the choice: %s\n", err.Error())
		return -1
	}

	return val

}

func Display_Main_Menu() {
	display_title()
	prettyfmt.Print("   Main Menu\n\n")
	prettyfmt.Print("1. Add Asset\n")
	prettyfmt.Print("2. Remove Asset\n")
	prettyfmt.Print("3. View Assets\n")
	prettyfmt.Print("4. Add to Blockchain (Just for testing as of now)\n")
	prettyfmt.Print("5. View Blockchain (Not working now)\n")
	prettyfmt.Print("6. Save\n")
	prettyfmt.Print("7. Exit\n")

}

func Add_Asset(w *wallet.Wallet) {
	Clear_Screen()
	var asset_name string
	var asset_init_loc string

	prettyfmt.Print("===== ADD ASSET =====\n\n")
	prettyfmt.Print("What is the new name of the asset you want to add?\n->")
	prettyfmt.Scanln(&asset_name)

	prettyfmt.Print("\n")

	prettyfmt.Print("Location of the file on your machine?\n->")
	prettyfmt.Scanln(&asset_init_loc)

	w.Add_Asset(asset_name, asset_init_loc)

	prettyfmt.Print("Press enter to go back to the main menu...\n")
	ScanChoice()
}

func Remove_Asset(w *wallet.Wallet) {
	Clear_Screen()
	var asset_name string

	prettyfmt.Print("===== REMOVE ASSET =====\n\n")

	assets := w.Get_All_Assets()

	prettyfmt.Print("Your assets:\n")

	for _, a := range assets {
		prettyfmt.Printf("  ---%s\n", a.Get_Name())
	}

	prettyfmt.Print("What asset do you want to remove?\n->")
	prettyfmt.Scanln(&asset_name)

	w.Remove_Asset(asset_name)

	prettyfmt.Print("Press enter to go back to the main menu...\n")
	ScanChoice()
}

func View_Assets(w *wallet.Wallet) {
	Clear_Screen()
	assets := w.Get_All_Assets()

	prettyfmt.Print("Your assets:\n")

	for _, a := range assets {
		prettyfmt.Printf("  -%s\n", a.Get_Name())
	}

	prettyfmt.Print("Press enter to go back to the main menu...\n")
	ScanChoice()
}

func Add_To_Blockchain_Test(w *wallet.Wallet, bc *blockchain.BlockChain) {
	Clear_Screen()
	var asset_name string
	var dest_user string

	prettyfmt.Print("===== ADD TO BLOCKCHAIN =====\n\n")

	prettyfmt.Print("What asset do you want to transaction?\n->")
	prettyfmt.Scanln(&asset_name)

	a := w.Get_Asset(asset_name)

	if a == nil {
		prettyfmt.ErrorF("There is no asset named \"%s\" in your wallet!\n", asset_name)
		prettyfmt.ErrorF("Failed to add \"%s\" to blockchain!\n", asset_name)
		prettyfmt.Print("Press enter to go back to the main menu...\n")
		ScanChoice()
		return
	}

	prettyfmt.Print("Who is the receiving user?\n->")
	prettyfmt.Scanln(&dest_user)

	bc.Add_Data_Test(w.Get_Username(), a, dest_user)

	prettyfmt.Print("Press enter to go back to the main menu...\n")
	ScanChoice()
}

func View_Blockchain_Test(bc *blockchain.BlockChain) {
	Clear_Screen()
	bc.View_Blockchain()

	prettyfmt.Print("Press enter to go back to the main menu...\n")
	ScanChoice()
}
