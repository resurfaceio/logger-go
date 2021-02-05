package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangesDefaultRules(t *testing.T) {
	httpRules := GetHttpRules()
	for {
		if !assert.Equal(t, newHttpRules("").strictRules, httpRules.GetDefaultRules()) {
			break
		}

		httpRules.SetDefaultRules("")
		if !assert.Equal(t, "", httpRules.defaultRules) {
			break
		}
		if !assert.Equal(t, 0, newHttpRules(httpRules.GetDefaultRules()).size) {
			break
		}

		httpRules.SetDefaultRules(" include default")
		if !assert.Equal(t, "", httpRules.GetDefaultRules()) {
			break
		}

		httpRules.SetDefaultRules("include default\n")
		if !assert.Equal(t, "", httpRules.GetDefaultRules()) {
			break
		}

		httpRules.SetDefaultRules("include default\ninclude default\n")
		if !assert.Equal(t, 0, newHttpRules(httpRules.GetDefaultRules()).size) {
			break
		}

		httpRules.SetDefaultRules("include default\ninclude default\nsample 42")
		rules := newHttpRules(httpRules.GetDefaultRules())
		if !assert.Equal(t, 1, rules.size) {
			break
		}
		if !assert.Equal(t, 1, len(rules.sample)) {
			break
		}
	}

	httpRules.SetDefaultRules(httpRules.GetStrictRules())

}
