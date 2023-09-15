package main

import (
	"github.com/TheJ0lly/GoChain/cli"
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

func main() {
	fv := cli.Init_Flags()

	if fv == nil {
		prettyfmt.ErrorF("Program has failed to parse arguments!\nTry again.\n")
		return
	}

	cli.Execute(fv)
}
