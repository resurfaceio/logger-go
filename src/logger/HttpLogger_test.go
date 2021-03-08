package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Testing HTTP Logger

func TestCreateInstance(t *testing.T) {

	//Creating a single instance
	httpLogger := NewHttpLogger()
	assert.NotNil(t, httpLogger)
	assert.Equal(t, HttpLogger.agent, httpLogger.Agent())
	assert.False(t, httpLogger.Enablable())
	assert.False(t, httpLogger.Enabled())
	assert.Nil(t, httpLogger.Queue())
	assert.Nil(t, httpLogger.Url())

}

func TestCreateMultipleInstances(t *testing.T) {
	//Creating multiple instances
	url1 := "https://resurface.io"
	url2 := "https://whatever.com"
	helper := GetTestHelper()

	logger1 := NewHttpLogger(url1)
	logger2 := NewHttpLogger(url2)
	logger3 := NewHttpLogger(helper.demoURL)

	//Logger 1
	assert.NotNil(t, logger1)
	assert.Equal(t, HttpLogger.agent, logger1.AGENT())
	assert.True(t, logger1.Enablable())
	assert.True(t, logger1.Enabled())
	assert.Equal(t, url1, logger1.Url())

	//Logger 2
	assert.NotNil(t, logger2)
	assert.Equal(t, HttpLogger.agent, logger2.Agent())
	assert.True(t, logger2.Enablable())
	assert.True(t, logger2.Enabled())
	assert.Equal(t, url2, logger2.Url())

	//Logger 3
	assert.NotNil(t, logger3)
	assert.Equal(t, HttpLogger.agent, logger3.Agent())
	assert.True(t, logger3.Enablable())
	assert.True(t, logger3.Enabled())
	assert.Equal(t, helper.demoURL, logger3.Url())

	//Testing Usage Logger
	//Disable
	UsageLoggers.disable()
	assert.False(t, UsageLoggers.isEnabled())
	assert.False(t, logger1.isEnabled())
	assert.False(t, logger2.isEnabled())
	assert.False(t, logger3.isEnabled())
	//Enable
	UsageLoggers.enable()
	assert.True(t, UsageLoggers.isEnabled())
	assert.True(t, logger1.isEnabled())
	assert.True(t, logger2.isEnabled())
	assert.True(t, logger3.isEnabled())
}

func TestHasValidAgent(t *testing.T) {
	//Has Valid Agent Test
	httpLogger := NewHttpLogger()

	agent := HttpLogger.agent
	assert.Greater(t, len(agent), 0)
	assert.Equal(t, ".go", agent[len(agent)-3:])
	assert.NotContains(t, agent, "\\")
	assert.NotContains(t, agent, "\"")
	assert.NotContains(t, agent, "'")
	assert.Equal(t, agent, httpLogger.Agent())

}
