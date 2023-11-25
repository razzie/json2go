package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"strings"
)

func toCamelCase(name string) string {
	parts := strings.Split(name, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

func generateStruct(jsonStr string, structName string) (string, error) {
	var data map[string]interface{}

	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}

	structDef := fmt.Sprintf("type %s struct {\n", structName)
	structDef = generateFields(structDef, data, "")
	structDef += "}\n"

	return structDef, nil
}

func generateFields(structDef string, data map[string]interface{}, prefix string) string {
	for key, value := range data {
		fieldName := fmt.Sprintf("%s%s", prefix, key)
		fieldName = toCamelCase(fieldName)
		switch v := value.(type) {
		case map[string]interface{}:
			structDef = generateFields(structDef, v, fieldName+"_")
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
	structDef = generateFields(structDef, data, "")
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

func main() {
	jsonStr := `{"name": "John", "age": 30, "isStudent": false, "grades": [90, 85, 88], "friends": [{"name": "Alice", "age": 28}, {"name": "Bob", "age": 32}]}`

	structName := "Person"
	structDef, err := generateStruct(jsonStr, structName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	formattedDef, err := format.Source([]byte(structDef))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(formattedDef))
}
