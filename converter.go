package GoToJava

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"GoToJava/constants"
	"GoToJava/ssa_helpers"
)

type Converter struct {
	DirPath string
	genName int
	ptrCnt  int
	ptrNow  bool
	currPtr uintptr

	used     map[string]bool
	usedPtr  map[uintptr]map[string]int
	inlineId map[string]int

	isJacoSupported bool
}

func CreateConverter(dirPath string, isJacoSupported bool) Converter {
	return Converter{
		DirPath:         dirPath,
		genName:         0,
		used:            map[string]bool{},
		usedPtr:         map[uintptr]map[string]int{},
		inlineId:        map[string]int{},
		isJacoSupported: isJacoSupported,
	}
}

func convertBaseType(goName string) string {
	switch goName {
	case "int", "int32", "uint16", "rune":
		return "Long"
	case "int16", "uint8", "byte":
		return "Long"
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
	return "Long"
}

func (conv *Converter) getInnerStructs(fieldType reflect.Type, kind reflect.Kind) string {
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
	case reflect.Slice, reflect.Array:
		fieldType = fieldType.Elem()
		kind = fieldType.Kind()

		name := conv.getInnerStructs(fieldType, kind)

		return "List<" + name + ">"
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

func (conv *Converter) convertStruct(structure interface{}) (string, error) {
	structVal := reflect.ValueOf(structure)
	structType := reflect.TypeOf(structure)
	structKind := structType.Kind()

	if structKind == reflect.Pointer {
		structType = structType.Elem()
		structVal = structVal.Elem()
	}

	name := structType.String()

	if strings.Contains(name, "struct") {
		id, ok := conv.inlineId[name]
		if !ok {
			id = conv.genName
			conv.genName++
			conv.inlineId[name] = id
		}
		name = fmt.Sprintf("generatedInlineStruct_%03d", id)
	}

	name = strings.ReplaceAll(name, ".", "_")

	if strings.Contains(name, "/") {
		return "", fmt.Errorf("name of the structure contains the '/' symbol")
	}
	if conv.used[name] {
		return name, nil
	}
	conv.used[name] = true

	filePath := filepath.Join(".", conv.DirPath, name+".kt")
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	structDef := constants.PackageLine + readerImports

	if conv.isJacoSupported {
		structDef = ssa_helpers.AddImportAndDefinition(structDef, name)
	} else {
		structDef += fmt.Sprintf(constants.StructDefinition, name)
	}

	deserializer := fmt.Sprintf(deserializeFunStart, name, name, name, name)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldType := field.Type
		//println(field.Type.String())

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

		fieldType = structVal.Field(i).Type()
		kind := fieldType.Kind()

		ktName := conv.getInnerStructs(fieldType, kind)

		if ktName == "" {
			// unsupported, ex functions
			continue
		}

		structDef += fmt.Sprintf(structField,
			field.Name, ktName)
		deserializer += fmt.Sprintf(deserializeField, field.Name, ktName)
	}

	if conv.isJacoSupported {
		structDef = ssa_helpers.AddInterfaceFunctions(structDef, name)
	}

	structDef += "}\n\n"
	file.Write([]byte(structDef))

	deserializer += deserializeEnd
	file.Write([]byte(deserializer))

	return name, nil
}

func getFieldString(conv *Converter, startString string) (string, bool) {
	skip := false

	if conv.ptrNow && conv.currPtr != 0 {
		var id int

		nameToID, ok := conv.usedPtr[conv.currPtr]

		if !ok {
			conv.usedPtr[conv.currPtr] = make(map[string]int)
			id = conv.ptrCnt
			conv.ptrCnt++

			conv.usedPtr[conv.currPtr][startString] = id
		} else {
			if id, ok = nameToID[startString]; ok {
				skip = true
			} else {
				id = conv.ptrCnt
				conv.ptrCnt++

				conv.usedPtr[conv.currPtr][startString] = id
			}
		}

		startString += " " + strconv.Itoa(id)

		conv.ptrNow = false
	}
	conv.ptrNow = false

	startString += "\n"
	return startString, skip
}

func (conv *Converter) fillInnerStructs(fieldType reflect.Type, fieldVal reflect.Value, kind reflect.Kind, fillerFile io.Writer) {
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

	case reflect.Slice, reflect.Array:
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

		fillerFile.Write([]byte("end\n"))
	case reflect.Struct:
		name := fieldType.String()

		if strings.Contains(name, "struct") {
			id, ok := conv.inlineId[name]
			if !ok {
				id = conv.genName
				conv.genName++
				conv.inlineId[name] = id
			}
			name = fmt.Sprintf("generatedInlineStruct_%03d", id)
		}

		name = strings.ReplaceAll(name, ".", "_")

		if _, ok := conv.used[name]; !ok {
			conv.convertStruct(reflect.Zero(fieldType).Interface())
		}

		structString, skip := getFieldString(conv, name)

		fillerFile.Write([]byte(structString))

		if skip {
			return
		}

		if fieldVal.Kind() != 0 {
			for i := 0; i < fieldType.NumField(); i++ {
				//fmt.Printf("%+v\n", fieldType.Field(i))
				field := fieldType.Field(i)
				innerFieldType := field.Type
				//println(field.Type.String())

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
				if strings.Contains(innerFieldType.String(), "/") {
					continue
				}

				innerFieldVal := fieldVal.Field(i)
				innerKind := innerFieldVal.Kind()

				if innerKind == reflect.Func {
					continue
				}

				fillerFile.Write([]byte(field.Name + " "))

				conv.fillInnerStructs(innerFieldType, innerFieldVal, innerKind, fillerFile)
			}
		}

		fillerFile.Write([]byte("end\n"))
	default:
		conv.ptrNow = false

		ktType := convertBaseType(kind.String())
		defaultVal := "0"
		if ktType == "Any" {
			defaultVal = "nil"
		}

		if fieldVal.IsValid() {
			switch kind {
			case reflect.String:
				fillerFile.Write([]byte(fmt.Sprintf("%s\n%q\n", ktType, fieldVal.String())))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				kek := fieldVal.Int()
				if kek > 1000000000000 {
					kek = 1000000000000
				}
				fillerFile.Write([]byte(fmt.Sprintf("%s\n%v\n", ktType, strconv.FormatInt(kek, 10))))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				kek := fieldVal.Uint()
				if kek > 1000000000000 {
					kek = 1000000000000
				}
				fillerFile.Write([]byte(fmt.Sprintf("%s\n%v\n", ktType, strconv.FormatUint(kek, 10))))
			case reflect.Bool:
				fillerFile.Write([]byte(fmt.Sprintf("%s\n%v\n", ktType, strconv.FormatBool(fieldVal.Bool()))))
			default:
				fillerFile.Write([]byte(fmt.Sprintf("%s\n%q\n", ktType, fieldVal)))
			}
		} else {
			fillerFile.Write([]byte(fmt.Sprintf("%s\n%s\n", ktType, defaultVal)))
		}
	}
}

