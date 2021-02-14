package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"strings"
)

func TestCreatesInstance(t *testing.T) {
	logger := NewBaseLoggerAgent(MOCK_AGENT)
	assert.NotNil(t, logger)
	assert.Equal(t, logger.getAgent(), MOCK_AGENT)
	assert.false(t, logger.getEnableable())
	assert.false(t, logger.getEnabled())
	assert.Nil(t, logger.getQueue())
	assert.Nil(t, logger.getUrl)
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

	assert.Equal(t, logger1.getAgent(), agent1)
	assert.True(t, logger1.getEnableable())
	assert.True(t, logger1.getEnabled())
	assert.Equal(t, logger1.getUrl(), url1)

	assert.Equal(t, logger2.getAgent(), agent2)
	assert.True(t, logger2.getEnableable())
	assert.True(t, logger2.getEnabled())
	assert.Equal(t, logger2.getUrl(), url2)

	assert.Equal(t, logger3.getAgent(), agent3)
	assert.True(t, logger3.getEnableable())
	assert.True(t, logger3.getEnabled())
	assert.Equal(t, logger3.getUrl(), helper.demoURL)

	//TODO: implement usage loggers
	//UsageLoggers.disable();
	//assert.False(t,UsageLoggers.GetEnabled())
	assert.False(t,logger1.getEnabled())
	assert.False(t,logger2.getEnabled())
	assert.False(t,logger3.getEnabled())
	//UsageLoggers.enable();
	//assert.True(t,UsageLoggers.GetEnabled())
	assert.True(t,logger1.getEnabled())
	assert.True(t,logger2.getEnabled())
	assert.True(t,logger3.getEnabled())
}

func TestHasValidHost(t *testing.T) {
	host := host_lookup()
	assert.NotNil(t, host)
	assert.Greater(t, len(host), 0)
	assert.Contains(t, host, "unknown")
	assert.Equal(t, host, NewBaseLoggerAgent(MOCK_AGENT).getHost())
}

func TestHasValidVersion(t *testing.T) {
	version := version_lookup()
	assert.NotNil(t, version)
	assert.Greater(t, len(version), 0)
	//replacement of the java "startswith" assertion
	//won't work rn since version is a dummy string
	assert.Equal(t, strings.Index(version, "2.0."), 0)
	assert.NotContains(t, version, "\\")
	assert.NotContains(t, version, "\"")
	assert.NotContains(t, version, "'")
	assert.Equal(t, NewBaseLoggerAgent(MOCK_AGENT).getVersion(), version)
}

func TestPerformsEnablingWhenExpected(t *testing.T) {
	helper := NewTestHelper()
	logger := NewBaseLogger(helper.mockAgent, helper.demoURL, false, nil)
	assert.True(t, logger.getEnableable())
	assert.False(t, logger.getEnabled())
	assert.Equal(t, logger.getUrl(), helper.demoURL)
	logger.enable()
	assert.True(t, logger.getEnabled())

	queue := []string{}
	logger = NewBaseLoggerAgentQueueEnabled(helper.mockAgent, queue, false)
	assert.True(t, logger.getEnableable())
	assert.False(t, logger.getEnabled())
	assert.Nil(t, logger.getUrl())
	logger.enable()
	logger.disable()
	logger.enable()
	assert.True(t, logger.getEnabled())
}

//needs some more stuff in the baselogger class for this to compile
func TestSkipsEnablingForInvalidUrls(t *testing.T) {
	helper := NewTestHelper()
	for _, url := range helper.mockURLSinvalid {
		logger := NewBaseLogger(MOCK_AGENT, url, false, nil)
		assert.False(t, logger.getEnableable())
		assert.False(t, logger.getEnabled())
		assert.Nil(t, logger.getUrl())
		logger.enable()
		logger.disable()
		logger.enable()
		assert.False(t, logger.getEnabled())
	}
}

func TestSkipsEnablingForNullUrls(t *testing.T) {
	url := ""
	helper := NewTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, url)
	assert.False(t, logger.getEnableable())
	assert.False(t, logger.getEnabled())
	assert.Nil(t, logger.getUrl())
	logger.enable()
	assert.False(t, logger.getEnabled())
}

