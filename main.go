package main

import (
	"github.com/TheJ0lly/GoChain/wallet"
)

func main() {
	// bc := blockchain.Initialize_BlockChain("C:\\Users\\eusic\\Desktop\\LucrareLicenta\\Database")

	// // bc.Add_Data("Matei")
	// // bc.Add_Data("Florin")
	// // bc.Add_Data("Brinzea")
	// // bc.Add_Data("Test1")
	// // bc.Add_Data("Test2")

	// bc.Save_State()

	w := wallet.Initialize_Wallet("matei", "test", "C:\\Users\\eusic\\Desktop\\LucrareLicenta\\Assets")

	if w == nil {
		return
	}

	w.Add_Asset("Second_PDF", "C:\\Users\\eusic\\Downloads\\writing an INTERPRETER in go.pdf")

	w.Save_State()

}
