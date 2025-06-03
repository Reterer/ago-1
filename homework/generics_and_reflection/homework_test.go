package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	v := reflect.ValueOf(person)
	t := reflect.TypeOf(person)

	var sb strings.Builder

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get("properties")
		if tag == "" {
			continue
		}

		var value string
		name, omitempty := strings.CutSuffix(tag, ",omitempty")

		switch fieldType.Type.Kind() {
		case reflect.Int:
			if omitempty && fieldValue.Int() == 0 {
				continue
			}
			value = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Bool:
			if omitempty && fieldValue.Bool() == false {
				continue
			}
			value = strconv.FormatBool(fieldValue.Bool())
		case reflect.String:
			if omitempty && fieldValue.String() == "" {
				continue
			}
			value = fieldValue.String()
		}

		if sb.Len() > 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(name)
		sb.WriteRune('=')
		sb.WriteString(value)
	}

	return sb.String()
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
