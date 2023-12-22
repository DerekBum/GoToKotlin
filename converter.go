package GoToJava

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func createDir(name string) error {
	return os.MkdirAll(name, os.ModePerm)
}

type converter struct {
	dirPath string
	used    map[string]bool
	genName int
}

func createConverter(dirPath string) converter {
	return converter{
		dirPath: dirPath,
		used:    map[string]bool{},
		genName: 0,
	}
}

func convertBaseType(goName string) string {
	switch goName {
	case "int", "int32", "uint16", "rune":
		return "int"
	case "int16", "uint8", "byte":
		return "short"
	case "int64", "uint32", "uint":
		return "long"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "string":
		return "String"
	case "bool":
		return "Boolean"
	case "interface {}":
		return "Object"
	}
	return ""
}

func (conv *converter) getInnerStructs(fieldType reflect.Type, kind reflect.Kind) string {
	switch kind {
	case reflect.Func:
		// skip
		return ""
	case reflect.Interface:
		fieldVal := reflect.Zero(fieldType)

		if !fieldVal.IsNil() {
			fieldType = fieldType.Elem()
			kind = fieldType.Kind()

			conv.getInnerStructs(fieldType, kind)
		}

		return "Object"
	case reflect.Pointer:
		fieldType = fieldType.Elem()
		kind = fieldType.Kind()

		name := conv.getInnerStructs(fieldType, kind)

		return name
	case reflect.Slice:
		fieldType = fieldType.Elem()
		kind = fieldType.Kind()

		name := conv.getInnerStructs(fieldType, kind)

		return name + "[]"
	case reflect.Map:
		keyType := fieldType.Key()
		keyKind := keyType.Kind()

		valType := fieldType.Elem()
		valKind := valType.Kind()

		keyName := conv.getInnerStructs(keyType, keyKind)
		valName := conv.getInnerStructs(valType, valKind)

		return fmt.Sprintf("Map<%s, %s>", keyName, valName)
	case reflect.Struct:
		sampleStruct := reflect.Zero(fieldType).Interface()
		name, _ := conv.convertStruct(sampleStruct)

		return name
	default:
		return convertBaseType(kind.String())
	}
}

func (conv *converter) convertStruct(structure interface{}) (string, error) {
	structVal := reflect.ValueOf(structure)
	structType := reflect.TypeOf(structure)
	structKind := structType.Kind()

	if structKind == reflect.Pointer {
		structType = structType.Elem()
		structVal = structVal.Elem()
	}

	name := structType.String()

	if strings.Contains(name, "struct") {
		name = fmt.Sprintf("generatedInlineStruct_%03d", conv.genName)
		conv.genName++
	}

	name = strings.ReplaceAll(name, ".", "_")

	if conv.used[name] {
		return name, nil
	}
	conv.used[name] = true

	filePath := filepath.Join(".", conv.dirPath, name+".java")
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	structDef := imports
	structDef += fmt.Sprintf(structDefinition, name)

	for i := 0; i < structType.NumField(); i++ {
		fmt.Printf("%+v\n", structType.Field(i))
		field := structType.Field(i)
		fieldType := field.Type
		println(field.Type.String())

		var javaName string

		fieldVal := structVal.Field(i)
		kind := fieldVal.Kind()

		javaName = conv.getInnerStructs(fieldType, kind)

		if javaName == "" {
			// unsupported, ex functions
			continue
		}

		structDef += fmt.Sprintf(structField,
			javaName, field.Name)
	}

	structDef += "}\n"

	file.Write([]byte(structDef))

	return name, nil
}

func RunConverter(dirPath string, structure interface{}) error {
	if err := createDir(dirPath); err != nil {
		return err
	}
	conv := createConverter(dirPath)

	_, err := conv.convertStruct(structure)
	return err
}
