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

}
