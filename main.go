package main

import (
	"Json-Parser/parser"
	"fmt"
)

func main() {
	parser := parser.NewParser("file.json")
	fmt.Println(parser.Parse())
}