func TestSubmitsToDemoUrl(t *testing.T) {
	helper := NewTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, helper.demoURL)
	message := [][]string{}
	message = append(message, []string{"agent", logger.getAgent()})
	message = append(message, []string{"version", logger.getVersion()})
	message = append(message, []string{"now", string(helper.mockNow)})
	message = append(message, []string{"protocol", "https"})
	//TODO: toimplement stringify method
	msg := "Json.stringify(message)"
	//TODO: implement parsable method
	assert.True(t, "parsable(msg)")
	logger.submit(msg)
	assert.Equal(t, logger.getSubmitFailues(), 0)
	assert.Equal(t, logger.getSubmitSuccesses(), 1)
}

func TestSubmitsToDemoUrlViaHttp(t *testing.T) {
	helper := NewTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, strings.Replace(helper.demoURL, "https://", "http://", 1))
	assert.Equal(t, strings.Index(logger.getUrl(), "http://"), 0)
	message := [][]string{}
	message = append(message, []string{"agent", logger.getAgent()})
	message = append(message, []string{"version", logger.getVersion()})
	message = append(message, []string{"now", string(helper.mockNow)})
	message = append(message, []string{"protocol", "http"})
	//TODO: toimplement stringify method
	msg := "Json.stringify(message)"
	//TODO: implement parsable method
	assert.True(t, "parsable(msg)")
	logger.submit(msg)
	assert.Equal(t, logger.getSubmitFailues(), 0)
	assert.Equal(t, logger.getSubmitSuccesses(), 1)
}

func TestSubmitsToDemoUrlWihoutCompression(t *testing.T) {
	helper := NewTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, helper.demoURL)
	logger.setSkipCompression(true)
	assert.True(t, logger.getSkipCompression())
	message := [][]string{}
	message = append(message, []string{"agent", logger.getAgent()})
	message = append(message, []string{"version", logger.getVersion()})
	message = append(message, []string{"now", string(helper.mockNow)})
	message = append(message, []string{"protocol", "https"})
	message = append(message, []string{"skip_compression", "true"})
	//TODO: toimplement stringify method
	msg := "Json.stringify(message)"
	//TODO: implement parsable method
	assert.True(t, "parsable(msg)")
	logger.submit(msg)
	assert.Equal(t, logger.getSubmitFailues(), 0)
	assert.Equal(t, logger.getSubmitSuccesses(), 1)
}

func TestSubmitsToDeniedUrl(t *testing.T) {
	helper := NewTestHelper()
	for _, url := range helper.mockURLSinvalid {
		logger := NewBaseLoggerAgentUrl(helper.mockAgent, url)
		assert.True(t, logger.getEnableable())
		assert.True(t, logger.getEnabled())
		logger.submit("{}")
		assert.Equal(t, logger.getSubmitFailues(), 1)
		assert.Equal(t, logger.getSubmitSuccesses(), 0)
	}
}

func TestSubmitsToQueue(t *testing.T) {
	helper := NewTestHelper()
	queue := []string{}
	logger := NewBaseLoggerAgentQueue(helper.mockAgent, queue)
	assert.Equal(t, logger.getQueue(), queue)
	assert.Nil(t, logger.getUrl)
	assert.True(t, logger.getEnableable())
	assert.True(t, logger.getEnabled())
	assert.Equal(t, len(queue), 0)
	logger.submit("{}")
	assert.Equal(t, len(queue), 1)
	logger.submit("{}")
	assert.Equal(t, len(queue), 2)
	assert.Equal(logger.getSubmitFailues(), 0)
	assert.Equal(logger.getSubmitSuccesses(), 0)
}

func TestUsesSkipOptions(t *testing.T) {
	helper := NewTestHelper()
	logger := NewBaseLoggerAgentUrl(helper.mockAgent, helper.demoURL)
	assert.False(t, logger.getSkipCompression())
	assert.False(t, logger.getSkipSubmission())

	logger.setSkipCompression(true)
	assert.True(t, logger.getSkipCompression())
	assert.False(t, logger.getSkipSubmission())

	logger.setSkipCompression(false)
	logger.setSkipSubmission(true)
	assert.False(t, logger.getSkipCompression())
	assert.True(t, logger.getSkipSubmission())
}
