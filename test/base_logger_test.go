package test

import (
	"testing"
	//"net/http/httptest"
	//"net/http"
)


func Test_Basic_Logger(t *testing.T) {
	testHelper := GetTestHelper()
	if testHelper.DEMO_URL == "" {
		t.Error("Helper DEMO_URL is empty")
	}
}
