// © 2016-2021 Resurface Labs Inc.

package logger

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"strings"

	"encoding/json"
)

func TestCreatesInstance(t *testing.T) {
	helper := newTestHelper()
	logger := newBaseLogger(helper.mockAgent, "", false, nil)
	assert.NotNil(t, logger)
	assert.Equal(t, helper.mockAgent, logger.agent)
	assert.False(t, logger.enableable)
	assert.False(t, logger.Enabled())
	assert.Nil(t, logger.queue)
	assert.Equal(t, "", logger.url)
}

func TestCreatesMultipleInstances(t *testing.T) {
	agent1 := "agent1"
	agent2 := "agent2"
	agent3 := "agent3"
	url1 := "https://resuface.io"
	url2 := "https://whatever.com"
	helper := newTestHelper()
	logger1 := newBaseLogger(agent1, url1, true, nil)
	logger2 := newBaseLogger(agent2, url2, true, nil)
	logger3 := newBaseLogger(agent3, helper.demoURL, true, nil)

	assert.Equal(t, agent1, logger1.agent)
	assert.True(t, logger1.enableable)
	assert.True(t, logger1.Enabled())
	assert.Equal(t, url1, logger1.url)

	assert.Equal(t, agent2, logger2.agent)
	assert.True(t, logger2.enableable)
	assert.True(t, logger2.Enabled())
	assert.Equal(t, url2, logger2.url)

	assert.Equal(t, agent3, logger3.agent)
	assert.True(t, logger3.enableable)
	assert.True(t, logger3.Enabled())
	assert.Equal(t, helper.demoURL, logger3.url)

	usageLoggers, _ := GetUsageLoggers()
	usageLoggers.Disable()
	assert.False(t, usageLoggers.IsEnabled())
	assert.False(t, logger1.Enabled())
	assert.False(t, logger2.Enabled())
	assert.False(t, logger3.Enabled())
	usageLoggers.Enable()
	assert.True(t, usageLoggers.IsEnabled())
	assert.True(t, logger1.Enabled())
	assert.True(t, logger2.Enabled())
	assert.True(t, logger3.Enabled())
}

func TestHasValidHost(t *testing.T) {
	helper := newTestHelper()
	baseLogger := newBaseLogger(helper.mockAgent, "", true, nil)
	host := baseLogger.host
	assert.NotNil(t, host)
	assert.Greater(t, len(host), 0)
	assert.NotContains(t, host, "unknown")
}

func TestHasValidVersion(t *testing.T) {
	helper := newTestHelper()
	baseLogger := newBaseLogger(helper.mockAgent, "", true, nil)
	version := versionLookup()
	assert.NotNil(t, version)
	assert.Greater(t, len(version), 0)
	//replacement of the java "startswith" assertion
	//won't work rn since version is a dummy string
	assert.True(t, strings.HasPrefix(version, "3."))
	assert.NotContains(t, version, "\\")
	assert.NotContains(t, version, "\"")
	assert.NotContains(t, version, "'")
	assert.Equal(t, baseLogger.version, version)
}

func TestPerformsEnablingWhenExpected(t *testing.T) {
	helper := newTestHelper()
	logger := newBaseLogger(helper.mockAgent, helper.demoURL, false, nil)
	assert.True(t, logger.enableable)
	assert.False(t, logger.Enabled())
	assert.Equal(t, helper.demoURL, logger.url)
	logger.Enable()
	assert.True(t, logger.enabled)
	assert.True(t, logger.Enabled())

	queue := []string{}
	logger = newBaseLogger(helper.mockAgent, "", false, queue)
	assert.True(t, logger.enableable)
	assert.False(t, logger.Enabled())
	assert.Equal(t, "", logger.url)
	logger.Enable()
	assert.True(t, logger.Enabled())
}

//needs some more stuff in the baselogger class for this to compile
func TestSkipsEnablingForInvalidUrls(t *testing.T) {
	helper := newTestHelper()
	for _, url := range helper.mockURLSinvalid {
		logger := newBaseLogger(helper.mockAgent, url, false, nil)
		assert.False(t, logger.enableable)
		assert.False(t, logger.Enabled())
		assert.Equal(t, "", logger.url)
		logger.Enable()
		assert.False(t, logger.Enabled())
	}
}

func TestSkipsEnablingForEmptyUrl(t *testing.T) {
	url := ""
	helper := newTestHelper()
	logger := newBaseLogger(helper.mockAgent, url, false, nil) // should this be false or true because it matters for the test with enabled bool
	assert.False(t, logger.enableable)
	assert.False(t, logger.Enabled())
	assert.Equal(t, "", logger.url)
	logger.Enable()
	assert.False(t, logger.Enabled())
}

