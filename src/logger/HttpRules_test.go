package logger

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangesDefaultRules(t *testing.T) {
	httpRules := GetHttpRules()
	for {
		rules, err := newHttpRules("")
		if err != nil {
			fmt.Println(err.Error())
		}
		if !assert.Equal(t, rules.StrictRules(), httpRules.DefaultRules()) {
			break
		}

		httpRules.SetDefaultRules("")
		if !assert.Equal(t, "", httpRules.DefaultRules()) {
			break
		}
		rules, _ = newHttpRules(httpRules.DefaultRules())
		if !assert.Equal(t, 0, rules.Size()) {
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
		rules, _ = newHttpRules(httpRules.DefaultRules())
		if !assert.Equal(t, 0, rules.Size()) {
			break
		}

		httpRules.SetDefaultRules("include default\ninclude default\nsample 42")
		rules, _ = newHttpRules(httpRules.DefaultRules())
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
	rules, _ := newHttpRules("include debug")
	assert.Equal(t, 2, rules.Size())
	assert.True(t, rules.AllowHttpUrl())
	assert.Equal(t, 1, len(rules.CopySessionField()))

	rules, _ = newHttpRules("include debug\n")
	assert.Equal(t, 2, rules.Size())
	rules, _ = newHttpRules("include debug\nsample 50")
	assert.Equal(t, 3, rules.Size())
	assert.Equal(t, 1, len(rules.Sample()))

	rules, _ = newHttpRules(" include debug\ninclude debug")
	assert.Equal(t, 4, rules.Size())
	rules, _ = newHttpRules("include debug\nsample 50\ninclude debug")
	assert.Equal(t, 5, rules.Size())

	httpRules := GetHttpRules()
	assert.Equal(t, httpRules.StrictRules(), httpRules.DefaultRules())
	httpRules.SetDefaultRules("include debug")
	rules, _ = newHttpRules("")
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
	rules, _ := newHttpRules("include standard")
	assert.Equal(t, 3, rules.Size())
	assert.Equal(t, 1, len(rules.Remove()))
	assert.Equal(t, 2, len(rules.Replace()))

	rules, _ = newHttpRules("include standard\n")
	assert.Equal(t, 3, rules.Size())
	rules, _ = newHttpRules("include standard\nsample 50")
	assert.Equal(t, 4, rules.Size())
	assert.Equal(t, 1, len(rules.Sample()))

	rules, _ = newHttpRules(" include standard\ninclude standard")
	assert.Equal(t, 6, rules.Size())
	rules, _ = newHttpRules("include standard\nsample 50\ninclude standard")
	assert.Equal(t, 7, rules.Size())

	httpRules := GetHttpRules()
	assert.Equal(t, httpRules.StrictRules(), httpRules.DefaultRules())
	for {
		httpRules.SetDefaultRules("include standard")
		rules, _ = newHttpRules("")
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
	rules, _ := newHttpRules("include strict")
	assert.Equal(t, 2, rules.Size())
	assert.Equal(t, 1, len(rules.Remove()))
	assert.Equal(t, 1, len(rules.Replace()))

	rules, _ = newHttpRules("include strict\n")
	assert.Equal(t, 2, rules.Size())
	rules, _ = newHttpRules("include strict\nsample 50")
	assert.Equal(t, 3, rules.Size())
	assert.Equal(t, 1, len(rules.Sample()))

	rules, _ = newHttpRules(" include strict\ninclude strict")
	assert.Equal(t, 4, rules.Size())
	rules, _ = newHttpRules(" include strict\nsample 50\ninclude strict")
	assert.Equal(t, 5, rules.Size())

	httpRules, _ := newHttpRules("")
	assert.Equal(t, httpRules.StrictRules(), httpRules.DefaultRules())
	for {
		httpRules.SetDefaultRules("include strict")
		rules, _ = newHttpRules("")
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
	rules, _ := newHttpRules("file://./rules1.txt")
	assert.Equal(t, 1, rules.Size())
	assert.Equal(t, 1, len(rules.Sample()))
	assert.Equal(t, 55, rules.Sample()[0].Param1())

	rules, _ = newHttpRules("file://./rules2.txt")
	assert.Equal(t, 3, rules.Size())
	assert.True(t, rules.AllowHttpUrl())
	assert.Equal(t, 1, len(rules.CopySessionField()))
	assert.Equal(t, 1, len(rules.Sample()))
	assert.Equal(t, 56, rules.Sample()[0].Param1())

	rules, _ = newHttpRules("file://./rules3.txt ")
	assert.Equal(t, 3, rules.Size())
	assert.Equal(t, 1, len(rules.Replace()))
	assert.Equal(t, 1, len(rules.sample))
	assert.Equal(t, 57, rules.Sample()[0].Param1())
}

func parseFail(t *testing.T, line string) {
	_, err := parseRule(line)
	assert.NotNil(t, err)
}

func parseOk(t *testing.T, line string, verb string,
	scope string, param1 interface{}, param2 interface{}) {

	rule, _ := parseRule(line)
	assert.Equal(t, verb, rule.Verb())

	if rule.Scope() == nil {
		assert.Nil(t, scope)
	} else {
		// this may need to change
		assert.Equal(t, scope, rule.Scope().String())
	}

	ruleParam1 := rule.Param1()
	_, notRegexp := ruleParam1.(regexp.Regexp)
	if ruleParam1 == nil {
		assert.Nil(t, param1)
	} else if !notRegexp {
		assert.Equal(t, ruleParam1.(*regexp.Regexp).String(), param1)
	} else {
		assert.Equal(t, ruleParam1, param1)
	}

	ruleParam2 := rule.Param2()
	_, notRegexp = ruleParam2.(*regexp.Regexp)
	if ruleParam2 == nil {
		assert.Nil(t, param2)
	} else if !notRegexp {
		assert.Equal(t, ruleParam2.(*regexp.Regexp).String(), param2)
	} else {
		assert.Equal(t, ruleParam2, param2)
	}
}

func TestParsesEmptyRules(t *testing.T) {
	rules, _ := newHttpRules("")
	assert.Equal(t, 2, rules.Size())
	rules, _ = newHttpRules(" ")
	assert.Equal(t, 2, rules.Size())
	rules, _ = newHttpRules("\t")
	assert.Equal(t, 2, rules.Size())
	rules, _ = newHttpRules("\n")
	assert.Equal(t, 2, rules.Size())

	parsedRule, _ := parseRule("")
	assert.Nil(t, parsedRule)
	parsedRule, _ = parseRule(" ")
	assert.Nil(t, parsedRule)
	parsedRule, _ = parseRule("\t")
	assert.Nil(t, parsedRule)
	parsedRule, _ = parseRule("\n")
	assert.Nil(t, parsedRule)
}

func TestParsesRulesWithBadVerbs(t *testing.T) {
	for _, verb := range []string{"b", "bozo", "*", ".*"} {
		parseFail(t, verb)
		parseFail(t, "!.*! "+verb)
		parseFail(t, "/.*/ "+verb)
		parseFail(t, "%request_body% "+verb)
		parseFail(t, "/^request_header:.*/ "+verb)
	}
}

func TestParsesRulesWithInvalidScopes(t *testing.T) {
	for _, s := range []string{"request_body", "*", ".*"} {
		parseFail(t, "/"+s)
		parseFail(t, "/"+s+" 1")
		parseFail(t, "/"+s+" # 1")
		parseFail(t, "/"+s+"/")
		parseFail(t, "/"+s+"/ # 1")
		parseFail(t, " / "+s)
		parseFail(t, "// "+s)
		parseFail(t, "/// "+s)
		parseFail(t, "/* "+s)
		parseFail(t, "/? "+s)
		parseFail(t, "/+ "+s)
		parseFail(t, "/( "+s)
		parseFail(t, "/(.* "+s)
		parseFail(t, "/(.*)) "+s)

		parseFail(t, "~"+s)
		parseFail(t, "!"+s+" 1")
		parseFail(t, "|"+s+" # 1")
		parseFail(t, "|"+s+"|")
		parseFail(t, "%"+s+"% # 1")
		parseFail(t, " % "+s)
		parseFail(t, "%% "+s)
		parseFail(t, "%%% "+s)
		parseFail(t, "%* "+s)
		parseFail(t, "%? "+s)
		parseFail(t, "%+ "+s)
		parseFail(t, "%( "+s)
		parseFail(t, "%(.* "+s)
		parseFail(t, "%(.*)) "+s)

		parseFail(t, "~"+s+"%")
		parseFail(t, "!"+s+"%# 1")
		parseFail(t, "|"+s+"% # 1")
		parseFail(t, "|"+s+"%")
		parseFail(t, "%"+s+"| # 1")
		parseFail(t, "~(.*! "+s)
		parseFail(t, "~(.*))! "+s)
		parseFail(t, "/(.*! "+s)
		parseFail(t, "/(.*))! "+s)
	}
}

func TestParsesAllowHttpRules(t *testing.T) {
	parseFail(t, "allow_http_url whaa")
	parseOk(t, "allow_http_url", "allow_http_url", "", nil, nil)
	parseOk(t, "allow_http_url # be safe bro!", "allow_http_url", "", nil, nil)
}

func TestParsesCopySessionFieldRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| copy_session_field %1%, %2%")
	parseFail(t, "!.*! copy_session_field /1/, 2")
	parseFail(t, "/.*/ copy_session_field /1/, /2")
	parseFail(t, "/.*/ copy_session_field /1/, /2/")
	parseFail(t, "/.*/ copy_session_field /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! copy_session_field %1%, %2%, %3%")
	parseFail(t, "/.*/ copy_session_field /1/, /2/, 3")
	parseFail(t, "/.*/ copy_session_field /1/, /2/, /3")
	parseFail(t, "/.*/ copy_session_field /1/, /2/, /3/")
	parseFail(t, "%.*% copy_session_field /1/, /2/, /3/ # blah")

	// with missing params
	parseFail(t, "!.*! copy_session_field")
	parseFail(t, "/.*/ copy_session_field")
	parseFail(t, "/.*/ copy_session_field /")
	parseFail(t, "/.*/ copy_session_field //")
	parseFail(t, "/.*/ copy_session_field blah")
	parseFail(t, "/.*/ copy_session_field # bleep")
	parseFail(t, "/.*/ copy_session_field blah # bleep")

	// with invalid params
	parseFail(t, "/.*/ copy_session_field /")
	parseFail(t, "/.*/ copy_session_field //")
	parseFail(t, "/.*/ copy_session_field ///")
	parseFail(t, "/.*/ copy_session_field /*/")
	parseFail(t, "/.*/ copy_session_field /?/")
	parseFail(t, "/.*/ copy_session_field /+/")
	parseFail(t, "/.*/ copy_session_field /(/")
	parseFail(t, "/.*/ copy_session_field /(.*/")
	parseFail(t, "/.*/ copy_session_field /(.*))/")

	// with valid regexes
	parseOk(t, "copy_session_field !.*!", "copy_session_field", "", "^.*$", nil)
	parseOk(t, "copy_session_field /.*/", "copy_session_field", "", "^.*$", nil)
	parseOk(t, "copy_session_field /^.*/", "copy_session_field", "", "^.*$", nil)
	parseOk(t, "copy_session_field /.*$/", "copy_session_field", "", "^.*$", nil)
	parseOk(t, "copy_session_field /^.*$/", "copy_session_field", "", "^.*$", nil)

	// with valid regexes and escape sequences
	parseOk(t, "copy_session_field !A\\!|B!", "copy_session_field", "", "^A!|B$", nil)
	parseOk(t, "copy_session_field |A\\|B|", "copy_session_field", "", "^A|B$", nil)
	parseOk(t, "copy_session_field |A\\|B\\|C|", "copy_session_field", "", "^A|B|C$", nil)
	parseOk(t, "copy_session_field /A\\/B\\/C/", "copy_session_field", "", "^A/B/C$", nil)
}

func TestParsesRemoveRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| remove %1%")
	parseFail(t, "~.*~ remove 1")
	parseFail(t, "/.*/ remove /1/")
	parseFail(t, "/.*/ remove 1 # bleep")
	parseFail(t, "|.*| remove %1%, %2%")
	parseFail(t, "!.*! remove /1/, 2")
	parseFail(t, "/.*/ remove /1/, /2")
	parseFail(t, "/.*/ remove /1/, /2/")
	parseFail(t, "/.*/ remove /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! remove %1%, %2%, %3%")
	parseFail(t, "/.*/ remove /1/, /2/, 3")
	parseFail(t, "/.*/ remove /1/, /2/, /3")
	parseFail(t, "/.*/ remove /1/, /2/, /3/")
	parseFail(t, "%.*% remove /1/, /2/, /3/ # blah")

	// with valid regexes
	parseOk(t, "%request_header:cookie|response_header:set-cookie% remove",
		"remove", "^request_header:cookie|response_header:set-cookie$", nil, nil)
	parseOk(t, "/request_header:cookie|response_header:set-cookie/ remove",
		"remove", "^request_header:cookie|response_header:set-cookie$", nil, nil)

	// with valid regexes and escape sequences
	parseOk(t, "!request_header\\!|response_header:set-cookie! remove",
		"remove", "^request_header!|response_header:set-cookie$", nil, nil)
	parseOk(t, "|request_header:cookie\\|response_header:set-cookie| remove",
		"remove", "^request_header:cookie|response_header:set-cookie$", nil, nil)
	parseOk(t, "|request_header:cookie\\|response_header:set-cookie\\|boo| remove",
		"remove", "^request_header:cookie|response_header:set-cookie|boo$", nil, nil)
	parseOk(t, "/request_header:cookie\\/response_header:set-cookie\\/boo/ remove",
		"remove", "^request_header:cookie/response_header:set-cookie/boo$", nil, nil)
}

func TestParsesRemoveIfRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| remove_if %1%, %2%")
	parseFail(t, "!.*! remove_if /1/, 2")
	parseFail(t, "/.*/ remove_if /1/, /2")
	parseFail(t, "/.*/ remove_if /1/, /2/")
	parseFail(t, "/.*/ remove_if /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! remove_if %1%, %2%, %3%")
	parseFail(t, "/.*/ remove_if /1/, /2/, 3")
	parseFail(t, "/.*/ remove_if /1/, /2/, /3")
	parseFail(t, "/.*/ remove_if /1/, /2/, /3/")
	parseFail(t, "%.*% remove_if /1/, /2/, /3/ # blah")

	// with missing params
	parseFail(t, "!.*! remove_if")
	parseFail(t, "/.*/ remove_if")
	parseFail(t, "/.*/ remove_if /")
	parseFail(t, "/.*/ remove_if //")
	parseFail(t, "/.*/ remove_if blah")
	parseFail(t, "/.*/ remove_if # bleep")
	parseFail(t, "/.*/ remove_if blah # bleep")

	// with invalid params
	parseFail(t, "/.*/ remove_if /")
	parseFail(t, "/.*/ remove_if //")
	parseFail(t, "/.*/ remove_if ///")
	parseFail(t, "/.*/ remove_if /*/")
	parseFail(t, "/.*/ remove_if /?/")
	parseFail(t, "/.*/ remove_if /+/")
	parseFail(t, "/.*/ remove_if /(/")
	parseFail(t, "/.*/ remove_if /(.*/")
	parseFail(t, "/.*/ remove_if /(.*))/")

	// with valid regexes
	parseOk(t, "%response_body% remove_if %<!--SKIP_BODY_LOGGING-->%",
		"remove_if", "^response_body$", "^<!--SKIP_BODY_LOGGING-->$", nil)
	parseOk(t, "/response_body/ remove_if /<!--SKIP_BODY_LOGGING-->/",
		"remove_if", "^response_body$", "^<!--SKIP_BODY_LOGGING-->$", nil)

	// with valid regexes and escape sequences
	parseOk(t, "!request_body|response_body! remove_if |<!--IGNORE_LOGGING-->\\|<!-SKIP-->|",
		"remove_if", "^request_body|response_body$", "^<!--IGNORE_LOGGING-->|<!-SKIP-->$", nil)
	parseOk(t, "|request_body\\|response_body| remove_if |<!--IGNORE_LOGGING-->\\|<!-SKIP-->|",
		"remove_if", "^request_body|response_body$", "^<!--IGNORE_LOGGING-->|<!-SKIP-->$", nil)
	parseOk(t, "|request_body\\|response_body\\|boo| remove_if |<!--IGNORE_LOGGING-->\\|<!-SKIP-->\\|asdf|",
		"remove_if", "^request_body|response_body|boo$", "^<!--IGNORE_LOGGING-->|<!-SKIP-->|asdf$", nil)
	parseOk(t, "/request_body\\/response_body\\/boo/ remove_if |<!--IGNORE_LOGGING-->\\|<!-SKIP-->\\|asdf|",
		"remove_if", "^request_body/response_body/boo$", "^<!--IGNORE_LOGGING-->|<!-SKIP-->|asdf$", nil)
}

func TestParsesRemoveIfFoundRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| remove_if_found %1%, %2%")
	parseFail(t, "!.*! remove_if_found /1/, 2")
	parseFail(t, "/.*/ remove_if_found /1/, /2")
	parseFail(t, "/.*/ remove_if_found /1/, /2/")
	parseFail(t, "/.*/ remove_if_found /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! remove_if_found %1%, %2%, %3%")
	parseFail(t, "/.*/ remove_if_found /1/, /2/, 3")
	parseFail(t, "/.*/ remove_if_found /1/, /2/, /3")
	parseFail(t, "/.*/ remove_if_found /1/, /2/, /3/")
	parseFail(t, "%.*% remove_if_found /1/, /2/, /3/ # blah")

	// with missing params
	parseFail(t, "!.*! remove_if_found")
	parseFail(t, "/.*/ remove_if_found")
	parseFail(t, "/.*/ remove_if_found /")
	parseFail(t, "/.*/ remove_if_found //")
	parseFail(t, "/.*/ remove_if_found blah")
	parseFail(t, "/.*/ remove_if_found # bleep")
	parseFail(t, "/.*/ remove_if_found blah # bleep")

	// with invalid params
	parseFail(t, "/.*/ remove_if_found /")
	parseFail(t, "/.*/ remove_if_found //")
	parseFail(t, "/.*/ remove_if_found ///")
	parseFail(t, "/.*/ remove_if_found /*/")
	parseFail(t, "/.*/ remove_if_found /?/")
	parseFail(t, "/.*/ remove_if_found /+/")
	parseFail(t, "/.*/ remove_if_found /(/")
	parseFail(t, "/.*/ remove_if_found /(.*/")
	parseFail(t, "/.*/ remove_if_found /(.*))/")

	// with valid regexes
	parseOk(t, "%response_body% remove_if_found %<!--SKIP_BODY_LOGGING-->%",
		"remove_if_found", "^response_body$", "<!--SKIP_BODY_LOGGING-->", nil)
	parseOk(t, "/response_body/ remove_if_found /<!--SKIP_BODY_LOGGING-->/",
		"remove_if_found", "^response_body$", "<!--SKIP_BODY_LOGGING-->", nil)

	// with valid regexes and escape sequences
	parseOk(t, "!request_body|response_body! remove_if_found |<!--IGNORE_LOGGING-->\\|<!-SKIP-->|",
		"remove_if_found", "^request_body|response_body$", "<!--IGNORE_LOGGING-->|<!-SKIP-->", nil)
	parseOk(t, "|request_body\\|response_body| remove_if_found |<!--IGNORE_LOGGING-->\\|<!-SKIP-->|",
		"remove_if_found", "^request_body|response_body$", "<!--IGNORE_LOGGING-->|<!-SKIP-->", nil)
	parseOk(t, "|request_body\\|response_body\\|boo| remove_if_found |<!--IGNORE_LOGGING-->\\|<!-SKIP-->\\|asdf|",
		"remove_if_found", "^request_body|response_body|boo$", "<!--IGNORE_LOGGING-->|<!-SKIP-->|asdf", nil)
	parseOk(t, "/request_body\\/response_body\\/boo/ remove_if_found |<!--IGNORE_LOGGING-->\\|<!-SKIP-->\\|asdf|",
		"remove_if_found", "^request_body/response_body/boo$", "<!--IGNORE_LOGGING-->|<!-SKIP-->|asdf", nil)
}

func TestParsesRemoveUnlessRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| remove_unless %1%, %2%")
	parseFail(t, "!.*! remove_unless /1/, 2")
	parseFail(t, "/.*/ remove_unless /1/, /2")
	parseFail(t, "/.*/ remove_unless /1/, /2/")
	parseFail(t, "/.*/ remove_unless /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! remove_unless %1%, %2%, %3%")
	parseFail(t, "/.*/ remove_unless /1/, /2/, 3")
	parseFail(t, "/.*/ remove_unless /1/, /2/, /3")
	parseFail(t, "/.*/ remove_unless /1/, /2/, /3/")
	parseFail(t, "%.*% remove_unless /1/, /2/, /3/ # blah")

	// with missing params
	parseFail(t, "!.*! remove_unless")
	parseFail(t, "/.*/ remove_unless")
	parseFail(t, "/.*/ remove_unless /")
	parseFail(t, "/.*/ remove_unless //")
	parseFail(t, "/.*/ remove_unless blah")
	parseFail(t, "/.*/ remove_unless # bleep")
	parseFail(t, "/.*/ remove_unless blah # bleep")

	// with invalid params
	parseFail(t, "/.*/ remove_unless /")
	parseFail(t, "/.*/ remove_unless //")
	parseFail(t, "/.*/ remove_unless ///")
	parseFail(t, "/.*/ remove_unless /*/")
	parseFail(t, "/.*/ remove_unless /?/")
	parseFail(t, "/.*/ remove_unless /+/")
	parseFail(t, "/.*/ remove_unless /(/")
	parseFail(t, "/.*/ remove_unless /(.*/")
	parseFail(t, "/.*/ remove_unless /(.*))/")

	// with valid regexes
	parseOk(t, "%response_body% remove_unless %<!--PERFORM_BODY_LOGGING-->%",
		"remove_unless", "^response_body$", "^<!--PERFORM_BODY_LOGGING-->$", nil)
	parseOk(t, "/response_body/ remove_unless /<!--PERFORM_BODY_LOGGING-->/",
		"remove_unless", "^response_body$", "^<!--PERFORM_BODY_LOGGING-->$", nil)

	// with valid regexes and escape sequences
	parseOk(t, "!request_body|response_body! remove_unless |<!--PERFORM_LOGGING-->\\|<!-SKIP-->|",
		"remove_unless", "^request_body|response_body$", "^<!--PERFORM_LOGGING-->|<!-SKIP-->$", nil)
	parseOk(t, "|request_body\\|response_body| remove_unless |<!--PERFORM_LOGGING-->\\|<!-SKIP-->|",
		"remove_unless", "^request_body|response_body$", "^<!--PERFORM_LOGGING-->|<!-SKIP-->$", nil)
	parseOk(t, "|request_body\\|response_body\\|boo| remove_unless |<!--PERFORM_LOGGING-->\\|<!-SKIP-->\\|skipit|",
		"remove_unless", "^request_body|response_body|boo$", "^<!--PERFORM_LOGGING-->|<!-SKIP-->|skipit$", nil)
	parseOk(t, "/request_body\\/response_body\\/boo/ remove_unless |<!--PERFORM_LOGGING-->\\|<!-SKIP-->\\|skipit|",
		"remove_unless", "^request_body/response_body/boo$", "^<!--PERFORM_LOGGING-->|<!-SKIP-->|skipit$", nil)
}

func TestParsesRemoveUnlessFoundRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| remove_unless_found %1%, %2%")
	parseFail(t, "!.*! remove_unless_found /1/, 2")
	parseFail(t, "/.*/ remove_unless_found /1/, /2")
	parseFail(t, "/.*/ remove_unless_found /1/, /2/")
	parseFail(t, "/.*/ remove_unless_found /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! remove_unless_found %1%, %2%, %3%")
	parseFail(t, "/.*/ remove_unless_found /1/, /2/, 3")
	parseFail(t, "/.*/ remove_unless_found /1/, /2/, /3")
	parseFail(t, "/.*/ remove_unless_found /1/, /2/, /3/")
	parseFail(t, "%.*% remove_unless_found /1/, /2/, /3/ # blah")

	// with missing params
	parseFail(t, "!.*! remove_unless_found")
	parseFail(t, "/.*/ remove_unless_found")
	parseFail(t, "/.*/ remove_unless_found /")
	parseFail(t, "/.*/ remove_unless_found //")
	parseFail(t, "/.*/ remove_unless_found blah")
	parseFail(t, "/.*/ remove_unless_found # bleep")
	parseFail(t, "/.*/ remove_unless_found blah # bleep")

	// with invalid params
	parseFail(t, "/.*/ remove_unless_found /")
	parseFail(t, "/.*/ remove_unless_found //")
	parseFail(t, "/.*/ remove_unless_found ///")
	parseFail(t, "/.*/ remove_unless_found /*/")
	parseFail(t, "/.*/ remove_unless_found /?/")
	parseFail(t, "/.*/ remove_unless_found /+/")
	parseFail(t, "/.*/ remove_unless_found /(/")
	parseFail(t, "/.*/ remove_unless_found /(.*/")
	parseFail(t, "/.*/ remove_unless_found /(.*))/")

	// with valid regexes
	parseOk(t, "%response_body% remove_unless_found %<!--PERFORM_BODY_LOGGING-->%",
		"remove_unless_found", "^response_body$", "<!--PERFORM_BODY_LOGGING-->", nil)
	parseOk(t, "/response_body/ remove_unless_found /<!--PERFORM_BODY_LOGGING-->/",
		"remove_unless_found", "^response_body$", "<!--PERFORM_BODY_LOGGING-->", nil)

	// with valid regexes and escape sequences
	parseOk(t, "!request_body|response_body! remove_unless_found |<!--PERFORM_LOGGING-->\\|<!-SKIP-->|",
		"remove_unless_found", "^request_body|response_body$", "<!--PERFORM_LOGGING-->|<!-SKIP-->", nil)
	parseOk(t, "|request_body\\|response_body| remove_unless_found |<!--PERFORM_LOGGING-->\\|<!-SKIP-->|",
		"remove_unless_found", "^request_body|response_body$", "<!--PERFORM_LOGGING-->|<!-SKIP-->", nil)
	parseOk(t, "|request_body\\|response_body\\|boo| remove_unless_found |<!--PERFORM_LOGGING-->\\|<!-SKIP-->\\|skipit|",
		"remove_unless_found", "^request_body|response_body|boo$", "<!--PERFORM_LOGGING-->|<!-SKIP-->|skipit", nil)
	parseOk(t, "/request_body\\/response_body\\/boo/ remove_unless_found |<!--PERFORM_LOGGING-->\\|<!-SKIP-->\\|skipit|",
		"remove_unless_found", "^request_body/response_body/boo$", "<!--PERFORM_LOGGING-->|<!-SKIP-->|skipit", nil)
}

