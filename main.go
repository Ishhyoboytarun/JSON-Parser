package main

import (
	"Json-Parser/parser"
	"fmt"
	"strings"
)

func unmarshalJSON(json string) (map[string]interface{}, error) {
	stack := make([]map[string]interface{}, 0)
	current := make(map[string]interface{})
	key := ""

	for i := 0; i < len(json); i++ {
		switch json[i] {
		case '{':
			newMap := make(map[string]interface{})
			if key != "" {
				current[key] = newMap
			} else {
				stack = append(stack, current)
				current = newMap
			}
			key = ""
		case '}':
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				parent[key] = current
				current = parent
			}
			key = ""
		case '"':
			i++
			start := i
			for i < len(json) && json[i] != '"' {
				i++
			}
			if i >= len(json) {
				return nil, fmt.Errorf("malformed JSON")
			}
			key = json[start:i]
		case ',':
			key = ""
		case ':':
			// Ignore colons in strings
			for i+1 < len(json) && json[i+1] == ' ' {
				i++
			}
			if i+1 < len(json) && json[i+1] != '"' {
				start := i
				for i < len(json) && json[i] != ',' && json[i] != '}' && json[i] != ']' {
					i++
				}
				value := json[start:i]
				
				return nil, fmt.Errorf("malformed JSON")
			}
		case ' ', '\t', '\n', '\r':
			// Ignore whitespaces
		default:
			start := i
			for i < len(json) && json[i] != ',' && json[i] != '}' && json[i] != ']' {
				i++
			}
			value := json[start:i]
			if key != "" {
				current[key] = strings.TrimSpace(value)
				key = ""
			}
		}
	}

	return current, nil
}

func main() {
	//jsonString := `{"name": "John", "age": 30, "city": "New York", "scores": [90, 85, 92]}`
	parser := parser.NewParser("file.json")
	jsonString, _ := parser.CreateJsonString()
	result, err := unmarshalJSON(jsonString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("%#v\n", result)
}
