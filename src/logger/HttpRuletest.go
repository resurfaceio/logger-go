package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangesDefaultRules(t *testing.T) {
	testRules := NewHttpRules()
	assert.Equal(t, testRules.strictRules, testRules.defaultRules,
		"When initialized default rules should equal strict rules")

	testRules.SetDefaultRules("")
	assert.Equal(t, "", testRules.defaultRules,
		"Rules should be empty after initialized with empty string")

	testRule.SetDefaultRules(" include default")
	assert.Equal(t, "", testRules.defaultRules, "")

	testRules.SetDefaultRules("include default\n")
	assert.Equal(t, "", testRules.defaultRules, "")

	testRules.SetDefaultRules("include default\ninclude default\n")
	assert.Equal(t, 0, testRules.size, "")

	testRules.SetDefaultRules("include default\ninclude default\nsample 42")
	testRules2 := NewHttpRules()
	testRules2.SetDefaultRules(testRules.defaultRules)
	assert.Equal(t, 1, testRules2.size, "")
	assert.Equal(t, 1, testRules2.sample.size, "")

}
