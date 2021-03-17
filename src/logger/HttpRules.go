package logger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var httpRules *HttpRules

type HttpRules struct {
	debugRules        string
	standardRules     string
	strictRules       string
	defaultRules      string
	allowHttpUrl      bool
	copySessionField  []*HttpRule
	remove            []*HttpRule
	removeIf          []*HttpRule
	removeIfFound     []*HttpRule
	removeUnless      []*HttpRule
	removeUnlessFound []*HttpRule
	replace           []*HttpRule
	sample            []*HttpRule
	skipCompression   bool
	skipSubmission    bool
	size              int
	stop              []*HttpRule
	stopIf            []*HttpRule
	stopIfFound       []*HttpRule
	stopUnless        []*HttpRule
	stopUnlessFound   []*HttpRule
	text              string
}

func GetHttpRules() *HttpRules {

}

// !!! this will need to return an error as well !!!
func NewHttpRules(rules string) *HttpRules {
	httpRules := GetHttpRules()
	if rules == "" {
		rules = httpRules.getDefaultRules()
	}

	//load rules from external files
	if strings.HasPrefix(rules, "file://") {
		// obtain file name
		rfile := strings.TrimPrefix(rules, "file://")
		rfile = strings.TrimSpace(rfile)

		// read rules from file
		buffer, err := ioutil.ReadFile(rfile)
		if err != nil {
			// error handling
		}
		rules = string(buffer)
	}

	//force default rules if necessary
	regex := regexp.MustCompile("(?m)^\\s*include default\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.defaultRules)
	if len(strings.TrimSpace(rules)) == 0 {
		rules = httpRules.defaultRules
	}

	//expand rule includes

	//include debug rules
	regex = regexp.MustCompile("(?m)^\\s*include debug\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.debugRules)
	// include standard rules
	regex = regexp.MustCompile("(?m)^\\s*include standard\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.standardRules)
	// include strict rules
	regex = regexp.MustCompile("(?m)^\\s*include strict\\s*$")
	rules = regex.ReplaceAllString(rules, httpRules.strictRules)

	_text := rules

	// parse all rules
	// Not sure about this part
	var prs []*HttpRule
	for _, rule := range regexp.MustCompile("\\r?\\n").Split(_text, -1) {
		parsed, err := parseRule(rule)
		if err == nil {
			prs = append(prs, parsed)
		}
	}

	_size := len(prs)

	_debugRules := "allow_http_url\ncopy_session_field /.*/\n"

	_standardRules := "/request_header:cookie|response_header:set-cookie/remove\n" +
		"/(request|response)_body|request_param/ replace /[a-zA-Z0-9.!#$%&’*+\\/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)/, /x@y.com/\n" +
		"/request_body|request_param|response_body/ replace /[0-9\\.\\-\\/]{9,}/, /xyxy/\n"

	_strictRules := "/request_url/ replace /([^\\?;]+).*/, /$1/\n" +
		"/request_body|response_body|request_param:.*|request_header:(?!user-agent).*|response_header:(?!(content-length)|(content-type)).*/ remove\n"

	_defaultRules := _strictRules

	// break out rules by verb
	_allowHttpUrl := len(ruleFilter(prs, "allow_http_url", ruleCompare)) > 0
	_copySessionField := ruleFilter(prs, "copy_session_field", ruleCompare)
	_remove := ruleFilter(prs, "remove", ruleCompare)
	_removeIf := ruleFilter(prs, "remove_if", ruleCompare)
	_removeIfFound := ruleFilter(prs, "remove_if_found", ruleCompare)
	_removeUnless := ruleFilter(prs, "remove_unless", ruleCompare)
	_removeUnlessFound := ruleFilter(prs, "remove_unless_found", ruleCompare)
	_replace := ruleFilter(prs, "replace", ruleCompare)
	_sample := ruleFilter(prs, "sample", ruleCompare)
	_skipCompression := len(ruleFilter(prs, "skip_compression", ruleCompare)) > 0
	_skipSubmission := len(ruleFilter(prs, "skip_submission", ruleCompare)) > 0
	_stop := ruleFilter(prs, "stop", ruleCompare)
	_stopIf := ruleFilter(prs, "stop_if", ruleCompare)
	_stopIfFound := ruleFilter(prs, "stop_if_found", ruleCompare)
	_stopUnless := ruleFilter(prs, "stop_unless", ruleCompare)
	_stopUnlessFound := ruleFilter(prs, "stop_unless_found", ruleCompare)

	if len(_sample) > 1 {
		// return error "Multiple sample rules"
	}

	return &HttpRules{
		debugRules:        _debugRules,
		standardRules:     _standardRules,
		strictRules:       _standardRules,
		defaultRules:      _defaultRules,
		allowHttpUrl:      _allowHttpUrl,
		copySessionField:  _copySessionField,
		remove:            _remove,
		removeIf:          _removeIf,
		removeIfFound:     _removeIfFound,
		removeUnless:      _removeUnless,
		removeUnlessFound: _removeUnlessFound,
		replace:           _replace,
		sample:            _sample,
		skipCompression:   _skipCompression,
		skipSubmission:    _skipSubmission,
		size:              _size,
		stop:              _stop,
		stopIf:            _stopIf,
		stopIfFound:       _stopIfFound,
		stopUnless:        _stopUnless,
		stopUnlessFound:   _stopUnlessFound,
		text:              _text,
	}
}

func (rules *HttpRules) getDefaultRules() string {
	return rules.defaultRules
}

func (rules *HttpRules) setDefaultRules(r string) {
	regex := regexp.MustCompile("(?m)^\\s*include default\\s*$")
	rules.defaultRules = regex.ReplaceAllString(r, "")
}

func (rules *HttpRules) getDebugRules() string {
	return rules.debugRules
}

func (rules *HttpRules) getStandardRules() string {
	return rules.standardRules
}

func (rules *HttpRules) getStrictRules() string {
	return rules.strictRules
}

// parse rule from single line
// !!! will return error eventually !!!
func parseRule(r string) (*HttpRule, error) {
	if r == "" || regexBlankOrComment.MatchString(r) {
		return nil, errors.New("Blank rule or comment")
	}
	if regexAllowHttpUrl.MatchString(r) {
		return NewHttpRule("allow_http_url", nil, nil, nil), nil
	}
	if regexCopySessionField.MatchString(r) {
		return NewHttpRule("copy_session_field", nil, nil, nil), nil
	}
	m := regexRemove.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("remove", parseRegex(r, m[0][1]), nil, nil), nil
	}
	m = regexRemoveIf.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("remove_if", parseRegex(r, m[0][1]), parseRegex(r, m[0][2]), nil), nil
	}
	m = regexRemoveIfFound.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("remove_if_found", parseRegex(r, m[0][1]), parseRegexFind(r, m[0][2]), nil), nil
	}
	m = regexRemoveUnless.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("remove_uless", parseRegex(r, m[0][1]), parseRegex(r, m[0][2]), nil), nil
	}
	m = regexRemoveUnlessFound.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("remove_uless_found", parseRegex(r, m[0][1]), parseRegexFind(r, m[0][2]), nil), nil
	}
	m = regexReplace.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("replace", parseRegex(r, m[0][1]), parseRegexFind(r, m[0][2]), parseString(r, m[0][3])), nil
	}
	m = regexSample.FindAllStringSubmatch(r, -1)
	if m != nil {
		m1, err := strconv.Atoi(m[0][1])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error parsing sample rule: %s", r))
		}
		if m1 < 1 || m1 > 99 {
			return nil, errors.New(fmt.Sprintf("Invalid sample percent: %d", m1))
		}
		return NewHttpRule("sample", nil, m1, nil), nil
	}
	m = regexSkipCompression.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("skip_compression", nil, nil, nil), nil
	}
	m = regexSkipSubmission.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("skip_submission", nil, nil, nil), nil
	}
	m = regexStop.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("stop", parseRegex(r, m[0][1]), nil, nil), nil
	}
	m = regexStopIf.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("stop_if", parseRegex(r, m[0][1]), parseRegex(r, m[0][2]), nil), nil
	}
	m = regexStopIfFound.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("stop_if_found", parseRegex(r, m[0][1]), parseRegexFind(r, m[0][2]), nil), nil
	}
	m = regexStopUnless.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("stop_unless", parseRegex(r, m[0][1]), parseRegex(r, m[0][2]), nil), nil
	}
	m = regexStopUnlessFound.FindAllStringSubmatch(r, -1)
	if m != nil {
		return NewHttpRule("stop_unless_found", parseRegex(r, m[0][1]), parseRegexFind(r, m[0][2]), nil), nil
	}
	return nil, errors.New(fmt.Sprintf("Invalid rule: %s", r))
}

