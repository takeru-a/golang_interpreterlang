package main

import (
	"fmt"
	"github.com/takeru-a/golang_interpreterlang/repl"
	"os"
)

func main() {

	fmt.Printf("Hello! This is the Aquamarine programming language!\n")
	fmt.Printf("\n")
	repl.Start(os.Stdin, os.Stdout)
}
