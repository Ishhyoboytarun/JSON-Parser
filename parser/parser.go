package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"sync"
)

var jsonMap JsonObject

type Parser struct {
	lock     sync.Mutex
	filePath string
}

func NewParser(path string) *Parser {
	return &Parser{
		lock:     sync.Mutex{},
		filePath: path,
	}
}

func (p *Parser) Unmarshal(result interface{}) error {

	p.Parse()
	// Check if result is a pointer to a struct
	structType := reflect.TypeOf(result)
	if structType.Kind() != reflect.Ptr || structType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("result must be a pointer to a struct")
	}

	structValue := reflect.ValueOf(result).Elem()

	for i := 0; i < structType.Elem().NumField(); i++ {
		field := structType.Elem().Field(i)
		tag := field.Tag.Get("json")

		// Check if the field has a corresponding key in the hashmap
		if value, ok := jsonMap[tag]; ok {
			// Convert the value to the type of the struct field
			fieldValue := reflect.ValueOf(value).Convert(field.Type)

			// Set the struct field value
			structValue.Field(i).Set(fieldValue)
		}
	}

	return nil
}

func (p *Parser) Parse() (JsonObject, error) {
	return p.parse()
}

type JsonObject map[string]interface{}

func (p *Parser) fileExists() bool {
	_, err := os.Stat(p.filePath)
	return !os.IsNotExist(err)
}

func (p *Parser) parse() (JsonObject, error) {

	p.lock.Lock()
	defer p.lock.Unlock()

	if !p.fileExists() {
		return nil, errors.New("file does not exist")
	}

	jsonString, err := p.createJsonString()
	if err != nil {
		return nil, err
	}

	jsonMap = JsonObject{}
	err = p.createJsonObject(jsonString)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}

func (p *Parser) createJsonString() (string, error) {

	file, err := os.Open(p.filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	jsonString := ""

	for scanner.Scan() {
		jsonString += scanner.Text()
	}
	return jsonString, nil
}

func (p *Parser) createJsonObject(jsonString string) error {

	if jsonString == "" {
		return nil
	}

	switch jsonString[0] {
	case '{':
		//Ignore the Paranthesis
		return p.createJsonObject(jsonString[1:])

	case ' ':
		return p.createJsonObject(jsonString[1:])

	case '}':
		//Ignore the Paranthesis
		return p.createJsonObject(jsonString[1:])

	case '"':
		index := 1

		// Fetch the token
		key := ""
		for jsonString[index] != '"' {
			key += string(jsonString[index])
			index++
		}
		index++

		// Ignore the spaces
		index = p.ignoreSpaces(index, jsonString)

		if jsonString[index] != ':' {
			return errors.New("malformed Json")
		}

		if jsonString[index] == ':' {
			index++
		}

		// Ignore the spaces
		index = p.ignoreSpaces(index, jsonString)

		val := ""
		if jsonString[index] == '"' {
			index++
		}
		for jsonString[index] != '"' && jsonString[index] != ',' && jsonString[index] != '}' {
			val += string(jsonString[index])
			index++
		}

		//convert the value according to the appropriate datatype
		jsonMap[key] = p.getValue(val)

		if jsonString[index] == '"' {
			index++
		}

		if jsonString[index] != ',' {
			index = p.ignoreSpaces(index, jsonString)
			if jsonString[index] != '}' {
				return errors.New("malformed Json")
			}
		}
		return p.createJsonObject(jsonString[index+1:])

	default:
		return errors.New("malformed Json")
	}

	return nil
}

func (p *Parser) ignoreSpaces(index int, jsonString string) int {

	for index < len(jsonString) && jsonString[index] == ' ' {
		index++
	}
	return index
}

func (p *Parser) getValue(val string) interface{} {

	if intValue, err := strconv.Atoi(val); err == nil {
		return intValue
	} else if floatValue, err := strconv.ParseFloat(val, 64); err == nil {
		return floatValue
	}
	return val
}
