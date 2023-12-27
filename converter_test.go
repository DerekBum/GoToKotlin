package GoToJava_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

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

type innerStruct struct {
	if1 int
	if2 bool
}

type BigStruct struct {
	innerStruct

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
}

var testDir = "./test_struct"
var genFiles = []string{
	"GoToJava_test_BigStruct.java",
	"GoToJava_test_innerStruct.java",
	"GoToJava_test_interfaceImpl.java",
	"GoToJava_test_someStruct.java",
}

func TestConvert(t *testing.T) {
	st := BigStruct{}

	err := GoToJava.RunConverter(testDir, st)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	dir, err := os.ReadDir(testDir)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(dir) != len(genFiles) {
		t.Errorf("wrong number of files, want %d, got %d", len(genFiles), len(dir))
	}
	for i, file := range dir {
		if file.Name() != genFiles[i] {
			t.Errorf("unexpected file, want %s, got %s", genFiles[i], file.Name())
		}
	}

	paths := genFiles
	for i := range paths {
		paths[i] = filepath.Join(testDir, paths[i])
	}

	cmd := exec.Command("javac", paths...)
	_, err = cmd.Output()

	if err != nil {
		t.Errorf("Java files did not compile")
	}
}

func TestMain(m *testing.M) {
	code := m.Run()

	os.RemoveAll(testDir)

	os.Exit(code)
}
