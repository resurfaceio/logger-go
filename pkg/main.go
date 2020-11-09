package main

import (
	"fmt"
	"../test"
)

func main() {
	fmt.Println("Hello Resurface!")

	testHelper := test.GetTestHelper()

	fmt.Println(testHelper)
}
