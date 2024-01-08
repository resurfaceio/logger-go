// © 2016-2024 Graylog, Inc.

package logger

import (
	"os"
	"sync"
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

func TestCreateInstanceWithURLByDefault(t *testing.T) {
	key := "USAGE_LOGGERS_URL"
	defer os.Setenv(key, os.Getenv(key))

	url := "http://whatever.com:8123/some/path"
	os.Setenv(key, url)
	//Creating a single instance
	HttpLogger, _ := NewHttpLogger(Options{})
	assert.NotNil(t, HttpLogger)
	assert.Equal(t, httpLoggerAgent, HttpLogger.agent)
	assert.False(t, HttpLogger.enableable)
	assert.False(t, HttpLogger.Enabled())
	assert.Nil(t, HttpLogger.queue)
	assert.Equal(t, url, HttpLogger.url)
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

	SendHttpMessage(logger, helper.MockResponse(), helper.MockRequestWithJson(), 0, 0, nil)

	assert.Contains(t, logger.queue[0], "[\"now", "SendHttpMessage did not append 'now' to message on null entry")
	assert.Contains(t, logger.queue[0], "[\"interval\",\"1", "SendHttpMessage did not appended 'floor interval' to message on null entry")

	logger, _ = NewHttpLogger(opt)

	now := time.Now()
	time.Sleep(100 * time.Millisecond)
	interval := time.Since(now).Milliseconds()

	SendHttpMessage(logger, helper.MockResponse(), helper.MockRequestWithJson(), now.Unix()*int64(time.Millisecond), interval, nil)

	assert.Contains(t, logger.queue[0], "[\"now", "SendHttpMessage did not append 'now' to message on manual entry")
	assert.Contains(t, logger.queue[0], "[\"interval\",\"", "SendHttpMessage did not append 'interval' to message on manual entry")
}

func TestStop(t *testing.T) {
	helper := newTestHelper()

	opt := Options{
		Queue:   make([]string, 0),
		Enabled: true,
	}

	logger, _ := NewHttpLogger(opt)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for logger.Enabled() {
				SendHttpMessage(logger, helper.MockResponse(), helper.MockRequestWithJson(), 0, 0, nil)
			}
		}()
	}

	logger.Stop()
	wg.Wait()

	assert.False(t, logger.Enabled())

	queueLen := len(logger.Queue())
	SendHttpMessage(logger, helper.MockResponse(), helper.MockRequestWithJson(), 0, 0, nil)
	assert.Equal(t, queueLen, len(logger.Queue()))
}
