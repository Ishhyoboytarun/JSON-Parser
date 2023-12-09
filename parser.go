package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
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

func (p *Parser) parse() (map[string]interface{}, error) {
	if !p.fileExists() {
		return nil, errors.New("file does not exist")
	}

	json, err := p.parseJson()
	if err != nil {
		return nil, err
	}
	return json, nil
}

func (p *Parser) parseJson() (map[string]interface{}, error) {
	p.Lock()
	defer p.Unlock()
	file, err := os.Open(p.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	response := make(map[string]interface{}, 1)

	scanner := bufio.NewScanner(file)
	skip := false
	for scanner.Scan() {
		if !skip {
			skip = true
			continue
		}
		line := scanner.Text()
		keyValues := strings.Split(line, ",")
		for i := 0; i < len(keyValues); i++ {
			key, value, err := p.extractKeyValue(line)
			if err != nil {
				return nil, err
			}
			if key == nil {
				continue
			}
			response[*key] = value
		}
	}
	return response, nil
}

func (p *Parser) extractKeyValue(line string) (*string, interface{}, error) {
	if line[0] == '{' || line[0] == '}' {
		return nil, nil, nil
	}

	key, err := p.extractKey(line)
	if err != nil {
		return nil, nil, err
	}

	value, err := p.extractValue(line)
	if err != nil {
		return nil, nil, err
	}

	return key, value, nil
}

func (p *Parser) extractKey(line string) (*string, error) {
	if !strings.Contains(line, ":") {
		return nil, errors.New("invalid JSON")
	}
}

func (p *Parser) extractValue(line string) (interface{}, error) {
	return nil, nil
}
