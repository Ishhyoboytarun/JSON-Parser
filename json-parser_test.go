package main

import (
	"Json-Parser/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonParser(t *testing.T) {
	t.Run("Test array-basic", func(t *testing.T) {
		expectedResponse := []interface{}{"Tarunn", "Gusain"}
		p := json.NewParser("tests/array-basic.json")
		object := make([]any, 0)
		err := p.Unmarshal(&object)
		assert.Nil(t, err)
		assert.Equal(t, "Gusain", object[1])
		assert.Equal(t, "Tarunn", object[0])
		assert.Equal(t, expectedResponse, object)
	})

	t.Run("Test basic", func(t *testing.T) {
		expectedResponse := &json.Person{Name: "Tarunn Gusain", Age: 25, GPA: 3.8, Company: "Oracle"}
		p := json.NewParser("tests/basic.json")
		object := new(json.Person)
		err := p.Unmarshal(object)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, object)
	})

	t.Run("Test basic-unformatted", func(t *testing.T) {
		expectedResponse := &json.Person{Name: "Tarunn Gusain", Age: 25, GPA: 3.8, Company: "Oracle"}
		p := json.NewParser("tests/basic.json")
		object := new(json.Person)
		err := p.Unmarshal(object)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, object)
	})

	t.Run("Test boolean-false", func(t *testing.T) {
		p := json.NewParser("tests/boolean-false.json")
		object := new(bool)
		err := p.Unmarshal(object)
		assert.Nil(t, err)
		assert.False(t, *object)
	})

	t.Run("Test boolean-true", func(t *testing.T) {
		p := json.NewParser("tests/boolean-true.json")
		object := new(bool)
		err := p.Unmarshal(object)
		assert.Nil(t, err)
		assert.True(t, *object)
	})

	t.Run("Test float", func(t *testing.T) {
		expectedResponse := 3.8
		p := json.NewParser("tests/float.json")
		object := new(float64)
		err := p.Unmarshal(object)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, *object)
	})

	t.Run("Test integer", func(t *testing.T) {
		expectedResponse := 24
		p := json.NewParser("tests/integer.json")
		object := new(int64)
		err := p.Unmarshal(object)
		assert.Nil(t, err)
		assert.Equal(t, int64(expectedResponse), *object)
	})

	t.Run("Test nested", func(t *testing.T) {
		expectedResponse := []interface{}{}
		p := json.NewParser("tests/nested.json")
		object := make([]interface{}, 0)
		err := p.Unmarshal(&object)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, object)
	})

	t.Run("Test null", func(t *testing.T) {
		p := json.NewParser("tests/null.json")
		object := new(int)
		err := p.Unmarshal(object)
		assert.Nil(t, err)
		assert.Equal(t, 0, *object)
	})

	t.Run("Test string", func(t *testing.T) {
		expectedResponse := "Hello, My name is Tarunn Gusain"
		p := json.NewParser("tests/string.json")
		object := new(string)
		err := p.Unmarshal(object)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, *object)
	})
}