// Parses regex for finding.
// !!! will return error eventually !!!
func parseRegex(r string, regex string) *regexp.Regexp /* error */ {
	s := parseString(r, regex)
	if "*" == s || "+" == s || "?" == s {
		// return error '"Invalid regex (%s) in rule: %s", regex, r'
	}
	if strings.HasPrefix(s, "^") {
		s = "^" + s
	}
	if strings.HasSuffix(s, "$") {
		s = s + "$"
	}

	regexp, err := regexp.Compile(s)
	if err != nil {
		// return error '"Invalid regex (%s) in rule: %s", regex, r'
	}
	return regexp
}

// Parses regex for finding.
// !!! will return error eventually !!!
func parseRegexFind(r string, regex string) *regexp.Regexp /* error */ {
	regexp, err := regexp.Compile(parseString(r, regex))
	if err != nil {
		// return error '"Invalid regex (%s) in rule: %s", regex, r'
	}
	return regexp
}

// Parses delimited string expression
// !!! will return an error eventually !!!
func parseString(r string, expr string) string /* error */ {
	separators := []string{"~", "!", "%", "|", "/"}
	for _, sep := range separators {
		regex := regexp.MustCompile(fmt.Sprintf("^[%s](.*)[%s]$", sep, sep))
		m := regex.FindAllStringSubmatch(expr, -1)
		if m != nil {
			m1 := m[0][1]
			regex = regexp.MustCompile(fmt.Sprintf("^[%s].*|.*[^\\\\][%s].*", sep, sep))
			if regex.MatchString(m1) {
				// return error '"Unescaped separator (%s) in rule: %s", sep, r'
			}
			return strings.Replace(m1, "\\"+sep, sep, -1)
		}
	}
	return "" /* error: '"Invalid expression (%s) in rule: %s", expr, r' */
}

