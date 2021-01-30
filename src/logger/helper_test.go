package logger

import (
	"testing"
)


func TestHelper(t *testing.T) {
	testHelper := NewTestHelper()
	if testHelper.DEMO_URL == "" {
		t.Error("Helper DEMO_URL is empty")
	}
}