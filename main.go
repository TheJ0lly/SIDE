package main

import (
	"fmt"
	"github.com/TheJ0lly/GoChain/cli"
)

//Change block data cap when ready
/*
Redo the TreeMatrix from TreeHash because:
	- if there are 2 or less MetaData, the tree will only create and store the root in the TreeMatrix.

Solution:
	All hashes need to be inside the tree including the tree root.

TreeMatrix[0]	H1, H2          or         H1, H1
TreeMatrix[1]    Root                       Root

If there are 3 hashes, for example:

TreeMatrix[0]  H1-H2  H3-H3
TreeMatrix[1]   H1  --  H2
TreeMatrix[2]      Root


And so on and so forth
*/

func main() {
	fv := cli.InitFlags()

	if fv == nil {
		fmt.Printf("Error: Program has failed to parse arguments!\nTry again.\n")
		return
	}

	cli.Execute(fv)

}
