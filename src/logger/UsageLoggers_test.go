package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvidesDefaultUrl(t *testing.T) {
	uLogger := GetUsageLoggers()
	//compare to empty string because there is no nil string in go
	assert.Equal(t, "", uLogger.UrlByDefault())
}
