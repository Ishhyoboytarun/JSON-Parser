package json

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Object any

type Map map[string]Object

type String *string

type Boolean *bool

type Array []*Object

type Float *float64

type Integer *int64

type Parser struct {
	lock     sync.Mutex
	filePath string
}

var jsonObject Object

func NewParser(path string) *Parser {
	return &Parser{
		lock:     sync.Mutex{},
		filePath: path,
	}
}

func (p *Parser) fileExists() bool {
	_, err := os.Stat(p.filePath)
	return !os.IsNotExist(err)
}

func (p *Parser) parse() (Object, error) {

	p.lock.Lock()
	defer p.lock.Unlock()

	if !p.fileExists() {
		return nil, errors.New("file does not exist")
	}

	jsonString, err := p.createJsonString()
	if err != nil {
		return nil, err
	}

	err = p.chooseStrategy(jsonString)

	if err != nil {
		return nil, err
	}
	return jsonObject, nil
}

func (p *Parser) chooseStrategy(jsonString string) error {

	switch jsonString[0] {
	case '{':
		if jsonString[len(jsonString)-1] != '}' {
			return errors.New("invalid Json")
		}
		var err error
		jsonObject = make(Map)
		jsonObject, err = p.parseObject(jsonString)
		return err

	case '"':
		if jsonString[len(jsonString)-1] != '"' {
			return errors.New("invalid Json")
		}
		val, err := p.parseString(jsonString)
		if err != nil {
			return err
		}
		jsonObject = String(val)
		return nil

	case 't':
		if jsonString == "true" {
			value := true
			jsonObject = Boolean(&value)
			return nil
		}
		return errors.New("invalid Json")

	case 'f':
		if jsonString == "false" {
			value := false
			jsonObject = Boolean(&value)
			return nil
		}
		return errors.New("invalid Json")

	case 'n':
		if jsonString == "null" {
			jsonObject = nil
			return nil
		}
		return errors.New("invalid Json")

	case '[':
		if jsonString[len(jsonString)-1] != ']' {
			return errors.New("invalid Json")
		}
		jsonObject = make([]any, 0)
		return p.parseSlice(jsonString[1 : len(jsonString)-1])

	default:
		return p.handleIntFloat(jsonString)
	}
}

func (p *Parser) handleIntFloat(jsonString string) error {
	ascii := int(jsonString[0])
	if ascii >= 48 && ascii <= 57 {
		if strings.Contains(jsonString, ".") {
			val, err := p.parseFloat(jsonString)
			if err != nil {
				return err
			}
			jsonObject = Float(val)
			return nil
		} else {
			val, err := p.parseInteger(jsonString)
			if err != nil {
				return err
			}
			jsonObject = Integer(val)
			return nil
		}
	}
	return errors.New("invalid Json")
}

func (p *Parser) parseString(jsonString string) (*string, error) {
	if strings.Contains(jsonString[1:len(jsonString)-1], `"`) {
		return nil, errors.New("invalid Json")
	}
	cval := jsonString[1 : len(jsonString)-1]
	return &cval, nil
}

func (p *Parser) parseFloat(jsonString string) (*float64, error) {
	if val, err := strconv.ParseFloat(jsonString, 0); err == nil {
		return &val, nil
	}
	return nil, errors.New("invalid Json")
}

func (p *Parser) parseInteger(jsonString string) (*int64, error) {
	if val, err := strconv.Atoi(jsonString); err == nil {
		cval := int64(val)
		return &cval, nil
	}
	return nil, errors.New("invalid Json")
}

