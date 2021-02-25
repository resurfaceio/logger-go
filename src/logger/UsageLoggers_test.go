package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvidesDefaultUrl(t *testing.T) {
	//I don't think this makes much sense, this string should be static
	//maybe implement singleton like with the helper struct
	uLogger := NewUsageLoggers()
	//compare to empty string because there is no nil string in go
	assert.Equal(t, "", uLogger.UrlByDefault())
}
