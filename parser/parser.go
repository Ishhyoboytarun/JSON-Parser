package parser

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type JsonObject interface{}

type JsonMap map[string]interface{}

type JsonString *string

type JsonBoolean bool

type JsonArray []interface{}

type JsonFloat *float64

type JsonInteger *int64

type Parser struct {
	lock     sync.Mutex
	filePath string
}

var jsonObject JsonObject

func NewParser(path string) *Parser {
	return &Parser{
		lock:     sync.Mutex{},
		filePath: path,
	}
}

func (p *Parser) Unmarshal(result interface{}) error {

	p.Parse()

	jsonObjectType := reflect.TypeOf(jsonObject)
	switch jsonObjectType.Elem().Kind() {
	case reflect.Map:
		return p.unmarshalMap(result)
	case reflect.String:
		return p.unmarshalString(result)
	}

	return nil
}

func (p *Parser) unmarshalMap(result interface{}) error {

	// Check if result is a pointer to a struct
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Ptr || resultType.Elem().Kind() != reflect.Struct {
		return errors.New("result must be a pointer to a struct")
	}

	structValue := reflect.ValueOf(result).Elem()
	for i := 0; i < resultType.Elem().NumField(); i++ {
		field := resultType.Elem().Field(i)
		tag := field.Tag.Get("json")
		// Check if the field has a corresponding key in the hashmap
		if value, ok := jsonObject.(JsonMap)[tag]; ok {
			// Convert the value to the type of the struct field
			fieldValue := reflect.ValueOf(value).Convert(field.Type)
			// Set the struct field value
			structValue.Field(i).Set(fieldValue)
		}
	}
	return nil
}

func (p *Parser) unmarshalString(result interface{}) error {
	// Check if result is a pointer to a struct
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Ptr || resultType.Elem().Kind() != reflect.String {
		return errors.New("result must be a pointer to a string")
	}
	value := reflect.ValueOf(jsonObject).Convert(resultType.Elem().Field(0).Type)
	reflect.ValueOf(result).Elem().Field(0).Set(value)
	return nil
}

func (p *Parser) Parse() (JsonObject, error) {
	return p.parse()
}

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

	err = p.chooseStrategy(jsonString)

	if err != nil {
		return nil, err
	}
	return jsonObject, nil
}

func (p *Parser) chooseStrategy(jsonString string) error {

	switch jsonString[0] {
	case '{':
		jsonObject = new(JsonMap)
		return p.createJsonObject(jsonString)

	case '"':
		val, err := p.convertJsonStringtoStringValue(jsonString)
		if err != nil {
			return err
		}
		jsonObject = JsonString(val)
		return nil

	case 't':
		if jsonString == "true" {
			jsonObject = JsonBoolean(true)
			return nil
		}
		return errors.New("invalid Json")

	case 'f':
		if jsonString == "false" {
			jsonObject = JsonBoolean(false)
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
		jsonObject = new(JsonArray)

	default:
		ascii := int(jsonString[0])
		if ascii >= 48 && ascii <= 57 {
			if strings.Contains(jsonString, ".") {
				val, err := p.convertJsonStringToFloatValue(jsonString)
				if err != nil {
					return err
				}
				jsonObject = JsonFloat(val)
				return nil
			} else {
				val, err := p.convertJsonStringToIntegerValue(jsonString)
				if err != nil {
					return err
				}
				jsonObject = JsonInteger(val)
				return nil
			}
		}
		return errors.New("invalid Json")
	}
	return nil
}

func (p *Parser) convertJsonStringtoStringValue(jsonString string) (*string, error) {
	if strings.Contains(jsonString[1:len(jsonString)-1], `"`) {
		return nil, errors.New("invalid Json")
	}
	cval := jsonString[1 : len(jsonString)-1]
	return &cval, nil
}

func (p *Parser) convertJsonStringToFloatValue(jsonString string) (*float64, error) {
	if val, err := strconv.ParseFloat(jsonString, 0); err == nil {
		return &val, nil
	}
	return nil, errors.New("invalid Json")
}

func (p *Parser) convertJsonStringToIntegerValue(jsonString string) (*int64, error) {
	if val, err := strconv.Atoi(jsonString); err == nil {
		cval := int64(val)
		return &cval, nil
	}
	return nil, errors.New("invalid Json")
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
		if jsonString[1] != '"' || jsonString[1] != '{' {
			return errors.New("malformed Json")
		}
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
		jsonObject.(JsonMap)[key] = p.getValue(val)

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
