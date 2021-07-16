package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Testing HTTP Logger

func TestCreateInstance(t *testing.T) {

	//Creating a single instance
	httpLogger, _ := newHttpLogger(Options{})
	assert.NotNil(t, httpLogger)
	assert.Equal(t, httpLoggerAgent, httpLogger.agent)
	assert.False(t, httpLogger.enableable)
	assert.False(t, httpLogger.Enabled())
	assert.Nil(t, httpLogger.queue)
	assert.Equal(t, "", httpLogger.url)

}

func TestCreateMultipleInstances(t *testing.T) {
	//Creating multiple instances
	url1 := "https://resurface.io"
	url2 := "https://whatever.com"
	helper := newTestHelper()
	options1 := Options{
		Url:     url1,
		Enabled: true,
	}
	options2 := Options{
		Url:     url2,
		Enabled: true,
	}
	options3 := Options{
		Url:     helper.demoURL,
		Enabled: true,
	}

	logger1, _ := newHttpLogger(options1)
	logger2, _ := newHttpLogger(options2)
	logger3, _ := newHttpLogger(options3)

	//Logger 1
	assert.NotNil(t, logger1)
	assert.Equal(t, httpLoggerAgent, logger1.agent)
	assert.True(t, logger1.enableable)
	assert.True(t, logger1.Enabled())
	assert.Equal(t, url1, logger1.url)

	//Logger 2
	assert.NotNil(t, logger2)
	assert.Equal(t, httpLoggerAgent, logger2.agent)
	assert.True(t, logger2.enableable)
	assert.True(t, logger2.Enabled())
	assert.Equal(t, url2, logger2.url)

	//Logger 3
	assert.NotNil(t, logger3)
	assert.Equal(t, httpLoggerAgent, logger3.agent)
	assert.True(t, logger3.enableable)
	assert.True(t, logger3.Enabled())
	assert.Equal(t, helper.demoURL, logger3.url)

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
	httpLogger, _ := newHttpLogger(Options{})

	assert.Greater(t, len(httpLoggerAgent), 0)
	assert.Equal(t, ".go", httpLoggerAgent[len(httpLoggerAgent)-3:])
	assert.NotContains(t, httpLoggerAgent, "\\")
	assert.NotContains(t, httpLoggerAgent, "\"")
	assert.NotContains(t, httpLoggerAgent, "'")
	assert.Equal(t, httpLoggerAgent, httpLogger.agent)

}
