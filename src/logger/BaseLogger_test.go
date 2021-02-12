package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	//TODO
}

func TestHasValidVersion(t *testing.T) {
	//TODO
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
