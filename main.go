package main

import (
	"fmt"
	"github.com/TheJ0lly/GoChain/cli"
)

//Change block data cap when ready

func main() {
	fv := cli.InitFlags()

	if fv == nil {
		fmt.Printf("Error: Program has failed to parse arguments!\nTry again.\n")
		return
	}

	cli.Execute(fv)

}
