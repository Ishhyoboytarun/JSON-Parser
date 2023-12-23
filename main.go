package main

import (
	"Json-Parser/json"
	"fmt"
)

func main() {
	p := json.NewParser("tests/basic.json")
	object := new(json.Person)
	err := p.Unmarshal(object)
	if err != nil {
		panic("invalid struct")
	}
	fmt.Println(object.Name)
	fmt.Println(object.GPA)
	fmt.Println(object.Age)
	fmt.Println(object.Company)
}
