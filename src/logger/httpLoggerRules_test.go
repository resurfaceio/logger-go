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
	testLogger1 := baseClientLogger{url: "http://mysote.com"}
	assert.Equal(t, testLogger1.enableable, false, "Logger enableable flag should be set to false")

	testLogger2 := baseClientLogger{url: "http://mysote.com", rules: ""}
	assert.Equal(t, testLogger2.enableable, false, "Logger enableable flag should be set to false")

	testLogger3 := baseClientLogger{url: "https://mysote.com"}
	assert.Equal(t, testLogger3.enableable, true, "Logger enableable flag should be set to true")

	testLogger4 := baseClientLogger{url: "http://mysote.com", rules: "allow_http_url"}
	assert.Equal(t, testLogger4.enableable, true, "Logger enableable flag should be set to true")

	testLogger5 := baseClientLogger{url: "http://mysote.com", rules: "allow_http_url\nallow_http_url"}
	assert.Equal(t, testLogger5.enableable, true, "Logger enableable flag should be set to true")
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
