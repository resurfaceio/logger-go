// © 2016-2024 Graylog, Inc.

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
	key := "USAGE_LOGGERS_URL"
	defer os.Setenv(key, os.Getenv(key))

	url := "http://thisurlisnotfrom.env/file"
	os.Setenv(key, url)
	uLogger, error := GetUsageLoggers()
	if error != nil {
		assert.Falsef(t, true, "GetUsageLoggers failed: %s", error.Error())
	}
	//compare to hard-coded string assigned to url above
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
	//compare to hard-coded string that must also be set inside a .env file
	assert.Equal(t, "http://thisurlisfrom.env/file", uLogger.UrlByDefault())
}
