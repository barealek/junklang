package main

import (
	"fmt"
	"os"

	"github.com/barealek/junklang"
)

func main() {

	args := os.Args
	if len(args) <= 1 {
		fmt.Println("Mangler en kommando: junk [run]")
		return
	}

	cmd := args[1]
	switch cmd {
	case "run":
		if len(args) != 3 {
			fmt.Println("Mangler en fil: junk run [fil]")
			return
		}
		script, err := readFile(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		evaluateScript(script)
	}
}

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	return string(b), err
}

func evaluateScript(s string) {

	global := junklang.NewScope(nil)

	lexer := junklang.NewLexer(s)
	tokens := lexer.Tokenize()

	parser := junklang.NewParser(tokens)
	nodes := parser.Parse()

	for _, node := range nodes {
		if node != nil {
			node.Call(global)
		}
	}
}
