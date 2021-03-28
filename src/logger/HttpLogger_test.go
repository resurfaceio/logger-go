package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Testing HTTP Logger

func TestCreateInstance(t *testing.T) {

	//Creating a single instance
	httpLogger := NewHttpLogger(Options{})
	assert.NotNil(t, httpLogger)
	assert.Equal(t, httpLoggerAgent, httpLogger.Agent())
	assert.False(t, httpLogger.Enableable())
	assert.False(t, httpLogger.Enabled())
	assert.Nil(t, httpLogger.Queue())
	assert.Nil(t, httpLogger.Url())

}

func TestCreateMultipleInstances(t *testing.T) {
	//Creating multiple instances
	url1 := "https://resurface.io"
	url2 := "https://whatever.com"
	helper := GetTestHelper()
	options1 := Options{
		url: url1,
	}
	options2 := Options{
		url: url2,
	}
	options3 := Options{
		url: helper.demoURL,
	}

	logger1 := NewHttpLogger(options1)
	logger2 := NewHttpLogger(options2)
	logger3 := NewHttpLogger(options3)

	//Logger 1
	assert.NotNil(t, logger1)
	assert.Equal(t, httpLoggerAgent, logger1.Agent())
	assert.True(t, logger1.Enableable())
	assert.True(t, logger1.Enabled())
	assert.Equal(t, url1, logger1.Url())

	//Logger 2
	assert.NotNil(t, logger2)
	assert.Equal(t, httpLoggerAgent, logger2.Agent())
	assert.True(t, logger2.Enableable())
	assert.True(t, logger2.Enabled())
	assert.Equal(t, url2, logger2.Url())

	//Logger 3
	assert.NotNil(t, logger3)
	assert.Equal(t, httpLoggerAgent, logger3.Agent())
	assert.True(t, logger3.Enableable())
	assert.True(t, logger3.Enabled())
	assert.Equal(t, helper.demoURL, logger3.Url())

	//Testing Usage Logger
	//Disable
	usageLoggers, _ := GetUsageLoggers()
	usageLoggers.Disable()
	assert.False(t, usageLoggers.IsEnabled())
	assert.False(t, logger1.Enabled())
	assert.False(t, logger2.Enabled())
	assert.False(t, logger3.Enabled())
	//Enable
	usageLoggers.Enable()
	assert.True(t, usageLoggers.IsEnabled())
	assert.True(t, logger1.Enabled())
	assert.True(t, logger2.Enabled())
	assert.True(t, logger3.Enabled())
}

func TestHasValidAgent(t *testing.T) {
	//Has Valid Agent Test
	httpLogger := NewHttpLogger(Options{})

	assert.Greater(t, len(httpLoggerAgent), 0)
	assert.Equal(t, ".go", httpLoggerAgent[len(httpLoggerAgent)-3:])
	assert.NotContains(t, httpLoggerAgent, "\\")
	assert.NotContains(t, httpLoggerAgent, "\"")
	assert.NotContains(t, httpLoggerAgent, "'")
	assert.Equal(t, httpLoggerAgent, httpLogger.Agent())

}
