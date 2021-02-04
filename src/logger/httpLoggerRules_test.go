package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Still need to import

// test override default rules

func TestOverrideDefaultRules(t *testing.T) {
	assert.Equal(t, httpRules.strictRules, httpRules.defaultRules, "HTTP default rules are not strict rules")

	logger := newLogger()
	logger.url = "https://mysite.com"
	assert.Equal(t, httpRules.strictRules, logger.rules.text, "logger rules are not set to default rules")
	logger = newLogger()
	logger.url = "https://mysite.com"
	logger.rules = "# 123"
	assert.Equal(t, "# 123", logger.rules.text, "logger default rules not overriden")

	httpRules.setDefaultRules("")
	logger = newLogger()
	logger.url = "https://mysite.com"
	assert.Equal(t, "", logger.rules.text, "logger default rules were not applied")
	logger = newLogger()
	logger.url = "https://mysite.com"
	logger.rules = "   "
	assert.Equal(t, "", logger.rules.text, "logger default rules not overriden or blank space not ignored")
	logger = newLogger()
	logger.url = "https://mysite.com"
	logger.rules = " sample 42"
	assert.Equal(t, " sample 42", logger.rules.text, "logger default rules not overriden")

	httpRules.setDefaultRules("skip_compression")
	logger = newLogger()
	logger.url = "https://mysite.com"
	assert.Equal(t, "skip_compression", logger.rules.text, "logger default rules not applied")
	logger = newLogger()
	logger.url = "https://mysite.com"
	logger.rules = "include default\nskip_submission\n"
	assert.Equal(t, "include default\nskip_submission\n", logger.rules.text, ":logger default rules not overriden")

	httpRules.setDefaultRules("sample 42\n")
	logger = newLogger()
	logger.url = "https://mysite.com"
	assert.Equal(t, "sample 42\n", logger.rules.text, "logger default rules not applied")
	logger = newLogger()
	logger.url = "https://mysite.com"
	logger.rules = "   "
	assert.Equal(t, "sample 42\n", logger.rules.text, "white space not ignored")
	logger = newLogger()
	logger.url = "https://mysite.com"
	logger.rules = "include default\nskip_submission\n"
	assert.Equal(t, "sample 42\n\nskip_submission\n", logger.rules.text, "logger rules not applied correctly")

	httpRules.setDefaultRules("inlude debug")
	logger = newLogger()
	logger.url = "https://mysite.com"
	logger.rules = httpRules.strictRules
	assert.Equal(t, httpRules.strictRules, logger.rules.text, "logger default rules not overriden")

	httpRules.setDefaultRules(httpRules.strictRules)
}

// test uses allow http url rules

func TestUsesAllowHttpUrlRules(t *testing.T) {
	// requires url, rules, and enableable to be in logger struct !!!
	logger := newLogger()
	logger.url = "http://mysite.com"
	assert.Equal(t, false, logger.enableable, "Logger enableable flag should be set to false")

	logger = newLogger()
	logger.url = "http://mysite.com"
	logger.rules = ""
	assert.Equal(t, false, logger.enableable, "Logger enableable flag should be set to false")

	logger = newLogger()
	logger.url = "https://mysite.com"
	assert.Equal(t, true, logger.enableable, "Logger enableable flag should be set to true")

	logger = newLogger()
	logger.url = "http://mysite.com"
	logger.rules = "allow_http_url"
	assert.Equal(t, true, logger.enableable, "Logger enableable flag should be set to true")

	logger = newLogger()
	logger.url = "http://mysite.com"
	logger.rules = "allow_http_url\nallow_http_url"
	assert.Equal(t, true, logger.enableable, "Logger enableable flag should be set to true")
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
