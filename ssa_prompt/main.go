package main

import (
	"GoToJava"
)

func GetSign(x int) (res int) {
	res = 123
	if x > 0 {
		return 1
	} else if x == 0 {
		return 0
	}
	res = -1
	return
}

func main() {
	var x int
	println(GetSign(x))

	conv := GoToJava.CreateConverter("ssaExample", true)

	conv.GenerateStructures(x)
	conv.FillStructures(nil, x)
}
