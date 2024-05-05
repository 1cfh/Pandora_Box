package main

import (
	"Pandora_Box/repl"
	"fmt"
	"os"
)

const banner = `	██████╗  █████╗ ███╗   ██╗██████╗  ██████╗ ██████╗  █████╗         ██████╗  ██████╗ ██╗  ██╗
	██╔══██╗██╔══██╗████╗  ██║██╔══██╗██╔═══██╗██╔══██╗██╔══██╗        ██╔══██╗██╔═══██╗╚██╗██╔╝
	██████╔╝███████║██╔██╗ ██║██║  ██║██║   ██║██████╔╝███████║        ██████╔╝██║   ██║ ╚███╔╝ 
	██╔═══╝ ██╔══██║██║╚██╗██║██║  ██║██║   ██║██╔══██╗██╔══██║        ██╔══██╗██║   ██║ ██╔██╗ 
	██║     ██║  ██║██║ ╚████║██████╔╝╚██████╔╝██║  ██║██║  ██║███████╗██████╔╝╚██████╔╝██╔╝ ██╗
	╚═╝     ╚═╝  ╚═╝╚═╝  ╚═══╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝`

func init() {
	// banner
	fmt.Println(banner)
}

func main() {

	// user, err := user2.Current()
	// if err != nil {
	//	panic(err)
	//}

	// fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)

	// fmt.Printf("Feel free to type in commands\n")

	repl.Start(os.Stdin, os.Stdout)

}
