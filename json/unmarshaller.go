package json

import (
	"errors"
	"reflect"
)

func (p *Parser) Unmarshal(result interface{}) error {

	p.parse()
	jsonObjectType := reflect.TypeOf(jsonObject)
	if jsonObjectType == nil {
		result = nil
		return nil
	}
	switch jsonObjectType.Kind() {
	case reflect.Map:
		return p.unmarshalMap(result)
	case reflect.Slice:
		return p.unmarshalSlice(result)
	default:
		switch jsonObjectType.Elem().Kind() {
		case reflect.String:
			return p.unmarshalString(result)
		case reflect.Int64:
			return p.unmarshalInteger(result)
		case reflect.Float64:
			return p.unmarshalFloat(result)
		case reflect.Bool:
			return p.unmarshalBoolean(result)
		default:
			return errors.New("invalid json")
		}
	}
}

func (p *Parser) unmarshalMap(result interface{}) error {

	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Ptr || resultType.Elem().Kind() != reflect.Struct {
		return errors.New("result must be a pointer to a struct")
	}

	structValue := reflect.ValueOf(result).Elem()
	for i := 0; i < resultType.Elem().NumField(); i++ {
		field := resultType.Elem().Field(i)
		tag := field.Tag.Get("json")
		// Check if the field has a corresponding key in the hashmap
		if value, ok := jsonObject.(Map)[tag]; ok {
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
	value := reflect.ValueOf(jsonObject).Convert(resultType).Elem()
	reflect.ValueOf(result).Elem().Set(value)
	return nil
}

func (p *Parser) unmarshalInteger(result interface{}) error {
	// Check if result is a pointer to a struct
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Ptr || resultType.Elem().Kind() != reflect.Int64 {
		return errors.New("result must be a pointer to a int64")
	}
	value := reflect.ValueOf(jsonObject).Convert(resultType).Elem()
	reflect.ValueOf(result).Elem().Set(value)
	return nil
}

func (p *Parser) unmarshalFloat(result interface{}) error {
	// Check if result is a pointer to a struct
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Ptr || resultType.Elem().Kind() != reflect.Float64 {
		return errors.New("result must be a pointer to a float32 or float64")
	}
	value := reflect.ValueOf(jsonObject).Convert(resultType).Elem()
	reflect.ValueOf(result).Elem().Set(value)
	return nil
}

func (p *Parser) unmarshalBoolean(result interface{}) error {
	// Check if result is a pointer to a struct
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Ptr || resultType.Elem().Kind() != reflect.Bool {
		return errors.New("result must be a pointer to a bool")
	}
	value := reflect.ValueOf(jsonObject).Convert(resultType).Elem()
	reflect.ValueOf(result).Elem().Set(value)
	return nil
}

func (p *Parser) unmarshalSlice(result any) error {
	// Check if result is a pointer to a struct
	resultType := reflect.TypeOf(result)
	resultValue := reflect.ValueOf(result).Elem()
	if resultType.Kind() != reflect.Ptr || resultType.Elem().Kind() != reflect.Slice || resultValue.Type().Elem().Kind() != reflect.Interface {
		return errors.New("result must be a pointer to an slice of interface")
	}
	jsonObjectValue := reflect.ValueOf(jsonObject)
	for i := 0; i < jsonObjectValue.Len(); i++ {
		resultValue.Set(reflect.Append(resultValue, jsonObjectValue.Index(i)))
	}
	return nil
}
