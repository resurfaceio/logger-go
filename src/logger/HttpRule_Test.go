package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangesDefaultRules(t *testing.T) {
	httpRules := GetHttpRules()
	for {
		if !assert.Equal(t, newHttpRules("").StrictRules(), httpRules.DefaultRules()) {
			break
		}

		httpRules.SetDefaultRules("")
		if !assert.Equal(t, "", httpRules.DefaultRules()) {
			break
		}
		if !assert.Equal(t, 0, newHttpRules(httpRules.DefaultRules()).size) {
			break
		}

		httpRules.SetDefaultRules(" include default")
		if !assert.Equal(t, "", httpRules.DefaultRules()) {
			break
		}

		httpRules.SetDefaultRules("include default\n")
		if !assert.Equal(t, "", httpRules.DefaultRules()) {
			break
		}

		httpRules.SetDefaultRules("include default\ninclude default\n")
		if !assert.Equal(t, 0, newHttpRules(httpRules.DefaultRules()).size) {
			break
		}

		httpRules.SetDefaultRules("include default\ninclude default\nsample 42")
		rules := newHttpRules(httpRules.DefaultRules())
		if !assert.Equal(t, 1, rules.size) {
			break
		}
		if !assert.Equal(t, 1, len(rules.sample)) {
			break
		}
	}

	httpRules.SetDefaultRules(httpRules.StrictRules())

}
