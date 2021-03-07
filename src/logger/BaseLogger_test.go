package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"strings"

	"encoding/json"
)

func TestCreatesInstance(t *testing.T) {
	helper := GetTestHelper()
	logger := NewBaseLoggerAgent(helper.mockAgent)
	assert.NotNil(t, logger)
	assert.Equal(t, helper.mockAgent, logger.Agent)
	assert.false(t, logger.Enableable())
	assert.false(t, logger.Enabled())
	assert.Nil(t, logger.Queue())
	assert.Nil(t, logger.Url())
}

func TestCreatesMultipleInstances(t *testing.T) {
	agent1 := "agent1"
	agent2 := "agent2"
	agent3 := "agent3"
	url1 := "https://resuface.io"
	url2 := "https://whatever.com"
	helper := GetTestHelper()
	logger1 := NewBaseLoggerAgentUrl(agent1, url1)
	logger2 := NewBaseLoggerAgentUrl(agent2, url2)
	logger3 := NewBaseLoggerAgentUrl(agent3, helper.demoURL)

	assert.Equal(t, agent1, logger1.Agent())
	assert.True(t, logger1.Enableable())
	assert.True(t, logger1.Enabled())
	assert.Equal(t, url1, logger1.Url())

	assert.Equal(t, agent2, logger2.Agent())
	assert.True(t, logger2.Enableable())
	assert.True(t, logger2.Enabled())
	assert.Equal(t, url2, logger2.Url())

	assert.Equal(t, agent3, logger3.Agent())
	assert.True(t, logger3.Enableable())
	assert.True(t, logger3.Enabled())
	assert.Equal(t, helper.demoURL, logger3.Url())

	//TODO: implement usage loggers
	//UsageLoggers.disable();
	//assert.False(t,UsageLoggers.Enabled())
	assert.False(t, logger1.Enabled())
	assert.False(t, logger2.Enabled())
	assert.False(t, logger3.Enabled())
	//UsageLoggers.enable();
	//assert.True(t,UsageLoggers.Enabled())
	assert.True(t, logger1.Enabled())
	assert.True(t, logger2.Enabled())
	assert.True(t, logger3.Enabled())
}

func TestHasValidHost(t *testing.T) {
	helper := GetTestHelper()
	baseLogger := NewBaseLoggerAgent(helper.mockAgent)
	host := baseLogger.hostLookup()
	assert.NotNil(t, host)
	assert.Greater(t, len(host), 0)
	assert.NotContains(t, host, "unknown")
}

func TestHasValidVersion(t *testing.T) {
	helper := GetTestHelper()
	baseLogger := NewBaseLoggerAgent(helper.mockAgent)
	version := versionLookup()
	assert.NotNil(t, version)
	assert.Greater(t, len(version), 0)
	//replacement of the java "startswith" assertion
	//won't work rn since version is a dummy string
	assert.Equal(t, strings.Index(version, "2.0."), 0)
	assert.NotContains(t, version, "\\")
	assert.NotContains(t, version, "\"")
	assert.NotContains(t, version, "'")
	assert.Equal(t, baseLogger.Version(), version)
}

func TestPerformsEnablingWhenExpected(t *testing.T) {
	helper := GetTestHelper()
	logger := NewBaseLogger(helper.mockAgent, helper.demoURL, false, nil)
	assert.True(t, logger.Enableable())
	assert.False(t, logger.Enabled())
	assert.Equal(t, helper.demoURL, logger.Url())
	logger.Enable()
	assert.True(t, logger.Enabled())

	queue := []string{}
	logger = NewBaseLoggerAgentQueueEnabled(helper.mockAgent, queue, false)
	assert.True(t, logger.Enableable())
	assert.False(t, logger.Enabled())
	assert.Nil(t, logger.Url())
	logger.Enable()
	logger.Disable()
	logger.Enable()
	assert.True(t, logger.Enabled())
}

//needs some more stuff in the baselogger class for this to compile
func TestSkipsEnablingForInvalidUrls(t *testing.T) {
	helper := GetTestHelper()
	for _, url := range helper.mockURLSinvalid {
		logger := NewBaseLogger(helper.mockAgent, url, false, nil)
		assert.False(t, logger.Enableable())
		assert.False(t, logger.Enabled())
		assert.Nil(t, logger.Url())
		logger.Enable()
		logger.Disable()
		logger.Enable()
		assert.False(t, logger.Enabled())
	}
}

