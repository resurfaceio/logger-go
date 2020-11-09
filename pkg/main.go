package main

import (
	"fmt"
	"./test/helper"
)

func main() {
	fmt.Println("Hello Resurface!")

	testHelper := helper.GetHelper()

	fmt.Println(testHelper)
}
