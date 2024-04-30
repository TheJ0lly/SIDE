package main

import (
	"github.com/TheJ0lly/GoChain/cli"
	"log"
)

//Change block data cap when ready

func main() {
	fv := cli.InitFlags()

	if fv == nil {
		log.Printf("ERROR: program has failed to parse arguments!\nTry again.\n")
		return
	}

	cli.Execute(fv)

}
