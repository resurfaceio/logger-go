package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"strings"
)

func TestBaseLogger(t *testing.T) {
	logger := newLogger()
	assert.NotNil(t, logger)
	//assert.Equal(logger.);
}

func TestCreatesInstance(t *testing.T) {
	logger := NewBaseLoggerAgent(MOCK_AGENT)
	assert.NotNil(t, logger)
	assert.Equal(t, logger.getAgent(), MOCK_AGENT)
	assert.false(t, logger.getEnableable())
	assert.false(t, logger.getEnableabled())
	assert.Nil(t, logger.getQueue())
	assert.Nil(t, logger.getUrl)
}

func TestCreatesMultipleInstances(t *testing.T) {
	//TODO
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
	assert.NotContains(t,version,"\\")
	assert.NotContains(t,version,"\"")
	assert.NotContains(t,version,"'")
	assert.Equal(t,NewBaseLoggerAgent(MOCK_AGENT).getVersion(),version)
}

func TestPerformsEnablingWhenExpected(t *testing.T) {
	//TODO
}

func TestSkipsEnablingForInvalidUrls(t *testing.T) {
	//TODO
}

func TestSkipsEnablingForNullUrls(t *testing.T) {
	//TODO
}

func TestSubmitsToDemoUrl(t *testing.T) {
	//TODO
}

func TestSubmitsToDemoUrlViaHttp(t *testing.T) {
	//TODO
}

func TestSubmitsToDemoUrlWihoutCompression(t *testing.T) {
	//TODO
}

func TestSubmitsToDeniedUrl(t *testing.T) {
	//TODO
}

func TestSubmitsToQueue(t *testing.T) {
	//TODO
}

func TestUsesSkipOptions(t *testing.T) {
	//TODO
}
