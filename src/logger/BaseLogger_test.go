package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"strings"
)

func TestCreatesInstance(t *testing.T) {
	logger := NewBaseLoggerAgent(MOCK_AGENT)
	assert.NotNil(t, logger)
	assert.Equal(t, logger.Agent, MOCK_AGENT)
	assert.false(t, logger.Enableable)
	assert.false(t, logger.Enabled)
	assert.Nil(t, logger.Queue)
	assert.Nil(t, logger.Url)
}

func TestCreatesMultipleInstances(t *testing.T) {
	agent1 := "agent1"
	agent2 := "agent2"
	agent3 := "agent3"
	url1 := "https://resuface.io"
	url2 := "https://whatever.com"
	helper := NewTestHelper()
	logger1 := NewBaseLoggerAgentUrl(agent1, url1)
	logger2 := NewBaseLoggerAgentUrl(agent2, url2)
	logger3 := NewBaseLoggerAgentUrl(agent3, helper.demoURL)

	assert.Equal(t, logger1.Agent, agent1)
	assert.True(t, logger1.Enableable)
	assert.True(t, logger1.Enabled)
	assert.Equal(t, logger1.Url, url1)

	assert.Equal(t, logger2.Agent, agent2)
	assert.True(t, logger2.Enableable)
	assert.True(t, logger2.Enabled)
	assert.Equal(t, logger2.Url, url2)

	assert.Equal(t, logger3.Agent, agent3)
	assert.True(t, logger3.Enableable)
	assert.True(t, logger3.Enabled)
	assert.Equal(t, logger3.Url, helper.demoURL)

	//TODO: implement usage loggers
	//UsageLoggers.disable();
	//assert.False(t,UsageLoggers.Enabled())
	assert.False(t, logger1.Enabled)
	assert.False(t, logger2.Enabled)
	assert.False(t, logger3.Enabled)
	//UsageLoggers.enable();
	//assert.True(t,UsageLoggers.Enabled())
	assert.True(t, logger1.Enabled)
	assert.True(t, logger2.Enabled)
	assert.True(t, logger3.Enabled)
}

func TestHasValidHost(t *testing.T) {
	helper := NewTestHelper()
	baseLogger := NewBaseLoggerAgent(helper.mockAgent)
	host := baseLogger.hostLookup()
	assert.NotNil(t, host)
	assert.Greater(t, len(host), 0)
	assert.NotContains(t, host, "unknown")
}

func TestHasValidVersion(t *testing.T) {
	helper := NewTestHelper()
	baseLogger := NewBaseLoggerAgent(helper.mockAgent)
	version := baseLogger.versionLookup()
	assert.NotNil(t, version)
	assert.Greater(t, len(version), 0)
	//replacement of the java "startswith" assertion
	//won't work rn since version is a dummy string
	assert.Equal(t, strings.Index(version, "2.0."), 0)
	assert.NotContains(t, version, "\\")
	assert.NotContains(t, version, "\"")
	assert.NotContains(t, version, "'")
	assert.Equal(t, baseLogger.Version, version)
}

func TestPerformsEnablingWhenExpected(t *testing.T) {
	helper := NewTestHelper()
	logger := NewBaseLogger(helper.mockAgent, helper.demoURL, false, nil)
	assert.True(t, logger.Enableable)
	assert.False(t, logger.Enabled)
	assert.Equal(t, logger.Url, helper.demoURL)
	logger.Enable()
	assert.True(t, logger.Enabled)

	queue := []string{}
	logger = NewBaseLoggerAgentQueueEnabled(helper.mockAgent, queue, false)
	assert.True(t, logger.Enableable)
	assert.False(t, logger.Enabled)
	assert.Nil(t, logger.Url)
	logger.Enable()
	logger.Disable()
	logger.Enable()
	assert.True(t, logger.Enabled)
}

//needs some more stuff in the baselogger class for this to compile
func TestSkipsEnablingForInvalidUrls(t *testing.T) {
	helper := NewTestHelper()
	for _, url := range helper.mockURLSinvalid {
		logger := NewBaseLogger(MOCK_AGENT, url, false, nil)
		assert.False(t, logger.Enableable)
		assert.False(t, logger.Enabled)
		assert.Nil(t, logger.Url)
		logger.Enable()
		logger.Disable()
		logger.Enable()
		assert.False(t, logger.Enabled)
	}
}

