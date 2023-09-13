package main

import (
	"os"
	"time"

	"github.com/TheJ0lly/GoChain/cli"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

func main() {
	bc, w := cli.Login_Or_Signup()

	cli.Display_Main_Menu()

	for {
		choice := cli.ScanChoice()

		switch choice {
		case 1: //Add an asset in the wallet.
			cli.Add_Asset(w)
		case 2: //Remove an asset from the wallet.
			cli.Remove_Asset(w)
		case 3: //View assets
			cli.View_Assets(w)
		case 4: //Add to Blockchain(Testing)
			cli.Add_To_Blockchain_Test(w, bc)
		case 5: //View Blockchain(Testing)
			cli.View_Blockchain_Test(bc)
		case 6: //Save state of the Wallet and the Blockchain.
			err := bc.Save_State()

			if err != nil {
				generalerrors.HandleError(err)
				break
			}

			err = w.Save_State()

			if err != nil {
				generalerrors.HandleError(err)
				break
			}

			prettyfmt.CPrint("Blockchain and Wallet state have been saved! Continuing in 3 seconds...",
				prettyfmt.GREEN)
			time.Sleep(time.Second * 3)
		case 7: //Exit the application.
			cli.Clear_Screen()
			os.Exit(0)
		}

		cli.Display_Main_Menu()
	}

}