func (p *Parser) splitNestedSlice(jsonString string) ([]string, error) {

	result := make([]string, 0)
	stack := make([]int, 0)
	for i := 0; i < len(jsonString); i++ {
		val := jsonString[i]
		if val == '[' {
			stack = append(stack, i)
		} else if val == ']' {
			if len(stack) == 0 {
				return nil, errors.New("invalid Json")
			}
			stack = stack[:len(stack)-1]
		} else if val == ' ' {
			continue
		} else if val == '{' {
			index := i
			for index < len(jsonString) && (jsonString[index] != ',' || jsonString[index] != '}' || jsonString[index] != '"') {
				index++
			}
			result = append(result, jsonString[i:index])
			i = index
		} else if val == '}' {
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return nil, errors.New("invalid Json")
			}
			result = append(result, jsonString[stack[len(stack)-1]:i])
			stack = stack[:len(stack)-1]
		} else if val == ',' {
			continue
		} else {
			if len(stack) > 0 {
				continue
			}
			if val == '"' {
				i++
			}
			index := i
			for index < len(jsonString) && jsonString[index] != ',' && jsonString[index] != ']' && jsonString[index] != '"' {
				index++
			}
			result = append(result, jsonString[i:index])
			i = index
		}
	}
	return result, nil
}

func (p *Parser) parseSlice(jsonString string) error {

	array, err := p.splitNestedSlice(jsonString)
	if err != nil {
		return err
	}
	var result []any
	for _, val := range array {
		switch val[0] {
		case ' ':
			p.parseSlice(val[1:])
		case '{':
			if val[len(val)-1] != '}' {
				return errors.New("invalid Json")
			}
			jsonObject = make(Map)
			val, err := p.parseObject(val[1 : len(val)-1])
			if err != nil {
				return err
			}
			result = append(result, val)
		case '[':
			if val[len(val)-1] != ']' {
				return errors.New("invalid Json")
			}
			p.parseSlice(val[1 : len(val)-1])
		default:
			i := 0
			for val[i] == ' ' {
				i++
			}
			ascii := int(val[i])
			if ascii >= 48 && ascii <= 57 {
				if strings.Contains(val, ".") {
					val, err := p.parseFloat(val[i:])
					if err != nil {
						return err
					}
					result = append(result, *val)
				} else {
					val, err := p.parseInteger(val[i:])
					if err != nil {
						return err
					}
					result = append(result, *val)
				}
			} else if ascii >= 65 && ascii <= 90 || ascii >= 98 && ascii <= 123 {
				result = append(result, val)
			} else {
				return errors.New("invalid Json")
			}
		}
	}
	jsonObject = append(jsonObject.([]any), result...)
	return nil
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

func (p *Parser) parseObject(jsonString string) (Map, error) {

	if jsonString == "" {
		return jsonObject.(Map), nil
	}

	switch jsonString[0] {
	case '{':
		i := p.ignoreSpaces(1, jsonString)
		if jsonString[i] != '"' || jsonString[i] == '{' {
			return nil, errors.New("malformed Json")
		}
		return p.parseObject(jsonString[1:])

	case ' ':
		return p.parseObject(jsonString[1:])

	case '}':
		//Ignore the Paranthesis
		return p.parseObject(jsonString[1:])

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
			return nil, errors.New("malformed Json")
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
		for index < len(jsonString) && jsonString[index] != '"' && jsonString[index] != ',' && jsonString[index] != '}' {
			val += string(jsonString[index])
			index++
		}

		//convert the value according to the appropriate datatype
		jsonObject.(Map)[key] = p.convertValue(strings.TrimSpace(val))

		if index >= len(jsonString) {
			return jsonObject.(Map), nil
		}

		if index < len(jsonString) && jsonString[index] == '"' {
			index++
		}

		if index < len(jsonString) && jsonString[index] != ',' {
			index = p.ignoreSpaces(index, jsonString)
			if jsonString[index] != '}' {
				return nil, errors.New("invalid Json")
			}
		}
		return p.parseObject(jsonString[index+1:])

	default:
		return nil, errors.New("malformed Json")
	}
}

func (p *Parser) ignoreSpaces(index int, jsonString string) int {

	for index < len(jsonString) && jsonString[index] == ' ' {
		index++
	}
	return index
}

func (p *Parser) convertValue(val string) Object {

	if intValue, err := strconv.Atoi(val); err == nil {
		return intValue
	} else if floatValue, err := strconv.ParseFloat(val, 64); err == nil {
		return floatValue
	}
	return val
}
