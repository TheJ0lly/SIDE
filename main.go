package main

import (
	"fmt"
	"os"
)

func main() {
	// fv := cli.Init_Flags()

	// if fv == nil {
	// 	prettyfmt.ErrorF("Program has failed to parse arguments!\nTry again.\n")
	// 	return
	// }

	// cli.Execute(fv)

	_, err := os.Stat("sadasd")

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
