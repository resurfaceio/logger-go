package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Testing HTTP Logger

func TestCreateInstance(t *testing.T) {

	//Creating a single instance
	httpLogger := NewHttpLogger()
	helper := GetTestHelper()
	assert.NotNil(t, httpLogger)
	assert.Equal(t, helper.mockAgent, httpLogger.Agent())
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
	assert.Equal(t, helper.mockAgent, logger1.AGENT())
	assert.False(t, logger1.Enablable())
	assert.False(t, logger1.Enabled())
	assert.Nil(t, logger1.Queue())
	assert.Equal(t, url1, logger1.Url())

	//Logger 2
	assert.NotNil(t, logger2)
	assert.Equal(t, helper.mockAgent, logger2.Agent())
	assert.False(t, logger2.Enablable())
	assert.False(t, logger2.Enabled())
	assert.Nil(t, logger2.Queue())
	assert.Equal(t, url2, logger2.Url())

	//Logger 3
	assert.NotNil(t, logger3)
	assert.Equal(t, helper.mockAgent, logger3.Agent())
	assert.False(t, logger3.Enablable())
	assert.False(t, logger3.Enabled())
	assert.Nil(t, logger3.Queue())
	assert.Nil(t, logger3.Url())
	assert.Equal(t, helper.demoURL, logger3.Url())

}

func TestHasValidAgent(t *testing.T) {
	//Has Valid Agent Test
	helper := GetTestHelper()
	logger1 := NewHttpLogger()

	agent := helper.mockAgent
	assert.Greater(t, len(agent), 0)
	assert.Equal(t, ".go", agent[len(agent)-3:])
	assert.NotContains(t, agent, "\\")
	assert.NotContains(t, agent, "\"")
	assert.NotContains(t, agent, "'")
	assert.Equal(t, agent, logger1.Agent())

}
