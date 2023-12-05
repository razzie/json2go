package main

import (
	"fmt"
	"go/format"
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
	for key, value := range data {
		fieldName := ToCamelCase(key)
		switch v := value.(type) {
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
			goType := getGoType(value)
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
	switch value.(type) {
	case float64:
		return "float64"
	case float32:
		return "float32"
	case int:
		return "int"
	case int32:
		return "int32"
	case int64:
		return "int64"
	case bool:
		return "bool"
	case string:
		return "string"
	default:
		return "interface{}"
	}
}
