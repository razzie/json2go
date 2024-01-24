package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/goccy/go-reflect"
)

func GenerateStruct(data map[string]interface{}, structName string) (string, error) {
	var structDef bytes.Buffer
	fmt.Fprintf(&structDef, "type %s struct {\n", structName)
	generateFields(&structDef, data)
	structDef.WriteString("}\n")

	formattedDef, err := format.Source(structDef.Bytes())
	if err != nil {
		return structDef.String(), err
	}

	return string(formattedDef), nil
}

func generateFields(w io.Writer, data map[string]interface{}) {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fieldName := ToCamelCase(key)
		switch v := data[key].(type) {
		case map[string]interface{}:
			fmt.Fprintf(w, "%s %s `json:\"%s\"`\n", fieldName, generateStruct(v), key)
		case []interface{}:
			fmt.Fprintf(w, "%s %s `json:\"%s\"`\n", fieldName, generateSlice(v), key)
		default:
			fmt.Fprintf(w, "%s %s `json:\"%s\"`\n", fieldName, getGoType(v), key)
		}
	}
}

func generateStruct(data map[string]interface{}) string {
	if len(data) == 0 {
		return "map[string]interface{}"
	}
	var structDef strings.Builder
	structDef.WriteString("struct {\n")
	generateFields(&structDef, data)
	structDef.WriteString("}")
	return structDef.String()
}

func generateSlice(data []interface{}) string {
	if len(data) == 0 {
		return "[]interface{}"
	}
	switch v := data[0].(type) {
	case map[string]interface{}:
		return "[]" + generateStruct(v)
	default:
		return "[]" + getGoType(v)
	}
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
