package logger

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

//Testing HTTP Logger

func TestCreateInstance(t *testing.T) {

	//Creating a single instance
	logger := newLogger()

	assert.NotNil(t, logger)
	assert.Equal(t, netHttpClientLogger.AGENT, logger.AGENT)
	assert.False(t, logger.isEnableable)
	assert.False(t, logger.LOGFLAG)
	assert.Nil(t, logger.GetQueue)
	assert.Nil(t, logger.GetURL)

}

func TestCreateMultipleInstances(T *testing.T) {
	//Creating multiple instances
	url1 := "https://resurface.io"
	url2 := "https://whatever.com"

	logger1 := newLogger(url1)
	logger2 := newLogger(url2)
	logger3 := newLogger(Helper.demoURL)

	//Logger 1
	assert.NotNil(t, logger1)
	assert.Equal(t, HttpLogger.AGENT, logger1.AGENT)
	assert.False(t, logger1.isEnableable)
	assert.False(t, logger1.LOGFLAG)
	assert.Nil(t, logger1.GetQueue)
	assert.Equal(t, url1, logger1.GetURL)

	//Logger 2
	assert.NotNil(t, logger2)
	assert.Equal(t, HttpLogger.AGENT, logger2.AGENT)
	assert.False(t, logger2.isEnableable)
	assert.False(t, logger2.LOGFLAG)
	assert.Nil(t, logger2.GetQueue)
	assert.Equal(t, url2, logger2.GetURL)

	//Logger 3
	assert.NotNil(t, logger3)
	assert.Equal(t, HttpLogger.AGENT, logger3.AGENT)
	assert.False(t, logger3.isEnableable)
	assert.False(t, logger3.LOGFLAG)
	assert.Nil(t, logger3.GetQueue)
	assert.Nil(t, logger3.GetURL)
	assert.Equal(t, testHelper.demoURL, logger3.GetURL)

}

func TestHasValidAgent(T *testing.T){
	//Has Valid Agent Test
	agent := logger.AGENT
	assert.Greater(t, len(agent),0)
	assert.Equal(t, ".go",agent[len(agent)-3:])
	assert.False(t, assert.Contains(t, agent, "\\" ))
	assert.False(t, assert.Contains(t, agent, "\"" ))
	assert.False(t, assert.Contains(t, agent, "'" ))
	assert.Equal(t, agent, newLogger().AGENT)

}

