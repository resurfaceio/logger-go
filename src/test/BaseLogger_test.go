package test

import (
	"../logger"
	"testing"
)

func BaseLoggerTest(t *testing.T) {
	testHelper := GetTestHelper()
	testHelper.MockCustomApp()

	testLogger := logger.NewLogger()
	testLogger.Get("https://github.com/resurfaceio/")
}
