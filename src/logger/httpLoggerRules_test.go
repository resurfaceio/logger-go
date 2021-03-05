package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test override default rules

func TestOverrideDefaultRules(t *testing.T) {
	assert.Equal(t, httpRules.StrictRules, httpRules.DefaultRules, "HTTP default rules are not strict rules")

	logger := NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	assert.Equal(t, httpRules.StrictRules, logger.rules.Text, "logger rules are not set to default rules")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules("# 123")
	assert.Equal(t, "# 123", logger.rules.Text, "logger default rules not overriden")

	httpRules.SetDefaultRules("")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	assert.Equal(t, "", logger.rules.Text, "logger default rules were not applied")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules("   ")
	assert.Equal(t, "", logger.rules.Text, "logger default rules not overriden or blank space not ignored")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules(" sample 42")
	assert.Equal(t, " sample 42", logger.rules.Text, "logger default rules not overriden")

	httpRules.SetDefaultRules("skip_compression")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	assert.Equal(t, "skip_compression", logger.rules.Text, "logger default rules not applied")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules("include default\nskip_submission\n")
	assert.Equal(t, "include default\nskip_submission\n", logger.rules.Text, ":logger default rules not overriden")

	httpRules.SetDefaultRules("sample 42\n")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	assert.Equal(t, "sample 42\n", logger.rules.Text, "logger default rules not applied")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules("   ")
	assert.Equal(t, "sample 42\n", logger.rules.Text, "white space not ignored")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules("include default\nskip_submission\n")
	assert.Equal(t, "sample 42\n\nskip_submission\n", logger.rules.Text, "logger rules not applied correctly")

	httpRules.SetDefaultRules("inlude debug")
	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	logger.SetRules(httpRules.StrictRules)
	assert.Equal(t, httpRules.StrictRules, logger.rules.text, "logger default rules not overriden")

	httpRules.SetDefaultRules(httpRules.StrictRules)
}

// test uses allow http url rules

func TestUsesAllowHttpUrlRules(t *testing.T) {
	// requires url, rules, and enableable to be in logger struct !!!
	logger := NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	assert.Equal(t, false, logger.Enableable(), "Logger enableable flag should be set to false")

	logger = NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	logger.SetRules("")
	assert.Equal(t, false, logger.Enableable(), "Logger enableable flag should be set to false")

	logger = NewHttpLogger()
	logger.SetUrl("https://mysite.com")
	assert.Equal(t, true, logger.Enableable(), "Logger enableable flag should be set to true")

	logger = NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	logger.SetRules("allow_http_url")
	assert.Equal(t, true, logger.Enableable(), "Logger enableable flag should be set to true")

	logger = NewHttpLogger()
	logger.SetUrl("http://mysite.com")
	logger.SetRules("allow_http_url\nallow_http_url")
	assert.Equal(t, true, logger.Enableable(), "Logger enableable flag should be set to true")
}

// test uses copy session field rules test

// func TestUsesCopySessionFieldRules(t *testing.T) {

// }

// test uses copy session field and remove rules test

// test uses copy session field and stop rules test

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
