package main

import "fmt"

type A1 struct {
	text1 string
	val1  int
}

type A2 struct {
	A1
	text2 string
	val2  int
}

func main() {
	a1 := A1{
		text1: "test",
		val1:  1,
	}
	fmt.Printf("a1 : %v\n", a1)

	a2 := A2{
		A1: A1{
			text1: "test",
			val1:  1,
		},
		text2: "test2",
		val2:  2,
	}
	fmt.Printf("a2 : %v\n", a2)
}
