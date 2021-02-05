package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncludeDebugRules(t *testing.T) {
	rules := newHttpRules()
	assert.Equal(t, 2, rules.Size)
	assert.True(t, rules.AllowHttpUrl)
	assert.Equal(t, 1, len(rules.CopySessionField))

	rules = newHttpRules("include debug\n")
	assert.Equal(t, 2, rules.Size)
	rules = newHttpRules("include debug\nsample 50")
	assert.Equal(t, 3, rules.Size)
	assert.Equal(t, 1, len(rules.Sample))

	rules = newHttpRules(" include debug\ninclude debug")
	assert.Equal(t, 4, rules.Size)
	rules = newHttpRules("include debug\nsample 50\ninclude debug")
	assert.Equal(t, 5, rules.Size)

	httpRules := GetHttpRules()
	assert.Equal(t, httpRules.StrictRules(), httpRules.DefaultRules())
	httpRules.SetDefaultRules("include debug")
	rules = newHttpRules("")
	for {
		if !assert.Equal(t, 2, rules.Size) {
			break
		}

		if !assert.True(t, rules.AllowHttpUrl) {
			break
		}

		if !assert.Equal(t, 1, len(rules.CopySessionField)) {
			break
		}

		break
	}

	httpRules.SetDefaultRules(httpRules.StrictRules())
}
