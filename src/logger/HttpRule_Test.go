package logger

import (
	"testing"
	"regexp"
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
		if !assert.Equal(t, 0, newHttpRules(httpRules.DefaultRules()).Size()) {
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
		if !assert.Equal(t, 0, newHttpRules(httpRules.DefaultRules()).Size()) {
			break
		}

		httpRules.SetDefaultRules("include default\ninclude default\nsample 42")
		rules := newHttpRules(httpRules.DefaultRules())
		if !assert.Equal(t, 1, rules.Size()) {
			break
		}
		if !assert.Equal(t, 1, len(rules.Sample())) {
			break
		}
		break
	}

	httpRules.SetDefaultRules(httpRules.StrictRules())

}

func TestIncludeDebugRules(t *testing.T) {
	rules := newHttpRules("include debug")
	assert.Equal(t, 2, rules.Size())
	assert.True(t, rules.AllowHttpUrl())
	assert.Equal(t, 1, len(rules.CopySessionField()))

	rules = newHttpRules("include debug\n")
	assert.Equal(t, 2, rules.Size())
	rules = newHttpRules("include debug\nsample 50")
	assert.Equal(t, 3, rules.Size())
	assert.Equal(t, 1, len(rules.Sample()))

	rules = newHttpRules(" include debug\ninclude debug")
	assert.Equal(t, 4, rules.Size())
	rules = newHttpRules("include debug\nsample 50\ninclude debug")
	assert.Equal(t, 5, rules.Size())

	httpRules := GetHttpRules()
	assert.Equal(t, httpRules.StrictRules(), httpRules.DefaultRules())
	httpRules.SetDefaultRules("include debug")
	rules = newHttpRules("")
	for {
		if !assert.Equal(t, 2, rules.Size()) {
			break
		}

		if !assert.True(t, rules.AllowHttpUrl()) {
			break
		}

		if !assert.Equal(t, 1, len(rules.CopySessionField())) {
			break
		}

		break
	}

	httpRules.SetDefaultRules(httpRules.StrictRules())
}

func TestIncludeStandardRules(t *testing.T) {
	rules := newHttpRules("include standard")
	assert.Equal(t, 3, rules.Size())
	assert.Equal(t, 1, len(rules.Remove()))
	assert.Equal(t, 2, len(rules.Replace()))

	rules = newHttpRules("include standard\n")
	assert.Equal(t, 3, rules.Size())
	rules = newHttpRules("include standard\nsample 50")
	assert.Equal(t, 4, rules.Size())
	assert.Equal(t, 1, len(rules.Sample()))

	rules = newHttpRules(" include standard\ninclude standard")
	assert.Equal(t, 6, rules.Size())
	rules = newHttpRules("inlcude standard\nsample 50\ninclude standard")
	assert.Equal(t, 7, rules.Size())

	httpRules := GetHttpRules()
	assert.Equal(t, httpRules.StrictRules(), httpRules.DefaultRules())
	for {
		httpRules.SetDefaultRules("inlcude standard")
		rules = newHttpRules("")
		if !assert.Equal(t, 3, rules.Size()) {
			break
		}

		if !assert.Equal(t, 1, len(rules.Remove())) {
			break
		}

		if !assert.Equal(t, 2, len(rules.Replace())) {
			break
		}

		break
	}

	httpRules.SetDefaultRules(httpRules.StrictRules())
}

func TestIncludeStrictRules(t *testing.T) {
	rules := newHttpRules("include strict")
	assert.Equal(t, 2, rules.Size())
	assert.Equal(t, 1, len(rules.Remove()))
	assert.Equal(t, 1, len(rules.Replace()))

	rules = newHttpRules("include strict\n")
	assert.Equal(t, 2, rules.Size())
	rules = newHttpRules("include strict\nsample 50")
	assert.Equal(t, 3, rules.Size())
	assert.Equal(t, 1, len(rules.Sample()))

	rules = newHttpRules(" include strict\ninclude strict")
	assert.Equal(t, 4, rules.Size())
	rules = newHttpRules(" include strict\nsample 50\ninclude strict")
	assert.Equal(t, 5, rules.Size())

	httpRules := GetHttpRules()
	assert.Equal(t, httpRules.GetStrictRules(), httpRules.GetDefaultRules())
	for {
		httpRules.SetDefaultRules("include strict")
		rules = newHttpRules("")
		if !assert.Equal(t, 2, rules.Size()) {
			break
		}

		if !assert.Equal(t, 1, len(rules.Remove())) {
			break
		}

		if !assert.Equal(t, 1, len(rules.Replace())) {
			break
		}

		break
	}

	httpRules.SetDefaultRules(httpRules.StrictRules())
}

func TestLoadsRulesFromFile(t *testing.T) {
	rules := newHttpRules("file://./test/rules1.txt")
	assert.Equal(t, 1, rules.Size())
	assert.Equal(t, 1, len(rules.Sample()))
	assert.Equal(t, 55, rules.Sample()[0].Param1())

	rules = newHttpRules("file://./test/rules2.txt")
	assert.Equal(t, 3, rules.Size())
	assert.True(t, rules.AllowHttpUrl())
	assert.Equal(t, 1, len(rules.CopySessionField()))
	assert.Equal(t, 1, len(rules.Sample()))
	assert.Equal(t, 56, rules.Sample()[0].Param1())

	rules = newHttpRules("file://./test/rules3.txt ")
	assert.Equal(t, 3, rules.Size())
	assert.Equal(t, 1, len(rules.Replace()))
	assert.Equal(t, 1, len(rules.Sample))
	assert.Equal(t, 57, rules.Sample()[0].Param1())
}

func parseFail(t *testing.T, line string) {
	httpRules := GetHttpRules()
	err := httpRules.ParseRule(line)
	if err != nil {
		return
	}

	assert.True(t, false)
}

func parseOk(t *testing.T, verb string, scope string,
	param1 interface{}, param2 interface{}) {

	httpRules := GetHttpRules()
	rule := httpRules.ParseRul(line)
	assert.Equal(t, verb, rule.Verb())
	
	if rule.Scope() == nil {
		assert.Nil(t, scope)
	} else {
		//not sure if this is right
		assert.Equal(t, scope, regexp.Compile(rule.Scope()))
	}

	rParam1 := rule.Param1()
	_, isPattern := rParam1.(regexp.Regexp)
	if rule.Param1() == nil {
		assert.Nil(t, param1)
	} else if !isPattern {
		assert
	}

        if (rule.param1 == null) {
            expect(param1).toBeNull();
        } else if (rule.param1 instanceof Pattern) {
            expect(param1).toEqual(((Pattern) (rule.param1)).pattern());
        } else {
            expect(param1).toEqual(rule.param1);
        }

        if (rule.param2 == null) {
            expect(param2).toBeNull();
        } else if (rule.param2 instanceof Pattern) {
            expect(param2).toEqual(((Pattern) (rule.param2)).pattern());
        } else {
            expect(param2).toEqual(rule.param2);
        }

}