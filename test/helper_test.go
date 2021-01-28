package test

import (
	"testing"
)

func TestHelper(t *testing.T) {
	testHelper := GetTestHelper()
	if testHelper.DEMO_URL == "" {
		t.Error("Helper DEMO_URL is empty")
	}
}



