package main

import (
	"fmt"
	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/blockchain"
)

//Change block data cap when ready

func main() {
	//fv := cli.InitFlags()
	//
	//if fv == nil {
	//	fmt.Printf("Error: Program has failed to parse arguments!\nTry again.\n")
	//	return
	//}
	//
	//cli.Execute(fv)

	bc, _ := blockchain.CreateNewBlockchain("C:\\Users\\eusic\\Desktop\\LicentaHelp\\Database")
	as := asset.CreateNewAsset("MateiPhoto", 0, nil)
	bc.AddData("Matei", as, "Ana")

	as = asset.CreateNewAsset("AnaPhoto", 0, nil)
	bc.AddData("Ana", as, "Matei")

	bie := blockchain.Block{}

	b := bie.Import("C:\\Users\\eusic\\Desktop\\LicentaHelp\\Database\\1C0FA3835C41477D826A8CE4823FD166F8BCFBBF02EA0BA3F8C9E0116BB85BCA")

	fmt.Printf("%T\n", b)
}