func (conv *Converter) fillValues(structure interface{}, fillerFile io.Writer) error {
	structVal := reflect.ValueOf(structure)
	structType := reflect.TypeOf(structure)
	structKind := structType.Kind()

	conv.fillInnerStructs(structType, structVal, structKind, fillerFile)

	return nil
}

func (conv *Converter) generateBaseDeserializers() error {
	filePath := filepath.Join(".", conv.DirPath, "baseDeserializers.kt")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	des := constants.PackageLine + readerImports + readBaseTypes
	_, err = file.Write([]byte(des))
	return err
}

func (conv *Converter) generateEntrypoint() error {
	filePath := filepath.Join(".", conv.DirPath, "Entrypoint.kt")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	start := constants.PackageLine + readerImports + kotlinConstants

	for name := range conv.used {
		start += ",\n"
		start += fmt.Sprintf(funcMapLine, name, name)
	}

	start += "\n)\n" + entrypoint

	_, err = file.Write([]byte(start))
	return err
}

func (conv *Converter) GenerateStructures(structure interface{}) error {
	if conv.isJacoSupported {
		err := ssa_helpers.GenerateJacoStructs(conv.DirPath)
		if err != nil {
			return err
		}
	}

	_, err := conv.convertStruct(structure)
	if err != nil {
		return err
	}

	err = conv.generateBaseDeserializers()
	if err != nil {
		return err
	}

	err = conv.generateEntrypoint()

	return err
}

func (conv *Converter) FillStructures(fillerFile io.Writer, structure interface{}) error {
	err := conv.fillValues(structure, fillerFile)
	if err != nil {
		return err
	}

	err = conv.generateBaseDeserializers()
	if err != nil {
		return err
	}

	err = conv.generateEntrypoint()

	return err
}
