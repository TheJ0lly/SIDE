package main

import "github.com/TheJ0lly/GoChain/cli"

func main() {
	bc, w := cli.Login_Or_Signup()

	bc.Save_State()
	w.Save_State()
}
