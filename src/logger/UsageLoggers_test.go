package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvidesDefaultUrl(t *testing.T) {
	uLogger, _ := GetUsageLoggers()

	/*
		without default url environment variable
		set the url should be empty
	*/
	assert.Equal(t, "", uLogger.UrlByDefault())
}
