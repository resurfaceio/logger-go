// © 2016-2021 Resurface Labs Inc.

package logger

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//Testing HTTP Logger

func TestCreateInstance(t *testing.T) {

	//Creating a single instance
	HttpLogger, _ := NewHttpLogger(Options{})
	assert.NotNil(t, HttpLogger)
	assert.Equal(t, httpLoggerAgent, HttpLogger.agent)
	assert.False(t, HttpLogger.enableable)
	assert.False(t, HttpLogger.Enabled())
	assert.Nil(t, HttpLogger.queue)
	assert.Equal(t, "", HttpLogger.url)

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

	logger1, _ := NewHttpLogger(options1)
	logger2, _ := NewHttpLogger(options2)
	logger3, _ := NewHttpLogger(options3)

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
	HttpLogger, _ := NewHttpLogger(Options{})

	assert.Greater(t, len(httpLoggerAgent), 0)
	assert.Equal(t, ".go", httpLoggerAgent[len(httpLoggerAgent)-3:])
	assert.NotContains(t, httpLoggerAgent, "\\")
	assert.NotContains(t, httpLoggerAgent, "\"")
	assert.NotContains(t, httpLoggerAgent, "'")
	assert.Equal(t, httpLoggerAgent, HttpLogger.agent)
}

func TestSetsNowAndInterval(t *testing.T) {
	helper := newTestHelper()

	opt := Options{
		Queue:   make([]string, 0),
		Enabled: true,
		Url:     helper.demoURL1,
	}

	logger, _ := NewHttpLogger(opt)

	SendHttpMessage(logger, helper.MockResponse(), helper.MockRequestWithJson(), 0, 0)

	assert.Contains(t, logger.queue[0], "[\"now", "SendHttpMessage did not append 'now' to message on null entry")
	assert.NotContains(t, logger.queue[0], "[\"interval", "SendHttpMessage appended 'interval' to message on null entry")

	logger, _ = NewHttpLogger(opt)

	now := time.Now()
	time.Sleep(100 * time.Millisecond)
	interval := time.Since(now).Milliseconds()

	SendHttpMessage(logger, helper.MockResponse(), helper.MockRequestWithJson(), (now.Unix() * int64(time.Millisecond)), interval)

	assert.Contains(t, logger.queue[0], "[\"now", "SendHttpMessage did not append 'now' to message on manual entry")
	assert.Contains(t, logger.queue[0], "[\"interval\",\"1", "SendHttpMessage did not append 'interval' to message on manual entry")
}
