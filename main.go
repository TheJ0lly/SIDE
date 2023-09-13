package main

import (
	"os"

	"github.com/TheJ0lly/GoChain/cli"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

func main() {
	bc, w := cli.Login_Or_Signup()

	cli.Display_Main_Menu()

	isRunning := true

	for isRunning {
		choice := cli.ScanChoice()

		switch choice {
		case 1:
			//Add an asset in the wallet.
			cli.Add_Asset(w)

			prettyfmt.Print("Press enter to go back to the main menu...\n")
			cli.ScanChoice()
		case 2:
			//Remove an asset from the wallet.
		case 3:
			err := bc.Save_State()

			if err != nil {
				generalerrors.HandleError(err)
			}

			err = w.Save_State()

			if err != nil {
				generalerrors.HandleError(err)
			}

		case 4:
			cli.Clear_Screen()
			os.Exit(0)
		}

		cli.Display_Main_Menu()
	}

}