func TestSkipsEnablingForNullUrl(t *testing.T) {
	url := ""
	helper := GetTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, url)
	assert.False(t, logger.Enableable())
	assert.False(t, logger.Enabled())
	assert.Nil(t, logger.Url())
	logger.Enable()
	assert.False(t, logger.Enabled())
}

func TestSubmitsToDemoUrl(t *testing.T) {
	helper := GetTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, helper.demoURL)
	message := [][]string{}
	message = append(message, []string{"agent", logger.Agent()})
	message = append(message, []string{"version", logger.Version()})
	message = append(message, []string{"now", string(helper.mockNow)})
	message = append(message, []string{"protocol", "https"})
	msg, err := json.Marshal(message)
	assert.True(t, err == nil)
	logger.Submit(string(msg))
	assert.Equal(t, 0, logger.SubmitFailures())
	assert.Equal(t, 1, logger.SubmitSuccesses())
}

func TestSubmitsToDemoUrlViaHttp(t *testing.T) {
	helper := GetTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, strings.Replace(helper.demoURL, "https://", "http://", 1))
	assert.Equal(t, 0, strings.Index(logger.Url(), "http://"))
	message := [][]string{}
	message = append(message, []string{"agent", logger.Agent()})
	message = append(message, []string{"version", logger.Version()})
	message = append(message, []string{"now", string(helper.mockNow)})
	message = append(message, []string{"protocol", "http"})
	msg, err := json.Marshal(message)
	assert.True(t, err == nil)
	logger.Submit(string(msg))
	assert.Equal(t, 0, logger.SubmitFailures())
	assert.Equal(t, 1, logger.SubmitSuccesses())
}

func TestSubmitsToDemoUrlWihoutCompression(t *testing.T) {
	helper := GetTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, helper.demoURL)
	logger.SetSkipCompression(true)
	assert.True(t, logger.SkipCompression())
	message := [][]string{}
	message = append(message, []string{"agent", logger.Agent()})
	message = append(message, []string{"version", logger.Version()})
	message = append(message, []string{"now", string(helper.mockNow)})
	message = append(message, []string{"protocol", "https"})
	message = append(message, []string{"skip_compression", "true"})
	msg, err := json.Marshal(message)
	assert.True(t, err == nil)
	logger.Submit(string(msg))
	assert.Equal(t, 0, logger.SubmitFailures)
	assert.Equal(t, 1, logger.SubmitSuccesses)
}

func TestSubmitsToDeniedUrl(t *testing.T) {
	helper := GetTestHelper()
	for _, url := range helper.mockURLSdenied {
		logger := NewBaseLoggerAgentUrl(helper.mockAgent, url)
		assert.True(t, logger.Enableable())
		assert.True(t, logger.Enabled())
		logger.Submit("{}")
		assert.Equal(t, 1, logger.SubmitFailures())
		assert.Equal(t, 0, logger.SubmitSuccesses())
	}
}

func TestSubmitsToQueue(t *testing.T) {
	helper := GetTestHelper()
	queue := []string{}
	logger := NewBaseLoggerAgentQueue(helper.mockAgent, queue)
	assert.Equal(t, queue, logger.Queue())
	assert.Nil(t, logger.Url())
	assert.True(t, logger.Enableable())
	assert.True(t, logger.Enabled())
	assert.Equal(t, 0, len(queue))
	logger.Submit("{}")
	assert.Equal(t, 1, len(queue))
	logger.Submit("{}")
	assert.Equal(t, 2, len(queue))
	assert.Equal(t, 0, logger.SubmitFailures())
	assert.Equal(t, 0, logger.SubmitSuccesses())
}

func TestUsesSkipOptions(t *testing.T) {
	helper := GetTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, helper.demoURL)
	assert.False(t, logger.SkipCompression())
	assert.False(t, logger.SkipSubmission())

	logger.SetSkipCompression(true)
	assert.True(t, logger.SkipCompression)
	assert.False(t, logger.SkipSubmission)

	logger.SetSkipCompression(false)
	logger.SetSkipSubmission(true)
	assert.False(t, logger.SkipCompression())
	assert.True(t, logger.SkipSubmission())
}