/*
The following unexported Regexps should be treat as constants
and remain unchanged throughout package usage
*/
var regexAllowHttpUrl *regexp.Regexp = regexp.MustCompile(`^\s*allow_http_url\s*(#.*)?$`)
var regexBlankOrComment *regexp.Regexp = regexp.MustCompile(`^\s*([#].*)*`)
var regexCopySessionField *regexp.Regexp = regexp.MustCompile(`^\s*copy_session_field\s+([~!%|\/].+[~!%|\/])\s*(#.*)?`)
var regexRemove *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove\s*(#.*)?`)
var regexRemoveIf *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove_if\s+([~!%|\/].+[~!%|\/])\s*(#.*)?`)
var regexRemoveIfFound *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove_if_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexRemoveUnless *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*remove_unless\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexRemoveUnlessFound *regexp.Regexp = regexp.MustCompile(`/^\s*([~!%|\/].+[~!%|\/])\s*remove_unless_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexReplace *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*replace[\s]+([~!%|\/].+[~!%|\/]),[\s]+([~!%|\/].*[~!%|\/])\s*(#.*)?$`)
var regexSample *regexp.Regexp = regexp.MustCompile(`^\s*sample\s+(\d+)\s*(#.*)?$`)
var regexSkipCompression *regexp.Regexp = regexp.MustCompile(`^\s*skip_compression\s*(#.*)?$`)
var regexSkipSubmission *regexp.Regexp = regexp.MustCompile(`^\s*skip_submission\s*(#.*)?$`)
var regexStop *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop\s*(#.*)?$`)
var regexStopIf *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_if\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexStopIfFound *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_if_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexStopUnless *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_unless\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)
var regexStopUnlessFound *regexp.Regexp = regexp.MustCompile(`^\s*([~!%|\/].+[~!%|\/])\s*stop_unless_found\s+([~!%|\/].+[~!%|\/])\s*(#.*)?$`)

/*
used in ruleFilter method to compare given rule string
and a parse rule's verb
*/
func ruleCompare(ruleString string, s string) bool {
	return ruleString == s
}

// filter a slice of HttpRules comparing a given rule string with an HttpRule's verb
func ruleFilter(parsedRules []*HttpRule, ruleString string, cond func(string, string) bool) []*HttpRule {
	result := []*HttpRule{}
	for _, rule := range parsedRules {
		if cond(ruleString, rule.verb) {
			result = append(result, rule)
		}
	}
	return result
}
