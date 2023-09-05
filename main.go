package main

import "github.com/TheJ0lly/GoChain/blockchain"

func main() {
	bc := blockchain.Initialize_BlockChain("C:\\Users\\eusic\\Desktop\\LucrareLicenta\\Database")

	bc.Add_Data("Matei")
	bc.Add_Data("Florin")
	bc.Add_Data("Brinzea")
	bc.Add_Data("Test1")
	bc.Add_Data("Test2")

	bc.Display_Block_Hashes()

	bc.Display_Data_Under_Block("B23B11B2BEB455284F5F16ECB82C772F773CBEE14837501065292B1FE5163FE6")

}
