package main

import (
	"os"

	"github.com/jeremi-traverse/monkey/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