func TestParsesReplaceRules(t *testing.T) {
	// with extra params
	parseFail(t, "!.*! replace %1%, %2%, %3%")
	parseFail(t, "/.*/ replace /1/, /2/, 3")
	parseFail(t, "/.*/ replace /1/, /2/, /3")
	parseFail(t, "/.*/ replace /1/, /2/, /3/")
	parseFail(t, "%.*% replace /1/, /2/, /3/ # blah")

	// with missing params
	parseFail(t, "!.*! replace")
	parseFail(t, "/.*/ replace")
	parseFail(t, "/.*/ replace /")
	parseFail(t, "/.*/ replace //")
	parseFail(t, "/.*/ replace blah")
	parseFail(t, "/.*/ replace # bleep")
	parseFail(t, "/.*/ replace blah # bleep")
	parseFail(t, "!.*! replace boo yah")
	parseFail(t, "/.*/ replace boo yah")
	parseFail(t, "/.*/ replace boo yah # bro")
	parseFail(t, "/.*/ replace /.*/ # bleep")
	parseFail(t, "/.*/ replace /.*/, # bleep")
	parseFail(t, "/.*/ replace /.*/, /# bleep")
	parseFail(t, "/.*/ replace // # bleep")
	parseFail(t, "/.*/ replace // // # bleep")

	// with invalid params
	parseFail(t, "/.*/ replace /")
	parseFail(t, "/.*/ replace //")
	parseFail(t, "/.*/ replace ///")
	parseFail(t, "/.*/ replace /*/")
	parseFail(t, "/.*/ replace /?/")
	parseFail(t, "/.*/ replace /+/")
	parseFail(t, "/.*/ replace /(/")
	parseFail(t, "/.*/ replace /(.*/")
	parseFail(t, "/.*/ replace /(.*))/")
	parseFail(t, "/.*/ replace /1/, ~")
	parseFail(t, "/.*/ replace /1/, !")
	parseFail(t, "/.*/ replace /1/, %")
	parseFail(t, "/.*/ replace /1/, |")
	parseFail(t, "/.*/ replace /1/, /")

	// with valid regexes
	parseOk(t, "%response_body% replace %kurt%, %vagner%", "replace", "^response_body$", "kurt", "vagner")
	parseOk(t, "/response_body/ replace /kurt/, /vagner/", "replace", "^response_body$", "kurt", "vagner")
	parseOk(t, "%response_body|.+_header:.+% replace %kurt%, %vagner%",
		"replace", "^response_body|.+_header:.+$", "kurt", "vagner")
	parseOk(t, "|response_body\\|.+_header:.+| replace |kurt|, |vagner\\|frazier|",
		"replace", "^response_body|.+_header:.+$", "kurt", "vagner|frazier")

	// with valid regexes and escape sequences
	parseOk(t, "|response_body\\|.+_header:.+| replace |kurt|, |vagner|",
		"replace", "^response_body|.+_header:.+$", "kurt", "vagner")
	parseOk(t, "|response_body\\|.+_header:.+\\|boo| replace |kurt|, |vagner|",
		"replace", "^response_body|.+_header:.+|boo$", "kurt", "vagner")
	parseOk(t, "|response_body| replace |kurt\\|bruce|, |vagner|",
		"replace", "^response_body$", "kurt|bruce", "vagner")
	parseOk(t, "|response_body| replace |kurt\\|bruce\\|kevin|, |vagner|",
		"replace", "^response_body$", "kurt|bruce|kevin", "vagner")
	parseOk(t, "|response_body| replace /kurt\\/bruce\\/kevin/, |vagner|",
		"replace", "^response_body$", "kurt/bruce/kevin", "vagner")
}

