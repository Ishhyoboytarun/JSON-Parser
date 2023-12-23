# Json Parser

[![Go Report Card](https://goreportcard.com/badge/github.com/tarunngusain08/Json-Parser)](https://goreportcard.com/report/github.com/tarunngusain08/Json-Parser)

A lightweight JSON parsing and unmarshalling library for Go.

## Features

- **JSON Parsing**: Parses JSON data and convert it into a structured format.
- **Unmarshalling**: Unmarshals JSON into Go structs, making it easy to work with the data.

## Installation

```bash
go get -u github.com/tarunngusain08/Json-Parser
```

## Usage
```json
{
  "Name": "Tarunn Gusain",
  "Age": 25,
  "GPA": 3.8,
  "Company": "Oracle"
}
```

### Unmarshalling

```go
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

```

## Contributions

Contributions are welcome! Feel free to open issues or pull requests for any improvements or bug fixes.


## Screenshots
<img width="1512" alt="Screenshot 2023-12-16 at 8 43 42 PM" src="https://github.com/Ishhyoboytarun/JSON-Parser/assets/36428256/f49ea925-6c9f-456b-9a51-144d5d670339">
<img width="1512" alt="Screenshot 2023-12-17 at 8 44 39 PM" src="https://github.com/Ishhyoboytarun/JSON-Parser/assets/36428256/9b6bb7e3-68c8-40da-ba5d-164ee3f1e37b">
<img width="1512" alt="Screenshot 2023-12-20 at 10 01 04 PM" src="https://github.com/tarunngusain08/Json-Parser/assets/36428256/bf35d075-7ab1-4050-b3c1-86cd85fea868">
<img width="1512" alt="Screenshot 2023-12-23 at 7 08 02 PM" src="https://github.com/tarunngusain08/Json-Parser/assets/36428256/e9b4b400-1f4b-40c6-992d-92cf8cd24321">
<img width="1512" alt="Screenshot 2023-12-23 at 7 08 22 PM" src="https://github.com/tarunngusain08/Json-Parser/assets/36428256/feb41ac9-70e9-404d-b11c-272bd144f43d">
<img width="1512" alt="Screenshot 2023-12-20 at 9 53 28 PM" src="https://github.com/tarunngusain08/Json-Parser/assets/36428256/74598143-320e-476a-967f-e0696a930432">

