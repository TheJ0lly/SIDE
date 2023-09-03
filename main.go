package main

import "github.com/TheJ0lly/GoChain/blockchain"

func main() {
	bc := blockchain.InitializeBlockChain()

	bc.AddData("Matei")
	bc.AddData("Camelia")
	bc.AddData("Florin")
	bc.AddData("Test1")
	bc.AddData("Test2")

	bc.PrintBlockChain()

}
