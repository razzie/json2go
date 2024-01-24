package main

import (
	"fmt"
	"go/format"
	"math"
	"sort"

	"github.com/goccy/go-reflect"
)

func GenerateStruct(data map[string]interface{}, structName string) (string, error) {
	structDef := fmt.Sprintf("type %s struct {\n", structName)
	structDef = generateFields(structDef, data)
	structDef += "}\n"

	formattedDef, err := format.Source([]byte(structDef))
	if err != nil {
		return structDef, err
	}

	return string(formattedDef), nil
}

func generateFields(structDef string, data map[string]interface{}) string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fieldName := ToCamelCase(key)
		switch v := data[key].(type) {
		case map[string]interface{}:
			structDef += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, generateStructFromMap(v), key)
		case []interface{}:
			if len(v) > 0 {
				// Handle non-empty slices
				switch e := v[0].(type) {
				case map[string]interface{}:
					// Nested slice of maps
					structDef += fmt.Sprintf("\t%s []%s `json:\"%s\"`\n", fieldName, generateStructFromMap(e), key)
				default:
					goType := getGoType(e)
					structDef += fmt.Sprintf("\t%s []%s `json:\"%s\"`\n", fieldName, goType, key)
				}
			} else {
				// Handle empty slices
				structDef += fmt.Sprintf("\t%s []interface{} `json:\"%s\"`\n", fieldName, key)
			}
		default:
			goType := getGoType(v)
			structDef += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, goType, key)
		}
	}
	return structDef
}

func generateStructFromMap(data map[string]interface{}) string {
	// Helper function to generate a struct from a map
	structDef := "struct {\n"
	structDef = generateFields(structDef, data)
	structDef += "}"
	return structDef
}

func getGoType(value interface{}) string {
	if value == nil {
		return "interface{}"
	}
	if value, ok := value.(float64); ok {
		return inferGoNumericType(value)
	}
	typ := reflect.TypeOf(value)
	tname := typ.Name()
	if len(tname) == 0 {
		return "interface{}"
	}
	return tname
}

func inferGoNumericType(value float64) string {
	if value != math.Floor(value) {
		return "float64"
	}

	if value >= math.MinInt && value <= math.MaxInt {
		return "int"
	}

	if value >= 0 {
		if value <= math.MaxUint64 {
			return "uint64"
		}
	} else {
		if value >= math.MinInt64 {
			return "int64"
		}
	}

	return "float64"
}
