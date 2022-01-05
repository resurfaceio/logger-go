// © 2016-2022 Resurface Labs Inc.

package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvidesDefaultUrl(t *testing.T) {
	uLogger, error := GetUsageLoggers()
	if error != nil {
		assert.Falsef(t, true, "GetUsageLoggers failed: %s", error.Error())
	}
	//compare to empty string because there is no nil string in go
	assert.Equal(t, "", uLogger.UrlByDefault())
}
