package GoToJava

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func createDir(name string) error {
	return os.MkdirAll(name, os.ModePerm)
}

type converter struct {
	dirPath string
	genName int
	ptrCnt  int
	ptrNow  bool
	currPtr uintptr

	used    map[string]bool
	usedPtr map[uintptr]int
}

func createConverter(dirPath string) converter {
	return converter{
		dirPath: dirPath,
		genName: 0,
		used:    map[string]bool{},
		usedPtr: map[uintptr]int{},
	}
}

func convertBaseType(goName string) string {
	switch goName {
	case "int", "int32", "uint16", "rune":
		return "Int"
	case "int16", "uint8", "byte":
		return "Short"
	case "int64", "uint32", "uint":
		return "Long"
	case "float32":
		return "Float"
	case "float64":
		return "Double"
	case "string":
		return "String"
	case "bool":
		return "Boolean"
	case "interface {}":
		return "Any"
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

		return "Any"
	case reflect.Pointer:
		fieldType = fieldType.Elem()
		kind = fieldType.Kind()

		name := conv.getInnerStructs(fieldType, kind)

		return name
	case reflect.Slice:
		fieldType = fieldType.Elem()
		kind = fieldType.Kind()

		name := conv.getInnerStructs(fieldType, kind)

		return "Array<" + name + ">"
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

	filePath := filepath.Join(".", conv.dirPath, name+".kt")
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	structDef := fmt.Sprintf(structDefinition, name)

	for i := 0; i < structType.NumField(); i++ {
		//fmt.Printf("%+v\n", structType.Field(i))
		field := structType.Field(i)
		fieldType := field.Type
		println(field.Type.String())

		if field.Name == "_" {
			// This is a blank identifier, no need to send.
			continue
		}
		if field.Name == "object" {
			// Invalid kotlin name.
			field.Name = "Object"
		}
		if field.Name == "val" {
			// Invalid kotlin name.
			field.Name = "Val"
		}

		fieldVal := structVal.Field(i)
		kind := fieldVal.Kind()

		ktName := conv.getInnerStructs(fieldType, kind)

		if ktName == "" {
			// unsupported, ex functions
			continue
		}

		structDef += fmt.Sprintf(structField,
			field.Name, ktName)
	}

	structDef += "}\n"

	file.Write([]byte(structDef))

	return name, nil
}

func getFieldString(conv *converter, startString string) (string, bool) {
	skip := false

	if conv.ptrNow && conv.currPtr != 0 {
		id, ok := conv.usedPtr[conv.currPtr]

		if !ok {
			id = conv.ptrCnt
			conv.ptrCnt++

			conv.usedPtr[conv.currPtr] = id
		} else {
			skip = true
		}

		startString += " " + strconv.Itoa(id)

		conv.ptrNow = false
	}
	conv.ptrNow = false

	startString += "\n"
	return startString, skip
}

func (conv *converter) fillInnerStructs(fieldType reflect.Type, fieldVal reflect.Value, kind reflect.Kind, fillerFile io.Writer) {
	switch kind {
	case reflect.Func:
		// skip
		return
	case reflect.Interface:
		var realVal reflect.Value

		if fieldVal.IsValid() {
			realVal = fieldVal.Elem()
		}

		if realVal.Kind() != 0 {
			fieldVal = realVal
			fieldType = fieldVal.Type()
			kind = fieldType.Kind()

			conv.fillInnerStructs(fieldType, fieldVal, kind, fillerFile)
		} else {
			conv.ptrNow = false
			fillerFile.Write([]byte("nil\n"))
		}
		return

	case reflect.Pointer:
		conv.ptrNow = true
		conv.currPtr = uintptr(fieldVal.UnsafePointer())

		fieldType = fieldType.Elem()
		fieldVal = fieldVal.Elem()
		kind = fieldType.Kind()

		conv.fillInnerStructs(fieldType, fieldVal, kind, fillerFile)

	case reflect.Slice:
		arrayString, skip := getFieldString(conv, "array")

		fillerFile.Write([]byte(arrayString))

		if skip {
			return
		}

		fieldType = fieldType.Elem()
		kind = fieldType.Kind()

		if fieldVal.Kind() != 0 {
			for i := 0; i < fieldVal.Len(); i++ {
				conv.fillInnerStructs(fieldType, fieldVal.Index(i), kind, fillerFile)
			}
		}

		fillerFile.Write([]byte("end\n"))

	case reflect.Map:
		mapString, skip := getFieldString(conv, "map")

		fillerFile.Write([]byte(mapString))

		if skip {
			return
		}

		keyType := fieldType.Key()
		keyKind := keyType.Kind()

		valType := fieldType.Elem()
		valKind := valType.Kind()

		if fieldVal.Kind() != 0 {
			for _, k := range fieldVal.MapKeys() {
				conv.fillInnerStructs(keyType, k, keyKind, fillerFile)
				conv.fillInnerStructs(valType, fieldVal.MapIndex(k), valKind, fillerFile)
			}
		}

		conv.getInnerStructs(keyType, keyKind)
		conv.getInnerStructs(valType, valKind)

		fillerFile.Write([]byte("end\n"))
	case reflect.Struct:
		name := fieldType.String()

		if strings.Contains(name, "struct") {
			name = fmt.Sprintf("generatedInlineStruct_%03d", conv.genName)
			conv.genName++
		}

		name = strings.ReplaceAll(name, ".", "_")

		structString, skip := getFieldString(conv, "struct "+name)

		fillerFile.Write([]byte(structString))

		if skip {
			return
		}

		if fieldVal.Kind() != 0 {
			for i := 0; i < fieldType.NumField(); i++ {
				fmt.Printf("%+v\n", fieldType.Field(i))
				field := fieldType.Field(i)
				innerFieldType := field.Type
				println(field.Type.String())

				if field.Name == "_" {
					// This is a blank identifier, no need to send.
					continue
				}

				innerFieldVal := fieldVal.Field(i)
				innerKind := innerFieldVal.Kind()

				fillerFile.Write([]byte(field.Name + " "))

				conv.fillInnerStructs(innerFieldType, innerFieldVal, innerKind, fillerFile)
			}
		}

		fillerFile.Write([]byte("end\n"))
	default:
		conv.ptrNow = false

		javaType := convertBaseType(kind.String())
		defaultVal := "0"
		if javaType == "Object" {
			defaultVal = "nil"
		}

		if fieldVal.IsValid() {
			fillerFile.Write([]byte(fmt.Sprintf("%s %v\n", javaType, fieldVal)))
		} else {
			fillerFile.Write([]byte(fmt.Sprintf("%s %s\n", javaType, defaultVal)))
		}
	}
}

func (conv *converter) fillValues(structure interface{}, fillerFile io.Writer) error {
	structVal := reflect.ValueOf(structure)
	structType := reflect.TypeOf(structure)
	structKind := structType.Kind()

	conv.fillInnerStructs(structType, structVal, structKind, fillerFile)

	return nil
}

// GENERATE PUBLIC CLASS IN ANOTHER FILE
// FILL IT WITH PUBLIC VARIABLES (GENERATED FROM GO)

func RunConverter(dirPath string, fillerFile io.Writer, structure interface{}) error {
	var err error

	if err = createDir(dirPath); err != nil {
		return err
	}
	conv := createConverter(dirPath)

	_, err = conv.convertStruct(structure)
	if err != nil {
		return err
	}

	conv.genName = 0
	err = conv.fillValues(structure, fillerFile)

	return err
}
