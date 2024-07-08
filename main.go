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

const description = `
	Hello! This is the Pandora_Box. Wish you happy! :)
`

func init() {
	// banner
	fmt.Println(banner)
	// description
	fmt.Println(description)
}

func main() {
	repl.Start(os.Stdin, os.Stdout)
	// doSomething(123)
}

func doSomething(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Println(v)
		fmt.Printf("整数: %d\n", v)
	case string:
		fmt.Println(v)
		fmt.Printf("字符串: %s\n", v)
	case bool:
		fmt.Println(v)
		fmt.Printf("布尔值: %t\n", v)
	default:
		fmt.Println(v)
		fmt.Printf("未知类型: %T\n", v)
	}
}
