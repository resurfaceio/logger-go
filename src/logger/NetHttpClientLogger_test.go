package logger

import (
	"fmt"
	"testing"
)

func TestBaseLogger(t *testing.T) {
	testHelper := NewTestHelper()

	testLogger := newLogger()

	testLogger.SetLogFlag(false)

	_, err := testLogger.Get(testHelper.DEMO_URL)
	if err == nil {
		fmt.Println("Get request success")
	}
	if testLogger.LOGFLAG {
		t.Error("LOG_FLAG true when set to false")
	}
}
