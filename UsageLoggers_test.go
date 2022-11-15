// © 2016-2022 Resurface Labs Inc.

package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvidesEmptyDefaultUrl(t *testing.T) {
	if _, err := os.Stat(".env"); err == nil {
		t.Skip(".env file exists")
	}
	uLogger, error := GetUsageLoggers()
	if error != nil {
		assert.Falsef(t, true, "GetUsageLoggers failed: %s", error.Error())
	}
	//compare to empty string because there is no nil string in go
	assert.Equal(t, "", uLogger.UrlByDefault())
}

func TestProvidesDefaultUrl(t *testing.T) {
	if _, err := os.Stat(".env"); err == nil {
		t.Skip(".env file exists")
	}
	url := "http://thisurlisnotfrom.env/file"
	os.Setenv("USAGE_LOGGERS_URL", url)
	uLogger, error := GetUsageLoggers()
	if error != nil {
		assert.Falsef(t, true, "GetUsageLoggers failed: %s", error.Error())
	}
	//compare to empty string because there is no nil string in go
	assert.Equal(t, url, uLogger.UrlByDefault())
}

func TestProvidesDefaultUrlFromFile(t *testing.T) {
	if _, err := os.Stat(".env"); err != nil {
		t.Skip(".env file does not exist")
	}
	uLogger, error := GetUsageLoggers()
	if error != nil {
		assert.Falsef(t, true, "GetUsageLoggers failed: %s", error.Error())
	}
	//compare to empty string because there is no nil string in go
	assert.Equal(t, "http://thisurlisfrom.env/file", uLogger.UrlByDefault())
}
