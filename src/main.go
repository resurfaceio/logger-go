package main

import (
	"./logger"
)

func main() {
	testLogger := logger.NewLogger()
	testLogger.Get("https://www.google.com/")
}
