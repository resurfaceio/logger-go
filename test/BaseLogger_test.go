package test

import (
	"../src/logger"
	"fmt"
	"testing"
)

func TestBaseLogger(t *testing.T) {
	testHelper := NewTestHelper()

	testLogger := logger.NewLogger()

	testLogger.SetLogFlag(false)

	_, err := testLogger.Get(testHelper.DEMO_URL)
	if err == nil {
		fmt.Println("Get request success")
	}
	if testLogger.LOG_FLAG {
		t.Error("LOG_FLAG true when set to false")
	}
}
