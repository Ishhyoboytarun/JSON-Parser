Certainly! Below is a template for a README.md file for your JSON parsing and unmarshalling library in Go:

---

# JSON Parser

[![Go Report Card](https://goreportcard.com/badge/github.com/Ishhyoboytarun/JSON-Parser)](https://goreportcard.com/report/github.com/Ishhyoboytarun/JSON-Parser)

A lightweight JSON parsing and unmarshalling library for Go.

## Features

- **JSON Parsing**: Efficiently parse JSON data and convert it into a structured format.
- **Unmarshalling**: Unmarshal JSON into Go structs, making it easy to work with the data.

## Installation

```bash
go get -u github.com/Ishhyoboytarun/JSON-Parser
```

## Usage
```json
{
  "Name": "Tarunn Gusain",
  "Age": 25,
  "GPA": 3.8
}
```

### Unmarshalling

```go
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
```

## Contributions

Contributions are welcome! Feel free to open issues or pull requests for any improvements or bug fixes.


## Screenshots
<img width="1512" alt="Screenshot 2023-12-16 at 8 43 42 PM" src="https://github.com/Ishhyoboytarun/JSON-Parser/assets/36428256/f49ea925-6c9f-456b-9a51-144d5d670339">
<img width="1512" alt="Screenshot 2023-12-17 at 8 44 39 PM" src="https://github.com/Ishhyoboytarun/JSON-Parser/assets/36428256/9b6bb7e3-68c8-40da-ba5d-164ee3f1e37b">