func TestSubmitsToDemoUrl(t *testing.T) {
	helper := newTestHelper()
	queue := []string{}
	logger := newBaseLogger(helper.mockAgent, helper.demoURL, true, queue)
	message := [][]string{}
	message = append(message, []string{"agent", logger.agent})
	message = append(message, []string{"version", logger.version})
	message = append(message, []string{"now", string(fmt.Sprint(helper.mockNow))})
	message = append(message, []string{"protocol", "https"})
	marshalled, _ := json.Marshal(message)
	logger.ndjsonHandler(string(marshalled))
	assert.Equal(t, int64(0), logger.submitFailures)
	assert.Equal(t, int64(1), logger.submitSuccesses)
}

func TestSubmitsToDemoUrlViaHttp(t *testing.T) {
	helper := newTestHelper()
	queue := []string{}
	logger := newBaseLogger(helper.mockAgent, strings.Replace(helper.demoURL, "https://", "http://", 1), true, queue)
	assert.Equal(t, 0, strings.Index(logger.url, "http://"))
	message := [][]string{}
	message = append(message, []string{"agent", logger.agent})
	message = append(message, []string{"version", logger.version})
	message = append(message, []string{"now", string(fmt.Sprint(helper.mockNow))})
	message = append(message, []string{"protocol", "http"})
	marshalled, _ := json.Marshal(message)
	logger.ndjsonHandler(string(marshalled))
	assert.Equal(t, int64(0), logger.submitFailures)
	assert.Equal(t, int64(1), logger.submitSuccesses)
}

func TestSubmitsToDemoUrlWihoutCompression(t *testing.T) {
	helper := newTestHelper()
	queue := []string{}
	logger := newBaseLogger(helper.mockAgent, helper.demoURL, true, queue)
	logger.skipCompression = true
	assert.True(t, logger.skipCompression)
	message := [][]string{}
	message = append(message, []string{"agent", logger.agent})
	message = append(message, []string{"version", logger.version})
	message = append(message, []string{"now", string(fmt.Sprint(helper.mockNow))})
	message = append(message, []string{"protocol", "https"})
	message = append(message, []string{"skip_compression", "true"})
	marshalled, _ := json.Marshal(message)
	logger.ndjsonHandler(string(marshalled))
	assert.Equal(t, int64(0), logger.submitFailures)
	assert.Equal(t, int64(1), logger.submitSuccesses)
}

func TestSubmitsToDeniedUrl(t *testing.T) {
	helper := newTestHelper()
	for _, url := range helper.mockURLSdenied {
		logger := newBaseLogger(helper.mockAgent, url, true, nil)
		assert.True(t, logger.enableable)
		assert.True(t, logger.Enabled())
		logger.ndjsonHandler("{}")
		time.Sleep(5 * time.Second) // Added because with async worker the test was checking for fail/success values before worker could actually try sending the request.
		fmt.Print("\n***************")
		fmt.Print(logger.submitFailures)
		fmt.Print(" | ")
		fmt.Print(logger.submitSuccesses)
		fmt.Print("***************\n")
		assert.Equal(t, int64(1), logger.submitFailures)
		assert.Equal(t, int64(0), logger.submitSuccesses)
	}
}

func TestSubmitsToQueue(t *testing.T) {
	helper := newTestHelper()
	queue := []string{}
	logger := newBaseLogger(helper.mockAgent, helper.mockURLSdenied[0], true, queue)
	assert.Equal(t, queue, logger.queue)
	assert.Equal(t, helper.mockURLSdenied[0], logger.url)
	assert.True(t, logger.enableable)
	assert.True(t, logger.Enabled())
	assert.Equal(t, 0, len(logger.queue))
	logger.ndjsonHandler("{}")
	assert.Equal(t, 1, len(logger.queue))
	logger.ndjsonHandler("{}")
	assert.Equal(t, 2, len(logger.queue))
	assert.Equal(t, int64(0), logger.submitFailures)
	assert.Equal(t, int64(2), logger.submitSuccesses)
}

func TestUsesSkipOptions(t *testing.T) {
	helper := newTestHelper()
	logger := newBaseLogger(helper.mockAgent, helper.demoURL, true, nil)
	assert.False(t, logger.skipCompression)
	assert.False(t, logger.skipSubmission)

	logger.skipCompression = true
	assert.True(t, logger.skipCompression)
	assert.False(t, logger.skipSubmission)

	logger.skipCompression = false
	logger.skipSubmission = true
	assert.False(t, logger.skipCompression)
	assert.True(t, logger.skipSubmission)
}
