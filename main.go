package main

import (
	parser2 "Json-Parser/parser"
	"fmt"
)

func main() {
	parser := parser2.NewParser("tests/basic.json")
	//fmt.Println(parser.Parse())

	p := new(Person)
	err := parser.Unmarshal(p)
	if err != nil {
		panic("Invalid struct")
	}

	fmt.Println(p.Name)
	fmt.Println(p.Age)
	fmt.Println(p.GPA)
}
