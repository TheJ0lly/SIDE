package main

import "github.com/TheJ0lly/GoChain/blockchain"

func main() {
	bc := blockchain.Initialize_BlockChain()

	bc.Add_Data("Matei")
	bc.Add_Data("Camelia")
	bc.Add_Data("Florin")
	bc.Add_Data("Test1")
	bc.Add_Data("Test2")

	bc.Print_BlockChain()

}
