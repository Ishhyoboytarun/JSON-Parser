package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
)

type Parser struct {
	sync.Mutex
	FilePath string
}

func NewParser(path string) *Parser {
	return &Parser{
		Mutex:    sync.Mutex{},
		FilePath: path,
	}
}

func (p *Parser) Parse() {
	p.parse()
}

func (p *Parser) fileExists() bool {
	_, err := os.Stat(p.FilePath)
	return !os.IsNotExist(err)
}

func (p *Parser) parse() (interface{}, error) {

	p.Lock()
	defer p.Unlock()

	if !p.fileExists() {
		return nil, errors.New("file does not exist")
	}

	jsonString, err := p.createJsonString()
	if err != nil {
		return nil, err
	}

	jsonMap, err := p.createJsonObject(jsonString, new(interface{}))
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}

func (p *Parser) createJsonString() (string, error) {

	file, err := os.Open(p.FilePath)
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

func (p *Parser) createJsonObject(jsonString string, response interface{}) (interface{}, error) {

	var resp interface{}
	var err error
	var index int
	for index = 0; jsonString[index] == ' '; index++ {
		val := jsonString[index]
		fmt.Println(val, jsonString[0])
	}
	switch jsonString[index] {
	case '"':
		str := ""
		for jsonString[index]!='"'{
			str += string(jsonString[index])
			index++
		}
	case '{':
		resp, err = p.createJsonObject(jsonString[1:], make(map[string]interface{}))
		if err != nil {
			return nil, err
		}
	case '[':
		resp, err = p.createJsonObject(jsonString[1:], make([]interface{}, 0))
		if err != nil {
			return nil, err
		}
	case '}':
		return response, nil
	case ']':
		return response, nil
	default:
		switch reflect.TypeOf(response).Kind() {
		case reflect.Map:
			resp, err = p.createJsonMap(jsonString)
			if err != nil {
				return nil, err
			}
		case reflect.Slice:
			resp, err = p.createJsonArray(jsonString)
			if err != nil {
				return nil, err
			}
		}
	}
	return resp, nil
}

func (p *Parser) createJsonMap(jsonString string) (map[string]interface{}, error) {

	response := make(map[string]interface{})
	for i := 0; jsonString[i] != '}'; i++ {
		if jsonString[i] == ' ' || jsonString[i] == '"' {
			continue
		} else {
			key := ""
			for i < len(jsonString) && jsonString[i] != '"' {
				key += string(jsonString[i])
				i++
			}
			var value interface{}
			var err error
			switch jsonString[0] {
			case '{':
				value, err = p.createJsonObject(jsonString[i:], make(map[string]interface{}))
				if err != nil {
					return nil, err
				}
			case '[':
				value, err = p.createJsonObject(jsonString[i:], make([]interface{}, 0))
				if err != nil {
					return nil, err
				}
			default:
				val := ""
				for i < len(jsonString) && jsonString[i] != ',' {
					val += string(jsonString[i])
					i++
				}
				value = val
			}
			response[key] = value
		}
	}
	return response, nil
}

func (p *Parser) createJsonArray(jsonString string) ([]interface{}, error) {

	response := make([]interface{}, 0)
	for i := 0; jsonString[i] != ']'; i++ {
		if jsonString[i] == ' ' || jsonString[i] == '"' {
			continue
		} else {
			var value interface{}
			var err error
			switch jsonString[0] {
			case '{':
				value, err = p.createJsonObject(jsonString[i:], make(map[string]interface{}))
				if err != nil {
					return nil, err
				}
			case '[':
				value, err = p.createJsonObject(jsonString[i:], make([]interface{}, 0))
				if err != nil {
					return nil, err
				}
			default:
				val := ""
				for i < len(jsonString) && jsonString[i] != ',' {
					val += string(jsonString[i])
					i++
				}
				value = val
			}
			response = append(response, value)
		}
	}
	return response, nil
}