func TestParsesSampleRules(t *testing.T) {
	parseFail(t, "sample")
	parseFail(t, "sample 50 50")
	parseFail(t, "sample 0")
	parseFail(t, "sample 100")
	parseFail(t, "sample 105")
	parseFail(t, "sample 10.5")
	parseFail(t, "sample blue")
	parseFail(t, "sample # bleep")
	parseFail(t, "sample blue # bleep")
	parseFail(t, "sample //")
	parseFail(t, "sample /42/")
}

func TestParsesSkipCompressionRules(t *testing.T) {
	parseFail(t, "skip_compression whaa")
	parseOk(t, "skip_compression", "skip_compression", "", nil, nil)
	parseOk(t, "skip_compression # slightly faster!", "skip_compression", "", nil, nil)
}

func TestParsesSkipSubmissionRules(t *testing.T) {
	parseFail(t, "skip_submission whaa")
	parseOk(t, "skip_submission", "skip_submission", "", nil, nil)
	parseOk(t, "skip_submission # slightly faster!", "skip_submission", "", nil, nil)
}

func TestParsesStopRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| stop %1%")
	parseFail(t, "~.*~ stop 1")
	parseFail(t, "/.*/ stop /1/")
	parseFail(t, "/.*/ stop 1 # bleep")
	parseFail(t, "|.*| stop %1%, %2%")
	parseFail(t, "!.*! stop /1/, 2")
	parseFail(t, "/.*/ stop /1/, /2")
	parseFail(t, "/.*/ stop /1/, /2/")
	parseFail(t, "/.*/ stop /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! stop %1%, %2%, %3%")
	parseFail(t, "/.*/ stop /1/, /2/, 3")
	parseFail(t, "/.*/ stop /1/, /2/, /3")
	parseFail(t, "/.*/ stop /1/, /2/, /3/")
	parseFail(t, "%.*% stop /1/, /2/, /3/ # blah")

	// with valid regexes
	parseOk(t, "%request_header:skip_usage_logging% stop", "stop", "^request_header:skip_usage_logging$", nil, nil)
	parseOk(t, "|request_header:skip_usage_logging| stop", "stop", "^request_header:skip_usage_logging$", nil, nil)
	parseOk(t, "/request_header:skip_usage_logging/ stop", "stop", "^request_header:skip_usage_logging$", nil, nil)

	// with valid regexes and escape sequences
	parseOk(t, "!request_header\\!! stop", "stop", "^request_header!$", nil, nil)
	parseOk(t, "|request_header\\|response_header| stop", "stop", "^request_header|response_header$", nil, nil)
	parseOk(t, "|request_header\\|response_header\\|boo| stop", "stop", "^request_header|response_header|boo$", nil, nil)
	parseOk(t, "/request_header\\/response_header\\/boo/ stop", "stop", "^request_header/response_header/boo$", nil, nil)
}

func TestParsesStopIfRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| stop_if %1%, %2%")
	parseFail(t, "!.*! stop_if /1/, 2")
	parseFail(t, "/.*/ stop_if /1/, /2")
	parseFail(t, "/.*/ stop_if /1/, /2/")
	parseFail(t, "/.*/ stop_if /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! stop_if %1%, %2%, %3%")
	parseFail(t, "/.*/ stop_if /1/, /2/, 3")
	parseFail(t, "/.*/ stop_if /1/, /2/, /3")
	parseFail(t, "/.*/ stop_if /1/, /2/, /3/")
	parseFail(t, "%.*% stop_if /1/, /2/, /3/ # blah")

	// with missing params
	parseFail(t, "!.*! stop_if")
	parseFail(t, "/.*/ stop_if")
	parseFail(t, "/.*/ stop_if /")
	parseFail(t, "/.*/ stop_if //")
	parseFail(t, "/.*/ stop_if blah")
	parseFail(t, "/.*/ stop_if # bleep")
	parseFail(t, "/.*/ stop_if blah # bleep")

	// with invalid params
	parseFail(t, "/.*/ stop_if /")
	parseFail(t, "/.*/ stop_if //")
	parseFail(t, "/.*/ stop_if ///")
	parseFail(t, "/.*/ stop_if /*/")
	parseFail(t, "/.*/ stop_if /?/")
	parseFail(t, "/.*/ stop_if /+/")
	parseFail(t, "/.*/ stop_if /(/")
	parseFail(t, "/.*/ stop_if /(.*/")
	parseFail(t, "/.*/ stop_if /(.*))/")

	// with valid regexes
	parseOk(t, "%response_body% stop_if %<!--IGNORE_LOGGING-->%", "stop_if", "^response_body$", "^<!--IGNORE_LOGGING-->$", nil)
	parseOk(t, "/response_body/ stop_if /<!--IGNORE_LOGGING-->/", "stop_if", "^response_body$", "^<!--IGNORE_LOGGING-->$", nil)

	// with valid regexes and escape sequences
	parseOk(t, "!request_body|response_body! stop_if |<!--IGNORE_LOGGING-->\\|<!-SKIP-->|",
		"stop_if", "^request_body|response_body$", "^<!--IGNORE_LOGGING-->|<!-SKIP-->$", nil)
	parseOk(t, "!request_body|response_body|boo\\!! stop_if |<!--IGNORE_LOGGING-->\\|<!-SKIP-->|",
		"stop_if", "^request_body|response_body|boo!$", "^<!--IGNORE_LOGGING-->|<!-SKIP-->$", nil)
	parseOk(t, "|request_body\\|response_body| stop_if |<!--IGNORE_LOGGING-->\\|<!-SKIP-->|",
		"stop_if", "^request_body|response_body$", "^<!--IGNORE_LOGGING-->|<!-SKIP-->$", nil)
	parseOk(t, "|request_body\\|response_body| stop_if |<!--IGNORE_LOGGING-->\\|<!-SKIP-->\\|pipe\\||",
		"stop_if", "^request_body|response_body$", "^<!--IGNORE_LOGGING-->|<!-SKIP-->|pipe|$", nil)
	parseOk(t, "/request_body\\/response_body/ stop_if |<!--IGNORE_LOGGING-->\\|<!-SKIP-->\\|pipe\\||",
		"stop_if", "^request_body/response_body$", "^<!--IGNORE_LOGGING-->|<!-SKIP-->|pipe|$", nil)
}

func TestParsesStopIfFoundRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| stop_if_found %1%, %2%")
	parseFail(t, "!.*! stop_if_found /1/, 2")
	parseFail(t, "/.*/ stop_if_found /1/, /2")
	parseFail(t, "/.*/ stop_if_found /1/, /2/")
	parseFail(t, "/.*/ stop_if_found /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! stop_if_found %1%, %2%, %3%")
	parseFail(t, "/.*/ stop_if_found /1/, /2/, 3")
	parseFail(t, "/.*/ stop_if_found /1/, /2/, /3")
	parseFail(t, "/.*/ stop_if_found /1/, /2/, /3/")
	parseFail(t, "%.*% stop_if_found /1/, /2/, /3/ # blah")

	// with missing params
	parseFail(t, "!.*! stop_if_found")
	parseFail(t, "/.*/ stop_if_found")
	parseFail(t, "/.*/ stop_if_found /")
	parseFail(t, "/.*/ stop_if_found //")
	parseFail(t, "/.*/ stop_if_found blah")
	parseFail(t, "/.*/ stop_if_found # bleep")
	parseFail(t, "/.*/ stop_if_found blah # bleep")

	// with invalid params
	parseFail(t, "/.*/ stop_if_found /")
	parseFail(t, "/.*/ stop_if_found //")
	parseFail(t, "/.*/ stop_if_found ///")
	parseFail(t, "/.*/ stop_if_found /*/")
	parseFail(t, "/.*/ stop_if_found /?/")
	parseFail(t, "/.*/ stop_if_found /+/")
	parseFail(t, "/.*/ stop_if_found /(/")
	parseFail(t, "/.*/ stop_if_found /(.*/")
	parseFail(t, "/.*/ stop_if_found /(.*))/")

	// with valid regexes
	parseOk(t, "%response_body% stop_if_found %<!--IGNORE_LOGGING-->%",
		"stop_if_found", "^response_body$", "<!--IGNORE_LOGGING-->", nil)
	parseOk(t, "/response_body/ stop_if_found /<!--IGNORE_LOGGING-->/",
		"stop_if_found", "^response_body$", "<!--IGNORE_LOGGING-->", nil)

	// with valid regexes and escape sequences
	parseOk(t, "!request_body|response_body! stop_if_found |<!--IGNORE_LOGGING-->\\|<!-SKIP-->|",
		"stop_if_found", "^request_body|response_body$", "<!--IGNORE_LOGGING-->|<!-SKIP-->", nil)
	parseOk(t, "!request_body|response_body|boo\\!! stop_if_found |<!--IGNORE_LOGGING-->\\|<!-SKIP-->|",
		"stop_if_found", "^request_body|response_body|boo!$", "<!--IGNORE_LOGGING-->|<!-SKIP-->", nil)
	parseOk(t, "|request_body\\|response_body| stop_if_found |<!--IGNORE_LOGGING-->\\|<!-SKIP-->|",
		"stop_if_found", "^request_body|response_body$", "<!--IGNORE_LOGGING-->|<!-SKIP-->", nil)
	parseOk(t, "|request_body\\|response_body| stop_if_found |<!--IGNORE_LOGGING-->\\|<!-SKIP-->\\|pipe\\||",
		"stop_if_found", "^request_body|response_body$", "<!--IGNORE_LOGGING-->|<!-SKIP-->|pipe|", nil)
	parseOk(t, "/request_body\\/response_body/ stop_if_found |<!--IGNORE_LOGGING-->\\|<!-SKIP-->\\|pipe\\||",
		"stop_if_found", "^request_body/response_body$", "<!--IGNORE_LOGGING-->|<!-SKIP-->|pipe|", nil)
}

func TestParsesStopUnlessRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| stop_unless %1%, %2%")
	parseFail(t, "!.*! stop_unless /1/, 2")
	parseFail(t, "/.*/ stop_unless /1/, /2")
	parseFail(t, "/.*/ stop_unless /1/, /2/")
	parseFail(t, "/.*/ stop_unless /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! stop_unless %1%, %2%, %3%")
	parseFail(t, "/.*/ stop_unless /1/, /2/, 3")
	parseFail(t, "/.*/ stop_unless /1/, /2/, /3")
	parseFail(t, "/.*/ stop_unless /1/, /2/, /3/")
	parseFail(t, "%.*% stop_unless /1/, /2/, /3/ # blah")

	// with missing params
	parseFail(t, "!.*! stop_unless")
	parseFail(t, "/.*/ stop_unless")
	parseFail(t, "/.*/ stop_unless /")
	parseFail(t, "/.*/ stop_unless //")
	parseFail(t, "/.*/ stop_unless blah")
	parseFail(t, "/.*/ stop_unless # bleep")
	parseFail(t, "/.*/ stop_unless blah # bleep")

	// with invalid params
	parseFail(t, "/.*/ stop_unless /")
	parseFail(t, "/.*/ stop_unless //")
	parseFail(t, "/.*/ stop_unless ///")
	parseFail(t, "/.*/ stop_unless /*/")
	parseFail(t, "/.*/ stop_unless /?/")
	parseFail(t, "/.*/ stop_unless /+/")
	parseFail(t, "/.*/ stop_unless /(/")
	parseFail(t, "/.*/ stop_unless /(.*/")
	parseFail(t, "/.*/ stop_unless /(.*))/")

	// with valid regexes
	parseOk(t, "%response_body% stop_unless %<!--DO_LOGGING-->%", "stop_unless", "^response_body$", "^<!--DO_LOGGING-->$", nil)
	parseOk(t, "/response_body/ stop_unless /<!--DO_LOGGING-->/", "stop_unless", "^response_body$", "^<!--DO_LOGGING-->$", nil)

	// with valid regexes and escape sequences
	parseOk(t, "!request_body|response_body! stop_unless |<!--DO_LOGGING-->\\|<!-NOSKIP-->|",
		"stop_unless", "^request_body|response_body$", "^<!--DO_LOGGING-->|<!-NOSKIP-->$", nil)
	parseOk(t, "!request_body|response_body|boo\\!! stop_unless |<!--DO_LOGGING-->\\|<!-NOSKIP-->|",
		"stop_unless", "^request_body|response_body|boo!$", "^<!--DO_LOGGING-->|<!-NOSKIP-->$", nil)
	parseOk(t, "|request_body\\|response_body| stop_unless |<!--DO_LOGGING-->\\|<!-NOSKIP-->|",
		"stop_unless", "^request_body|response_body$", "^<!--DO_LOGGING-->|<!-NOSKIP-->$", nil)
	parseOk(t, "|request_body\\|response_body| stop_unless |<!--DO_LOGGING-->\\|<!-NOSKIP-->\\|pipe\\||",
		"stop_unless", "^request_body|response_body$", "^<!--DO_LOGGING-->|<!-NOSKIP-->|pipe|$", nil)
	parseOk(t, "/request_body\\/response_body/ stop_unless |<!--DO_LOGGING-->\\|<!-NOSKIP-->\\|pipe\\||",
		"stop_unless", "^request_body/response_body$", "^<!--DO_LOGGING-->|<!-NOSKIP-->|pipe|$", nil)
}