func TestSkipsEnablingForNullUrls(t *testing.T) {
	url := ""
	helper := NewTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, url)
	assert.False(t, logger.Enableable)
	assert.False(t, logger.Enabled)
	assert.Nil(t, logger.Url)
	logger.Enable()
	assert.False(t, logger.Enabled)
}

func TestSubmitsToDemoUrl(t *testing.T) {
	helper := NewTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, helper.demoURL)
	message := [][]string{}
	message = append(message, []string{"agent", logger.Agent})
	message = append(message, []string{"version", logger.Version})
	message = append(message, []string{"now", string(helper.mockNow)})
	message = append(message, []string{"protocol", "https"})
	//TODO: toimplement stringify method
	msg := "Json.stringify(message)"
	//TODO: implement parsable method
	assert.True(t, "parsable(msg)")
	logger.Submit(msg)
	assert.Equal(t, logger.SubmitFailures, 0)
	assert.Equal(t, logger.SubmitSuccesses, 1)
}

func TestSubmitsToDemoUrlViaHttp(t *testing.T) {
	helper := NewTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, strings.Replace(helper.demoURL, "https://", "http://", 1))
	assert.Equal(t, strings.Index(logger.Url, "http://"), 0)
	message := [][]string{}
	message = append(message, []string{"agent", logger.Agent})
	message = append(message, []string{"version", logger.Version})
	message = append(message, []string{"now", string(helper.mockNow)})
	message = append(message, []string{"protocol", "http"})
	//TODO: toimplement stringify method
	msg := "Json.stringify(message)"
	//TODO: implement parsable method
	assert.True(t, "parsable(msg)")
	logger.Submit(msg)
	assert.Equal(t, logger.SubmitFailures, 0)
	assert.Equal(t, logger.SubmitSuccesses, 1)
}

func TestSubmitsToDemoUrlWihoutCompression(t *testing.T) {
	helper := NewTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, helper.demoURL)
	logger.SkipCompression = true
	assert.True(t, logger.SkipCompression)
	message := [][]string{}
	message = append(message, []string{"agent", logger.Agent})
	message = append(message, []string{"version", logger.Version})
	message = append(message, []string{"now", string(helper.mockNow)})
	message = append(message, []string{"protocol", "https"})
	message = append(message, []string{"skip_compression", "true"})
	//TODO: toimplement stringify method
	msg := "Json.stringify(message)"
	//TODO: implement parsable method
	assert.True(t, "parsable(msg)")
	logger.Submit(msg)
	assert.Equal(t, logger.SubmitFailures, 0)
	assert.Equal(t, logger.SubmitSuccesses, 1)
}

func TestSubmitsToDeniedUrl(t *testing.T) {
	helper := NewTestHelper()
	for _, url := range helper.mockURLSinvalid {
		logger := NewBaseLoggerAgentUrl(helper.mockAgent, url)
		assert.True(t, logger.Enableable)
		assert.True(t, logger.Enabled)
		logger.Submit("{}")
		assert.Equal(t, logger.SubmitFailures, 1)
		assert.Equal(t, logger.SubmitSuccesses, 0)
	}
}

func TestSubmitsToQueue(t *testing.T) {
	helper := NewTestHelper()
	queue := []string{}
	logger := NewBaseLoggerAgentQueue(helper.mockAgent, queue)
	assert.Equal(t, logger.Queue, queue)
	assert.Nil(t, logger.Url)
	assert.True(t, logger.Enableable)
	assert.True(t, logger.Enabled)
	assert.Equal(t, len(queue), 0)
	logger.Submit("{}")
	assert.Equal(t, len(queue), 1)
	logger.Submit("{}")
	assert.Equal(t, len(queue), 2)
	assert.Equal(logger.SubmitFailures, 0)
	assert.Equal(logger.SubmitSuccesses, 0)
}

func TestUsesSkipOptions(t *testing.T) {
	helper := NewTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, helper.demoURL)
	assert.False(t, logger.SkipCompression)
	assert.False(t, logger.SkipSubmission)

	logger.SkipCompression = true
	assert.True(t, logger.SkipCompression)
	assert.False(t, logger.SkipSubmission)

	logger.SkipCompression = false
	logger.SkipSubmission = true
	assert.False(t, logger.SkipCompression)
	assert.True(t, logger.SkipSubmission)
}
