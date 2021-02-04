package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Still need to import logger above

// test override default rules

func TestOverrideDefaultRules() {

}

// test uses allow http url rules

func TestUsesAllowHttpUrlRules(t *testing.T) {
	// requires url, rules, and enableable to be in logger struct !!!
	testLogger := newLogger()
	testLogger.url = "http://mysite.com"
	assert.Equal(t, false, testLogger.enableable, "Logger enableable flag should be set to false")

	testLogger = newLogger()
	testLogger.url = "http://mysite.com"
	testLogger.rules = ""
	assert.Equal(t, false, testLogger.enableable, "Logger enableable flag should be set to false")

	testLogger = newLogger()
	testLogger.url = "https://mysite.com"
	assert.Equal(t, true, testLogger.enableable, "Logger enableable flag should be set to true")

	testLogger = newLogger()
	testLogger.url = "http://mysite.com"
	testLogger.rules = "allow_http_url"
	assert.Equal(t, true, testLogger.enableable, "Logger enableable flag should be set to true")

	testLogger = newLogger()
	testLogger.url = "http://mysite.com"
	testLogger.rules = "allow_http_url\nallow_http_url"
	assert.Equal(t, true, testLogger.enableable, "Logger enableable flag should be set to true")
}

// test uses remove rules

// test uses remove if rules

// test uses remove if found rules

// test uses remove unless rules

// test uses remove unless found rules

// test uses replace rules

// test uses replace rules with complex expressions

// test uses sample rules

// test uses skip compression rules

// test uses skip submission rules

// test uses stop rules

// test uses stop if rules

// test uses stop if found rules

// test uses stop unless rules

// test uses stop unless found rules