func TestParsesStopUnlessFoundRules(t *testing.T) {
	// with extra params
	parseFail(t, "|.*| stop_unless_found %1%, %2%")
	parseFail(t, "!.*! stop_unless_found /1/, 2")
	parseFail(t, "/.*/ stop_unless_found /1/, /2")
	parseFail(t, "/.*/ stop_unless_found /1/, /2/")
	parseFail(t, "/.*/ stop_unless_found /1/, /2/, /3/ # blah")
	parseFail(t, "!.*! stop_unless_found %1%, %2%, %3%")
	parseFail(t, "/.*/ stop_unless_found /1/, /2/, 3")
	parseFail(t, "/.*/ stop_unless_found /1/, /2/, /3")
	parseFail(t, "/.*/ stop_unless_found /1/, /2/, /3/")
	parseFail(t, "%.*% stop_unless_found /1/, /2/, /3/ # blah")

	// with missing params
	parseFail(t, "!.*! stop_unless_found")
	parseFail(t, "/.*/ stop_unless_found")
	parseFail(t, "/.*/ stop_unless_found /")
	parseFail(t, "/.*/ stop_unless_found //")
	parseFail(t, "/.*/ stop_unless_found blah")
	parseFail(t, "/.*/ stop_unless_found # bleep")
	parseFail(t, "/.*/ stop_unless_found blah # bleep")

	// with invalid params
	parseFail(t, "/.*/ stop_unless_found /")
	parseFail(t, "/.*/ stop_unless_found //")
	parseFail(t, "/.*/ stop_unless_found ///")
	parseFail(t, "/.*/ stop_unless_found /*/")
	parseFail(t, "/.*/ stop_unless_found /?/")
	parseFail(t, "/.*/ stop_unless_found /+/")
	parseFail(t, "/.*/ stop_unless_found /(/")
	parseFail(t, "/.*/ stop_unless_found /(.*/")
	parseFail(t, "/.*/ stop_unless_found /(.*))/")

	// with valid regexes
	parseOk(t, "%response_body% stop_unless_found %<!--DO_LOGGING-->%",
		"stop_unless_found", "^response_body$", "<!--DO_LOGGING-->", nil)
	parseOk(t, "/response_body/ stop_unless_found /<!--DO_LOGGING-->/",
		"stop_unless_found", "^response_body$", "<!--DO_LOGGING-->", nil)

	// with valid regexes and escape sequences
	parseOk(t, "!request_body|response_body! stop_unless_found |<!--DO_LOGGING-->\\|<!-NOSKIP-->|",
		"stop_unless_found", "^request_body|response_body$", "<!--DO_LOGGING-->|<!-NOSKIP-->", nil)
	parseOk(t, "!request_body|response_body|boo\\!! stop_unless_found |<!--DO_LOGGING-->\\|<!-NOSKIP-->|",
		"stop_unless_found", "^request_body|response_body|boo!$", "<!--DO_LOGGING-->|<!-NOSKIP-->", nil)
	parseOk(t, "|request_body\\|response_body| stop_unless_found |<!--DO_LOGGING-->\\|<!-NOSKIP-->|",
		"stop_unless_found", "^request_body|response_body$", "<!--DO_LOGGING-->|<!-NOSKIP-->", nil)
	parseOk(t, "|request_body\\|response_body| stop_unless_found |<!--DO_LOGGING-->\\|<!-NOSKIP-->\\|pipe\\||",
		"stop_unless_found", "^request_body|response_body$", "<!--DO_LOGGING-->|<!-NOSKIP-->|pipe|", nil)
	parseOk(t, "/request_body\\/response_body/ stop_unless_found |<!--DO_LOGGING-->\\|<!-NOSKIP-->\\|pipe\\||",
		"stop_unless_found", "^request_body/response_body$", "<!--DO_LOGGING-->|<!-NOSKIP-->|pipe|", nil)
}

func TestReturnsExpectedErrors(t *testing.T) {

	_, err := newHttpRules("file://~/bleepblorpbleepblorp12345")
	assert.Equal(t, "Failed to load rules: ~/bleepblorpbleepblorp12345", err.Error())

	_, err = newHttpRules("/*! stop")
	assert.Equal(t, "Invalid expression (/*!) in rule: /*! stop", err.Error())

	_, err = newHttpRules("/*/ stop")
	assert.Equal(t, "Invalid expression (/*!) in rule: /*! stop", err.Error())

	_, err = newHttpRules("/boo")
	assert.Equal(t, "Invalid rule: /boo", err.Error())

	_, err = newHttpRules("sample 123")
	assert.Equal(t, "Invalid sample percent: 123", err.Error())

	_, err = newHttpRules("!!! stop")
	assert.Equal(t, "Unescaped separator (!) in rule: !!! stop", err.Error())

}
