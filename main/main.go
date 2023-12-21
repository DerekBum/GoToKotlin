package main

import (
	"GoToJava"
	"fmt"
)

type smth interface {
	SomeFunc() bool
}

type lol struct {
	fl1 string
}

type lel struct {
	fl21 []string
}

func (l lel) SomeFunc() bool {
	return true
}

type kekes struct {
	Field1      int
	FieldString string
	field3      bool
	F4          byte
	f5          *int16
	f6          []uint
	f7          *[]uint
	f8          interface{}
	F9          lol
	//f10         smth
	f11 map[string]bool
}

func main() {
	//lele := lel{}
	//k := kekes{f10: lele}
	k := kekes{}
	fmt.Printf("%v", GoToJava.RunConverter("example", k))
}
