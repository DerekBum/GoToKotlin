package main

import (
	"fmt"
	"os"

	"GoToJava"
)

type nonEmptyInterface interface {
	SomeFunc() bool
}

type someStruct struct {
	fl1 string
}

type interfaceImpl struct {
	fl21 []string
}

func (l interfaceImpl) SomeFunc() bool {
	return true
}

type customType int16

type BigStruct struct {
	_ interfaceImpl

	Field1      int
	FieldString string
	field3      bool
	F4          byte
	f5          *int16
	f6          []uint
	f7          *[]uint
	f8          interface{}
	f85         interface{}
	f9          someStruct
	f10         nonEmptyInterface
	f11         map[string]bool
	f12         map[interface{}]interface{}
	f13         *map[someStruct]interfaceImpl
	f14         *map[*interfaceImpl]*someStruct
	f15         []someStruct
	f16         *[]*someStruct
	f17         *nonEmptyInterface
	f18         *someStruct
	f19         customType
	f20         *customType
	f21         []*customType
	f22         map[customType]*customType
	f23         *BigStruct
	f24         interface{}
}

func main() {
	/*// Replace interface{} with any for this test.
	// Parse the source files.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./ssa_prompt/main.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Print(err) // parse error
	}
	files := []*ast.File{f}
	// Create the type-checker's package.
	pkg := types.NewPackage("main", "")
	// Type-check the package, load dependencies.
	// Create and build the SSA program.
	pkgBuild, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()}, fset, pkg, files, ssa.PrintFunctions)
	if err != nil {
		fmt.Print(err) // type error in some package
	}*/

	impl := interfaceImpl{}
	k := BigStruct{f85: 123, f10: impl}
	k.f9.fl1 = "123"
	k.f11 = make(map[string]bool)
	k.f11["lolol"] = false

	mp := make(map[someStruct]interfaceImpl)
	mp[someStruct{fl1: "key"}] = interfaceImpl{fl21: []string{"some", "value"}}
	k.f13 = &mp

	f10 := interfaceImpl{fl21: []string{"hello world!", "hola", "привет"}}
	k.f10 = f10
	k.f23 = &k
	k.f24 = &k

	os.Mkdir("example", os.ModePerm)
	file, _ := os.Create("example/filled.txt")

	fmt.Printf("%v", GoToJava.RunConverter("example", file, &k))
}
