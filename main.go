package main

import (
	"github.com/TheJ0lly/GoChain/cli"
)

func main() {
	cli.Login_Or_Signup()

	cli.BC.Save_State()
	cli.W.Save_State()

}
