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
}

func createConverter(path string) converter {
	return converter{
		dirPath: path,
		used:    map[string]bool{},
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

func convertType(goName string) string {
	// Base types
	base := convertBaseType(goName)
	if base != "" {
		return base
	}

	if strings.HasPrefix(goName, "*") {
		// A pointer
		goVal := strings.TrimPrefix(goName, "*")
		javaVal := convertType(goVal)
		return javaVal //TODO
	}

	if strings.HasPrefix(goName, "map[") {
		splited := strings.FieldsFunc(goName, func(r rune) bool {
			return r == '[' || r == ']'
		})
		goValKey := splited[1]
		goValVal := splited[2]

		javaValKey := convertType(goValKey)
		javaValVal := convertType(goValVal)

		return fmt.Sprintf("Map<%s, %s>", javaValKey, javaValVal)
	}

	if strings.HasPrefix(goName, "[]") {
		// An array
		goVal := strings.TrimPrefix(goName, "[]")
		javaVal := convertType(goVal)
		return javaVal + "[]"
	}

	return "" // Other struct
}

func (conv *converter) convertStruct(structure interface{}) error {
	//lol := reflect.ValueOf(&structure).Elem()
	structType := reflect.TypeOf(structure)
	name := structType.Name()

	if conv.used[name] {
		return nil
	}
	conv.used[name] = true

	filePath := filepath.Join(".", conv.dirPath, name+".java")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	structDef := fmt.Sprintf(structDefinition, name)

	for i := 0; i < structType.NumField(); i++ {
		fmt.Printf("%+v\n", structType.Field(i))
		field := structType.Field(i)
		println(field.Type.String())
		javaName := convertType(field.Type.String())

		if javaName == "" {
			/*val := reflect.ValueOf(structure).Field(i).Interface()

			f := reflect.TypeOf(val)
			fmt.Printf("%+v\n", f)*/
			javaName = field.Type.String()
			if strings.Contains(javaName, ".") {
				splited := strings.Split(javaName, ".")
				javaName = splited[len(splited)-1]
			}
			sampleStruct := reflect.Zero(field.Type).Interface()
			conv.convertStruct(sampleStruct)
		}

		structDef += fmt.Sprintf(structField,
			javaName, field.Name)
	}

	structDef += "}\n"

	file.Write([]byte(structDef))

	return nil
}

func RunConverter(name string, structure interface{}) error {
	if err := createDir(name); err != nil {
		return err
	}
	conv := createConverter(name)

	return conv.convertStruct(structure)
}
